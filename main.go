package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ANSI escape codes for resetting text attributes
const reset = "\033[0m"

func main() {
	// Parse command-line arguments for month and year
	month, year := parseArgs()

	// If no arguments, display current month/year
	if month == 0 || year == 0 {
		now := time.Now()
		month = int(now.Month())
		year = now.Year()
	}

	drawCalendar(month, year)
}

// parseArgs parses command-line arguments for month and year.
// It expects "go run your_program.go [month] [year]".
func parseArgs() (int, int) {
	args := make([]string, 0)
	// Skip the program name itself
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	month, year := 0, 0
	if len(args) >= 1 {
		m, err := strconv.Atoi(args[0])
		if err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}
	if len(args) >= 2 {
		y, err := strconv.Atoi(args[1])
		if err == nil && y >= 1900 { // Arbitrary sensible lower bound for year
			year = y
		}
	}
	return month, year
}

// drawCalendar draws a calendar for the given month and year in the terminal.
func drawCalendar(month int, year int) {
	// Get the first day of the month
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	// Get the last day of the month
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// Print calendar header
	fmt.Println()
	currentMonthYear := fmt.Sprintf("%s %d", firstOfMonth.Format("January"), year)
	leftPadding := (23 - len(currentMonthYear)) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", leftPadding), currentMonthYear)

	fmt.Printf("")
	fmt.Println(" Su Mo Tu We Th Fr Sa")

	// Print leading spaces for the first day of the week
	for i := 0; i < int(firstOfMonth.Weekday()); i++ {
		fmt.Printf("   ")
	}

	// Print days
	for day := 1; day <= lastOfMonth.Day(); day++ {
		currentDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		dayOfWeek := currentDay.Weekday()

		// Highlight current day
		if currentDay.Year() == time.Now().Year() &&
			currentDay.Month() == time.Now().Month() &&
			currentDay.Day() == time.Now().Day() {
			color.New(color.FgRed, color.Bold).Printf("%3d", day)
		} else {
			fmt.Printf("%3d", day)
		}

		if dayOfWeek == time.Saturday {
			fmt.Println() // New line after Saturday
		}
	}
	fmt.Println("\n") // Extra new line at the end
}
