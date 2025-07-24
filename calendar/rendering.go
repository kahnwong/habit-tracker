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

func renderCalendar(numWeeks int, monthDaysBlock [][]string) {
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
}
