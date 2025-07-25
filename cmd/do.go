package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:               "do",
	Short:             "Track a habit",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		habit.Do(args)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
