package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

var show6mCmd = &cobra.Command{
	Use:               "show-6m",
	Short:             "Show habit stats for the last 6 months",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		habit.ShowHabitActivity(6, args)
	},
}

func init() {
	rootCmd.AddCommand(show6mCmd)
}
