package screens

import (
	"fmt"
	"teapot/internal/models"
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddAnotherAppModel struct {
	options        []AddAnotherOption
	cursor         int
	appCount       int
	architectureType models.ArchitectureType
}

type AddAnotherOption struct {
	Key         string
	Name        string
	Description string
}

func NewAddAnotherAppModel(appCount int, architecture models.ArchitectureType) AddAnotherAppModel {
	options := []AddAnotherOption{
		{"add", "Add Another App", "Add another application to your monorepo"},
		{"continue", "Continue to Dev Tools", "Proceed with the current applications"},
	}

	// For single app architecture, don't allow adding more apps
	if architecture == models.ArchitectureSingle {
		options = []AddAnotherOption{
			{"continue", "Continue to Dev Tools", "Proceed with your single application"},
		}
	}

	return AddAnotherAppModel{
		options:          options,
		cursor:           0,
		appCount:         appCount,
		architectureType: architecture,
	}
}

func (m AddAnotherAppModel) Init() tea.Cmd {
	return nil
}

func (m AddAnotherAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.options[m.cursor].Key == "continue" {
				// Continue to dev tools
				return m, func() tea.Msg {
					return AddAnotherAppSelectedMsg{
						Action: "continue",
					}
				}
			} else {
				// Add another app
				return m, func() tea.Msg {
					return AddAnotherAppSelectedMsg{
						Action: "add",
					}
				}
			}
		}
	}
	return m, nil
}

func (m AddAnotherAppModel) View() string {
	subtitle := components.RenderSubtitle("Application Configuration Complete!")

	// Show current apps summary
	appsSection := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render(fmt.Sprintf("✓ %d application(s) configured", m.appCount))

	// Architecture note
	var architectureNote string
	if m.architectureType == models.ArchitectureSingle {
		architectureNote = lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Italic(true).
			Margin(1, 0, 1, 0).
			Render("Note: Single app architecture selected - only one app allowed")
	} else {
		architectureNote = lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Italic(true).
			Margin(1, 0, 1, 0).
			Render("You can add multiple applications to your monorepo")
	}

	// Options
	var choices string
	for i, option := range m.options {
		cursor := " "
		
		var checked string
		var optionStyle lipgloss.Style
		
		if option.Key == "continue" {
			// Continue option styling
			checked = "→"
		} else {
			// Add another app option
			checked = "+"
		}
		
		if m.cursor == i {
			optionStyle = styles.FocusedStyle
		} else {
			optionStyle = styles.UnselectedStyle
		}

		// Unified choice rendering
		choiceText := cursor + " " + checked + " " + option.Name
		choice := optionStyle.Render(choiceText)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(option.Description)

		choices += choice + "\n" + description + "\n\n"
	}

	return subtitle + "\n\n" + 
		   appsSection + "\n" + 
		   architectureNote + "\n" + 
		   choices
}

type AddAnotherAppSelectedMsg struct {
	Action string
}