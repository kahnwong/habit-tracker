package habit

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/kahnwong/habit-tracker/calendar"

	"github.com/rs/zerolog/log"
)

func Create(args []string) {
	if len(args) < 1 {
		log.Fatal().Msg("You must specify a habit")
	} else {
		err := Habit.CreateHabit(args[0])
		if err != nil {
			log.Error().Msgf("Habit %s already exists", args[0])
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
			log.Info().Msgf("No activities found for habit: %s", args[0])
			os.Exit(0)
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

func ShowPeriodActivity(period string, args []string) {
	// fetch activities
	var activities []Activity
	var err error
	activities, err = Habit.GetPeriodActivity(period)
	if err != nil {
		log.Info().Msgf("No activities found for period: %s", period)
		os.Exit(0)
	}

	if len(activities) == 0 {
		log.Info().Msgf("No activities found for period: %s", period)
		os.Exit(0)
	}

	// do something
	fmt.Println(activities)
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
