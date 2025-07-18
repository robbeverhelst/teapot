package screens

import (
	"fmt"
	"strings"
	"time"

	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GeneratingModel struct {
	progress    int
	currentStep string
	steps       []GenerationStep
	completed   []bool
	done        bool
}

type GenerationStep struct {
	Name        string
	Description string
}

func NewGeneratingModel() GeneratingModel {
	steps := []GenerationStep{
		{"Base project initialized", "Creating project structure"},
		{"Monorepo workspace configured", "Setting up workspace"},
		{"Apps scaffolded", "Generating applications"},
		{"Packages created", "Setting up shared packages"},
		{"Installing dependencies", "Running package manager"},
		{"Setting up Git hooks", "Configuring development tools"},
		{"Configuring CI/CD", "Setting up automation"},
	}

	return GeneratingModel{
		progress:    0,
		currentStep: "Creating project structure...",
		steps:       steps,
		completed:   make([]bool, len(steps)),
		done:        false,
	}
}

func (m GeneratingModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return ProgressMsg{}
	})
}

func (m GeneratingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "enter":
			if m.done {
				return m, func() tea.Msg {
					return GenerationCompleteMsg{}
				}
			}
		}

	case ProgressMsg:
		if !m.done {
			if m.progress < 100 {
				m.progress += 2
				stepIndex := (m.progress * len(m.steps)) / 100
				if stepIndex < len(m.steps) {
					m.currentStep = m.steps[stepIndex].Description
					if stepIndex > 0 {
						m.completed[stepIndex-1] = true
					}
				}
			}

			if m.progress >= 100 {
				m.done = true
				m.completed[len(m.completed)-1] = true
				m.currentStep = "Project generation complete!"
				return m, func() tea.Msg {
					return GenerationCompleteMsg{}
				}
			}

			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return ProgressMsg{}
			})
		}
	}

	return m, nil
}

func (m GeneratingModel) View() string {
	subtitle := components.RenderSubtitle("🫖 Generating your project...")

	progressBar := m.renderProgressBar()
	statusText := fmt.Sprintf("%s %d%%", m.currentStep, m.progress)

	var stepsList strings.Builder
	for i, step := range m.steps {
		icon := "○"
		style := styles.UnselectedStyle
		
		if m.completed[i] {
			icon = "✓"
			style = styles.CheckedStyle
		} else if i == (m.progress*len(m.steps))/100 {
			icon = "⟳"
			style = styles.SelectedStyle
		}

		stepsList.WriteString(style.Render(icon + " " + step.Name) + "\n")
	}

	content := subtitle + "\n\n" + progressBar + "\n" + statusText + "\n\n" + stepsList.String()

	if m.done {
		content += "\n" + components.RenderHelp("enter: continue")
	} else {
		content += "\n" + components.RenderHelp("esc: cancel")
	}

	return content
}

func (m GeneratingModel) renderProgressBar() string {
	width := 31
	filled := (m.progress * width) / 100
	
	var bar strings.Builder
	bar.WriteString("┌")
	for i := 0; i < width; i++ {
		if i < filled {
			bar.WriteString("█")
		} else {
			bar.WriteString("░")
		}
	}
	bar.WriteString("┐")

	return lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Render(bar.String())
}

type ProgressMsg struct{}

type GenerationCompleteMsg struct{}