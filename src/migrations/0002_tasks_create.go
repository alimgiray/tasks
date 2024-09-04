package migrations

import "github.com/alimgiray/tasks/src/database"

type MigrationCreateTasks struct {
	MigrationBase
}

func NewMigrationCreateTasks(db *database.Database) Migration {
	return &MigrationCreateTasks{
		MigrationBase{
			Name: "0002_tasks_create",
			Query: `
				CREATE TABLE IF NOT EXISTS tasks (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					key TEXT NOT NULL,
					deadline DATETIME NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					assignee TEXT NOT NULL,
					title TEXT NOT NULL,
					status TEXT NOT NULL,
					description TEXT,
					project_id INTEGER,
					FOREIGN KEY (project_id) REFERENCES projects(id)
				);

				CREATE TRIGGER IF NOT EXISTS update_tasks_updated_at
				AFTER UPDATE ON tasks
				FOR EACH ROW
				BEGIN
					UPDATE tasks
					SET updated_at = CURRENT_TIMESTAMP
					WHERE id = OLD.id;
				END;
			`,
			DB: db,
		},
	}
}
