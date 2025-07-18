package screens

import (
	"strings"
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AIToolsModel struct {
	editors []EditorOption
	cursor  int
}

type EditorOption struct {
	Key         string
	Name        string
	Description string
	Extensions  []string
	Selected    bool
	IsContinue  bool
}

func NewAIToolsModel() AIToolsModel {
	return AIToolsModel{
		editors: []EditorOption{
			{
				Key:         "claude-code",
				Name:        "Claude Code",
				Description: "Anthropic's official CLI for Claude",
				Extensions:  []string{"claude-code", "VS Code integration"},
				Selected:    false,
				IsContinue:  false,
			},
			{
				Key:         "cursor",
				Name:        "Cursor",
				Description: "AI-powered code editor",
				Extensions:  []string{"Built-in AI", "VS Code fork"},
				Selected:    false,
				IsContinue:  false,
			},
			{
				Key:         "windsurf",
				Name:        "Windsurf",
				Description: "Codeium's AI-native IDE",
				Extensions:  []string{"AI chat", "Code generation"},
				Selected:    false,
				IsContinue:  false,
			},
			{
				Key:         "continue-dev",
				Name:        "Continue.dev",
				Description: "Open-source AI coding assistant",
				Extensions:  []string{"VS Code extension", "JetBrains plugin"},
				Selected:    false,
				IsContinue:  false,
			},
			{
				Key:         "none",
				Name:        "Skip AI Tools",
				Description: "Set up AI tools later",
				Extensions:  []string{},
				Selected:    false,
				IsContinue:  false,
			},
			{
				Key:         "continue",
				Name:        "Continue",
				Description: "Proceed with selected AI tools",
				Extensions:  []string{},
				Selected:    false,
				IsContinue:  true,
			},
		},
		cursor: 0,
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
		case " ":
			// Space key toggles selection
			if !m.editors[m.cursor].IsContinue {
				m.editors[m.cursor].Selected = !m.editors[m.cursor].Selected
			}
			return m, nil
		case "enter":
			if m.editors[m.cursor].IsContinue {
				// Continue with selected tools
				selectedTools := []string{}
				selectedExtensions := []string{}
				for _, editor := range m.editors {
					if !editor.IsContinue && editor.Selected {
						selectedTools = append(selectedTools, editor.Key)
						selectedExtensions = append(selectedExtensions, editor.Extensions...)
					}
				}
				
				return m, func() tea.Msg {
					return AIToolsSelectedMsg{
						Editor:     strings.Join(selectedTools, ","),
						Extensions: selectedExtensions,
					}
				}
			} else {
				// Toggle selection with Enter too
				m.editors[m.cursor].Selected = !m.editors[m.cursor].Selected
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
		
		if editor.IsContinue {
			// Continue option styling
			checked = "→"
			if m.cursor == i {
				editorStyle = styles.FocusedStyle
			} else {
				editorStyle = styles.UnselectedStyle
			}
		} else {
			// Regular editor styling (multiselect)
			if editor.Selected {
				checked = "☑" // Checked box
				editorStyle = styles.CheckedStyle
			} else {
				checked = "☐" // Unchecked box
				editorStyle = styles.UnselectedStyle
			}
			
			if m.cursor == i {
				editorStyle = styles.FocusedStyle
			}
		}

		// Unified choice rendering
		choiceText := cursor + " " + checked + " " + editor.Name
		choice := editorStyle.Render(choiceText)

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