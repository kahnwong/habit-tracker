package calendar

import (
	"fmt"
	"time"
)

func generateMonths(lookbackMonths int) []time.Time {
	allMonths := make([]time.Time, lookbackMonths)
	for i := 0; i < lookbackMonths; i++ {
		allMonths[lookbackMonths-1-i] = now.AddDate(0, -i, 0)
	}

	return removePreviousYearDates(allMonths)
}

func removePreviousYearDates(dates []time.Time) []time.Time {
	var currentYearDates []time.Time
	currentYear := now.Year()

	for _, date := range dates {
		if date.Year() == currentYear {
			currentYearDates = append(currentYearDates, date)
		}
	}
	return currentYearDates
}

func getMaxWeeks(currentBlockMonths []time.Time) int {
	// --- Calculate Max Number of Weeks for the Current Block ---
	numWeeks := 0
	for _, m := range currentBlockMonths {
		firstDayOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
		offset := int(firstDayOfMonth.Weekday())

		daysInMonth := getDaysInMonth(m.Year(), m.Month())

		weeks := (offset + daysInMonth + 6) / 7
		if weeks > numWeeks {
			numWeeks = weeks
		}
	}
	return numWeeks
}

func createMonthDaysBlock(currentBlockMonths []time.Time, numWeeks int, highlightDates []time.Time) [][]string {
	monthDaysBlock := make([][]string, len(currentBlockMonths))
	for j, m := range currentBlockMonths {
		firstDayOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
		offset := int(firstDayOfMonth.Weekday())

		daysInMonth := getDaysInMonth(m.Year(), m.Month())

		currentMonthDays := make([]string, numWeeks*7)
		for k := 0; k < offset; k++ {
			currentMonthDays[k] = "   " // Each "day slot" is 3 chars wide
		}
		for day := 1; day <= daysInMonth; day++ {
			currentDay := time.Date(m.Year(), m.Month(), day, 0, 0, 0, 0, time.Local)
			dayStr := fmt.Sprintf("%2d", day) // Format as " 1" or "10"
			if isHighlighted(currentDay, highlightDates) {
				currentMonthDays[offset+day-1] = fmt.Sprintf("\033[31m%s\033[0m", dayStr) // Store just the formatted date with color codes
			} else {
				currentMonthDays[offset+day-1] = dayStr // Store just the formatted date
			}
		}
		monthDaysBlock[j] = currentMonthDays
	}
	return monthDaysBlock
}

// helpers
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

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

func isHighlighted(date time.Time, highlightDates []time.Time) bool {
	for _, hd := range highlightDates {
		// Compare Year, Month, and Day to ignore time components
		if date.Year() == hd.Year() && date.Month() == hd.Month() && date.Day() == hd.Day() {
			return true
		}
	}
	return false
}
