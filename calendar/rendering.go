package calendar

import (
	"fmt"
	"strings"
	"time"
)

func printHeaders(currentBlockMonths []time.Time) {
	for idx, m := range currentBlockMonths {
		title := m.Format("January 2006")
		titlePadding := calendarElementRenderedWidth - len(title)
		leftPad := titlePadding / 2
		rightPad := titlePadding - leftPad
		fmt.Printf("%s%s%s", strings.Repeat(" ", leftPad), title, strings.Repeat(" ", rightPad))

		if idx < len(currentBlockMonths)-1 {
			fmt.Printf("%s", strings.Repeat(" ", paddingBetweenCalendars))
		}
	}
	fmt.Println()
}

func printWeekdayHeaders(currentBlockMonths []time.Time) {
	for idx := range currentBlockMonths {
		fmt.Printf("Su Mo Tu We Th Fr Sa  ") // This is exactly 22 chars
		if idx < len(currentBlockMonths)-1 {
			fmt.Printf("%s", strings.Repeat(" ", paddingBetweenCalendars))
		}
	}
	fmt.Println()
}
