package screens

import (
	"teapot/internal/ui/components"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CIPipelineModel struct {
	providers     []ProviderOption
	features      []FeatureOption
	providerIdx   int
	featureIdx    int
	currentArea   int // 0 = providers, 1 = features
	selectedProvider int
}

type ProviderOption struct {
	Key         string
	Name        string
	Description string
}

type FeatureOption struct {
	Key         string
	Name        string
	Description string
	Selected    bool
	IsContinue  bool // Special flag for continue option
}

func NewCIPipelineModel() CIPipelineModel {
	return CIPipelineModel{
		providers: []ProviderOption{
			{"github", "GitHub Actions", "Integrated with GitHub repositories"},
			{"gitlab", "GitLab CI", "GitLab's built-in CI/CD system"},
			{"jenkins", "Jenkins", "Self-hosted automation server"},
			{"skip", "Skip CI/CD", "Set up CI/CD later manually"},
		},
		features: []FeatureOption{
			{"testing", "Testing", "Run tests on every push", true, false},
			{"linting", "Linting & Formatting", "Code quality checks", true, false},
			{"docker", "Docker Image Build", "Build and push container images", false, false},
			{"deployment", "Automatic Deployment", "Deploy on successful builds", false, false},
			{"security", "Security Scanning", "Vulnerability and dependency checks", false, false},
			{"continue", "Continue", "Proceed with selected configuration", false, true},
		},
		providerIdx:      0,
		featureIdx:       0,
		currentArea:      0,
		selectedProvider: 0, // Default to GitHub Actions
	}
}

func (m CIPipelineModel) Init() tea.Cmd {
	return nil
}

func (m CIPipelineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.currentArea = (m.currentArea + 1) % 2
		case "shift+tab":
			m.currentArea = (m.currentArea - 1 + 2) % 2
		case "j", "down":
			if m.currentArea == 0 && m.providerIdx < len(m.providers)-1 {
				m.providerIdx++
			} else if m.currentArea == 1 && m.featureIdx < len(m.features)-1 {
				m.featureIdx++
			}
		case "k", "up":
			if m.currentArea == 0 && m.providerIdx > 0 {
				m.providerIdx--
			} else if m.currentArea == 1 && m.featureIdx > 0 {
				m.featureIdx--
			}
		case " ":
			if m.currentArea == 0 {
				m.selectedProvider = m.providerIdx
			} else if m.currentArea == 1 {
				// Space only works on non-continue features
				if !m.features[m.featureIdx].IsContinue {
					m.features[m.featureIdx].Selected = !m.features[m.featureIdx].Selected
				}
			}
		case "enter":
			if m.currentArea == 0 {
				// In provider area, select current provider and move to features
				m.selectedProvider = m.providerIdx
				m.currentArea = 1 // Move focus to features section
				return m, nil
			} else {
				// In features area, either select feature or continue
				if m.features[m.featureIdx].IsContinue {
					// Continue to next screen
					selectedFeatures := []string{}
					for _, feature := range m.features {
						if !feature.IsContinue && feature.Selected {
							selectedFeatures = append(selectedFeatures, feature.Key)
						}
					}
					
					return m, func() tea.Msg {
						return CIPipelineSelectedMsg{
							Provider: m.providers[m.selectedProvider].Key,
							Features: selectedFeatures,
						}
					}
				} else {
					// Select current feature
					m.features[m.featureIdx].Selected = !m.features[m.featureIdx].Selected
					return m, nil
				}
			}
		case "s":
			// Skip CI/CD setup
			return m, func() tea.Msg {
				return CIPipelineSelectedMsg{
					Provider: "skip",
					Features: []string{},
				}
			}
		}
	}
	return m, nil
}

func (m CIPipelineModel) View() string {
	subtitle := components.RenderSubtitle("CI/CD Pipeline Setup")

	// Provider selection
	providerLabel := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Margin(0, 0, 1, 0).
		Render("🚀 Choose CI/CD Provider:")

	var providerChoices string
	for i, provider := range m.providers {
		cursor := " "
		if m.currentArea == 0 && m.providerIdx == i {
			cursor = ">"
		}

		checked := " "
		if m.selectedProvider == i {
			checked = "●"
		}

		var providerStyle lipgloss.Style
		if m.currentArea == 0 && m.providerIdx == i {
			providerStyle = styles.FocusedStyle
		} else {
			providerStyle = styles.UnselectedStyle
		}

		// Render the entire line together for proper alignment
		choiceText := cursor + " " + checked + " " + provider.Name
		choice := providerStyle.Render(choiceText)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(provider.Description)

		providerChoices += choice + "\n" + description + "\n"
	}

	// Features selection
	featuresLabel := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Margin(2, 0, 1, 0).
		Render("⚙️ Pipeline Features:")

	var featureChoices string
	for i, feature := range m.features {
		cursor := " "
		if m.currentArea == 1 && m.featureIdx == i {
			cursor = ">"
		}

		var checked string
		var featureStyle lipgloss.Style
		
		if feature.IsContinue {
			// Continue option styling
			checked = "→"
		} else {
			// Regular feature styling
			checked = "☐"
			if feature.Selected {
				checked = "☑"
			}
		}

		if m.currentArea == 1 && m.featureIdx == i {
			featureStyle = styles.FocusedStyle
		} else if feature.Selected {
			featureStyle = styles.CheckedStyle
		} else {
			featureStyle = styles.UnselectedStyle
		}

		// Render the entire line together for proper alignment
		choiceText := cursor + " " + checked + " " + feature.Name
		choice := featureStyle.Render(choiceText)

		description := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Margin(0, 0, 0, 4).
			Render(feature.Description)

		featureChoices += choice + "\n" + description + "\n"
	}

	// Skip option
	skipNote := lipgloss.NewStyle().
		Foreground(styles.ColorWarning).
		Bold(true).
		Margin(1, 0, 0, 0).
		Render("Press 's' to skip CI/CD setup")

	return subtitle + "\n\n" + 
		   providerLabel + "\n" + providerChoices + "\n" +
		   featuresLabel + "\n" + featureChoices + "\n" +
		   skipNote
}

type CIPipelineSelectedMsg struct {
	Provider string
	Features []string
}