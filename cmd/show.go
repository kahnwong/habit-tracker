package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

// [TODO] autocomplete habit names
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show habit stats",
	Run: func(cmd *cobra.Command, args []string) {
		lookbackMonths := 3
		habit.Show(lookbackMonths, args)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
