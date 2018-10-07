package main

import (
	"log"
)

func main() {

	config, err := ConfigOpen("config.yaml")
	if err != nil {
		log.Fatal("ERROR ON CREATING CONFIG: ", err)
	}

	mysql, err = NewMySql(config.MySql)
	if err != nil {
		log.Fatal("ERROR CONNECTING TO MYSQL DATABASE: ", err)
	}
	log.Println("Connected to MySql Database")

	err = OpenWebServer(config, mysql)
	if err != nil {
		log.Fatal("ERROR ON OPENING WEBSERVER: ", err)
	}
}
