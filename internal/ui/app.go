package ui

import (
	"teapot/internal/models"
	"teapot/internal/ui/components"
	"teapot/internal/ui/screens"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ScreenSizeAware interface for screens that need to respond to size changes
type ScreenSizeAware interface {
	SetSize(width, height int) (tea.Model, tea.Cmd)
}

type Model struct {
	state         models.AppState
	screenModels  map[models.Screen]tea.Model
	windowWidth   int
	windowHeight  int
}

func NewModel() Model {
	state := models.AppState{
		CurrentScreen: models.WelcomeScreen,
		Project:       models.ProjectConfig{},
		Quitting:      false,
	}

	screenModels := make(map[models.Screen]tea.Model)
	screenModels[models.WelcomeScreen] = screens.NewWelcomeModel()

	return Model{
		state:        state,
		screenModels: screenModels,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		
		// Update all screen models with new dimensions
		for screen, screenModel := range m.screenModels {
			if sizeAwareModel, ok := screenModel.(ScreenSizeAware); ok {
				updatedModel, cmd := sizeAwareModel.SetSize(msg.Width, msg.Height)
				m.screenModels[screen] = updatedModel
				if cmd != nil {
					return m, cmd
				}
			}
		}
		
		return m, nil

	case tea.KeyMsg:
		// Handle global quit keys first before passing to individual screens
		switch msg.String() {
		case "ctrl+c":
			// Allow Ctrl+C to quit from any screen
			m.state.Quitting = true
			return m, tea.Quit
		case "esc":
			// Allow Esc to quit from any screen except generation/complete screens
			if m.state.CurrentScreen != models.GeneratingScreen && m.state.CurrentScreen != models.CompleteScreen {
				m.state.Quitting = true
				return m, tea.Quit
			}
		case "backspace":
			// Handle back navigation
			return m.handleBackNavigation()
		}
		// If quit keys didn't trigger, pass the key to the current screen
		if screenModel, exists := m.screenModels[m.state.CurrentScreen]; exists {
			updatedModel, cmd := screenModel.Update(msg)
			m.screenModels[m.state.CurrentScreen] = updatedModel
			return m, cmd
		}

	case screens.WelcomeCompleteMsg:
		if m.state.CurrentScreen == models.WelcomeScreen {
			m.state.CurrentScreen = models.ProjectSetupScreen
			if _, exists := m.screenModels[models.ProjectSetupScreen]; !exists {
				m.screenModels[models.ProjectSetupScreen] = screens.NewProjectSetupModel()
			}
		}
		return m, nil

	case screens.ProjectSetupCompleteMsg:
		if m.state.CurrentScreen == models.ProjectSetupScreen {
			m.state.Project.Name = msg.ProjectName
			m.state.Project.Description = msg.Description
			m.state.CurrentScreen = models.ArchitectureScreen
			if _, exists := m.screenModels[models.ArchitectureScreen]; !exists {
				m.screenModels[models.ArchitectureScreen] = screens.NewArchitectureModel()
			}
		}
		return m, nil

	case screens.ArchitectureSelectedMsg:
		if m.state.CurrentScreen == models.ArchitectureScreen {
			m.state.Project.Architecture = msg.Architecture
			m.state.CurrentScreen = models.AddAppsScreen
			if _, exists := m.screenModels[models.AddAppsScreen]; !exists {
				m.screenModels[models.AddAppsScreen] = screens.NewAddAppsModel()
			}
		}
		return m, nil

	case screens.AppTypeSelectedMsg:
		if m.state.CurrentScreen == models.AddAppsScreen {
			// Create new application and start configuration
			newApp := models.Application{
				ID:      "app-" + string(msg.AppType),
				Name:    string(msg.AppType),
				Type:    msg.AppType,
				Options: make(map[string]interface{}),
			}
			m.state.CurrentApp = &newApp
			m.state.CurrentScreen = models.AppConfigScreen
			// Always create a new AppConfigModel for the new app type
			m.screenModels[models.AppConfigScreen] = screens.NewAppConfigModel(msg.AppType)
		}
		return m, nil

	case screens.AppConfigCompleteMsg:
		if m.state.CurrentScreen == models.AppConfigScreen && m.state.CurrentApp != nil {
			// Save the configured app
			m.state.CurrentApp.Name = msg.AppName
			m.state.CurrentApp.Options = msg.Options
			m.state.Project.Applications = append(m.state.Project.Applications, *m.state.CurrentApp)
			m.state.CurrentApp = nil
			
			// Move to "Add Another App?" screen
			m.state.CurrentScreen = models.AddAnotherAppScreen
			if _, exists := m.screenModels[models.AddAnotherAppScreen]; !exists {
				m.screenModels[models.AddAnotherAppScreen] = screens.NewAddAnotherAppModel(
					len(m.state.Project.Applications), 
					m.state.Project.Architecture,
				)
			}
		}
		return m, nil

	case screens.AddAnotherAppSelectedMsg:
		if m.state.CurrentScreen == models.AddAnotherAppScreen {
			if msg.Action == "add" {
				// Go back to app selection
				m.state.CurrentScreen = models.AddAppsScreen
				// Refresh the screen model for new app selection
				m.screenModels[models.AddAppsScreen] = screens.NewAddAppsModel()
			} else {
				// Continue to dev tools
				m.state.CurrentScreen = models.DevToolsScreen
				if _, exists := m.screenModels[models.DevToolsScreen]; !exists {
					m.screenModels[models.DevToolsScreen] = screens.NewDevToolsModel()
				}
			}
		}
		return m, nil

	case screens.DevToolsSelectedMsg:
		if m.state.CurrentScreen == models.DevToolsScreen {
			m.state.Project.DevTools.Linting = msg.LintingTool
			m.state.Project.DevTools.TypeScript = true
			m.state.Project.DevTools.Husky = true
			m.state.Project.DevTools.LintStaged = true
			
			m.state.CurrentScreen = models.InfrastructureScreen
			if _, exists := m.screenModels[models.InfrastructureScreen]; !exists {
				m.screenModels[models.InfrastructureScreen] = screens.NewInfrastructureModel()
			}
		}
		return m, nil

	case screens.InfrastructureSelectedMsg:
		if m.state.CurrentScreen == models.InfrastructureScreen {
			m.state.Project.Infrastructure.Docker = msg.Options["docker"]
			m.state.Project.Infrastructure.DockerCompose = msg.Options["docker-compose"]
			m.state.Project.Infrastructure.Pulumi = msg.Options["pulumi"]
			m.state.Project.Infrastructure.Terraform = msg.Options["terraform"]
			
			m.state.CurrentScreen = models.CIPipelineScreen
			if _, exists := m.screenModels[models.CIPipelineScreen]; !exists {
				m.screenModels[models.CIPipelineScreen] = screens.NewCIPipelineModel()
			}
		}
		return m, nil

	case screens.CIPipelineSelectedMsg:
		if m.state.CurrentScreen == models.CIPipelineScreen {
			m.state.Project.CIPipeline.Provider = msg.Provider
			m.state.Project.CIPipeline.Features = msg.Features
			
			m.state.CurrentScreen = models.AIToolsScreen
			if _, exists := m.screenModels[models.AIToolsScreen]; !exists {
				m.screenModels[models.AIToolsScreen] = screens.NewAIToolsModel()
			}
		}
		return m, nil

	case screens.AIToolsSelectedMsg:
		if m.state.CurrentScreen == models.AIToolsScreen {
			m.state.Project.AITools.Editor = msg.Editor
			m.state.Project.AITools.Extensions = msg.Extensions
			
			// Move to generation
			m.state.CurrentScreen = models.GeneratingScreen
			if _, exists := m.screenModels[models.GeneratingScreen]; !exists {
				m.screenModels[models.GeneratingScreen] = screens.NewGeneratingModel()
			}
		}
		return m, nil

	case screens.GenerationCompleteMsg:
		if m.state.CurrentScreen == models.GeneratingScreen {
			m.state.CurrentScreen = models.CompleteScreen
			if _, exists := m.screenModels[models.CompleteScreen]; !exists {
				m.screenModels[models.CompleteScreen] = screens.NewCompleteModel(m.state.Project.Name)
			}
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleBackNavigation() (Model, tea.Cmd) {
	switch m.state.CurrentScreen {
	case models.ProjectSetupScreen:
		m.state.CurrentScreen = models.WelcomeScreen
	case models.ArchitectureScreen:
		m.state.CurrentScreen = models.ProjectSetupScreen
	case models.AddAppsScreen:
		if len(m.state.Project.Applications) > 0 {
			// If we have apps, go back to "add another" screen
			m.state.CurrentScreen = models.AddAnotherAppScreen
			if _, exists := m.screenModels[models.AddAnotherAppScreen]; !exists {
				m.screenModels[models.AddAnotherAppScreen] = screens.NewAddAnotherAppModel(
					len(m.state.Project.Applications), 
					m.state.Project.Architecture,
				)
			}
		} else {
			// No apps yet, go back to architecture
			m.state.CurrentScreen = models.ArchitectureScreen
		}
	case models.AppConfigScreen:
		m.state.CurrentScreen = models.AddAppsScreen
		// Clear the current app being configured
		m.state.CurrentApp = nil
	case models.AddAnotherAppScreen:
		// Go back to the last app's configuration
		if len(m.state.Project.Applications) > 0 {
			// Remove the last app and go back to its config
			lastApp := m.state.Project.Applications[len(m.state.Project.Applications)-1]
			m.state.Project.Applications = m.state.Project.Applications[:len(m.state.Project.Applications)-1]
			m.state.CurrentApp = &lastApp
			m.state.CurrentScreen = models.AppConfigScreen
			m.screenModels[models.AppConfigScreen] = screens.NewAppConfigModel(lastApp.Type)
		} else {
			m.state.CurrentScreen = models.AddAppsScreen
		}
	case models.DevToolsScreen:
		m.state.CurrentScreen = models.AddAnotherAppScreen
		if _, exists := m.screenModels[models.AddAnotherAppScreen]; !exists {
			m.screenModels[models.AddAnotherAppScreen] = screens.NewAddAnotherAppModel(
				len(m.state.Project.Applications), 
				m.state.Project.Architecture,
			)
		}
	case models.InfrastructureScreen:
		m.state.CurrentScreen = models.DevToolsScreen
	case models.CIPipelineScreen:
		m.state.CurrentScreen = models.InfrastructureScreen
	case models.AIToolsScreen:
		m.state.CurrentScreen = models.CIPipelineScreen
	default:
		// For other screens, do nothing or default behavior
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.state.Quitting {
		return ""
	}

	if m.windowWidth == 0 || m.windowHeight == 0 {
		// Beautiful loading screen
		loading := lipgloss.NewStyle().
			Foreground(styles.ColorPrimary).
			Bold(true).
			Align(lipgloss.Center).
			Render("🫖 Loading Teapot...")
		
		return lipgloss.NewStyle().
			Background(styles.ColorBgPrimary).
			Align(lipgloss.Center).
			Render(loading)
	}

	leftPanel := m.renderLeftPanel()
	rightPanel := components.RenderProjectStructure(m.state.Project, m.windowWidth, m.windowHeight)

	// Create proper spacing between panels with subtle shadow effect
	leftPanelWithShadow := lipgloss.NewStyle().
		MarginRight(1).
		Render(leftPanel)
	
	rightPanelWithShadow := lipgloss.NewStyle().
		MarginLeft(1).
		Render(rightPanel)
	
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanelWithShadow,
		lipgloss.NewStyle().Width(2).Render(" "), // Spacer between panels
		rightPanelWithShadow,
	)

	// Create fullscreen background with centered content
	return lipgloss.NewStyle().
		Background(styles.ColorBgPrimary).
		Width(m.windowWidth).
		Height(m.windowHeight).
		Render(
			lipgloss.Place(
				m.windowWidth,
				m.windowHeight,
				lipgloss.Center,
				lipgloss.Center,
				content,
			),
		)
}

func (m Model) renderLeftPanel() string {
	var content string

	// Title section with compact spacing
	content += components.RenderTitle() + "\n"

	// Main content area with compact spacing
	if screenModel, exists := m.screenModels[m.state.CurrentScreen]; exists {
		content += screenModel.View()
	}

	// Progress section with compact space
	content += "\n" + components.RenderProgressIndicator(m.state.CurrentScreen)
	
	// Help text at the bottom with compact spacing
	content += "\n" + m.getHelpText()

	return styles.GetLeftPanelStyle(m.windowWidth, m.windowHeight).Render(content)
}

func (m Model) getHelpText() string {
	switch m.state.CurrentScreen {
	case models.WelcomeScreen:
		return components.RenderHelp("enter/space: get started • ctrl+c: quit")
	case models.ProjectSetupScreen:
		return components.RenderHelp("tab: switch fields • enter: continue • backspace: back • esc: quit")
	case models.ArchitectureScreen:
		return components.RenderHelp("↑↓: navigate • enter: select • backspace: back • esc: quit")
	case models.AddAppsScreen:
		return components.RenderHelp("↑↓: navigate • enter: select • backspace: back • esc: quit")
	case models.AppConfigScreen:
		return components.RenderHelp("↑↓: navigate • space: toggle • enter: continue • backspace: back • esc: quit")
	case models.AddAnotherAppScreen:
		return components.RenderHelp("↑↓: navigate • enter: select • backspace: back • esc: quit")
	case models.DevToolsScreen:
		return components.RenderHelp("↑↓: navigate • enter: select • backspace: back • esc: quit")
	case models.InfrastructureScreen:
		return components.RenderHelp("↑↓: navigate • space/enter: select • s: skip • backspace: back • esc: quit")
	case models.CIPipelineScreen:
		return components.RenderHelp("↑↓: navigate • space/enter: select • tab: switch areas • s: skip • backspace: back • esc: quit")
	case models.AIToolsScreen:
		return components.RenderHelp("↑↓: navigate • enter: select • s: skip • backspace: back • esc: quit")
	case models.GeneratingScreen:
		return ""
	case models.CompleteScreen:
		return ""
	default:
		return components.RenderHelp("j/k: navigate • space: select • enter: continue • backspace: back • esc: quit")
	}
}

