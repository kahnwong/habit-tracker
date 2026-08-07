CREATE TABLE habit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    is_completed INTEGER NOT NULL,
    habit_id INTEGER NOT NULL,
    FOREIGN KEY (habit_id) REFERENCES habit(id) ON DELETE CASCADE,
    UNIQUE (habit_id, date, is_completed)
);

CREATE INDEX idx_activity_habit_id ON activity (habit_id);
CREATE INDEX idx_activity_date ON activity (date);
