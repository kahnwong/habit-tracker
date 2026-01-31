package cmd

import (
	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a habit.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := habit.Create(args); err != nil {
			log.Fatal().Err(err).Msg("failed to create habit")
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
