package screens

import (
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type CompleteModel struct {
	projectName string
}

func NewCompleteModel(projectName string) CompleteModel {
	return CompleteModel{
		projectName: projectName,
	}
}

func (m CompleteModel) Init() tea.Cmd {
	return nil
}

func (m CompleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m CompleteModel) View() string {
	successTitle := styles.CheckedStyle.Render("✨ Project created successfully!")
	
	pathText := "Your project is ready at:\n./" + m.projectName

	nextSteps := "Next steps:\n" +
		"1. cd " + m.projectName + "\n" +
		"2. pnpm install\n" +
		"3. pnpm dev"

	commands := "Available commands:\n" +
		"• pnpm dev      - Start all apps\n" +
		"• pnpm build    - Build all apps\n" +
		"• pnpm lint     - Lint all packages\n" +
		"• pnpm test     - Run tests"

	footer := styles.CheckedStyle.Render("Happy coding! 🚀")

	content := successTitle + "\n\n" + pathText + "\n\n" + nextSteps + "\n\n" + commands + "\n\n" + footer

	content += "\n\n" + components.RenderHelp("enter: exit")

	return content
}