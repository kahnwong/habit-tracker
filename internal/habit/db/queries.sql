-- name: CreateHabit :exec
INSERT INTO habit (name) VALUES (sqlc.arg(name));

-- name: ListHabits :many
SELECT name FROM habit ORDER BY name;

-- name: CreateActivity :exec
INSERT OR IGNORE INTO activity (date, is_completed, habit_id)
SELECT sqlc.arg(date), sqlc.arg(is_completed), id
FROM habit
WHERE name = sqlc.arg(habit_name);

-- name: DeleteActivity :exec
DELETE FROM activity
WHERE date = sqlc.arg(date)
  AND habit_id = (SELECT id FROM habit WHERE name = sqlc.arg(habit_name));

-- name: ListCompletedHabitActivities :many
SELECT a.date, a.is_completed, h.name AS habit_name
FROM activity AS a
JOIN habit AS h ON h.id = a.habit_id
WHERE a.is_completed = 1
  AND a.date >= sqlc.arg(lookback_start)
  AND h.name = sqlc.arg(habit_name)
ORDER BY date;
