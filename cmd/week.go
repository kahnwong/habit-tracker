package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show habit stats for week",
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.ShowPeriodActivity("week"); err != nil {
			log.Fatal().Err(err).Msg("failed to show period activity")
		}
	},
}

func init() {
	rootCmd.AddCommand(weekCmd)
}
