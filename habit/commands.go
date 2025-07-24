package habit

import (
	"os"
	"slices"
	"time"

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

func GetActivities(args []string) []Activity {
	var activities []Activity
	var err error
	if validateHabit(args) {
		activities, err = Habit.GetActivity(args[0], 3)
		if err != nil {
			log.Info().Msgf("No activities found for habit: %s", args[0])
			os.Exit(1)
		}
	}

	return activities
}

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
