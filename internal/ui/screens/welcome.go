package screens

import (
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeModel struct {
	ready bool
}

func NewWelcomeModel() WelcomeModel {
	return WelcomeModel{
		ready: false,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return m, func() tea.Msg {
				return WelcomeCompleteMsg{}
			}
		}
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	// Welcome header
	welcomeText := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Align(lipgloss.Center).
		Width(40).
		Render("Welcome to Teapot! 🫖")

	// Description
	description := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Align(lipgloss.Center).
		Width(40).
		Render("The modern monorepo builder for full-stack developers")

	// Features list
	features := []string{
		"🚀 Modern tech stack (React, Next.js, TanStack, Expo)",
		"🛠️ Development tools (ESLint, Prettier, Biome)",
		"🐳 Infrastructure setup (Docker, Pulumi, Terraform)",
		"⚡ CI/CD pipelines (GitHub Actions, GitLab CI)",
		"🤖 AI tooling integration (Claude, Cursor, Windsurf)",
	}

	var featureList string
	for _, feature := range features {
		featureList += lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Margin(0, 0, 0, 2).
			Render(feature) + "\n"
	}

	// Call to action
	ctaText := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Bold(true).
		Align(lipgloss.Center).
		Width(40).
		Margin(1, 0, 0, 0).
		Render("Ready to build something amazing?")

	return welcomeText + "\n\n" + 
		   description + "\n\n" + 
		   featureList + "\n" + 
		   ctaText
}

type WelcomeCompleteMsg struct{}