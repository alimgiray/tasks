package views

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/alimgiray/tasks/src/database"
	"github.com/alimgiray/tasks/src/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxKeyLength      = 10
	maxAssigneeLength = 10
	maxTitleLength    = 40
)

type DashboardModel struct {
	tasks      []*models.Task
	currentRow int

	width  int
	height int
	cursor Cursor

	// Creating new task related fields
	newTask      *models.Task // a new task currently being created
	creatingTask bool         // whether we are in task creation mode
	fieldIndex   int          // which column we are creating or editing?
	projectID    int

	// Editing related fields
	isEditing   bool
	inputBuffer string
}

func InitialDashboard(width, height int, projectID int) DashboardModel {
	return DashboardModel{
		tasks:      database.LoadTasksFromDB(projectID),
		currentRow: 0,
		width:      width,
		height:     height,
		cursor:     NewCursor("|"),
		projectID:  projectID,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return m.cursor.Blink()
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle key presses
	case tea.KeyMsg:
		// Task creating mode allows users to enter text
		// so it is different from "view only" mode
		if m.creatingTask {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.creatingTask = false
				m.inputBuffer = ""
				m.fieldIndex = 0
			case "enter":
				switch m.fieldIndex {
				case 0: // Task Key
					m.newTask.Key = m.inputBuffer[0:int(math.Min(float64(len(m.inputBuffer)), maxKeyLength))]
				case 1: // Deadline
					deadline, _ := time.Parse("02/01/2006", m.inputBuffer)
					m.newTask.Deadline = deadline
				case 2: // Assignee
					m.newTask.Assignee = m.inputBuffer[0:int(math.Min(float64(len(m.inputBuffer)), maxAssigneeLength))]
				case 3: // Title
					m.newTask.Title = m.inputBuffer[0:int(math.Min(float64(len(m.inputBuffer)), maxTitleLength))]

					// Set default fields
					m.newTask.Status = "Planning"    // Set default status
					m.newTask.CreatedAt = time.Now() // Set created date
					m.newTask.ProjectID = m.projectID

					err := database.DB.CreateTask(m.newTask)
					if err != nil {
						// TODO find a way to handle errors
						log.Fatalf("Failed to create task: %v", err)
					}
					m.tasks = append(m.tasks, m.newTask)

					// End task creation
					m.creatingTask = false
					m.currentRow = len(m.tasks) - 1
					m.inputBuffer = ""
					m.fieldIndex = 0 // Reset field index, because it is end of creating
				}
				m.fieldIndex++
				m.inputBuffer = ""
			case "backspace":
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			default:
				m.inputBuffer += msg.String()
			}
		} else if m.isEditing {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.isEditing = false
				m.inputBuffer = ""
			case "enter":
				if m.fieldIndex == 0 { // Task key
					m.tasks[m.currentRow].Key = m.inputBuffer[0:int(math.Min(float64(len(m.inputBuffer)), maxKeyLength))]
				}
				if m.fieldIndex == 1 { // Deadline
					deadline, _ := time.Parse("02/01/2006", m.inputBuffer)
					m.tasks[m.currentRow].Deadline = deadline
				}
				if m.fieldIndex == 2 { // Assignee
					m.tasks[m.currentRow].Assignee = m.inputBuffer[0:int(math.Min(float64(len(m.inputBuffer)), maxAssigneeLength))]
				}
				if m.fieldIndex == 3 { // Title
					m.tasks[m.currentRow].Title = m.inputBuffer[0:int(math.Min(float64(len(m.inputBuffer)), maxTitleLength))]
				}
				if m.fieldIndex == 4 { // Description
					m.tasks[m.currentRow].Description = m.inputBuffer
				}

				err := database.DB.UpdateTask(m.tasks[m.currentRow])
				if err != nil {
					// TODO find a way to handle errors
					log.Fatalf("Failed to create task: %v", err)
				}

				m.isEditing = false
				m.inputBuffer = ""
			case "backspace":
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			default:
				m.inputBuffer += msg.String()
			}
		} else {
			// Viewing tasks
			switch msg.String() {
			case "ctrl+c", "esc", "q":
				return m, tea.Quit
			case "up":
				if m.currentRow > 0 {
					m.currentRow--
				}
			case "down":
				if m.currentRow < len(m.tasks)-1 {
					m.currentRow++
				}
			case "n":
				m.creatingTask = true
				m.newTask = &models.Task{CreatedAt: time.Now()}
				m.fieldIndex = 0
				m.inputBuffer = ""
				m.currentRow = len(m.tasks) // Move the cursor to the new task row
			case "right", "left":
				m.tasks[m.currentRow].Expanded = !m.tasks[m.currentRow].Expanded
			case "p":
				return InitialProjectPage(m.width, m.height), tea.ClearScreen
			case "s":
				m.tasks[m.currentRow].Status = NextStatus(m.tasks[m.currentRow].Status)
			case "k":
				// Edit key field
				m.isEditing = true
				m.fieldIndex = 0
			case "d":
				// Edit deadline field
				m.isEditing = true
				m.fieldIndex = 1
			case "a":
				// Edit assignee field
				m.isEditing = true
				m.fieldIndex = 2
			case "t":
				// Edit title field
				m.isEditing = true
				m.fieldIndex = 3
			case "e":
				// Edit description field
				m.isEditing = true
				m.fieldIndex = 4
				// Expand the task if it's not already expanded, then start editing
				if !m.tasks[m.currentRow].Expanded {
					m.tasks[m.currentRow].Expanded = true
				}
				m.inputBuffer = m.tasks[m.currentRow].Description
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}

	return m, m.cursor.Update(msg)
}

func (m DashboardModel) View() string {
	var b strings.Builder

	// Render header
	b.WriteString(m.generateHeader() + "\n")

	// Render tasks
	for i, task := range m.tasks {
		style := m.getCurrentRowStyle(i)

		// Just displaying tasks in view mode
		// Or in editing, but this is not the row we are editing
		if !m.isEditing || m.currentRow != i {
			taskInfo := m.createTaskLine(task, -1)
			b.WriteString(" ")
			b.WriteString(style.Render(taskInfo) + "\n")
			// If expanded, render description
			if task.Expanded {
				descriptionStyle := m.getDescriptionStyle()
				b.WriteString(descriptionStyle.Render(task.Description) + "\n")
			}
		}

		// Rendering currently edited row
		if m.isEditing && m.currentRow == i {
			taskInfo := m.createTaskLine(task, m.fieldIndex)
			b.WriteString(" ")
			b.WriteString(style.Render(taskInfo) + "\n")
			// If expanded, render description
			if task.Expanded {
				descriptionStyle := m.getDescriptionStyle()
				if m.fieldIndex == 4 { // Editing description
					b.WriteString(descriptionStyle.Render(m.inputBuffer+m.cursor.Render()) + "\n")
				} else {
					b.WriteString(descriptionStyle.Render(task.Description) + "\n")
				}
			}
		}
	}

	// Creating new task
	if m.creatingTask {
		m.renderNewTaskLine(&b)
	}

	// Render footer
	emptyRows := m.getEmptyRows(b.String())
	b.WriteString(emptyRows) // Fill gaps
	b.WriteString(m.generateFooter())

	return b.String()
}

func (m DashboardModel) getRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Width(m.width - 5)
}

func (m DashboardModel) getSelectedRowStyle() lipgloss.Style {
	return m.getRowStyle().Background(lipgloss.Color("#333333"))
}

func (m DashboardModel) getCurrentRowStyle(i int) lipgloss.Style {
	style := m.getRowStyle()

	if m.currentRow == i && !m.creatingTask {
		style = m.getSelectedRowStyle()
	}

	return style
}

func (m DashboardModel) getDescriptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().Margin(0, 3)
}

func (m DashboardModel) createTaskLine(task *models.Task, index int) string {
	columns := m.getColumns()

	format := ""
	args := []interface{}{}

	for i, column := range columns {
		isEditable := i >= 0 && i <= 3
		if m.isEditing && index == i && isEditable { // If editing and it's the current editable index
			args = append(args, column.Width, m.inputBuffer+m.cursor.Render())
		} else {
			switch column.Name {
			case "Task Key":
				args = append(args, column.Width, task.Key)
			case "Deadline":
				args = append(args, column.Width, task.Deadline.Format("02 Jan 2006"))
			case "Assignee":
				args = append(args, column.Width, task.Assignee)
			case "Title":
				args = append(args, column.Width, task.Title)
			case "Created":
				args = append(args, column.Width, task.CreatedAt.Format("02 Jan 2006"))
			case "Status":
				statusStyle := statusStyles[task.Status]
				statusLabel := statusStyle.Width(column.Width).Align(lipgloss.Center).Render(task.Status)
				args = append(args, column.Width, statusLabel)
			}
		}
		format += "%-*s "
	}

	format = format[:len(format)-1]
	return fmt.Sprintf(format, args...)
}

func (m DashboardModel) renderNewTaskLine(b *strings.Builder) {
	if m.creatingTask {
		// Render the pink indicator
		pinkIndicator := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Render("> ")
		b.WriteString(pinkIndicator)

		columns := m.getColumns()
		fields := []string{
			m.newTask.Key,
			" dd/mm/yyyy",
			"  " + m.newTask.Assignee,
			m.newTask.Title,
		}

		// Format the deadline if it's set
		if !m.newTask.Deadline.IsZero() {
			fields[1] = " " + m.newTask.Deadline.Format("02 Jan 2006")
		}

		// Loop through the fields and render them
		for idx, field := range fields {
			if idx == m.fieldIndex {
				spacing := ""
				switch columns[idx].Name {
				case "Deadline":
					spacing = " " // Move cursor one space to the right for the deadline
				case "Assignee":
					spacing = "  " // Move cursor two spaces to the right for the assignee
				case "Title":
					spacing = "   " // Move cursor three spaces to the right for the title
				}
				b.WriteString(spacing + fmt.Sprintf("%-*s", columns[idx].Width, m.inputBuffer+m.cursor.Render()))
			} else {
				b.WriteString(fmt.Sprintf("%-*s", columns[idx].Width, field))
			}
		}

		// Render the "Created" and "Status" columns
		createdColumn := columns[len(columns)-2] // Second last column should be "Created"
		statusColumn := columns[len(columns)-1]  // Last column should be "Status"

		b.WriteString(fmt.Sprintf(" %-*s %-*s\n", createdColumn.Width, "", statusColumn.Width, ""))
	}
}
