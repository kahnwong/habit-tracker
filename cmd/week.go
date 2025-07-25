package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

// [TODO] autocomplete period: today, week
// [TODO] sql should return missing habits
var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show habit stats for week",
	Run: func(cmd *cobra.Command, args []string) {
		habit.ShowPeriodActivity("week")
	},
}

func init() {
	rootCmd.AddCommand(weekCmd)
}
