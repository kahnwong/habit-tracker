-- name: CreateHabit :exec
INSERT INTO habit (name) VALUES (sqlc.arg(name));

-- name: ListHabits :many
SELECT name FROM habit ORDER BY name;

-- name: CreateActivity :exec
INSERT OR IGNORE INTO activity (date, is_completed, habit_name)
VALUES (sqlc.arg(date), sqlc.arg(is_completed), sqlc.arg(habit_name));

-- name: DeleteActivity :exec
DELETE FROM activity
WHERE date = sqlc.arg(date)
  AND habit_name = sqlc.arg(habit_name);

-- name: ListCompletedHabitActivities :many
SELECT date, is_completed, habit_name
FROM activity
WHERE is_completed = 1
  AND date >= sqlc.arg(lookback_start)
  AND habit_name = sqlc.arg(habit_name)
ORDER BY date;
