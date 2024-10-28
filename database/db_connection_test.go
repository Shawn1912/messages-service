package database

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDB_Success(t *testing.T) {
	// Use a valid connection string for the test database
	connStr := "user=postgres password=postgres dbname=messages_test sslmode=disable"

	// Call InitDB with the test connection string
	InitDB(connStr)
	defer DB.Close()

	// Check that DB is not nil
	if DB == nil {
		t.Fatal("Expected DB to be initialized, but got nil")
	}

	// Ping the database to ensure the connection is active
	err := DB.Ping()
	if err != nil {
		t.Fatalf("Expected to ping DB successfully, but got error: %v", err)
	}
}
