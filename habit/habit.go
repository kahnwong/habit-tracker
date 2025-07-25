package habit

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Habit *Application

type Application struct {
	DB *sqlx.DB
}

type periodActivityRow map[string]interface{} // for periodActivity

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
	lookbackStart := time.Now().AddDate(0, -lookbackMonths, 0)
	var completedActivities []Activity

	query := `
	SELECT date, is_completed, habit_name
	FROM activity
	WHERE is_completed = 1 AND date >= ? AND habit_name = ?
	ORDER BY date;`

	err := Habit.DB.Select(&completedActivities, query, lookbackStart, habitName)
	if err != nil {
		return completedActivities, fmt.Errorf("error fetching activity for habit '%s'", habitName)
	}

	return completedActivities, nil
}

func (Habit *Application) GetPeriodActivity(period string) ([]periodActivityRow, []string, error) {
	var lookbackStart time.Time
	var now = time.Now()
	var dates []string

	lookbackStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch period {
	case "today":
		dates = []string{lookbackStart.Format("2006-01-02")}
	case "week":
		lookbackStart = lookbackStart.AddDate(0, 0, -7)
		for date := lookbackStart; !date.After(now); date = date.AddDate(0, 0, 1) {
			dates = append(dates, date.Format("2006-01-02"))
		}
	}

	// prep query
	selectClauses := []string{"h.name AS habit_name"}
	for _, date := range dates {
		selectClauses = append(selectClauses, fmt.Sprintf("SUM(CASE WHEN a.date = '%s' THEN a.is_completed ELSE 0 END) AS \"%s\"", date, date))
	}
	selectStmt := strings.Join(selectClauses, ",\n    ")

	baseQuery := fmt.Sprintf(`
	SELECT
	   %s
	FROM
	   habit AS h
	LEFT JOIN
		activity AS a ON h.name = a.habit_name AND a.date IN (?)
	GROUP BY
	   h.name
	ORDER BY
	   h.name;`, selectStmt)

	query, args, err := sqlx.In(baseQuery, dates)
	if err != nil {
		return nil, dates, fmt.Errorf("error preparing query with sqlx.In: %w", err)
	}

	// execute query
	query = Habit.DB.Rebind(query)
	rows, err := Habit.DB.Queryx(query, args...)
	if err != nil {
		return nil, dates, fmt.Errorf("error executing query: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing rows")
		}
	}(rows)

	// parse
	var completedActivities []periodActivityRow
	for rows.Next() {
		row := make(periodActivityRow)
		err = rows.MapScan(row)
		if err != nil {
			return nil, dates, fmt.Errorf("error scanning row: %w", err)
		}

		completedActivities = append(completedActivities, row)
	}

	return completedActivities, dates, nil
}
