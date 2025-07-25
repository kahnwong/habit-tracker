package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

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
