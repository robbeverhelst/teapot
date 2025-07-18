package screens

import (
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AIToolsModel struct {
	editors  []EditorOption
	cursor   int
	selected int
}

type EditorOption struct {
	Key         string
	Name        string
	Description string
	Extensions  []string
}

func NewAIToolsModel() AIToolsModel {
	return AIToolsModel{
		editors: []EditorOption{
			{
				"claude-code",
				"Claude Code",
				"Anthropic's official CLI for Claude",
				[]string{"claude-code", "VS Code integration"},
			},
			{
				"cursor",
				"Cursor",
				"AI-powered code editor",
				[]string{"Built-in AI", "VS Code fork"},
			},
			{
				"windsurf",
				"Windsurf",
				"Codeium's AI-native IDE",
				[]string{"AI chat", "Code generation"},
			},
			{
				"continue-dev",
				"Continue.dev",
				"Open-source AI coding assistant",
				[]string{"VS Code extension", "JetBrains plugin"},
			},
			{
				"none",
				"Skip AI Tools",
				"Set up AI tools later",
				[]string{},
			},
			{
				"continue",
				"Continue",
				"Proceed with selected AI tools",
				[]string{},
			},
		},
		cursor:   0,
		selected: -1,
	}
}

func (m AIToolsModel) Init() tea.Cmd {
	return nil
}

func (m AIToolsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.editors)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.editors[m.cursor].Key == "continue" {
				// Continue with selected editor
				if m.selected >= 0 {
					return m, func() tea.Msg {
						return AIToolsSelectedMsg{
							Editor:     m.editors[m.selected].Key,
							Extensions: m.editors[m.selected].Extensions,
						}
					}
				}
				return m, nil // No selection made
			} else {
				// Select current editor
				m.selected = m.cursor
				return m, nil
			}
		case "s":
			// Skip AI tools setup
			return m, func() tea.Msg {
				return AIToolsSelectedMsg{
					Editor:     "none",
					Extensions: []string{},
				}
			}
		}
	}
	return m, nil
}

func (m AIToolsModel) View() string {
	subtitle := components.RenderSubtitle("AI Development Environment")

	var choices string
	for i, editor := range m.editors {
		cursor := " "
		
		var checked string
		var editorStyle lipgloss.Style
		
		if editor.Key == "continue" {
			// Continue option styling
			checked = "→"
			if m.cursor == i {
				editorStyle = styles.FocusedStyle
			} else {
				editorStyle = styles.UnselectedStyle
			}
		} else {
			// Regular editor styling
			if m.selected == i {
				checked = "●"
				editorStyle = styles.CheckedStyle
			} else {
				checked = " "
				editorStyle = styles.UnselectedStyle
			}
			
			if m.cursor == i {
				editorStyle = styles.FocusedStyle
			}
		}

		choice := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render(cursor) + " " +
			editorStyle.Render(checked+" "+editor.Name)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(editor.Description)

		// Extensions/features (only for non-continue options)
		var extensionsList string
		if editor.Key != "continue" && len(editor.Extensions) > 0 {
			for _, ext := range editor.Extensions {
				extensionsList += lipgloss.NewStyle().
					Foreground(styles.ColorAccent).
					Margin(0, 0, 0, 6).
					Render("• " + ext) + "\n"
			}
		}

		choices += choice + "\n" + description
		if extensionsList != "" {
			choices += "\n" + extensionsList
		}
		choices += "\n"
	}

	// Compact benefits and notes
	benefitsNote := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Margin(1, 0, 0, 0).
		Render("🤖 Includes: AI coding, debugging, and documentation tools")

	// Skip option
	skipNote := lipgloss.NewStyle().
		Foreground(styles.ColorWarning).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("Press 's' to skip AI tools setup")

	return subtitle + "\n\n" + choices + benefitsNote + "\n" + skipNote
}

type AIToolsSelectedMsg struct {
	Editor     string
	Extensions []string
}