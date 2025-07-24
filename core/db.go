package core

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func Foo() {
	fmt.Println("foo")
}

func init() {
	// Define the database file name
	dbFileName := "habits.sqlite"

	// Remove the database file if it already exists for a clean start
	if err := os.Remove(dbFileName); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing database file: %v", err)
	}
	log.Printf("Removed existing database file (if any): %s", dbFileName)

	// Open the database connection using sqlx.Connect
	// The first argument is the driver name ("sqlite3")
	// The second argument is the data source name (the file path for SQLite)
	db, err := sqlx.Connect("sqlite3", dbFileName)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	// Ensure the database connection is closed when the main function exits
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	log.Println("Successfully connected to SQLite database.")

	// Define the SQL statement to create the 'users' table
	// This table will have an 'id' (primary key, auto-increment), 'name', and 'email'
	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`

	// Execute the schema to create the table
	// db.ExecContext is generally preferred for production, but for simple
	// schema creation, db.Exec is sufficient.
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Table 'users' created successfully!")

	// You can optionally verify the table exists by trying to insert data
	// or querying schema information, but for this example, we'll keep it simple.
}
