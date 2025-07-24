package main

import (
	"fmt"
	"github.com/kahnwong/habit-tracker/calendar"
	"time"
)

// [TODO] month/year color
// [TODO] replace highlighted date color
// [TODO] add generated code disclaimer

func main() {
	// Define dates to highlight (example: today, tomorrow, and a date in the past)
	highlightDates := []time.Time{
		time.Date(2025, time.June, 15, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.May, 10, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.July, 10, 0, 0, 0, 0, time.Local),
	}

	calendar.RenderCalendarView(3, highlightDates)

	fmt.Printf("----%s\n", time.Now()) // debug
}
