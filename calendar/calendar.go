// originally generated via gemini
// refactored by `Karn Wong <karn@karnwong.me>`

package calendar

import (
	"fmt"
	"time"
)

var (
	now = time.Now()
)

const (
	monthsPerRow                 = 3
	calendarElementRenderedWidth = 22
	paddingBetweenCalendars      = 3
)

func RenderCalendarView(lookbackMonths int, highlightDates []time.Time) {
	allMonths := generateMonths(lookbackMonths)

	for i := 0; i < len(allMonths); i += monthsPerRow {
		endIdx := i + monthsPerRow
		if endIdx > len(allMonths) {
			endIdx = len(allMonths)
		}
		currentBlockMonths := allMonths[i:endIdx]

		printHeaders(currentBlockMonths)
		printWeekdayHeaders(currentBlockMonths)

		numWeeks := getMaxWeeks(currentBlockMonths)
		monthDaysBlock := createMonthDaysBlock(currentBlockMonths, numWeeks, highlightDates)

		renderCalendar(numWeeks, monthDaysBlock)
		fmt.Println()
	}
}
