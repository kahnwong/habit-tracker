package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:               "undo",
	Short:             "Untrack a habit",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		habit.Undo(args)
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
