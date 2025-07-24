package cmd

import (
	"fmt"
	"github.com/kahnwong/habit-tracker/core"
	"time"

	"github.com/kahnwong/habit-tracker/calendar"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show habit stats",
	Run: func(cmd *cobra.Command, args []string) {
		_ = core.Foo()
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
