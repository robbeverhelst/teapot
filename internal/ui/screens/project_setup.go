package screens

import (
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectSetupModel struct {
	projectName   string
	description   string
	nameCursor    int
	descCursor    int
	focusedField  int // 0 = name, 1 = description
}

func NewProjectSetupModel() ProjectSetupModel {
	return ProjectSetupModel{
		projectName:  "",
		description:  "",
		nameCursor:   0,
		descCursor:   0,
		focusedField: 0,
	}
}

func (m ProjectSetupModel) Init() tea.Cmd {
	return nil
}

func (m ProjectSetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.projectName != "" {
				return m, func() tea.Msg {
					return ProjectSetupCompleteMsg{
						ProjectName: m.projectName,
						Description: m.description,
					}
				}
			}
		case "tab":
			m.focusedField = (m.focusedField + 1) % 2
		case "shift+tab":
			m.focusedField = (m.focusedField - 1 + 2) % 2
		case "backspace":
			if m.focusedField == 0 && m.nameCursor > 0 {
				m.projectName = m.projectName[:m.nameCursor-1] + m.projectName[m.nameCursor:]
				m.nameCursor--
			} else if m.focusedField == 1 && m.descCursor > 0 {
				m.description = m.description[:m.descCursor-1] + m.description[m.descCursor:]
				m.descCursor--
			}
		case "left":
			if m.focusedField == 0 && m.nameCursor > 0 {
				m.nameCursor--
			} else if m.focusedField == 1 && m.descCursor > 0 {
				m.descCursor--
			}
		case "right":
			if m.focusedField == 0 && m.nameCursor < len(m.projectName) {
				m.nameCursor++
			} else if m.focusedField == 1 && m.descCursor < len(m.description) {
				m.descCursor++
			}
		case "home":
			if m.focusedField == 0 {
				m.nameCursor = 0
			} else {
				m.descCursor = 0
			}
		case "end":
			if m.focusedField == 0 {
				m.nameCursor = len(m.projectName)
			} else {
				m.descCursor = len(m.description)
			}
		default:
			if len(msg.String()) == 1 {
				char := msg.String()
				if m.focusedField == 0 && isValidProjectNameCharSetup(char) {
					m.projectName = m.projectName[:m.nameCursor] + char + m.projectName[m.nameCursor:]
					m.nameCursor++
				} else if m.focusedField == 1 {
					m.description = m.description[:m.descCursor] + char + m.description[m.descCursor:]
					m.descCursor++
				}
			}
		}
	}
	return m, nil
}

func (m ProjectSetupModel) View() string {
	subtitle := components.RenderSubtitle("Let's set up your project!")

	// Project name field
	nameLabel := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("📝 Project name: (required)")

	nameValue := m.projectName
	if m.focusedField == 0 {
		if m.nameCursor < len(nameValue) {
			nameValue = nameValue[:m.nameCursor] + "_" + nameValue[m.nameCursor+1:]
		} else {
			nameValue += "_"
		}
	}

	nameBorderColor := styles.ColorBorderPrimary
	if m.focusedField == 0 {
		nameBorderColor = styles.ColorBorderAccent
	}

	nameBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(nameBorderColor).
		Foreground(styles.ColorTextPrimary).
		Padding(0, 1).
		Width(35).
		Margin(1, 0, 0, 0).
		Render(nameValue)

	// Description field
	descLabel := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("📄 Description: (optional)")

	descValue := m.description
	if m.focusedField == 1 {
		if m.descCursor < len(descValue) {
			descValue = descValue[:m.descCursor] + "_" + descValue[m.descCursor+1:]
		} else {
			descValue += "_"
		}
	}

	descBorderColor := styles.ColorBorderPrimary
	if m.focusedField == 1 {
		descBorderColor = styles.ColorBorderAccent
	}

	descBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(descBorderColor).
		Foreground(styles.ColorTextPrimary).
		Padding(0, 1).
		Width(35).
		Margin(1, 0, 0, 0).
		Render(descValue)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Margin(1, 0, 0, 0).
		Render("Tab: switch fields")

	return subtitle + "\n\n" + 
		   nameLabel + "\n" + nameBox + "\n\n" + 
		   descLabel + "\n" + descBox + "\n\n" + 
		   instructions
}

func isValidProjectNameCharSetup(char string) bool {
	if len(char) != 1 {
		return false
	}
	c := char[0]
	return (c >= 'a' && c <= 'z') || 
		   (c >= 'A' && c <= 'Z') || 
		   (c >= '0' && c <= '9') || 
		   c == '-' || c == '_'
}

type ProjectSetupCompleteMsg struct {
	ProjectName string
	Description string
}