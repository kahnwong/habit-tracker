package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show habit stats for today",
	Run: func(cmd *cobra.Command, args []string) {
		habit.ShowPeriodActivity("today")
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
