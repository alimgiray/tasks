package migrations

func Run(migrations []Migration) {
	if len(migrations) > 0 {
		migrations = append([]Migration{NewMigrationCreateMigration(migrations[0].GetDB())}, migrations...) // Mandatory. It creates migrations table.
	}
	for _, m := range migrations {
		if m.ShouldRun() {
			if err := m.Run(); err != nil {
				panic(err)
			}
			if err := m.Update(); err != nil {
				panic(err)
			}
		}
	}
}
