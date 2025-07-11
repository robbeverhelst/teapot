package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateWelcome state = iota
	stateBrewing
	stateReady
)

type model struct {
	state    state
	spinner  spinner.Model
	ticks    int
	quitting bool
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginTop(1).
			MarginBottom(1)

	teapotStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(1)

	hotTeaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		state:   stateWelcome,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter", " ":
			if m.state == stateWelcome {
				m.state = stateBrewing
				m.ticks = 0
				return m, tea.Batch(
					m.spinner.Tick,
					tickCmd(),
				)
			}
		}

	case tickMsg:
		m.ticks++
		if m.ticks >= 30 { // ~3 seconds
			m.state = stateReady
			return m, nil
		}
		return m, tickCmd()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	teapot := `
    ___
   |___|
  /     \
 | ~~~ ~~|
  \     /
   '---'
`

	var s strings.Builder

	switch m.state {
	case stateWelcome:
		s.WriteString(titleStyle.Render("☕ Welcome to Teapot!"))
		s.WriteString("\n")
		s.WriteString(teapotStyle.Render(teapot))
		s.WriteString("\n")
		s.WriteString(textStyle.Render("Press ENTER to brew some tea..."))
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press q to quit"))

	case stateBrewing:
		s.WriteString(titleStyle.Render("☕ Brewing your tea..."))
		s.WriteString("\n")
		s.WriteString(teapotStyle.Render(teapot))
		s.WriteString("\n")
		s.WriteString(m.spinner.View() + " Please wait while we prepare your perfect cup")
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press q to quit"))

	case stateReady:
		s.WriteString(titleStyle.Render("☕ Your tea is ready!"))
		s.WriteString("\n")
		s.WriteString(teapotStyle.Render(teapot))
		s.WriteString("\n")
		s.WriteString(hotTeaStyle.Render("✨ Enjoy your delicious tea! ✨"))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press q to quit"))
	}

	return lipgloss.NewStyle().Margin(1, 2).Render(s.String())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}