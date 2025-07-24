package habit

type Activity struct {
	Date        string `db:"date"`
	IsCompleted int    `db:"is_completed"` // 0 for false, 1 for true
	HabitName   string `db:"habit_name"`
}
