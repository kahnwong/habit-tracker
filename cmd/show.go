package cmd

import (
	"os"
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
		lookbackMonths := 4

		// fetch activities
		activities := habit.GetActivities(lookbackMonths, args)
		if len(activities) == 0 {
			log.Info().Msgf("No activities found for habit: %s", args[0])
			os.Exit(0)
		}

		// convert to time object
		var dates []time.Time
		layout := "2006-01-02"
		for _, a := range activities {
			t, err := time.Parse(layout, a.Date)
			if err != nil {
				log.Error().Msgf("Error parsing date: %s", a.Date)
				continue
			}
			dates = append(dates, t)
		}

		calendar.RenderCalendarView(lookbackMonths, dates)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
