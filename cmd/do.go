package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:               "do",
	Short:             "Track a habit",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.Do(args); err != nil {
			log.Fatal().Err(err).Msg("failed to track habit")
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
