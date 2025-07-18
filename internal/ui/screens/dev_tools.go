package screens

import (
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DevToolsModel struct {
	options  []DevToolOption
	cursor   int
	selected int
}

type DevToolOption struct {
	Key         string
	Name        string
	Description string
}

func NewDevToolsModel() DevToolsModel {
	return DevToolsModel{
		options: []DevToolOption{
			{"prettier-eslint", "Prettier + ESLint", "Traditional formatting and linting setup"},
			{"biome", "Biome", "Fast, modern toolchain for web projects"},
			{"custom", "Custom Setup", "Configure your own linting and formatting"},
			{"continue", "Continue", "Proceed with selected dev tools"},
		},
		cursor:   0,
		selected: -1,
	}
}

func (m DevToolsModel) Init() tea.Cmd {
	return nil
}

func (m DevToolsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				// Continue with selected tool (default to first option if none selected)
				selectedTool := "prettier-eslint" // default
				if m.selected >= 0 {
					selectedTool = m.options[m.selected].Key
				}
				return m, func() tea.Msg {
					return DevToolsSelectedMsg{
						LintingTool: selectedTool,
					}
				}
			} else {
				// Select current tool
				m.selected = m.cursor
				return m, nil
			}
		}
	}
	return m, nil
}

func (m DevToolsModel) View() string {
	subtitle := components.RenderSubtitle("Code Quality & Formatting")

	var choices string
	for i, option := range m.options {
		cursor := " "
		
		var checked string
		var optionStyle lipgloss.Style
		var extraInfo string
		
		if option.Key == "continue" {
			// Continue option styling
			checked = "→"
			extraInfo = ""
			if m.cursor == i {
				optionStyle = styles.FocusedStyle
			} else {
				optionStyle = styles.UnselectedStyle
			}
		} else {
			// Regular option styling
			if m.selected == i {
				checked = "●"
				optionStyle = styles.CheckedStyle
			} else {
				checked = " "
				optionStyle = styles.UnselectedStyle
			}
			
			if m.cursor == i {
				optionStyle = styles.FocusedStyle
			}
			
			switch option.Key {
			case "prettier-eslint":
				extraInfo = "Includes: Prettier, ESLint, Husky, lint-staged"
			case "biome":
				extraInfo = "Includes: Formatting, linting, import sorting (faster)"
			case "custom":
				extraInfo = "You'll configure your own setup"
			}
		}

		// Unified choice rendering
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

		choices += choice + "\n" + description + "\n" + extra + "\n\n"
	}

	// Additional tools section
	additionalTools := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("📦 Additional Tools (will be included):")

	tools := []string{
		"✓ TypeScript configuration",
		"✓ Git hooks setup (Husky)",
		"✓ Pre-commit linting (lint-staged)",
		"✓ Editor configuration (.editorconfig)",
	}

	var toolsList string
	for _, tool := range tools {
		toolsList += lipgloss.NewStyle().
			Foreground(styles.ColorSuccess).
			Margin(0, 0, 0, 2).
			Render(tool) + "\n"
	}

	return subtitle + "\n\n" + choices + additionalTools + "\n" + toolsList
}

type DevToolsSelectedMsg struct {
	LintingTool string
}