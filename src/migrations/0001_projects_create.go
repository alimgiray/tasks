package migrations

import "github.com/alimgiray/tasks/src/database"

type MigrationCreateProjects struct {
	MigrationBase
}

func NewMigrationCreateProjects(db *database.Database) Migration {
	return &MigrationCreateProjects{
		MigrationBase{
			Name: "0001_projects_create",
			Query: `
				CREATE TABLE IF NOT EXISTS projects (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
				
				CREATE TRIGGER IF NOT EXISTS update_projects_updated_at
				AFTER UPDATE ON projects
				FOR EACH ROW
				BEGIN
					UPDATE projects
					SET updated_at = CURRENT_TIMESTAMP
					WHERE id = OLD.id;
				END;
				`,
			DB: db,
		},
	}
}
