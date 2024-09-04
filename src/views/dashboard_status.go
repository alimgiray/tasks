package views

import "github.com/charmbracelet/lipgloss"

// TODO maybe move this into DB and make them dynamic?
const (
	STATUS_PLANNING    = "Planning"
	STATUS_WAITING     = "Waiting"
	STATUS_IN_PROGRESS = "In Progress"
	STATUS_TESTING     = "Testing"
	STATUS_DONE        = "Done"
)

var statuses = []string{STATUS_PLANNING, STATUS_WAITING, STATUS_IN_PROGRESS, STATUS_TESTING, STATUS_DONE}

var statusStyles = map[string]lipgloss.Style{
	STATUS_PLANNING: lipgloss.NewStyle().
		Background(lipgloss.Color("#FF5F87")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Bold(true),

	STATUS_WAITING: lipgloss.NewStyle().
		Background(lipgloss.Color("#87D7FF")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Bold(true),

	STATUS_IN_PROGRESS: lipgloss.NewStyle().
		Background(lipgloss.Color("#FFD700")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Bold(true),

	STATUS_TESTING: lipgloss.NewStyle().
		Background(lipgloss.Color("#AF87FF")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Bold(true),

	STATUS_DONE: lipgloss.NewStyle().
		Background(lipgloss.Color("#5FD700")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Bold(true),
}

func NextStatus(current string) string {
	for i, status := range statuses {
		if status == current {
			return statuses[(i+1)%len(statuses)]
		}
	}
	return statuses[0]
}
