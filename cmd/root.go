package cmd

import (
	"os"
	"strings"

	"github.com/kahnwong/habit-tracker/internal/habit"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

func HabitsGet(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var autocomplete []string

	if len(args) == 0 {
		if err := habit.Init(); err != nil {
			return autocomplete, cobra.ShellCompDirectiveNoFileComp
		}

		autocomplete, _ = habit.Habit.GetHabits()
	}

	return autocomplete, cobra.ShellCompDirectiveNoFileComp
}

var rootCmd = &cobra.Command{
	Use:              "habit-tracker",
	Short:            "Display habits activity in tui",
	PersistentPreRun: initHabitDB,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func initHabitDB(cmd *cobra.Command, args []string) {
	commandPath := cmd.CommandPath()
	if commandPath == "habit-tracker" || strings.HasPrefix(commandPath, "habit-tracker completion") {
		return
	}

	if err := habit.Init(); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
