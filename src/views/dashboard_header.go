package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var headerStyle = lipgloss.NewStyle().Bold(true).Padding(0, 2)

func (m DashboardModel) generateHeader() string {
	header := ""
	columns := m.getColumns()

	args := []interface{}{}

	for _, column := range columns {
		header += "%-*s "
		args = append(args, column.Width, column.Name)
	}

	header = header[:len(header)-1]
	return headerStyle.Render(fmt.Sprintf(header, args...))
}
