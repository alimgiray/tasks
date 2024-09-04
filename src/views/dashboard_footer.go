package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m DashboardModel) getFooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#282A36")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Align(lipgloss.Center).
		Width(m.width)
}

func (m DashboardModel) getEmptyRows(currentScreen string) string {
	remainingLines := m.height - lipgloss.Height(currentScreen) - 1
	emptyLines := ""
	if remainingLines > 0 {
		emptyLines = (strings.Repeat("\n", remainingLines))
	}
	return emptyLines
}

func (m DashboardModel) generateFooter() string {
	footerStyle := m.getFooterStyle()
	footer := footerStyle.Render("p: Projects | n: New Task | arrows: Navigation | s: Status | a: Assign | d: Deadline | e: Description | q: Quit")
	if m.creatingTask {
		footer = footerStyle.Render("enter: Next Field | ctrl+c: Abort")
	}
	if m.isEditing {
		footer = footerStyle.Render("enter: Complete Editing | ctrl+c: Abort")
	}
	return footer
}
