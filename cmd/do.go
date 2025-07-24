package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

// [TODO] autocomplete habit names

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Track a habit",
	Run: func(cmd *cobra.Command, args []string) {
		habit.Do(args)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
