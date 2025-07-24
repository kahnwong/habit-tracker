package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func Foo() {
	fmt.Println("foo")
}

// Define expected schemas for all tables
var tableSchemas = map[string]string{
	"users": `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`,
	"habit": `
	CREATE TABLE habit (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`,
	"activity": `
	CREATE TABLE activity (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,  -- YYYY-MM-DD format (e.g., '2023-10-27')
		is_completed INTEGER NOT NULL, -- 0 for false, 1 for true (boolean)
		habit_name TEXT NOT NULL,
		FOREIGN KEY (habit_name) REFERENCES habit(name) ON DELETE CASCADE
	);
	CREATE INDEX idx_activity_habit_name ON activity (habit_name);
	CREATE INDEX idx_activity_date ON activity (date);`,
}

// Define expected column definitions for schema validation for each table
var allExpectedColumns = map[string]map[string]string{
	"users": {
		"id":    "INTEGER",
		"name":  "TEXT",
		"email": "TEXT",
	},
	"habit": {
		"id":   "INTEGER",
		"name": "TEXT",
	},
	"activity": {
		"id":           "INTEGER",
		"date":         "TEXT",
		"is_completed": "INTEGER",
		"habit_name":   "TEXT",
	},
}

func init() {
	// Define the database file name
	dbFileName := "habits.sqlite"

	// Check if the database file exists
	dbExists := true
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		dbExists = false
		log.Printf("Database file '%s' not found. It will be created.", dbFileName)
	} else if err != nil {
		// Handle other potential errors when checking file status
		log.Fatalf("Error checking database file status: %v", err)
	}

	// Open the database connection using sqlx.Connect
	// If the file does not exist, the sqlite3 driver will create it upon connection.
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

	// Iterate through all defined table schemas
	for tableName, schema := range tableSchemas {
		if !dbExists {
			// Database file did not exist, so create the table
			log.Printf("Creating table '%s'...", tableName)
			_, err = db.Exec(schema)
			if err != nil {
				log.Fatalf("Error creating table '%s': %v", tableName, err)
			}
			log.Printf("Table '%s' created successfully!", tableName)
		} else {
			// Database file existed, validate its schema
			log.Printf("Database file '%s' found. Validating schema for table '%s'...", dbFileName, tableName)
			expectedCols, ok := allExpectedColumns[tableName]
			if !ok {
				log.Printf("Warning: No expected column definition found for table '%s'. Skipping schema validation for this table.", tableName)
				continue
			}
			if err := validateSchema(db, tableName, expectedCols); err != nil {
				log.Fatalf("Schema validation failed for table '%s': %v", tableName, err)
			}
			log.Printf("Schema for table '%s' validated successfully.", tableName)
		}
	}
	log.Println("All tables processed successfully.")
}

// tableExists checks if a given table exists in the database.
func tableExists(db *sqlx.DB, tableName string) (bool, error) {
	var count int
	// Query sqlite_master to check for the table's existence
	query := `SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?`
	err := db.Get(&count, query, tableName)
	if err != nil {
		return false, fmt.Errorf("error checking if table '%s' exists: %w", tableName, err)
	}
	return count > 0, nil
}

// validateSchema validates the schema of a specified table against a map of expected columns.
// It checks if the table exists and if it has the expected columns with correct types.
func validateSchema(db *sqlx.DB, tableName string, expectedColumns map[string]string) error {
	// First, check if the table itself exists
	exists, err := tableExists(db, tableName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("table '%s' does not exist in the database", tableName)
	}

	// Query table info using PRAGMA to get column details
	rows, err := db.Queryx(fmt.Sprintf("PRAGMA table_info(%s);", tableName))
	if err != nil {
		return fmt.Errorf("error querying table info for '%s': %w", tableName, err)
	}
	defer rows.Close()

	// Map to store found columns and their types
	foundColumns := make(map[string]string)
	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string // column type (e.g., TEXT, INTEGER)
			notnull    int
			dflt_value sql.NullString // Default value, can be NULL
			pk         int            // Primary key flag
		)
		// Scan the results from PRAGMA table_info
		if err := rows.Scan(&cid, &name, &columnType, &notnull, &dflt_value, &pk); err != nil {
			return fmt.Errorf("error scanning table info row: %w", err)
		}
		foundColumns[name] = columnType
	}

	// Validate each expected column against the found columns
	for colName, expectedType := range expectedColumns {
		foundType, ok := foundColumns[colName]
		if !ok {
			return fmt.Errorf("missing expected column: '%s'", colName)
		}
		// For simplicity, we'll check for an exact type match.
		// SQLite's type affinity can sometimes return slightly different names
		// (e.g., VARCHAR instead of TEXT), but for basic types, this is usually sufficient.
		if foundType != expectedType {
			return fmt.Errorf("column '%s' has unexpected type: expected '%s', got '%s'", colName, expectedType, foundType)
		}
	}

	// Optionally, you might want to check for extra columns not in expectedColumns,
	// but for now, we only ensure all expected columns are present and correct.

	return nil // Schema is valid
}
