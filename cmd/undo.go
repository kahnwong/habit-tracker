package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:               "undo",
	Short:             "Untrack a habit",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.Undo(args); err != nil {
			log.Fatal().Err(err).Msg("failed to untrack habit")
		}
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
