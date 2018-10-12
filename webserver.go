package main

import (
	"encoding/json"
	"log"
	"net/http"

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

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/login.html")
	})

	router.HandleFunc("/_manage", func(w http.ResponseWriter, r *http.Request) {
		LinkHandleManagePage(w, r, config)
	})

	router.HandleFunc("/deleteShortUrl", func(w http.ResponseWriter, r *http.Request) {
		tokenhash := r.FormValue("tokenhash")
		shortlink := r.FormValue("shortlink")
		frommanagepage := r.FormValue("frommanagepage")
		if tokenhash != GetSHA256Hash(config.CreationToken) {
			WriteError(w, 403, "unauthorized")
			return
		}
		_, err := mysql.Query("DELETE FROM shortlinks WHERE shortlink = ?", shortlink)
		if err != nil {
			WriteError(w, 500, err.Error())
			return
		}
		if frommanagepage == "1" {
			r.Form.Set("token", config.CreationToken)
			LinkHandleManagePage(w, r, config)
			return
		}
		WriteResponse(w, 200, "removed", nil)
	})

	router.HandleFunc("/createShortUrl", func(w http.ResponseWriter, r *http.Request) {
		LinkHandlerCreateRequest(w, r, mysql, config)
	})

	router.HandleFunc("/{shortlink}", func(w http.ResponseWriter, r *http.Request) {
		LinkHandler(w, r, mysql)
	})

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets/static"))))

	if certEnabled {
		log.Printf("Running Webserver in TLS mode on port %s", config.Port)
		return http.ListenAndServeTLS(":"+config.Port, config.Cert.CertFile, config.Cert.KeyFile, nil)
	} else {
		log.Printf("Running Webserver in NON TLS mode on port %s", config.Port)
		return http.ListenAndServe(":"+config.Port, nil)
	}
}
