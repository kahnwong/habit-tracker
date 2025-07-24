package habit

import (
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
			log.Error().Msgf("Habit %s already exists\n", args[0])
		} else {
			log.Info().Msgf("Habit %s created\n", args[0])
		}
	}
}

func Do(args []string) {
	if len(args) == 0 {
		log.Fatal().Msg("Habit must be specified")
	}

	habitName := args[0]

	// validate habit
	habits, err := Habit.GetHabits()
	if err != nil {
		log.Fatal().Msg("Error getting habits")
	}

	// track habit
	isValidHabit := slices.Contains(habits, habitName)
	if isValidHabit {
		today := time.Now().Format("2006-01-02")

		activity := Activity{
			Date: today, IsCompleted: 1, HabitName: habitName}

		err = Habit.Do(activity)
		if err != nil {
			log.Error().Msg("Error logging a habit")
		} else {
			log.Info().Msgf("Logged %s for today\n", args[0])
		}
	} else {
		log.Error().Msgf("Invalid habit: %s\n", args[0])
	}
}

func Undo(args []string) { // some chunks are duplicated from `Do()`
	if len(args) == 0 {
		log.Fatal().Msg("Habit must be specified")
	}

	habitName := args[0]

	// validate habit
	habits, err := Habit.GetHabits()
	if err != nil {
		log.Fatal().Msg("Error getting habits")
	}

	// untrack habit
	isValidHabit := slices.Contains(habits, habitName)
	if isValidHabit {
		today := time.Now().Format("2006-01-02")

		activity := Activity{
			Date: today, IsCompleted: 0, HabitName: habitName}

		err = Habit.Undo(activity)
		if err != nil {
			log.Error().Msg("Error undoing a habit")
		} else {
			log.Info().Msgf("Undo %s for today\n", args[0])
		}
	} else {
		log.Error().Msgf("Invalid habit: %s\n", args[0])
	}
}
