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

// isLeapYear checks if a year is a leap year.
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// getDaysInMonth returns the number of days in a given month and year.
func getDaysInMonth(year int, month time.Month) int {
	if month == time.February {
		if isLeapYear(year) {
			return 29
		}
		return 28
	} else if month == time.April || month == time.June ||
		month == time.September || month == time.November {
		return 30
	}
	return 31
}

// isHighlighted checks if a given date is in the list of highlight dates.
func isHighlighted(date time.Time, highlightDates []time.Time) bool {
	for _, hd := range highlightDates {
		// Compare Year, Month, and Day to ignore time components
		if date.Year() == hd.Year() && date.Month() == hd.Month() && date.Day() == hd.Day() {
			return true
		}
	}
	return false
}

// calculateDisplayedWidth calculates the visible width of a string, ignoring ANSI escape codes.
func calculateDisplayedWidth(s string) int {
	inEscape := false
	width := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			width++
		}
	}
	return width
}
