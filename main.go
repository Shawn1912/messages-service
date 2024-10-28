package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shawn1912/messages-service/database"
	"github.com/shawn1912/messages-service/handlers"
)

func main() {
	// TODO: Use environment variables
	database.InitDB("user=postgres password=postgres dbname=messages sslmode=disable")

	router := mux.NewRouter()

	router.HandleFunc("/message", handlers.CreateMessage).Methods("POST")
	router.HandleFunc("/message/{id:[0-9]+}", handlers.GetMessage).Methods("GET")
	router.HandleFunc("/message/{id:[0-9]+}", handlers.UpdateMessage).Methods("PATCH")
	router.HandleFunc("/message/{id:[0-9]+}", handlers.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/messages", handlers.ListMessages).Methods("GET")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
