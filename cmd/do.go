package cmd

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

// [TODO] autocomplete habit names

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Track a habit",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a habit")
			os.Exit(1)
		}

		habitName := args[0]

		// validate habit
		habits, err := Habit.GetHabits()
		if err != nil {
			fmt.Println("Error getting habits")
			os.Exit(1)
		}

		isValidHabit := slices.Contains(habits, habitName)
		if isValidHabit {
			today := time.Now().Format("2006-01-02")

			activity := habit.Activity{
				Date: today, IsCompleted: 1, HabitName: habitName}

			err = Habit.Do(activity)
			if err != nil {
				fmt.Println("Error logging a habit")
			} else {
				fmt.Printf("Logged %s for today\n", args[0])
			}
		} else {
			fmt.Printf("Invalid habit: %s\n", habitName)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
