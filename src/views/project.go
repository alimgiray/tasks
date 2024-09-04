package views

import (
	"log"
	"strings"

	"github.com/alimgiray/tasks/src/database"
	"github.com/alimgiray/tasks/src/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectModel struct {
	projects []*models.Project

	currentRow int

	width  int
	height int
	cursor Cursor

	newProject      *models.Project // a new project currently being created
	creatingProject bool            // whether we are in project creation mode

	// Editing related fields
	isEditing   bool
	inputBuffer string
}

func InitialProjectPage(width, height int) ProjectModel {
	return ProjectModel{
		projects:   database.LoadProjectsFromDB(),
		currentRow: 0,
		width:      width,
		height:     height,
		cursor:     NewCursor("|"),
	}
}

func (m ProjectModel) Init() tea.Cmd {
	return m.cursor.Blink()
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle key presses
	case tea.KeyMsg:
		if m.creatingProject {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.creatingProject = false
				m.inputBuffer = ""
			case "enter":
				m.newProject.Name = m.inputBuffer
				projectID, err := database.DB.CreateProject(m.newProject)
				if err != nil {
					log.Fatalf("Failed to create project: %v", err)
				}
				m.newProject.ProjectID = projectID
				m.projects = append(m.projects, m.newProject)

				// End project creation
				m.creatingProject = false
				m.currentRow = len(m.projects) - 1
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
				m.projects[m.currentRow].Name = m.inputBuffer
				err := database.DB.UpdateProject(m.projects[m.currentRow])
				if err != nil {
					log.Fatalf("Failed to update project: %v", err)
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
			switch msg.String() {
			case "ctrl+c", "esc", "q":
				return m, tea.Quit
			case "up":
				if m.currentRow > 0 {
					m.currentRow--
				}
			case "down":
				if m.currentRow < len(m.projects)-1 {
					m.currentRow++
				}
			case "enter":
				projectID := m.projects[m.currentRow].ProjectID
				return InitialDashboard(m.width, m.height, projectID), tea.ClearScreen
			case "n":
				m.creatingProject = true
				m.newProject = &models.Project{Name: ""}
				m.inputBuffer = ""
				m.currentRow = len(m.projects)
			case "e":
				m.isEditing = true
				m.inputBuffer = m.projects[m.currentRow].Name
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}

	return m, m.cursor.Update(msg)
}

func (m ProjectModel) View() string {
	var b strings.Builder

	// Render header
	b.WriteString(m.generateHeader() + "\n")

	for i, project := range m.projects {
		style := m.getCurrentRowStyle(i)

		// Just displaying projects in view mode
		// Or in editing, but this is not the row we are editing
		if !m.isEditing || m.currentRow != i {
			b.WriteString(" ")
			b.WriteString(style.Render(project.Name) + "\n")
		}

		// Rendering currently edited row
		if m.isEditing && m.currentRow == i {
			b.WriteString(" ")
			b.WriteString(style.Render(m.inputBuffer+m.cursor.Render()) + "\n")
		}

	}

	// Creating new project
	if m.creatingProject {
		m.renderNewProjectLine(&b)
	}

	// Render footer
	emptyRows := m.getEmptyRows(b.String())
	b.WriteString(emptyRows) // Fill gaps
	b.WriteString(m.generateFooter())

	return b.String()
}

func (m ProjectModel) getRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Width(m.width - 2)
}

func (m ProjectModel) getSelectedRowStyle() lipgloss.Style {
	return m.getRowStyle().Background(lipgloss.Color("#333333"))
}

func (m ProjectModel) getCurrentRowStyle(i int) lipgloss.Style {
	style := m.getRowStyle()

	if m.currentRow == i && !m.creatingProject {
		style = m.getSelectedRowStyle()
	}

	return style
}

func (m ProjectModel) renderNewProjectLine(b *strings.Builder) {
	if m.creatingProject {
		// Render the pink indicator
		pinkIndicator := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Render("> ")
		b.WriteString(pinkIndicator + m.inputBuffer + m.cursor.Render())
	}
}
