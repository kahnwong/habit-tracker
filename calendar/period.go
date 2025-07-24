package calendar

import "time"

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
