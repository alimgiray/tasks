package migrations

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/alimgiray/tasks/src/database"
)

type Migration interface {
	ShouldRun() bool
	Run() error
	Update() error
	GetDB() *database.Database
}

type MigrationBase struct {
	Name  string
	Query string
	DB    *database.Database
}

func (m *MigrationBase) ShouldRun() bool {
	_, err := m.DB.GetMigrationName(m.Name)
	if err != nil {
		if err == sql.ErrNoRows || strings.Contains(err.Error(), "no rows") {
			// The migration has not been run yet
			return true
		}
		// An actual error occurred
		log.Println(err)
		return false
	}
	// The migration has already been run
	return false
}

func (m *MigrationBase) Run() error {
	return m.DB.RunMigrationQuery(m.Query)
}

func (m *MigrationBase) Update() error {
	return m.DB.UpdateMigration(m.Name, time.Now())
}

func (m *MigrationBase) GetDB() *database.Database {
	return m.DB
}
