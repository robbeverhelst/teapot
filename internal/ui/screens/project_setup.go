package screens

import (
	"fmt"
	
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"
	"teapot/internal/validation"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectSetupModel struct {
	nameInput    textinput.Model
	descInput    textinput.Model
	focusedField int // 0 = name, 1 = description
	validationErr error
}

func NewProjectSetupModel() ProjectSetupModel {
	// Create name input with validation
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter project name (e.g., my-awesome-project)"
	nameInput.Focus()
	nameInput.CharLimit = 50
	nameInput.Width = 35

	// Create description input
	descInput := textinput.New()
	descInput.Placeholder = "Enter description (optional)"
	descInput.CharLimit = 200
	descInput.Width = 35

	return ProjectSetupModel{
		nameInput:     nameInput,
		descInput:     descInput,
		focusedField:  0,
		validationErr: nil,
	}
}

func (m ProjectSetupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ProjectSetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusedField == 0 {
				// Validate project name before continuing
				projectName := m.nameInput.Value()
				if projectName == "" {
					m.validationErr = validation.ValidateProjectName(projectName)
					return m, nil
				}
				
				if err := validation.ValidateProjectName(projectName); err != nil {
					m.validationErr = err
					return m, nil
				}
				
				// If name is valid, move to description field
				m.nameInput.Blur()
				m.descInput.Focus()
				m.focusedField = 1
				m.validationErr = nil
				return m, nil
			} else {
				// Validate description and complete setup
				description := m.descInput.Value()
				if err := validation.ValidateProjectDescription(description); err != nil {
					m.validationErr = err
					return m, nil
				}
				
				// Both fields are valid, proceed
				m.validationErr = nil
				return m, func() tea.Msg {
					return ProjectSetupCompleteMsg{
						ProjectName: m.nameInput.Value(),
						Description: description,
					}
				}
			}
		case "tab":
			if m.focusedField == 0 {
				m.nameInput.Blur()
				m.descInput.Focus()
				m.focusedField = 1
			} else {
				m.descInput.Blur()
				m.nameInput.Focus()
				m.focusedField = 0
			}
			m.validationErr = nil
		case "shift+tab":
			if m.focusedField == 1 {
				m.descInput.Blur()
				m.nameInput.Focus()
				m.focusedField = 0
			} else {
				m.nameInput.Blur()
				m.descInput.Focus()
				m.focusedField = 1
			}
			m.validationErr = nil
		}
	}

	// Clear validation error when user types
	if msg, ok := msg.(tea.KeyMsg); ok {
		if len(msg.String()) == 1 || msg.String() == "backspace" {
			m.validationErr = nil
		}
	}

	// Update the focused input
	if m.focusedField == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.descInput, cmd = m.descInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ProjectSetupModel) View() string {
	subtitle := components.RenderSubtitle("Let's set up your project!")

	// Project name field
	nameLabel := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("📝 Project name: (required)")

	// Style the name input based on focus
	nameInputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderPrimary).
		Padding(0, 1).
		Margin(1, 0, 0, 0)

	if m.focusedField == 0 {
		nameInputStyle = nameInputStyle.BorderForeground(styles.ColorBorderAccent)
	}

	nameBox := nameInputStyle.Render(m.nameInput.View())

	// Description field
	descLabel := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("📄 Description: (optional)")

	// Style the description input based on focus
	descInputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderPrimary).
		Padding(0, 1).
		Margin(1, 0, 0, 0)

	if m.focusedField == 1 {
		descInputStyle = descInputStyle.BorderForeground(styles.ColorBorderAccent)
	}

	descBox := descInputStyle.Render(m.descInput.View())

	// Error message if validation fails
	var errorMsg string
	if m.validationErr != nil {
		errorMsg = lipgloss.NewStyle().
			Foreground(styles.ColorError).
			Margin(1, 0, 0, 0).
			Render("⚠️  " + m.validationErr.Error())
	}

	// Instructions
	var instructions string
	if m.focusedField == 0 {
		instructions = lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(1, 0, 0, 0).
			Render("Enter: next field • Tab: switch fields")
	} else {
		instructions = lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(1, 0, 0, 0).
			Render("Enter: continue • Tab: switch fields")
	}

	// Character count helpers
	nameCount := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Margin(0, 0, 0, 1).
		Render("(" + fmt.Sprintf("%d/50", len(m.nameInput.Value())) + ")")

	descCount := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Margin(0, 0, 0, 1).
		Render("(" + fmt.Sprintf("%d/200", len(m.descInput.Value())) + ")")

	content := subtitle + "\n\n" + 
		      nameLabel + nameCount + "\n" + nameBox + "\n\n" + 
		      descLabel + descCount + "\n" + descBox + "\n\n" + 
		      instructions

	if errorMsg != "" {
		content += "\n" + errorMsg
	}

	return content
}

type ProjectSetupCompleteMsg struct {
	ProjectName string
	Description string
}