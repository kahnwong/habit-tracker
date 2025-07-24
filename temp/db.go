package temp

import (
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// [TODO] undo activity (today only)

func Foo() error {

	//habits := []string{"Drink Water", "Exercise", "Read Book", "Meditate"}
	//
	//log.Println("\n--- Inserting sample habits ---")
	//for _, habitName := range habits {
	//
	//}
	//log.Println("--- Sample habits insertion complete ---")

	//////// activity
	log.Println("\n--- Inserting sample activities ---")

	// Get today's date and yesterday's date for sample activities
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Sample activities data: {date, is_completed, habit_name}
	activities := []struct {
		Date        string
		IsCompleted int // 0 for false, 1 for true
		HabitName   string
	}{
		{today, 0, "Drink Water"},
		{today, 0, "Exercise"},
		{yesterday, 1, "Drink Water"},
		{yesterday, 1, "Read Book"},
	}

	for _, activity := range activities {
		query := `INSERT INTO activity (date, is_completed, habit_name) VALUES (?, ?, ?)`
		result, err := db.Exec(query, activity.Date, activity.IsCompleted, activity.HabitName)
		if err != nil {
			// Check if the error is due to a foreign key constraint violation
			if err.Error() == "FOREIGN KEY constraint failed" {
				log.Printf("Warning: Could not insert activity for habit '%s' on '%s'. Habit might not exist.", activity.HabitName, activity.Date)
				continue // Skip to the next activity
			}
			return fmt.Errorf("error inserting activity for habit '%s' on '%s': %w", activity.HabitName, activity.Date, err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected for activity: %w", err)
		}
		if rowsAffected > 0 {
			log.Printf("Inserted activity: Habit '%s' on %s, Completed: %t", activity.HabitName, activity.Date, activity.IsCompleted == 1)
		}
	}
	log.Println("--- Sample activities insertion complete ---")
	return nil
}

func main() {

	// Insert a user
	err := app.CreateUser(ctx, "Alice", "alice@example.com")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}

	// Get a user
	user, err := app.GetUserByID(ctx, 1)
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("Fetched user: %+v\n", user)
	}

	// Try inserting a duplicate email to see UNIQUE constraint in action
	err = app.CreateUser(ctx, "Bob", "alice@example.com")
	if err != nil {
		log.Printf("Expected error inserting duplicate email: %v", err)
	}
}
