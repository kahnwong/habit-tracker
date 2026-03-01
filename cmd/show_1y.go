package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var showYearCmd = &cobra.Command{
	Use:               "show-1y",
	Short:             "Show habit stats for current year",
	ValidArgsFunction: HabitsGet,
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.ShowHabitActivity(12, args); err != nil {
			log.Fatal().Err(err).Msg("failed to show habit activity")
		}
	},
}

func init() {
	rootCmd.AddCommand(showYearCmd)
}
