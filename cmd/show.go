package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:               "show",
	Short:             "Show habit stats for the last 3 months",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		habit.ShowHabitActivity(3, args)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
