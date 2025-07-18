package screens

import (
	"fmt"
	"strings"

	"teapot/internal/generator"
	"teapot/internal/models"
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type YAMLPreviewModel struct {
	project    models.ProjectConfig
	yamlContent string
	scrollOffset int
	maxLines    int
	options     []string
	cursor      int
}

func NewYAMLPreviewModel(project models.ProjectConfig) YAMLPreviewModel {
	yamlContent, err := generator.GenerateTeapotYAML(project)
	if err != nil {
		yamlContent = fmt.Sprintf("Error generating YAML: %v", err)
	}

	return YAMLPreviewModel{
		project:     project,
		yamlContent: yamlContent,
		scrollOffset: 0,
		maxLines:    20, // Will be adjusted based on terminal height
		options: []string{
			"Save teapot.yml",
			"Continue to Generation",
			"Back to Edit",
		},
		cursor: 1, // Default to "Continue to Generation"
	}
}

func (m YAMLPreviewModel) Init() tea.Cmd {
	return nil
}

func (m YAMLPreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "ctrl+j":
			// Scroll YAML content down
			yamlLines := strings.Split(m.yamlContent, "\n")
			if m.scrollOffset < len(yamlLines)-m.maxLines {
				m.scrollOffset++
			}
		case "ctrl+k":
			// Scroll YAML content up
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "enter":
			switch m.cursor {
			case 0: // Save teapot.yml
				return m, func() tea.Msg {
					return YAMLSaveMsg{Project: m.project}
				}
			case 1: // Continue to Generation
				return m, func() tea.Msg {
					return YAMLContinueMsg{Project: m.project}
				}
			case 2: // Back to Edit
				return m, func() tea.Msg {
					return YAMLBackMsg{}
				}
			}
		}
	}
	return m, nil
}

func (m YAMLPreviewModel) View() string {
	subtitle := components.RenderSubtitle("📄 Configuration Preview")

	// Project summary
	summary := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render(fmt.Sprintf("✓ Project: %s", m.project.Name))

	// YAML content box
	yamlLines := strings.Split(m.yamlContent, "\n")
	displayLines := yamlLines
	if len(yamlLines) > m.maxLines {
		end := m.scrollOffset + m.maxLines
		if end > len(yamlLines) {
			end = len(yamlLines)
		}
		displayLines = yamlLines[m.scrollOffset:end]
	}

	yamlBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderNeon).
		Background(styles.ColorBgSecondary).
		Foreground(styles.ColorTextSecondary).
		Padding(1, 2).
		Width(70).
		Height(m.maxLines + 2).
		Margin(1, 0).
		Render(strings.Join(displayLines, "\n"))

	// Scroll indicator
	scrollInfo := ""
	if len(yamlLines) > m.maxLines {
		scrollInfo = lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Italic(true).
			Render(fmt.Sprintf("(Showing %d-%d of %d lines - Use Ctrl+J/K to scroll)", 
				m.scrollOffset+1, 
				min(m.scrollOffset+m.maxLines, len(yamlLines)), 
				len(yamlLines)))
	}

	// Options
	var choices string
	for i, option := range m.options {
		cursor := " "
		var optionStyle lipgloss.Style
		var icon string

		switch i {
		case 0: // Save
			icon = "💾"
		case 1: // Continue
			icon = "🚀"
		case 2: // Back
			icon = "←"
		}

		if m.cursor == i {
			optionStyle = styles.FocusedStyle
		} else {
			optionStyle = styles.UnselectedStyle
		}

		choiceText := cursor + " " + icon + " " + option
		choice := optionStyle.Render(choiceText)
		choices += choice + "\n"
	}

	// Help text
	helpText := components.RenderHelp("↑↓: navigate • enter: select • ctrl+j/k: scroll YAML • backspace: back")

	return subtitle + "\n\n" + 
		   summary + "\n" + 
		   yamlBox + "\n" + 
		   scrollInfo + "\n\n" + 
		   choices + "\n\n" + 
		   helpText
}

// Message types for YAML preview actions
type YAMLSaveMsg struct {
	Project models.ProjectConfig
}

type YAMLContinueMsg struct {
	Project models.ProjectConfig
}

type YAMLBackMsg struct{}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}