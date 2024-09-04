package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m ProjectModel) getFooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#282A36")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Align(lipgloss.Center).
		Width(m.width)
}

func (m ProjectModel) getEmptyRows(currentScreen string) string {
	remainingLines := m.height - lipgloss.Height(currentScreen) - 1
	emptyLines := ""
	if remainingLines > 0 {
		emptyLines = (strings.Repeat("\n", remainingLines))
	}
	return emptyLines
}

func (m ProjectModel) generateFooter() string {
	footerStyle := m.getFooterStyle()
	footer := footerStyle.Render("n: New Project | arrows: Navigation | q: Quit")
	if m.creatingProject {
		footer = footerStyle.Render("enter: Complete | ctrl+c: Abort")
	}
	if m.isEditing {
		footer = footerStyle.Render("enter: Complete Editing | ctrl+c: Abort")
	}
	return footer
}
