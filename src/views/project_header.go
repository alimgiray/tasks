package views

import (
	"github.com/charmbracelet/lipgloss"
)

func (m ProjectModel) generateHeader() string {
	projectHeaderStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#282A36")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Align(lipgloss.Center).
		Width(m.width)

	return projectHeaderStyle.Render("Projects")
}
