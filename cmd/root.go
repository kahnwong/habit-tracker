package cmd

import (
	"os"

	"github.com/kahnwong/habit-tracker/habit"

	"github.com/spf13/cobra"
)

func HabitsGet(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var autocomplete []string

	if len(args) == 0 {
		autocomplete, _ = habit.Habit.GetHabits()
	}

	return autocomplete, cobra.ShellCompDirectiveNoFileComp
}

var rootCmd = &cobra.Command{
	Use:   "habit-tracker",
	Short: "Display habits activity in tui",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
