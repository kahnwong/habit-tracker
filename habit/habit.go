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

// [TODO] undo activity (today only)

func (Habit *Application) CreateHabit(habit string) error {
	query := `INSERT INTO habit (name) VALUES (?)`
	_, err := Habit.DB.Exec(query, habit)
	if err != nil {
		return fmt.Errorf("error inserting habit '%s': %w", habit, err)
	}

	return nil
}

func (Habit *Application) GetHabits() ([]string, error) {
	var habits []string
	query := "SELECT name FROM habit"

	err := Habit.DB.Select(&habits, query)
	if err != nil {
		return habits, fmt.Errorf("error fetching habits")
	}

	return habits, nil
}

func (Habit *Application) Do(activity Activity) error {
	query := `INSERT OR IGNORE INTO activity (date, is_completed, habit_name) VALUES (?, ?, ?)`
	_, err := Habit.DB.Exec(query, activity.Date, activity.IsCompleted, activity.HabitName)
	if err != nil {
		return fmt.Errorf("error inserting activity for habit '%s' on '%s': %w", activity.HabitName, activity.Date, err)
	}

	return nil
}
