package cmd

import (
	"github.com/kahnwong/habit-tracker/habit"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a habit.",
	Run: func(cmd *cobra.Command, args []string) {
		habit.Create(args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
