package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	// Get the last three months
	months := make([]time.Time, 3)
	for i := 0; i < 3; i++ {
		months[2-i] = now.AddDate(0, -i, 0)
	}

	// Print headers
	for _, m := range months {
		fmt.Printf("%20s", m.Format("January 2006"))
	}
	fmt.Println()

	for _, _ = range months {
		fmt.Printf("Su Mo Tu We Th Fr Sa  ")
	}
	fmt.Println()

	// Calculate and print days
	numWeeks := 0
	for _, m := range months {
		firstDayOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
		// Go's Weekday starts Sunday=0, Monday=1, ..., Saturday=6
		// We want to align the first day correctly under its weekday column.
		offset := int(firstDayOfMonth.Weekday())

		daysInMonth := 0
		if m.Month() == time.February {
			if isLeapYear(m.Year()) {
				daysInMonth = 29
			} else {
				daysInMonth = 28
			}
		} else if m.Month() == time.April || m.Month() == time.June ||
			m.Month() == time.September || m.Month() == time.November {
			daysInMonth = 30
		} else {
			daysInMonth = 31
		}

		weeks := (offset + daysInMonth + 6) / 7 // Calculate total number of rows needed
		if weeks > numWeeks {
			numWeeks = weeks
		}
	}

	// Create a 2D slice to hold the days for each month, padding with empty strings
	monthDays := make([][]string, len(months))
	for i, m := range months {
		firstDayOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
		offset := int(firstDayOfMonth.Weekday())

		daysInMonth := 0
		if m.Month() == time.February {
			if isLeapYear(m.Year()) {
				daysInMonth = 29
			} else {
				daysInMonth = 28
			}
		} else if m.Month() == time.April || m.Month() == time.June ||
			m.Month() == time.September || m.Month() == time.November {
			daysInMonth = 30
		} else {
			daysInMonth = 31
		}

		currentMonthDays := make([]string, numWeeks*7)
		for j := 0; j < offset; j++ {
			currentMonthDays[j] = "  " // Pad initial empty spaces
		}
		for day := 1; day <= daysInMonth; day++ {
			currentMonthDays[offset+day-1] = fmt.Sprintf("%2d", day)
		}
		monthDays[i] = currentMonthDays
	}

	// Print the calendar row by row
	for week := 0; week < numWeeks; week++ {
		for _, days := range monthDays {
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
}

// isLeapYear checks if a year is a leap year.
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
