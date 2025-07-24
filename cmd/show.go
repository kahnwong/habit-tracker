package cmd

import (
	"fmt"
	"time"

	"github.com/kahnwong/habit-tracker/habit"
	"github.com/rs/zerolog/log"

	"github.com/kahnwong/habit-tracker/calendar"

	"github.com/spf13/cobra"
)

// [TODO] autocomplete habit names
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show habit stats",
	Run: func(cmd *cobra.Command, args []string) {
		activities := habit.GetActivities(args)
		if len(activities) == 0 {
			log.Info().Msgf("No activities found for habit: %s", args[0])
		}
		fmt.Println(activities) // [TODO] WIP

		///
		highlightDates := []time.Time{
			time.Date(2025, time.June, 15, 0, 0, 0, 0, time.Local),
			time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local),
			time.Date(2025, time.May, 10, 0, 0, 0, 0, time.Local),
			time.Date(2025, time.July, 10, 0, 0, 0, 0, time.Local),
		}

		calendar.RenderCalendarView(3, highlightDates)

		fmt.Printf("----%s\n", time.Now()) // debug
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
