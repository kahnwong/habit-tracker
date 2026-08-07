package habit

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"testing"
	"time"

	cliBase "github.com/kahnwong/cli-base"
	"github.com/kahnwong/habit-tracker/internal/habit/store"
	sqliteBase "github.com/kahnwong/sqlite-base"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Path string `yaml:"PATH"`
}

var config *Config
var dbFileName string

var Habit *Application

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Application struct {
	DB      *sql.DB
	Queries *store.Queries
}

type Activity struct {
	Date        string `db:"date"`
	IsCompleted int    `db:"is_completed"` // 0 for false, 1 for true
	HabitName   string `db:"habit_name"`
}

type gooseErrorLogger struct{}

func (gooseErrorLogger) Printf(string, ...interface{}) {}

func (gooseErrorLogger) Fatalf(format string, v ...interface{}) {
	stdlog.Fatalf(format, v...)
}

type periodActivityRow map[string]interface{} // for periodActivity

func (Habit *Application) CreateHabit(habit string) error {
	err := Habit.Queries.CreateHabit(context.Background(), habit)
	if err != nil {
		return fmt.Errorf("error inserting habit '%s': %w", habit, err)
	}

	return nil
}

func (Habit *Application) GetHabits() ([]string, error) {
	habits, err := Habit.Queries.ListHabits(context.Background())
	if err != nil {
		return habits, fmt.Errorf("error fetching habits")
	}

	return habits, nil
}

func (Habit *Application) Do(activity Activity) error {
	err := Habit.Queries.CreateActivity(context.Background(), store.CreateActivityParams{
		Date:        activity.Date,
		IsCompleted: int64(activity.IsCompleted),
		HabitName:   activity.HabitName,
	})
	if err != nil {
		return fmt.Errorf("error inserting activity for habit '%s' on '%s': %w", activity.HabitName, activity.Date, err)
	}

	return nil
}

func (Habit *Application) Undo(activity Activity) error {
	err := Habit.Queries.DeleteActivity(context.Background(), store.DeleteActivityParams{
		Date:      activity.Date,
		HabitName: activity.HabitName,
	})
	if err != nil {
		return fmt.Errorf("error deleting habit '%s' on '%s': %w", activity.HabitName, activity.Date, err)
	}

	return nil
}

func (Habit *Application) GetHabitActivity(habitName string, lookbackMonths int) ([]Activity, error) {
	lookbackStart := time.Now().AddDate(0, -lookbackMonths, 0)

	rows, err := Habit.Queries.ListCompletedHabitActivities(context.Background(), store.ListCompletedHabitActivitiesParams{
		LookbackStart: lookbackStart.Format("2006-01-02"),
		HabitName:     habitName,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching activity for habit '%s'", habitName)
	}

	completedActivities := make([]Activity, 0, len(rows))
	for _, row := range rows {
		completedActivities = append(completedActivities, Activity{
			Date:        row.Date,
			IsCompleted: int(row.IsCompleted),
			HabitName:   row.HabitName,
		})
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

	placeholders := make([]string, 0, len(dates))
	args := make([]any, 0, len(dates))
	for _, date := range dates {
		placeholders = append(placeholders, "?")
		args = append(args, date)
	}

	query := fmt.Sprintf(`
	SELECT
	   %s
	FROM
	   habit AS h
	LEFT JOIN
		activity AS a ON h.id = a.habit_id AND a.date IN (%s)
	GROUP BY
	   h.name
	ORDER BY
	   h.name;`, selectStmt, strings.Join(placeholders, ", "))

	// execute query
	rows, err := Habit.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, dates, fmt.Errorf("error executing query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing rows")
		}
	}(rows)

	columns, err := rows.Columns()
	if err != nil {
		return nil, dates, fmt.Errorf("error reading columns: %w", err)
	}

	// parse
	var completedActivities []periodActivityRow
	for rows.Next() {
		values := make([]any, len(columns))
		scanArgs := make([]any, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dates, fmt.Errorf("error scanning row: %w", err)
		}

		row := make(periodActivityRow)
		for i, column := range columns {
			switch value := values[i].(type) {
			case []byte:
				row[column] = string(value)
			default:
				row[column] = value
			}
		}

		completedActivities = append(completedActivities, row)
	}
	if err := rows.Err(); err != nil {
		return nil, dates, fmt.Errorf("error reading rows: %w", err)
	}

	return completedActivities, dates, nil
}

func init() {
	// set logs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	goose.SetLogger(gooseErrorLogger{})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("MODE") == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// init config
	var err error
	config, err = cliBase.ReadYaml[Config]("~/.config/habit-tracker/config.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && testing.Testing() {
			return
		}
		log.Fatal().Err(err).Msg("failed to read config")
	}

	dbFileName, err = cliBase.ExpandHome(config.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to expand home path")
	}

	// init app
	db, err := sqliteBase.Open(sqliteBase.Config{
		Path:         dbFileName,
		MigrationDir: "migrations",
		MigrationFS:  migrationFiles,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}

	Habit = &Application{
		DB:      db,
		Queries: store.New(db),
	}
}
