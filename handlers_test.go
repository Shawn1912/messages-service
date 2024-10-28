package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Set up the test database connection
	// TODO: Set up environment variables
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=messages_test sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to test database: %v", err)
	}
	defer testDB.Close()

	// Ensure the database is accessible
	if err = testDB.Ping(); err != nil {
		log.Fatalf("Error pinging test database: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Clean up the test database
	teardownTestDatabase()

	os.Exit(code)
}

func teardownTestDatabase() {
	testDB.Exec("TRUNCATE TABLE messages RESTART IDENTITY CASCADE;")
}

// Tests POST /message
func TestCreateMessage(t *testing.T) {
	// Prepare the request body
	payload := map[string]string{"content": "Racecar"}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/messages", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Use httptest to record the response
	rr := httptest.NewRecorder()

	// Set up the router and handler
	router := mux.NewRouter()
	router.HandleFunc("/messages", CreateMessage).Methods("POST")

	// Assign the test database to the global DB variable
	DB = testDB

	// Call the handler
	router.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	var respMsg Message
	err = json.Unmarshal(rr.Body.Bytes(), &respMsg)
	if err != nil {
		t.Fatal(err)
	}

	if respMsg.Content != "Racecar" {
		t.Errorf("Expected content 'Racecar', got '%s'", respMsg.Content)
	}
	if !respMsg.IsPalindrome {
		t.Error("Expected IsPalindrome to be true")
	}
	if respMsg.ID == 0 {
		t.Error("Expected a valid ID")
	}
}

// Tests GET /message/{id}
func TestGetMessage(t *testing.T) {
	// Insert a test message into the test database
	var msgID int64
	err := testDB.QueryRow(
		"INSERT INTO messages (content, is_palindrome, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id",
		"Madam", true).Scan(&msgID)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the request
	req, err := http.NewRequest("GET", "/messages/"+strconv.FormatInt(msgID, 10), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Use httptest to record the response
	rr := httptest.NewRecorder()

	// Set up the router and handler
	router := mux.NewRouter()
	router.HandleFunc("/messages/{id:[0-9]+}", GetMessage).Methods("GET")

	// Assign the test database to the global DB variable
	DB = testDB

	// Call the handler
	router.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var respMsg Message
	err = json.Unmarshal(rr.Body.Bytes(), &respMsg)
	if err != nil {
		t.Fatal(err)
	}

	if respMsg.ID != msgID {
		t.Errorf("Expected ID %d, got %d", msgID, respMsg.ID)
	}
	if respMsg.Content != "Madam" {
		t.Errorf("Expected content 'Madam', got '%s'", respMsg.Content)
	}
	if !respMsg.IsPalindrome {
		t.Error("Expected IsPalindrome to be true")
	}
}

// Tests GET /messages?limit={}&page={}
func TestListMessages(t *testing.T) {
	// Clean the database before the test starts
	teardownTestDatabase()

	// Insert multiple test messages into the test database
	for i := 1; i <= 25; i++ {
		content := fmt.Sprintf("Test message %d", i)
		isPalindrome := false
		if i%5 == 0 {
			content = "Madam"
			isPalindrome = true
		}
		_, err := testDB.Exec(
			"INSERT INTO messages (content, is_palindrome, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())",
			content, isPalindrome)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Prepare the request with pagination parameters
	req, err := http.NewRequest("GET", "/messages?page=2&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Use httptest to record the response
	rr := httptest.NewRecorder()

	// Set up the router and handler
	router := mux.NewRouter()
	router.HandleFunc("/messages", ListMessages).Methods("GET")

	// Assign the test database to the global DB variable
	DB = testDB

	// Call the handler
	router.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Parse the response
	var response struct {
		Messages   []Message `json:"messages"`
		Pagination struct {
			CurrentPage   int `json:"currentPage"`
			PageSize      int `json:"pageSize"`
			TotalPages    int `json:"totalPages"`
			TotalMessages int `json:"totalMessages"`
		} `json:"pagination"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Validate pagination metadata
	if response.Pagination.CurrentPage != 2 {
		t.Errorf("Expected CurrentPage 2, got %d", response.Pagination.CurrentPage)
	}
	if response.Pagination.PageSize != 10 {
		t.Errorf("Expected PageSize 10, got %d", response.Pagination.PageSize)
	}
	if response.Pagination.TotalPages != 3 {
		t.Errorf("Expected TotalPages 3, got %d", response.Pagination.TotalPages)
	}
	if response.Pagination.TotalMessages != 25 {
		t.Errorf("Expected TotalMessages 25, got %d", response.Pagination.TotalMessages)
	}

	// Validate messages count
	if len(response.Messages) != 10 {
		t.Errorf("Expected 10 messages, got %d", len(response.Messages))
	}
}
