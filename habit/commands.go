package habit

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/kahnwong/habit-tracker/calendar"

	"github.com/rs/zerolog/log"
)

func Create(args []string) {
	if len(args) < 1 {
		log.Fatal().Msg("You must specify a habit")
	} else {
		err := Habit.CreateHabit(args[0])
		if err != nil {
			log.Info().Msgf("Habit %s already exists", args[0])
		} else {
			log.Info().Msgf("Habit %s created", args[0])
		}
	}
}

func Do(args []string) {
	if validateHabit(args) {
		today := time.Now().Format("2006-01-02")

		activity := Activity{
			Date: today, IsCompleted: 1, HabitName: args[0]}

		err := Habit.Do(activity)
		if err != nil {
			log.Error().Msg("Error logging a habit")
		} else {
			log.Info().Msgf("Logged %s for today", args[0])
		}
	} else {
		log.Error().Msgf("Invalid habit: %s", args[0])
	}
}

func Undo(args []string) { // some chunks are duplicated from `Do()`
	if validateHabit(args) {
		today := time.Now().Format("2006-01-02")

		activity := Activity{
			Date: today, IsCompleted: 0, HabitName: args[0]}

		err := Habit.Undo(activity)
		if err != nil {
			log.Error().Msg("Error undoing a habit")
		} else {
			log.Info().Msgf("Undo %s for today", args[0])
		}
	} else {
		log.Error().Msgf("Invalid habit: %s", args[0])
	}
}

func ShowHabitActivity(lookbackMonths int, args []string) {
	// fetch activities
	var activities []Activity
	var err error
	if validateHabit(args) {
		activities, err = Habit.GetHabitActivity(args[0], lookbackMonths)
		if err != nil {
			log.Fatal().Msgf("Error fetching activities for habit: %s", args[0])
		}
	}

	if len(activities) == 0 {
		log.Info().Msgf("No activities found for habit: %s", args[0])
		os.Exit(0)
	}

	// parse date
	var dates []time.Time
	layout := "2006-01-02"
	for _, a := range activities {
		t, err := time.Parse(layout, a.Date)
		if err != nil {
			log.Error().Msgf("Error parsing date: %s", a.Date)
			continue
		}
		dates = append(dates, t)
	}

	// render calendar
	calendar.RenderCalendarView(lookbackMonths, dates)
}

func ShowPeriodActivity(period string) {
	// fetch activities
	activities, dates, err := Habit.GetPeriodActivity(period)
	if err != nil {
		log.Fatal().Msgf("Error fetching activities for period: %s", period)
	}

	if len(activities) == 0 {
		log.Info().Msgf("No activities found for period: %s", period)
		os.Exit(0)
	}

	// render table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	//// set header
	headers := append([]string{""}, dates...)
	var headerRow table.Row
	for _, h := range headers {
		headerRow = append(headerRow, fmt.Sprintf("%2s  ", h))
	}
	t.AppendHeader(headerRow)

	/// unwind data
	isCompletedIcon := map[int64]string{0: "", 1: "✓"} // [TODO] apply color
	for _, activity := range activities {
		var elems []interface{}
		elems = append(elems, fmt.Sprintf("%-6s", activity["habit_name"])) // %-6s for left-alignment and padding

		for _, date := range dates {
			if intVal, ok := activity[date].(int64); ok {
				elems = append(elems, fmt.Sprintf("%6s", isCompletedIcon[intVal])) // %6s for center alignment
			} else {
				log.Error().Err(err).Msgf("Failed to cast to int. Value is of type %T\n", activity[date])
			}
		}

		t.AppendRows([]table.Row{
			elems,
		})
	}

	//// styling
	t.SetStyle(table.Style{
		Options: table.OptionsNoBordersAndSeparators,
	})

	// render
	t.Render()
}

// utils
func validateHabit(args []string) bool {
	if len(args) == 0 {
		log.Fatal().Msg("Habit must be specified")
	}

	habitName := args[0]
	habits, err := Habit.GetHabits()
	if err != nil {
		log.Fatal().Msg("Error getting habits")
	}

	return slices.Contains(habits, habitName)
}
