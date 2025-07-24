package main

import (
	"fmt"
	"strings"
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

	allMonths := generateMonths(3)

	// Process months in blocks of 'monthsPerRow'
	for i := 0; i < len(allMonths); i += monthsPerRow {
		endIdx := i + monthsPerRow
		if endIdx > len(allMonths) {
			endIdx = len(allMonths)
		}
		currentBlockMonths := allMonths[i:endIdx]

		printHeaders(currentBlockMonths)

		// --- Print Weekday Headers ---
		for idx := range currentBlockMonths {
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

		// --- Print the calendar row by row for the current block ---
		for week := 0; week < numWeeks; week++ {
			for idx, days := range monthDaysBlock {
				var weekLineBuilder strings.Builder
				for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
					idxInDays := week*7 + dayOfWeek
					dayContent := ""
					if idxInDays < len(days) {
						dayContent = days[idxInDays]
					}

					// Explicitly ensure each day slot consumes 3 characters.
					// ANSI escape codes add length to the string but not to the displayed width.
					// We need to account for the *displayed* width.
					displayLength := calculateDisplayedWidth(dayContent)

					// Prepend spaces to center the number within its 3-character slot.
					// Or just ensure it's left-aligned and padded right. "%2s" already handles leading space.
					// We just need a trailing space if the content itself (like " 1" or "10") is 2 chars.
					if displayLength == 2 { // This is a " 1" or "10"
						weekLineBuilder.WriteString(dayContent + " ")
					} else if displayLength == 1 { // This shouldn't happen with "%2d" but as a fallback
						weekLineBuilder.WriteString(" " + dayContent + " ") // Center single digit if it somehow appears
					} else { // This is either an empty slot "   " or already correctly formatted
						weekLineBuilder.WriteString(dayContent) // Use as is, assuming it's "   " or already 3 chars
					}
				}
				// After building the 7-day string for the week, pad it to the full calendar element width.
				renderedWeekLine := weekLineBuilder.String()
				// The intended displayed width of 7 days * 3 chars/day = 21 chars.
				// Our overall calendarElementRenderedWidth is 22 (to match "Su Mo Tu We Th Fr Sa  ").
				// So, we need to ensure the end of the line has one extra space.
				// However, if the last day of the month falls on a Saturday and there's no trailing " " in its slot,
				// the renderedWeekLine might be 20 chars long instead of 21 (e.g., " 1  2  3  4  5  6 7").
				// Let's ensure the total displayed width for the day grid part is 21, then add the final space.

				// Calculate actual displayed width of the built string, ignoring ANSI codes
				actualDisplayedWidth := calculateDisplayedWidth(renderedWeekLine)

				fmt.Print(renderedWeekLine)
				// Pad the remaining space up to 21 (for the 7 days) and then add the final fixed space.
				if actualDisplayedWidth < calendarElementRenderedWidth-1 { // calendarElementRenderedWidth - 1 because the header ends with "  "
					fmt.Print(strings.Repeat(" ", (calendarElementRenderedWidth-1)-actualDisplayedWidth))
				}
				fmt.Print(" ") // The fixed trailing space to match header width

				// Add padding between calendar elements, but not after the last one in the row
				if idx < len(monthDaysBlock)-1 {
					fmt.Printf("%s", strings.Repeat(" ", paddingBetweenCalendars))
				}
			}
			fmt.Println()
		}

		fmt.Println() // Add an extra newline between blocks for better separation
	}
	fmt.Printf("----%s\n", now) // debug
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

// time period
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

// rendering
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
