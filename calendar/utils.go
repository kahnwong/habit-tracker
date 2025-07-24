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
