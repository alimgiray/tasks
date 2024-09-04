package views

type Column struct {
	Name  string
	Width int
}

func (m DashboardModel) getColumns() []Column {
	return []Column{
		{
			Name:  "Task Key",
			Width: 15,
		},
		{
			Name:  "Deadline",
			Width: 20,
		},
		{
			Name:  "Assignee",
			Width: 22,
		},
		{
			Name: "Title",
			Width: func(m DashboardModel) int {
				return m.width -
					15 - // Task column
					20 - // Deadline column
					22 - // Assignee column
					15 - // Created column
					15 - // Status column
					12 // ???
			}(m),
		},
		{
			Name:  "Created",
			Width: 15,
		},
		{
			Name:  "Status",
			Width: 15, // "In Progress" is 11 characters long, so we use 15 to give it a little padding
		},
	}
}
