package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shawn1912/messages-service/utils"
)

// CreateMessage creates a new message.
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var msg Message
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(msg.Content) > 1000 {
		http.Error(w, "Message content exceeds 1000 characters", http.StatusBadRequest)
		return
	}

	msg.IsPalindrome = utils.IsPalindrome(msg.Content)

	err = DB.QueryRow(
		"INSERT INTO messages (content, is_palindrome) VALUES ($1, $2) RETURNING id",
		msg.Content, msg.IsPalindrome).Scan(&msg.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// GetMessage retrieves a message by its ID.
func GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	var msg Message
	err = DB.QueryRow("SELECT id, content, is_palindrome FROM messages WHERE id = $1", id).Scan(&msg.ID, &msg.Content, &msg.IsPalindrome)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}
