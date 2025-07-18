package screens

import (
	"teapot/internal/models"
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddAppsModel struct {
	apps     []models.AppType
	cursor   int
	selected models.AppType
}

func NewAddAppsModel() AddAppsModel {
	return AddAppsModel{
		apps: []models.AppType{
			models.AppTypeReact,
			models.AppTypeNext,
			models.AppTypeTanStack,
			models.AppTypeExpo,
			models.AppTypeNest,
			models.AppTypeBasicNode,
			"continue", // Special continue option
		},
		cursor:   0,
		selected: "",
	}
}

func (m AddAppsModel) Init() tea.Cmd {
	return nil
}

func (m AddAppsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.apps)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.apps[m.cursor] == "continue" {
				// Continue with selected app
				if m.selected != "" {
					return m, func() tea.Msg {
						return AppTypeSelectedMsg{
							AppType: m.selected,
						}
					}
				}
				return m, nil // No selection made
			} else {
				// Select current app
				m.selected = m.apps[m.cursor]
				return m, nil
			}
		}
	}
	return m, nil
}

func (m AddAppsModel) View() string {
	subtitle := components.RenderSubtitle("Choose your application type:")

	var choices string
	for i, app := range m.apps {
		cursor := " "
		
		var checked string
		var optionStyle lipgloss.Style
		var descriptionText string
		var name string
		
		if app == "continue" {
			// Continue option styling
			checked = "→"
			name = "Continue"
			descriptionText = "Proceed with selected application"
			if m.cursor == i {
				optionStyle = styles.FocusedStyle
			} else {
				optionStyle = styles.UnselectedStyle
			}
		} else {
			// Regular app option styling
			if m.selected == app {
				checked = "●"
				optionStyle = styles.CheckedStyle
			} else {
				checked = " "
				optionStyle = styles.UnselectedStyle
			}
			
			if m.cursor == i {
				optionStyle = styles.FocusedStyle
			}
			
			name = models.AppTypeNames[app]
			
			switch app {
			case models.AppTypeReact:
				descriptionText = "Client-side React application"
			case models.AppTypeNext:
				descriptionText = "Full-stack React framework"
			case models.AppTypeTanStack:
				descriptionText = "Modern full-stack React framework"
			case models.AppTypeExpo:
				descriptionText = "React Native mobile application"
			case models.AppTypeNest:
				descriptionText = "Scalable Node.js server framework"
			case models.AppTypeBasicNode:
				descriptionText = "Basic Node.js application"
			}
		}

		choice := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render(cursor) + " " +
			optionStyle.Render(checked+" "+name)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(descriptionText)

		choices += choice + "\n" + description + "\n"
	}

	return subtitle + "\n\n" + choices
}

type AppTypeSelectedMsg struct {
	AppType models.AppType
}