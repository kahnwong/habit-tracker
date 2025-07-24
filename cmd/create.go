package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a habit.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must specify a habit")
		} else {
			err := Habit.CreateHabit(args[0])
			if err != nil {
				fmt.Printf("Habit %s already exists\n", args[0])
			} else {
				fmt.Printf("Created habit: %s\n", args[0])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
