package sqlite

import "context"

func (db *DB) MigrateJobTable() error {
	return db.db.Exec(context.Background(), `
		CREATE TABLE Job 
			(
				id INTEGER PRIMARY KEY AUTOINCREMENT,

				repositoryName TEXT NOT NULL,
				repositoryUrl  TEXT NOT NULL,
				findings       TEXT,

				status      INT,
				queuedAt    DATE,
				scanningAt  DATE,
				finnishedAt DATE
			)
	`)
}
