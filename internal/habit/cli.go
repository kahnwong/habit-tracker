package habit

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kahnwong/habit-tracker/internal/calendar"

	"github.com/rs/zerolog/log"
)

func Create(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("you must specify a habit")
	}

	err := Habit.CreateHabit(args[0])
	if err != nil {
		log.Info().Msgf("Habit %s already exists", args[0])
	} else {
		log.Info().Msgf("Habit %s created", args[0])
	}

	return nil
}

func Do(args []string) error {
	if err := validateHabit(args); err != nil {
		return err
	}

	date := time.Now().Format("2006-01-02")
	if len(args) == 2 {
		date = args[1]
	}

	activity := Activity{
		Date: date, IsCompleted: 1, HabitName: args[0]}

	err := Habit.Do(activity)
	if err != nil {
		return fmt.Errorf("error logging habit: %w", err)
	}

	log.Info().Msgf("Logged %s for %s", args[0], date)
	return nil
}

func Undo(args []string) error {
	if err := validateHabit(args); err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")

	activity := Activity{
		Date: today, IsCompleted: 0, HabitName: args[0]}

	err := Habit.Undo(activity)
	if err != nil {
		return fmt.Errorf("error undoing habit: %w", err)
	}

	log.Info().Msgf("Undo %s for today", args[0])
	return nil
}

func ShowHabitActivity(lookbackMonths int, args []string) error {
	if err := validateHabit(args); err != nil {
		return err
	}

	// fetch activities
	activities, err := Habit.GetHabitActivity(args[0], lookbackMonths)
	if err != nil {
		return fmt.Errorf("error fetching activities for habit %s: %w", args[0], err)
	}

	if len(activities) == 0 {
		log.Info().Msgf("No activities found for habit: %s", args[0])
		return nil
	}

	// parse date
	var dates []time.Time
	layout := "2006-01-02"
	for _, a := range activities {
		t, err := time.Parse(layout, a.Date)
		if err != nil {
			return fmt.Errorf("error parsing date %s: %w", a.Date, err)
		}
		dates = append(dates, t)
	}

	// render calendar
	calendar.RenderCalendarView(lookbackMonths, dates)
	return nil
}

func ShowPeriodActivity(period string) error {
	// fetch activities
	activities, dates, err := Habit.GetPeriodActivity(period)
	if err != nil {
		return fmt.Errorf("error fetching activities for period %s: %w", period, err)
	}

	if len(activities) == 0 {
		log.Info().Msgf("No activities found for period: %s", period)
		return nil
	}

	// render table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	//// set header
	var dateFormatted []string
	for _, d := range dates {
		v, err := time.Parse("2006-01-02", d)
		if err != nil {
			return fmt.Errorf("error parsing date %s: %w", d, err)
		}
		dateFormatted = append(dateFormatted, v.Format("Mon"))
	}

	headers := append([]string{""}, dateFormatted...)

	var headerRow table.Row
	for _, h := range headers {
		headerRow = append(headerRow, h)
	}
	t.AppendHeader(headerRow)

	//// unwind data
	////// %4s for center alignment, has to stay here because color package has fixed bytes
	isCompletedIcon := map[int64]string{0: fmt.Sprintf(" %s", " "), 1: fmt.Sprintf(" %s", "✓")}
	for _, activity := range activities {
		var elems []interface{}
		elems = append(elems, fmt.Sprintf("%-4s", activity["habit_name"])) // %-6s for left-alignment and padding

		for _, date := range dates {
			if intVal, ok := activity[date].(int64); ok {
				elems = append(elems, isCompletedIcon[intVal])
			} else {
				return fmt.Errorf("failed to cast to int, value is of type %T", activity[date])
			}
		}

		t.AppendRows([]table.Row{
			elems,
		})
	}

	////// styling
	//t.SetStyle(table.Style{
	//	Options: table.OptionsNoBordersAndSeparators,
	//})

	t.SetStyle(table.StyleColoredBlackOnCyanWhite)

	// render
	t.Render()
	return nil
}

// utils
func validateHabit(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("habit must be specified")
	}

	habitName := args[0]
	habits, err := Habit.GetHabits()
	if err != nil {
		return fmt.Errorf("error getting habits: %w", err)
	}

	if !slices.Contains(habits, habitName) {
		return fmt.Errorf("invalid habit: %s", habitName)
	}

	return nil
}
