package database

import "time"

func (db *Database) GetMigrationName(migrationName string) (string, error) {
	var name string
	err := db.db.QueryRow("SELECT name FROM migrations WHERE name = $1", migrationName).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (db *Database) RunMigrationQuery(query string) error {
	_, err := db.db.Exec(query)
	return err
}

func (db *Database) UpdateMigration(name string, updatedAt time.Time) error {
	_, err := db.db.Exec("INSERT INTO migrations (name, completed_at) VALUES ($1, $2)", name, updatedAt)
	return err
}
