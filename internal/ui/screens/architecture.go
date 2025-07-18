package screens

import (
	"teapot/internal/models"
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ArchitectureModel struct {
	options []models.ArchitectureType
	cursor  int
}

func NewArchitectureModel() ArchitectureModel {
	return ArchitectureModel{
		options: []models.ArchitectureType{
			models.ArchitectureTurborepo,
			models.ArchitectureSingle,
			models.ArchitectureNx,
			"continue", // Special continue option
		},
		cursor: 0,
	}
}

func (m ArchitectureModel) Init() tea.Cmd {
	return nil
}

func (m ArchitectureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.options[m.cursor] != models.ArchitectureNx { // Nx is coming soon
				return m, func() tea.Msg {
					return ArchitectureSelectedMsg{
						Architecture: m.options[m.cursor],
					}
				}
			}
		}
	}
	return m, nil
}

func (m ArchitectureModel) View() string {
	subtitle := components.RenderSubtitle("Choose your project architecture:")

	var choices string
	for i, option := range m.options {
		cursor := " "
		
		checked := " "
		if m.cursor == i {
			checked = "●"
		}

		var optionStyle lipgloss.Style
		var descriptionText string
		
		switch option {
		case models.ArchitectureTurborepo:
			optionStyle = styles.SelectedStyle
			descriptionText = "Fast, incremental builds with caching"
		case models.ArchitectureSingle:
			optionStyle = styles.UnselectedStyle
			descriptionText = "Simple single application setup"
		case models.ArchitectureNx:
			optionStyle = styles.UnselectedStyle.Copy().Foreground(styles.ColorTextMuted)
			descriptionText = "Enterprise-grade monorepo tools"
		}

		if m.cursor == i {
			optionStyle = styles.FocusedStyle
		}

		name := models.ArchitectureNames[option]
		if option == models.ArchitectureNx {
			name += " 🚧"
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

type ArchitectureSelectedMsg struct {
	Architecture models.ArchitectureType
}