CREATE TABLE habit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    is_completed INTEGER NOT NULL,
    habit_name TEXT NOT NULL,
    FOREIGN KEY (habit_name) REFERENCES habit(name) ON DELETE CASCADE,
    UNIQUE (habit_name, date, is_completed)
);

CREATE INDEX idx_activity_habit_name ON activity (habit_name);
CREATE INDEX idx_activity_date ON activity (date);
