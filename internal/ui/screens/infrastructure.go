package screens

import (
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InfrastructureModel struct {
	options []InfraOption
	cursor  int
}

type InfraOption struct {
	Key         string
	Name        string
	Description string
	Selected    bool
	IsContinue  bool // Special flag for continue option
}

func NewInfrastructureModel() InfrastructureModel {
	return InfrastructureModel{
		options: []InfraOption{
			{"docker", "Docker", "Containerize your applications", false, false},
			{"docker-compose", "Docker Compose", "Multi-container development setup", false, false},
			{"pulumi", "Pulumi", "Infrastructure as Code for Kubernetes", false, false},
			{"terraform", "Terraform", "Infrastructure as Code for cloud resources", false, false},
			{"continue", "Continue", "Proceed with selected infrastructure", false, true},
		},
		cursor: 0,
	}
}

func (m InfrastructureModel) Init() tea.Cmd {
	return nil
}

func (m InfrastructureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case " ":
			// Space only works on non-continue options
			if !m.options[m.cursor].IsContinue {
				m.options[m.cursor].Selected = !m.options[m.cursor].Selected
			}
		case "enter":
			// Enter either selects current item or continues if on continue option
			if m.options[m.cursor].IsContinue {
				// Continue to next screen
				selectedOptions := make(map[string]bool)
				for _, option := range m.options {
					if !option.IsContinue {
						selectedOptions[option.Key] = option.Selected
					}
				}
				
				return m, func() tea.Msg {
					return InfrastructureSelectedMsg{
						Options: selectedOptions,
					}
				}
			} else {
				// Select current item
				m.options[m.cursor].Selected = !m.options[m.cursor].Selected
				return m, nil
			}
		case "s":
			// Skip infrastructure setup
			return m, func() tea.Msg {
				return InfrastructureSelectedMsg{
					Options: make(map[string]bool), // Empty options = skip
				}
			}
		}
	}
	return m, nil
}

func (m InfrastructureModel) View() string {
	subtitle := components.RenderSubtitle("Infrastructure & Deployment")

	var choices string
	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		var checked string
		var optionStyle lipgloss.Style
		var extraInfo string
		
		if option.IsContinue {
			// Continue option styling
			checked = "→"
			extraInfo = ""
		} else {
			// Regular option styling
			checked = "☐"
			if option.Selected {
				checked = "☑"
			}
			
			switch option.Key {
			case "docker":
				extraInfo = "Single container setup with Dockerfile"
			case "docker-compose":
				extraInfo = "Multi-service development environment"
			case "pulumi":
				extraInfo = "Deploy to Kubernetes with type-safe code"
			case "terraform":
				extraInfo = "Deploy to AWS, GCP, Azure with HCL"
			}
		}

		if m.cursor == i {
			optionStyle = styles.FocusedStyle
		} else if option.Selected {
			optionStyle = styles.CheckedStyle
		} else {
			optionStyle = styles.UnselectedStyle
		}

		// Render the entire line together for proper alignment
		choiceText := cursor + " " + checked + " " + option.Name
		choice := optionStyle.Render(choiceText)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(option.Description)

		extra := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Italic(true).
			Render(extraInfo)

		choices += choice + "\n" + description + "\n" + extra + "\n"
	}

	// Compact cloud providers note
	cloudNote := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Margin(1, 0, 0, 0).
		Render("☁️ Auto-detects: AWS, Vercel, Railway")

	// Skip option
	skipNote := lipgloss.NewStyle().
		Foreground(styles.ColorWarning).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("Press 's' to skip infrastructure setup")

	return subtitle + "\n\n" + choices + cloudNote + "\n" + skipNote
}

type InfrastructureSelectedMsg struct {
	Options map[string]bool
}