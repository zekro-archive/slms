package main

import (
	"database/sql"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func rowsEmpty(rows *sql.Rows) bool {
	status := true
	for rows.Next() {
		status = false
	}
	return status
}

func LinkHandler(w http.ResponseWriter, r *http.Request, mysql *MySql) {
	vars := mux.Vars(r)
	shortlink := vars["shortlink"]

	rows, err := mysql.Query("SELECT rootlink FROM shortlinks WHERE shortlink = ?", shortlink)
	if err != nil {
		WriteError(w, 500, err.Error())
		return
	}
	rootlink := ""
	rows.Next()
	rows.Scan(&rootlink)
	if rootlink == "" {
		http.ServeFile(w, r, "./assets/invalid.html")
	}
	mysql.Query("UPDATE shortlinks SET accesses = accesses + 1, lastaccess = ? WHERE shortlink = ?", time.Now(), shortlink)
	http.Redirect(w, r, rootlink, http.StatusSeeOther)
}

func LinkHandlerCreate(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/createlink.html")
}

func LinkHandleManagePage(w http.ResponseWriter, r *http.Request, config *Config) {
	token := r.FormValue("token")
	if token != config.CreationToken {
		WriteError(w, 403, "unauthorized")
		return
	}

	shortLinks := make([]*ShortLink, 0)

	rows, err := mysql.Query("SELECT * FROM shortlinks")
	if err != nil {
		WriteError(w, 500, err.Error())
	}
	for rows.Next() {
		sl := new(ShortLink)
		rows.Scan(&sl.RootLink, &sl.ShortLink, &sl.Created, &sl.Accesses, &sl.LastAccess)
		shortLinks = append(shortLinks, sl)
	}

	tokenHash := GetSHA256Hash(token)
	t := template.New("manage.html")
	t, _ = t.ParseFiles("./assets/manage.html")
	t.Execute(w, struct {
		ShortLinks []*ShortLink
		TokenHash  string
		AppVersion string
		AppCommit  string
		AppDate    string
	}{
		ShortLinks: shortLinks,
		TokenHash:  tokenHash,
		AppVersion: AppVersion,
		AppCommit:  AppCommit,
		AppDate:    AppDate,
	})
}

func LinkHandlerCreateRequest(w http.ResponseWriter, r *http.Request, mysql *MySql, config *Config) {
	rooturl := r.FormValue("rooturl")
	shortlink := r.FormValue("shortlink")
	enteredtoken := r.FormValue("creationtoken")
	frommanagepage := r.FormValue("frommanagepage")

	if enteredtoken != GetSHA256Hash(config.CreationToken) {
		WriteError(w, 401, "unauthorized")
		return
	}

	for _, bl := range BLOCKED_SHORTLINKS {
		if shortlink == bl {
			WriteError(w, 400, "blacklisted shortlink")
			return
		}
	}

	if shortlink == "" {
		shortlink = RandomString(config.RandShortLen)
	}

	rows, err := mysql.Query("SELECT * FROM shortlinks WHERE shortlink = ?", shortlink)
	if err != nil {
		log.Println("ERROR GETTING DATABASE VALUES: ", err)
		WriteError(w, 500, err.Error())
		return
	}
	if !rowsEmpty(rows) {
		fmt.Println("UPDATE")
		_, err = mysql.Query("UPDATE shortlinks SET rootlink = ?, created = ? WHERE shortlink = ?", rooturl, time.Now(), shortlink)
		if err != nil {
			log.Println("ERROR UPDATING SHORTLINK DB ENTRY: ", err)
			WriteError(w, 500, err.Error())
			return
		}
		if frommanagepage == "1" {
			r.Form.Set("token", config.CreationToken)
			LinkHandleManagePage(w, r, config)
			return
		}
		WriteResponse(w, 201, "updated shortlink", map[string]string{
			"mode":     "UPDATED",
			"rooturl":  rooturl,
			"shorturl": fmt.Sprintf("%s//%s/%s", r.URL.Scheme, r.URL.Host, shortlink),
		})

	} else {
		_, err = mysql.Query("INSERT INTO shortlinks (rootlink, shortlink) VALUES (?, ?)", rooturl, shortlink)
		if err != nil {
			log.Println("ERROR CREATING SHORTLINK DB ENTRY: ", err)
			WriteError(w, 500, err.Error())
			return
		}
		if frommanagepage == "1" {
			r.Form.Set("token", config.CreationToken)
			LinkHandleManagePage(w, r, config)
			return
		}
		WriteResponse(w, 201, "created shortlink", map[string]string{
			"mode":     "CREATED",
			"rooturl":  rooturl,
			"shorturl": fmt.Sprintf("%s//%s/%s", r.URL.Scheme, r.URL.Host, shortlink),
		})
	}
}
