package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitDB("user=postgres dbname=messages sslmode=disable")

	router := mux.NewRouter()

	router.HandleFunc("/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/messages/{id:[0-9]+}", GetMessage).Methods("GET")
	router.HandleFunc("/messages/{id:[0-9]+}", UpdateMessage).Methods("PATCH")
	router.HandleFunc("/messages/{id:[0-9]+}", DeleteMessage).Methods("DELETE")
	router.HandleFunc("/messages", ListMessages).Methods("GET")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
