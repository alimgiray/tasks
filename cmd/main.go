package main

import (
	"fmt"

	"github.com/alimgiray/tasks/src/database"
	"github.com/alimgiray/tasks/src/migrations"
	"github.com/alimgiray/tasks/src/views"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	db := database.Get()

	migrations.Run([]migrations.Migration{
		migrations.NewMigrationCreateProjects(db),
		migrations.NewMigrationCreateTasks(db),
	})
}

func main() {
	p := tea.NewProgram(views.InitialProjectPage(1, 0), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
