package main

import (
	"encoding/json"
	"log"
	"net/http"
	"html/template"

	"github.com/gorilla/mux"
)

var mysql *MySql
var filePaths []string

type WebServerCert struct {
	CertFile string `yaml:"certfile"`
	KeyFile  string `yaml:"keyfile"`
}

type WebServerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, code int, message string) {
	err := &WebServerError{code, message}
	bdata, _ := json.MarshalIndent(err, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bdata)
}

func WriteResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	bdata, _ := json.MarshalIndent(&map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bdata)
}

func OpenWebServer(config *Config, mysql *MySql) error {
	certEnabled := config.Cert.KeyFile != "" && config.Cert.CertFile != ""

	router := mux.NewRouter()

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	// router.Methods("GET").HandlerFunc(fileServerHandler)

	router.HandleFunc("/_create", func(w http.ResponseWriter, r *http.Request) {
		LinkHandlerCreate(w, r)
	})

	router.HandleFunc("/_manage", func(w http.ResponseWriter, r *http.Request) {
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

		t := template.New("manage.html")
		t, _ = t.ParseFiles("./assets/manage.html")
		t.Execute(w, struct{
			ShortLinks []*ShortLink
		}{
			ShortLinks: shortLinks,
		})
	})

	router.HandleFunc("/createShortUrl", func(w http.ResponseWriter, r *http.Request) {
		LinkHandlerCreateRequest(w, r, mysql, config.CreationToken, config.RandShortLen)
	})

	router.HandleFunc("/{shortlink}", func(w http.ResponseWriter, r *http.Request) {
		LinkHandler(w, r, mysql)
	})

	http.Handle("/", router)

	if certEnabled {
		log.Printf("Running Webserver in TLS mode on port %s", config.Port)
		return http.ListenAndServeTLS(":"+config.Port, config.Cert.CertFile, config.Cert.KeyFile, nil)
	} else {
		log.Printf("Running Webserver in NON TLS mode on port %s", config.Port)
		return http.ListenAndServe(":"+config.Port, nil)
	}
}
