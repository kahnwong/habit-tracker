package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

// [TODO] autocomplete habit names
var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Untrack a habit",
	Run: func(cmd *cobra.Command, args []string) {
		habit.Undo(args)
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
