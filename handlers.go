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

	query := `
        INSERT INTO messages (content, is_palindrome)
        VALUES ($1, $2)
        RETURNING id, created_at, updated_at
    `

	err = DB.QueryRow(query, msg.Content, msg.IsPalindrome).
		Scan(&msg.ID, &msg.CreatedAt, &msg.UpdatedAt)
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

	query := `
        SELECT id, content, is_palindrome, created_at, updated_at 
		FROM messages 
		WHERE id = $1
    `

	err = DB.QueryRow(query, id).
		Scan(&msg.ID, &msg.Content, &msg.IsPalindrome, &msg.CreatedAt, &msg.UpdatedAt)
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

// UpdateMessage updates an existing message by its ID.
func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Retrieve existing message from the database
	var existingMsg Message
	err = DB.QueryRow("SELECT id, content, is_palindrome, created_at, updated_at FROM messages WHERE id = $1", id).
		Scan(&existingMsg.ID, &existingMsg.Content, &existingMsg.IsPalindrome, &existingMsg.CreatedAt, &existingMsg.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Read and parse the request body
	var msgUpdates struct {
		Content *string `json:"text"`
	}
	err = json.NewDecoder(r.Body).Decode(&msgUpdates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update fields if they are provided
	if msgUpdates.Content != nil {
		if len(*msgUpdates.Content) > 1000 {
			http.Error(w, "Message content exceeds 1000 characters", http.StatusBadRequest)
			return
		}
		existingMsg.Content = *msgUpdates.Content
		existingMsg.IsPalindrome = utils.IsPalindrome(existingMsg.Content)
	}

	// Update the message in the database
	query := `
        UPDATE messages
        SET content = $1, is_palindrome = $2, updated_at = NOW()
        WHERE id = $3
        RETURNING created_at, updated_at
    `

	err = DB.QueryRow(query, existingMsg.Content, existingMsg.IsPalindrome, id).
		Scan(&existingMsg.CreatedAt, &existingMsg.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingMsg)
}

// DeleteMessage deletes a message by its ID.
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	result, err := DB.Exec("DELETE FROM messages WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListMessages returns a list of messages, up to a maximum of 100.
// TODO: Add Pagination
func ListMessages(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, content, is_palindrome, created_at, updated_at FROM messages LIMIT 100"

	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Content, &msg.IsPalindrome, &msg.CreatedAt, &msg.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
