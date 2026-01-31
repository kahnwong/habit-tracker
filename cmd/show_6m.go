package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var show6mCmd = &cobra.Command{
	Use:               "show-6m",
	Short:             "Show habit stats for the last 6 months",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.ShowHabitActivity(6, args); err != nil {
			log.Fatal().Err(err).Msg("failed to show habit activity")
		}
	},
}

func init() {
	rootCmd.AddCommand(show6mCmd)
}
