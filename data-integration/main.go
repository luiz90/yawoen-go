package main

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	//util := Util{}
	//util.loadCsvData()
}

func main() {

	router := NewRouter()
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	log.Fatal(http.ListenAndServe(":9000",
		handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
