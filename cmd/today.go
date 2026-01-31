package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show habit stats for today",
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.ShowPeriodActivity("today"); err != nil {
			log.Fatal().Err(err).Msg("failed to show period activity")
		}
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
