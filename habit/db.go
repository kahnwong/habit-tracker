package habit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
)

var dbFileName = "habits.sqlite" // [TODO] set via config (plain yaml, not sops)

func init() {
	//dbExists := isDBExists()
	_ = isDBExists()

	app := &Application{
		DB: initDB(),
	}

	fmt.Println(app)
}

func isDBExists() bool {
	dbExists := true
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		dbExists = false
		log.Printf("Database file '%s' not found. It will be created.", dbFileName) // [TODO] replace
	} else if err != nil {
		log.Fatalf("Error checking database file status: %v", err) // [TODO] replace
	}
	return dbExists
}

func initDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", dbFileName)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err) // [TODO] replace
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	return db
}
