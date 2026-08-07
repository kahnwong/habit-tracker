-- +goose Up
CREATE TABLE activity_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    is_completed INTEGER NOT NULL,
    habit_id INTEGER NOT NULL,
    FOREIGN KEY (habit_id) REFERENCES habit(id) ON DELETE CASCADE,
    UNIQUE (habit_id, date, is_completed)
);

INSERT INTO activity_new (id, date, is_completed, habit_id)
SELECT a.id, a.date, a.is_completed, h.id
FROM activity AS a
JOIN habit AS h ON h.name = a.habit_name;

DROP TABLE activity;
ALTER TABLE activity_new RENAME TO activity;

CREATE INDEX idx_activity_habit_id ON activity (habit_id);
CREATE INDEX idx_activity_date ON activity (date);

-- +goose Down
CREATE TABLE activity_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    is_completed INTEGER NOT NULL,
    habit_name TEXT NOT NULL,
    FOREIGN KEY (habit_name) REFERENCES habit(name) ON DELETE CASCADE,
    UNIQUE (habit_name, date, is_completed)
);

INSERT INTO activity_old (id, date, is_completed, habit_name)
SELECT a.id, a.date, a.is_completed, h.name
FROM activity AS a
JOIN habit AS h ON h.id = a.habit_id;

DROP TABLE activity;
ALTER TABLE activity_old RENAME TO activity;

CREATE INDEX idx_activity_habit_name ON activity (habit_name);
CREATE INDEX idx_activity_date ON activity (date);
