package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/shawn1912/messages-service/database"
)

func TestMain(m *testing.M) {
	// Initialize the test database
	database.InitDB("user=postgres password=postgres dbname=messages_test sslmode=disable")
	defer database.DB.Close()

	// Run tests
	code := m.Run()

	// Clean up the test database
	database.DB.Exec("TRUNCATE TABLE messages RESTART IDENTITY CASCADE;")

	os.Exit(code)
}

func setupTestDatabase() {
	// Clean the messages table before each test
	database.DB.Exec("TRUNCATE TABLE messages RESTART IDENTITY CASCADE;")
}

func TestCreateMessageRoute(t *testing.T) {
	setupTestDatabase()

	router := setupRouter()

	// Prepare the request
	payload := map[string]string{"content": "Racecar"}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/message", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	// Parse the response body
	var resp database.Message
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the response
	if resp.Content != "Racecar" {
		t.Errorf("Expected content 'Racecar', got '%s'", resp.Content)
	}
	if !resp.IsPalindrome {
		t.Error("Expected IsPalindrome to be true")
	}
}
