package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	now := time.Now() // July 24, 2025

	// Define dates to highlight (example: today, tomorrow, and a date in the past)
	highlightDates := []time.Time{
		now,
		now.AddDate(0, 0, 1), // July 25, 2025
		time.Date(2025, time.June, 15, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), // Another highlight for a 4th month example
		time.Date(2025, time.May, 10, 0, 0, 0, 0, time.Local),    // Highlight for the 3rd last month
	}

	// Get the last four months for demonstration (current month + 3 previous = 4 months total)
	// If now is July, this will be April, May, June, July.
	numMonthsToDisplay := 12
	allMonths := make([]time.Time, numMonthsToDisplay)
	for i := 0; i < numMonthsToDisplay; i++ {
		// Populate in reverse to get chronological order (earliest to latest)
		allMonths[numMonthsToDisplay-1-i] = now.AddDate(0, -i, 0)
	}

	// Define the number of months per row
	monthsPerRow := 3

	// Define the fixed width of each calendar element (7 days * 3 chars/day = 21 chars for "XX " days.
	// We'll aim for a total width that comfortably holds "Su Mo Tu We Th Fr Sa  " (22 chars)
	const calendarElementRenderedWidth = 22 // "Su Mo Tu We Th Fr Sa  " is 22 chars
	const paddingBetweenCalendars = 3

	// Process months in blocks of 'monthsPerRow'
	for i := 0; i < len(allMonths); i += monthsPerRow {
		endIdx := i + monthsPerRow
		if endIdx > len(allMonths) {
			endIdx = len(allMonths)
		}
		currentBlockMonths := allMonths[i:endIdx]

		// --- Print Headers for the Current Block ---
		for idx, m := range currentBlockMonths {
			title := m.Format("January 2006")
			titlePadding := calendarElementRenderedWidth - len(title)
			leftPad := titlePadding / 2
			rightPad := titlePadding - leftPad
			fmt.Printf("%s%s%s", strings.Repeat(" ", leftPad), title, strings.Repeat(" ", rightPad))

			if idx < len(currentBlockMonths)-1 { // Don't add padding after the last month in the row
				fmt.Printf("%s", strings.Repeat(" ", paddingBetweenCalendars))
			}
		}
		fmt.Println()

		// --- Print Weekday Headers ---
		for idx, _ := range currentBlockMonths {
			fmt.Printf("Su Mo Tu We Th Fr Sa  ") // This is exactly 22 chars
			if idx < len(currentBlockMonths)-1 {
				fmt.Printf("%s", strings.Repeat(" ", paddingBetweenCalendars))
			}
		}
		fmt.Println()

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

		// --- Create a 2D slice to hold the formatted days for each month in the current block ---
		monthDaysBlock := make([][]string, len(currentBlockMonths))
		for j, m := range currentBlockMonths {
			firstDayOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
			offset := int(firstDayOfMonth.Weekday())

			daysInMonth := getDaysInMonth(m.Year(), m.Month())

			// This inner slice will hold all formatted day strings for this month, including leading/trailing blanks
			currentMonthDays := make([]string, numWeeks*7)
			for k := 0; k < offset; k++ {
				currentMonthDays[k] = "  " // Pad initial empty spaces. Each "  " is 2 chars. Day itself is 2 chars, then space. So day takes 3. "  " means "  " + " " = 3
			}
			for day := 1; day <= daysInMonth; day++ {
				currentDay := time.Date(m.Year(), m.Month(), day, 0, 0, 0, 0, time.Local)
				if isHighlighted(currentDay, highlightDates) {
					currentMonthDays[offset+day-1] = fmt.Sprintf("\033[31m%2d\033[0m", day) // Red color, takes 2 chars for day, ANSI codes don't count
				} else {
					currentMonthDays[offset+day-1] = fmt.Sprintf("%2d", day) // Takes 2 chars
				}
			}
			monthDaysBlock[j] = currentMonthDays
		}

		// --- Print the calendar row by row for the current block ---
		for week := 0; week < numWeeks; week++ {
			for idx, days := range monthDaysBlock {
				var weekLineBuilder strings.Builder
				for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
					idxInDays := week*7 + dayOfWeek
					if idxInDays < len(days) && days[idxInDays] != "" {
						// Each day slot, whether it's a number or "  ", needs to be 3 characters wide
						// If highlighted, the ANSI codes add invisible characters, but the displayed content is 2 chars.
						// So, we need to add the trailing space explicitly.
						// Example: " 1 ", "10 ", "   "
						if strings.Contains(days[idxInDays], "\033[") { // Is it a highlighted string?
							weekLineBuilder.WriteString(days[idxInDays]) // Add the color codes and the number
							// The displayed width is 2 chars. We need 1 more space for the 3-char slot.
							weekLineBuilder.WriteString(" ")
						} else {
							weekLineBuilder.WriteString(fmt.Sprintf("%2s ", days[idxInDays])) // Ensure 2-char content + 1 space
						}
					} else {
						weekLineBuilder.WriteString("   ") // 3 spaces for an empty day slot
					}
				}
				// Pad the entire week string to ensure it's calendarElementRenderedWidth (22) characters long
				// before adding the inter-calendar padding.
				// The actual displayed length of the string from weekLineBuilder.String() will be calendarElementRenderedWidth-1 (21)
				// because it does not include the last space that "Su Mo Tu We Th Fr Sa  " includes.
				// Let's ensure it's always 22 by adding one more space at the end.
				renderedWeekLine := weekLineBuilder.String()
				fmt.Print(renderedWeekLine)
				// The "Su Mo Tu We Th Fr Sa  " header is 22 chars. Our days are 7 * 3 = 21. Plus the final space makes 22.
				// If renderedWeekLine.Len() is 21, and we need 22, add one more space.
				if len(renderedWeekLine) < calendarElementRenderedWidth {
					fmt.Print(strings.Repeat(" ", calendarElementRenderedWidth-len(renderedWeekLine)))
				}

				// Add padding between calendar elements, but not after the last one in the row
				if idx < len(monthDaysBlock)-1 {
					fmt.Printf("%s", strings.Repeat(" ", paddingBetweenCalendars))
				}
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
