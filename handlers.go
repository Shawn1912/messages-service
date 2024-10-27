package main

import (
	"encoding/json"
	"io"
	"net/http"

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
