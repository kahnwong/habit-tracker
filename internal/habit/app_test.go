package habit

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory database for testing
func setupTestDB(t *testing.T) *Application {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Create tables
	for _, schema := range tableSchemas {
		_, err := db.Exec(schema)
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
	}

	return &Application{DB: db}
}

func TestCreateHabit(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	tests := []struct {
		name      string
		habit     string
		wantError bool
	}{
		{"create valid habit", "exercise", false},
		{"create another habit", "reading", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.CreateHabit(tt.habit)
			if (err != nil) != tt.wantError {
				t.Errorf("CreateHabit() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCreateHabitDuplicate(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	// Create first habit
	err := app.CreateHabit("exercise")
	if err != nil {
		t.Fatalf("failed to create first habit: %v", err)
	}

	// Try to create duplicate
	err = app.CreateHabit("exercise")
	if err == nil {
		t.Error("expected error when creating duplicate habit, got nil")
	}
}

func TestGetHabits(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	// Initially empty
	habits, err := app.GetHabits()
	if err != nil {
		t.Fatalf("GetHabits() error = %v", err)
	}
	if len(habits) != 0 {
		t.Errorf("expected 0 habits, got %d", len(habits))
	}

	// Add habits
	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit 'exercise': %v", err)
	}
	if err := app.CreateHabit("reading"); err != nil {
		t.Fatalf("failed to create habit 'reading': %v", err)
	}
	if err := app.CreateHabit("meditation"); err != nil {
		t.Fatalf("failed to create habit 'meditation': %v", err)
	}

	habits, err = app.GetHabits()
	if err != nil {
		t.Fatalf("GetHabits() error = %v", err)
	}
	if len(habits) != 3 {
		t.Errorf("expected 3 habits, got %d", len(habits))
	}
}

func TestDo(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	// Create habit first
	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	activity := Activity{
		Date:        "2024-01-15",
		IsCompleted: 1,
		HabitName:   "exercise",
	}

	err := app.Do(activity)
	if err != nil {
		t.Errorf("Do() error = %v", err)
	}

	// Verify activity was recorded
	var count int
	err = app.DB.Get(&count, "SELECT COUNT(*) FROM activity WHERE habit_name = ? AND date = ?", "exercise", "2024-01-15")
	if err != nil {
		t.Fatalf("failed to query activity: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 activity record, got %d", count)
	}
}

func TestDoIdempotent(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	activity := Activity{
		Date:        "2024-01-15",
		IsCompleted: 1,
		HabitName:   "exercise",
	}

	// Insert twice
	if err := app.Do(activity); err != nil {
		t.Fatalf("failed to do activity (first call): %v", err)
	}
	if err := app.Do(activity); err != nil {
		t.Fatalf("failed to do activity (second call): %v", err)
	}

	// Should only have one record due to INSERT OR IGNORE
	var count int
	err := app.DB.Get(&count, "SELECT COUNT(*) FROM activity WHERE habit_name = ? AND date = ?", "exercise", "2024-01-15")
	if err != nil {
		t.Fatalf("failed to query activity: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 activity record (idempotent), got %d", count)
	}
}

func TestUndo(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	activity := Activity{
		Date:        "2024-01-15",
		IsCompleted: 1,
		HabitName:   "exercise",
	}

	// Add activity
	if err := app.Do(activity); err != nil {
		t.Fatalf("failed to do activity: %v", err)
	}

	// Undo activity
	err := app.Undo(activity)
	if err != nil {
		t.Errorf("Undo() error = %v", err)
	}

	// Verify activity was deleted
	var count int
	err = app.DB.Get(&count, "SELECT COUNT(*) FROM activity WHERE habit_name = ? AND date = ?", "exercise", "2024-01-15")
	if err != nil {
		t.Fatalf("failed to query activity: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 activity records after undo, got %d", count)
	}
}

func TestGetHabitActivity(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	// Add activities using recent dates (within last month)
	now := time.Now()
	date1 := now.AddDate(0, 0, -5).Format("2006-01-02")
	date2 := now.AddDate(0, 0, -4).Format("2006-01-02")
	date3 := now.AddDate(0, 0, -3).Format("2006-01-02")

	activities := []Activity{
		{Date: date1, IsCompleted: 1, HabitName: "exercise"},
		{Date: date2, IsCompleted: 1, HabitName: "exercise"},
		{Date: date3, IsCompleted: 1, HabitName: "exercise"},
	}

	for _, activity := range activities {
		if err := app.Do(activity); err != nil {
			t.Fatalf("failed to do activity: %v", err)
		}
	}

	// Get activity with 1 month lookback
	result, err := app.GetHabitActivity("exercise", 1)
	if err != nil {
		t.Fatalf("GetHabitActivity() error = %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 activities, got %d", len(result))
	}
}

func TestGetPeriodActivityToday(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit 'exercise': %v", err)
	}
	if err := app.CreateHabit("reading"); err != nil {
		t.Fatalf("failed to create habit 'reading': %v", err)
	}

	today := time.Now().Format("2006-01-02")

	// Add today's activities
	if err := app.Do(Activity{Date: today, IsCompleted: 1, HabitName: "exercise"}); err != nil {
		t.Fatalf("failed to do activity: %v", err)
	}

	rows, dates, err := app.GetPeriodActivity("today")
	if err != nil {
		t.Fatalf("GetPeriodActivity() error = %v", err)
	}

	if len(dates) != 1 {
		t.Errorf("expected 1 date for today, got %d", len(dates))
	}

	if len(rows) != 2 {
		t.Errorf("expected 2 habit rows, got %d", len(rows))
	}
}

func TestGetPeriodActivityWeek(t *testing.T) {
	app := setupTestDB(t)
	defer func() {
		if err := app.DB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	if err := app.CreateHabit("exercise"); err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	rows, dates, err := app.GetPeriodActivity("week")
	if err != nil {
		t.Fatalf("GetPeriodActivity() error = %v", err)
	}

	if len(dates) != 8 {
		t.Errorf("expected 8 dates for week, got %d", len(dates))
	}

	if len(rows) != 1 {
		t.Errorf("expected 1 habit row, got %d", len(rows))
	}
}
