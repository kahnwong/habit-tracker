package habit

import (
	"fmt"
	"time"

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

func (Habit *Application) Undo(activity Activity) error {
	deleteQuery := "DELETE FROM activity WHERE date = ? AND habit_name = ?"
	_, err := Habit.DB.Exec(deleteQuery, activity.Date, activity.HabitName)
	if err != nil {
		return fmt.Errorf("error deleting habit '%s' on '%s': %w", activity.HabitName, activity.Date, err)
	}

	return nil
}

func (Habit *Application) GetHabitActivity(habitName string, lookbackMonths int) ([]Activity, error) {
	query := `
	SELECT date, is_completed, habit_name
	FROM activity
	WHERE is_completed = 1 AND date >= ? AND habit_name = ?
	ORDER BY date;`

	lookbackStart := time.Now().AddDate(0, -lookbackMonths, 0)
	var completedActivities []Activity
	err := Habit.DB.Select(&completedActivities, query, lookbackStart, habitName)
	if err != nil {
		return completedActivities, fmt.Errorf("error fetching activity for habit '%s'", habitName)
	}

	return completedActivities, nil
}

func (Habit *Application) GetPeriodActivity(period string) ([]Activity, error) {
	query := `
	SELECT date, is_completed, habit_name
	FROM activity
	WHERE is_completed = 1 AND date >= ?
	ORDER BY date;`

	var lookbackStart time.Time
	var now = time.Now()
	switch period {
	case "today":
		lookbackStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	lookbackStartStr := lookbackStart.Format("2006-01-02")

	var completedActivities []Activity
	err := Habit.DB.Select(&completedActivities, query, lookbackStartStr)
	if err != nil {
		return completedActivities, fmt.Errorf("error fetching activity for period '%s'", period)
	}

	return completedActivities, nil
}
