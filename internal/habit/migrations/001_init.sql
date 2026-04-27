-- +goose Up
CREATE TABLE habit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,  -- YYYY-MM-DD format (e.g., '2023-10-27')
    is_completed INTEGER NOT NULL, -- 0 for false, 1 for true (boolean)
    habit_name TEXT NOT NULL,
    FOREIGN KEY (habit_name) REFERENCES habit(name) ON DELETE CASCADE,
    UNIQUE (habit_name, date, is_completed)
);
CREATE INDEX idx_activity_habit_name ON activity (habit_name);
CREATE INDEX idx_activity_date ON activity (date);

-- +goose Down
DROP TABLE activity;
DROP TABLE habit;
