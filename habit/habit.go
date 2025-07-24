package habit

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Habit *Application

type Application struct {
	DB *sqlx.DB
}

func (Habit *Application) CreateHabit(habit string) error {
	query := `INSERT INTO habit (name) VALUES (?)`
	_, err := Habit.DB.Exec(query, habit)
	if err != nil {
		return fmt.Errorf("error inserting habit '%s': %w", habit, err)
	}

	return nil
}
