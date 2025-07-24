package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

// [TODO] autocomplete habit names
var showYearCmd = &cobra.Command{
	Use:   "show-year",
	Short: "Show habit stats for current year",
	Run: func(cmd *cobra.Command, args []string) {
		habit.ShowHabitActivity(12, args)
	},
}

func init() {
	rootCmd.AddCommand(showYearCmd)
}
