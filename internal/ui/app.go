// Package ui provides the main user interface components for the Teapot CLI application.
// It implements the Bubble Tea Model-View-Update pattern for managing the terminal UI,
// including screen navigation, window resizing, and state management.
package ui

import (
	"time"
	
	"teapot/internal/errors"
	"teapot/internal/generator"
	"teapot/internal/models"
	"teapot/internal/navigation"
	"teapot/internal/ui/components"
	"teapot/internal/ui/screens"
	"teapot/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ScreenSizeAware interface for screens that need to respond to terminal size changes.
// Screens implementing this interface will receive size updates when the terminal is resized.
type ScreenSizeAware interface {
	// SetSize updates the screen model with new terminal dimensions
	SetSize(width, height int) (tea.Model, tea.Cmd)
}

// Model represents the main application model that manages the overall UI state.
// It coordinates between different screens and handles navigation, window resizing,
// and global key bindings.
type Model struct {
	// state holds the current application state including screen and project config
	state         models.AppState
	// screenModels contains the instantiated screen models for each screen type
	screenModels  map[models.Screen]tea.Model
	// windowWidth stores the current terminal width
	windowWidth   int
	// windowHeight stores the current terminal height
	windowHeight  int
	// navigationFlow manages the navigation state machine
	navigationFlow *navigation.NavigationFlow
	// errorRecovery manages error handling and recovery
	errorRecovery *errors.ErrorRecovery
	// errorDisplay manages error display in the UI
	errorDisplay  *components.ErrorDisplay
}

// NewModel creates a new main application model with initial state.
// It initializes the application with the welcome screen and empty project configuration.
func NewModel() Model {
	state := models.AppState{
		CurrentScreen: models.WelcomeScreen,
		Project:       models.ProjectConfig{},
		Quitting:      false,
	}

	screenModels := make(map[models.Screen]tea.Model)
	screenModels[models.WelcomeScreen] = screens.NewWelcomeModel()

	return Model{
		state:          state,
		screenModels:   screenModels,
		navigationFlow: navigation.NewNavigationFlow(),
		errorRecovery:  errors.NewErrorRecovery(50, true), // 50 errors max, panic recovery enabled
		errorDisplay:   components.NewErrorDisplay(components.ErrorDisplayInline, true, 10*time.Second),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// addScreenModel adds a new screen model and applies current window dimensions if available
func (m *Model) addScreenModel(screen models.Screen, model tea.Model) {
	// Update size for new screen if window dimensions are available
	if m.windowWidth > 0 && m.windowHeight > 0 {
		if sizeAwareModel, ok := model.(ScreenSizeAware); ok {
			model, _ = sizeAwareModel.SetSize(m.windowWidth, m.windowHeight)
		}
	}
	m.screenModels[screen] = model
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Wrap the update function with error recovery
	defer func() {
		if r := recover(); r != nil {
			panicErr := errors.NewTeapotError(
				errors.ErrorTypePanic,
				"Application panic during update",
				nil,
			)
			m.errorDisplay.ShowError(panicErr)
			cmd := m.errorRecovery.HandleError(panicErr)
			if cmd != nil {
				// Execute the recovery command
				_ = cmd()
			}
		}
	}()
	
	switch msg := msg.(type) {
	case errors.ErrorOccurredMsg:
		// Handle non-recoverable errors
		m.errorDisplay.ShowError(msg.Error)
		if msg.Error.Type == errors.ErrorTypeSystem {
			// System errors are fatal
			m.state.Quitting = true
			return m, tea.Quit
		}
		return m, nil
		
	case errors.ErrorRecoveredMsg:
		// Handle recoverable errors
		m.errorDisplay.ShowError(msg.Error)
		// Apply recovery action if needed
		return m.handleErrorRecovery(msg.Error)
		
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		
		// Only update current screen for performance - other screens will be updated when navigated to
		if screenModel, exists := m.screenModels[m.state.CurrentScreen]; exists {
			if sizeAwareModel, ok := screenModel.(ScreenSizeAware); ok {
				updatedModel, cmd := sizeAwareModel.SetSize(msg.Width, msg.Height)
				m.screenModels[m.state.CurrentScreen] = updatedModel
				return m, cmd
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
			// Handle ESC key for dismissing errors first
			if m.errorDisplay.HasError() {
				m.errorDisplay.ClearError()
				return m, nil
			}
			// Allow Esc to quit from any screen except generation/complete screens
			if m.state.CurrentScreen != models.GeneratingScreen && m.state.CurrentScreen != models.CompleteScreen {
				m.state.Quitting = true
				return m, tea.Quit
			}
		case "backspace":
			// Check if we're in a screen with text input that might need backspace
			if m.state.CurrentScreen == models.ProjectSetupScreen || m.state.CurrentScreen == models.AppConfigScreen {
				// Let the screen handle backspace first
				if screenModel, exists := m.screenModels[m.state.CurrentScreen]; exists {
					updatedModel, cmd := screenModel.Update(msg)
					m.screenModels[m.state.CurrentScreen] = updatedModel
					// If the screen returns a command, it handled the backspace
					if cmd != nil {
						return m, cmd
					}
				}
			}
			// If not handled by screen, do back navigation
			return m.navigateBack()
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
				m.addScreenModel(models.ProjectSetupScreen, screens.NewProjectSetupModel())
			}
		}
		return m, nil

	case screens.ProjectSetupCompleteMsg:
		if m.state.CurrentScreen == models.ProjectSetupScreen {
			m.state.Project.Name = msg.ProjectName
			m.state.Project.Description = msg.Description
			m.state.CurrentScreen = models.ArchitectureScreen
			if _, exists := m.screenModels[models.ArchitectureScreen]; !exists {
				m.addScreenModel(models.ArchitectureScreen, screens.NewArchitectureModel())
			}
		}
		return m, nil

	case screens.ArchitectureSelectedMsg:
		if m.state.CurrentScreen == models.ArchitectureScreen {
			m.state.Project.Architecture = msg.Architecture
			m.state.CurrentScreen = models.AddAppsScreen
			if _, exists := m.screenModels[models.AddAppsScreen]; !exists {
				m.addScreenModel(models.AddAppsScreen, screens.NewAddAppsModel())
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
			
			// Move to YAML preview
			m.state.CurrentScreen = models.YAMLPreviewScreen
			if _, exists := m.screenModels[models.YAMLPreviewScreen]; !exists {
				m.screenModels[models.YAMLPreviewScreen] = screens.NewYAMLPreviewModel(m.state.Project)
			}
		}
		return m, nil

	case screens.YAMLSaveMsg:
		if m.state.CurrentScreen == models.YAMLPreviewScreen {
			// Save teapot.yml to current directory
			err := generator.SaveTeapotYAML(msg.Project, ".")
			if err != nil {
				// Handle error (could show error message)
				return m, nil
			}
			// Could show success message
		}
		return m, nil

	case screens.YAMLContinueMsg:
		if m.state.CurrentScreen == models.YAMLPreviewScreen {
			// Continue to generation
			m.state.CurrentScreen = models.GeneratingScreen
			if _, exists := m.screenModels[models.GeneratingScreen]; !exists {
				m.screenModels[models.GeneratingScreen] = screens.NewGeneratingModel()
			}
		}
		return m, nil

	case screens.YAMLBackMsg:
		if m.state.CurrentScreen == models.YAMLPreviewScreen {
			// Go back to AI tools screen
			m.state.CurrentScreen = models.AIToolsScreen
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

// navigateBack handles back navigation using the navigation state machine
func (m Model) navigateBack() (Model, tea.Cmd) {
	// Use the navigation flow to handle back navigation
	newState, newScreen, screenModel := m.navigationFlow.NavigateBack(m.state.CurrentScreen, &m.state)
	
	// Update the state
	m.state = *newState
	m.state.CurrentScreen = newScreen
	
	// If a specific screen model was returned, use it
	if screenModel != nil {
		if teaModel, ok := screenModel.(tea.Model); ok {
			m.addScreenModel(newScreen, teaModel)
		}
	} else {
		// Ensure the target screen model exists
		if _, exists := m.screenModels[newScreen]; !exists {
			// Create the screen model using the navigation flow factory
			factory := m.navigationFlow.GetScreenFactory(newScreen)
			if newScreen == models.AddAnotherAppScreen {
				// Special case: AddAnotherAppScreen needs app count and architecture
				screenModel := factory(len(m.state.Project.Applications), m.state.Project.Architecture)
				if teaModel, ok := screenModel.(tea.Model); ok {
					m.addScreenModel(newScreen, teaModel)
				}
			} else {
				screenModel := factory()
				if teaModel, ok := screenModel.(tea.Model); ok {
					m.addScreenModel(newScreen, teaModel)
				}
			}
		}
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

	// Create proper spacing between panels with enhanced visual separation
	leftPanelWithShadow := lipgloss.NewStyle().
		MarginRight(1).
		Render(leftPanel)
	
	rightPanelWithShadow := lipgloss.NewStyle().
		MarginLeft(1).
		Render(rightPanel)
	
	// Enhanced spacer with neon accent
	spacer := lipgloss.NewStyle().
		Width(3).
		Background(styles.ColorBgSecondary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderNeon).
		Render(" ")
	
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanelWithShadow,
		spacer,
		rightPanelWithShadow,
	)

	// Add error display overlay if there's an error
	var finalContent string
	if m.errorDisplay.HasError() {
		errorOverlay := m.errorDisplay.Render(m.windowWidth, m.windowHeight)
		if errorOverlay != "" {
			// For modal errors, overlay on top of content
			if m.errorDisplay.HasError() {
				finalContent = lipgloss.Place(
					m.windowWidth,
					m.windowHeight,
					lipgloss.Center,
					lipgloss.Center,
					content,
				)
				finalContent = lipgloss.Place(
					m.windowWidth,
					m.windowHeight,
					lipgloss.Center,
					lipgloss.Center,
					finalContent+"\n"+errorOverlay,
				)
			} else {
				finalContent = lipgloss.Place(
					m.windowWidth,
					m.windowHeight,
					lipgloss.Center,
					lipgloss.Center,
					content,
				)
			}
		} else {
			finalContent = lipgloss.Place(
				m.windowWidth,
				m.windowHeight,
				lipgloss.Center,
				lipgloss.Center,
				content,
			)
		}
	} else {
		finalContent = lipgloss.Place(
			m.windowWidth,
			m.windowHeight,
			lipgloss.Center,
			lipgloss.Center,
			content,
		)
	}

	// Create fullscreen background with centered content
	return lipgloss.NewStyle().
		Background(styles.ColorBgPrimary).
		Width(m.windowWidth).
		Height(m.windowHeight).
		Render(finalContent)
}

func (m Model) renderLeftPanel() string {
	var content string

	// Title section with improved spacing
	content += components.RenderTitle() + "\n\n"

	// Main content area with better spacing
	if screenModel, exists := m.screenModels[m.state.CurrentScreen]; exists {
		content += screenModel.View()
	}

	// Progress section with improved spacing
	content += "\n\n" + components.RenderProgressIndicator(m.state.CurrentScreen)
	
	// Help text at the bottom with improved spacing
	content += "\n\n" + m.getHelpText()
	
	// Add inline error display if there's an error
	if m.errorDisplay.HasError() {
		errorContent := m.errorDisplay.Render(m.windowWidth, m.windowHeight)
		if errorContent != "" {
			content += "\n\n" + errorContent
		}
	}

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
		return components.RenderHelp("↑↓: navigate • space/enter: select • s: skip • backspace: back • esc: quit")
	case models.YAMLPreviewScreen:
		return components.RenderHelp("↑↓: navigate • enter: select • ctrl+j/k: scroll • backspace: back • esc: quit")
	case models.GeneratingScreen:
		return ""
	case models.CompleteScreen:
		return ""
	default:
		return components.RenderHelp("j/k: navigate • space: select • enter: continue • backspace: back • esc: quit")
	}
}

// handleErrorRecovery applies recovery actions for recoverable errors
func (m Model) handleErrorRecovery(err *errors.TeapotError) (Model, tea.Cmd) {
	switch err.Type {
	case errors.ErrorTypeNavigation:
		// Navigation errors: go back to welcome screen
		m.state.CurrentScreen = models.WelcomeScreen
		return m, nil
		
	case errors.ErrorTypeUI:
		// UI errors: refresh current screen
		if screenModel, exists := m.screenModels[m.state.CurrentScreen]; exists {
			// Reinitialize the current screen
			if initCmd := screenModel.Init(); initCmd != nil {
				return m, m.errorRecovery.WrapCommand(initCmd)
			}
		}
		return m, nil
		
	case errors.ErrorTypePanic:
		// Panic errors: reset to safe state
		m.state.CurrentScreen = models.WelcomeScreen
		m.state.CurrentApp = nil
		m.state.Quitting = false
		return m, nil
		
	case errors.ErrorTypeValidation:
		// Validation errors: stay on current screen, error already displayed
		return m, nil
		
	default:
		// Unknown error type: stay on current screen
		return m, nil
	}
}

