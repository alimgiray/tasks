package migrations

import "github.com/alimgiray/tasks/src/database"

type MigrationCreateMigration struct {
	MigrationBase
}

// Always run this one because it is the first migration
func (m *MigrationCreateMigration) ShouldRun() bool {
	return true
}

// No need to update, it creates the table.
func (m *MigrationCreateMigration) Update() error {
	return nil
}

func NewMigrationCreateMigration(db *database.Database) Migration {
	return &MigrationCreateMigration{
		MigrationBase{
			Name: "0000_migration_create",
			Query: `
				CREATE TABLE IF NOT EXISTS migrations (
					name TEXT PRIMARY KEY,
					completed_at TIMESTAMP
				);`,
			DB: db,
		},
	}
}
