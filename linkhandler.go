package main

import (
	"database/sql"
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

func LinkHandlerCreateRequest(w http.ResponseWriter, r *http.Request, mysql *MySql, creationtoken string, randshortlen int) {
	rooturl := r.FormValue("rooturl")
	shortlink := r.FormValue("shortlink")
	enteredtoken := r.FormValue("creationtoken")

	if enteredtoken != GetSHA256Hash(creationtoken) {
		WriteError(w, 401, "unauthorized")
		return
	}

	if shortlink == "" {
		shortlink = RandomString(randshortlen)
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
		WriteResponse(w, 201, "created shortlink", map[string]string{
			"mode":     "CREATED",
			"rooturl":  rooturl,
			"shorturl": fmt.Sprintf("%s//%s/%s", r.URL.Scheme, r.URL.Host, shortlink),
		})
	}
}
