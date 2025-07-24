package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	// Define dates to highlight (example: today, tomorrow, and a date in the past)
	highlightDates := []time.Time{
		now,
		now.AddDate(0, 0, 1),
		time.Date(2025, time.June, 15, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), // Another highlight for a 4th month example
	}

	// Get the last four months for demonstration (you can change this number)
	// Example: now.AddDate(0, -3, 0) gives 4 months total (current + 3 previous)
	numMonthsToDisplay := 12
	allMonths := make([]time.Time, numMonthsToDisplay)
	for i := 0; i < numMonthsToDisplay; i++ {
		allMonths[numMonthsToDisplay-1-i] = now.AddDate(0, -i, 0)
	}

	// Define the number of months per row
	monthsPerRow := 3

	// Process months in blocks of 'monthsPerRow'
	for i := 0; i < len(allMonths); i += monthsPerRow {
		endIdx := i + monthsPerRow
		if endIdx > len(allMonths) {
			endIdx = len(allMonths)
		}
		currentBlockMonths := allMonths[i:endIdx]

		// Print headers for the current block
		for _, m := range currentBlockMonths {
			fmt.Printf("%20s", m.Format("January 2006"))
		}
		fmt.Println()

		for _, _ = range currentBlockMonths {
			fmt.Printf("Su Mo Tu We Th Fr Sa  ")
		}
		fmt.Println()

		// Calculate max number of weeks for the current block
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

		// Create a 2D slice to hold the formatted days for each month in the current block
		monthDaysBlock := make([][]string, len(currentBlockMonths))
		for j, m := range currentBlockMonths {
			firstDayOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
			offset := int(firstDayOfMonth.Weekday())

			daysInMonth := getDaysInMonth(m.Year(), m.Month())

			currentMonthDays := make([]string, numWeeks*7) // Use numWeeks from the block
			for k := 0; k < offset; k++ {
				currentMonthDays[k] = "  " // Pad initial empty spaces
			}
			for day := 1; day <= daysInMonth; day++ {
				currentDay := time.Date(m.Year(), m.Month(), day, 0, 0, 0, 0, time.Local)
				if isHighlighted(currentDay, highlightDates) {
					currentMonthDays[offset+day-1] = fmt.Sprintf("\033[31m%2d\033[0m", day) // Red color
				} else {
					currentMonthDays[offset+day-1] = fmt.Sprintf("%2d", day)
				}
			}
			monthDaysBlock[j] = currentMonthDays
		}

		// Print the calendar row by row for the current block
		for week := 0; week < numWeeks; week++ {
			for _, days := range monthDaysBlock {
				for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
					idx := week*7 + dayOfWeek
					if idx < len(days) {
						fmt.Printf("%s ", days[idx])
					} else {
						fmt.Print("   ") // Pad if month is shorter
					}
				}
				fmt.Print(" ") // Space between months
			}
			fmt.Println()
		}
		fmt.Println() // Add an extra newline between blocks for better separation
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
