// Package navigation provides state machine logic for screen navigation in the Teapot CLI.
// It handles the flow between screens, back navigation, and conditional navigation paths
// based on application state.
package navigation

import (
	"teapot/internal/models"
	"teapot/internal/ui/screens"
)

// NavigationFlow manages the navigation state machine for the Teapot CLI application.
// It provides clean abstractions for screen transitions and back navigation logic.
type NavigationFlow struct {
	// transitions maps each screen to its previous screen for back navigation
	transitions map[models.Screen]models.Screen
	// conditionalTransitions handles complex navigation logic based on app state
	conditionalTransitions map[models.Screen]func(*models.AppState) models.Screen
}

// NewNavigationFlow creates a new navigation flow state machine.
// It initializes the basic navigation transitions and conditional logic.
func NewNavigationFlow() *NavigationFlow {
	nf := &NavigationFlow{
		transitions: make(map[models.Screen]models.Screen),
		conditionalTransitions: make(map[models.Screen]func(*models.AppState) models.Screen),
	}
	
	// Define basic navigation transitions
	nf.transitions[models.ProjectSetupScreen] = models.WelcomeScreen
	nf.transitions[models.ArchitectureScreen] = models.ProjectSetupScreen
	nf.transitions[models.AppConfigScreen] = models.AddAppsScreen
	nf.transitions[models.DevToolsScreen] = models.AddAnotherAppScreen
	nf.transitions[models.InfrastructureScreen] = models.DevToolsScreen
	nf.transitions[models.CIPipelineScreen] = models.InfrastructureScreen
	nf.transitions[models.AIToolsScreen] = models.CIPipelineScreen
	
	// Define conditional navigation logic
	nf.conditionalTransitions[models.AddAppsScreen] = func(state *models.AppState) models.Screen {
		if len(state.Project.Applications) > 0 {
			// If we have apps, go back to "add another" screen
			return models.AddAnotherAppScreen
		}
		// No apps yet, go back to architecture
		return models.ArchitectureScreen
	}
	
	nf.conditionalTransitions[models.AddAnotherAppScreen] = func(state *models.AppState) models.Screen {
		// This is handled specially in the main navigation logic
		// as it involves removing the last app and going back to its config
		return models.AppConfigScreen
	}
	
	return nf
}

// GetPreviousScreen returns the previous screen for back navigation.
// It handles both simple transitions and conditional logic based on app state.
func (nf *NavigationFlow) GetPreviousScreen(current models.Screen, state *models.AppState) models.Screen {
	// Check for conditional transitions first
	if conditionFunc, exists := nf.conditionalTransitions[current]; exists {
		return conditionFunc(state)
	}
	
	// Fall back to simple transitions
	if previous, exists := nf.transitions[current]; exists {
		return previous
	}
	
	// No transition defined, stay on current screen
	return current
}

// CanNavigateBack returns true if the user can navigate back from the current screen.
// Some screens (like generation and complete) don't allow back navigation.
func (nf *NavigationFlow) CanNavigateBack(current models.Screen) bool {
	switch current {
	case models.WelcomeScreen, models.GeneratingScreen, models.CompleteScreen:
		return false
	default:
		return true
	}
}

// NavigateBack performs the back navigation operation and returns any necessary state changes.
// It handles complex cases like removing the last app when navigating back from AddAnotherAppScreen.
func (nf *NavigationFlow) NavigateBack(current models.Screen, state *models.AppState) (*models.AppState, models.Screen, interface{}) {
	if !nf.CanNavigateBack(current) {
		return state, current, nil
	}
	
	// Handle special case: navigating back from AddAnotherAppScreen
	if current == models.AddAnotherAppScreen {
		if len(state.Project.Applications) > 0 {
			// Remove the last app and go back to its configuration
			lastApp := state.Project.Applications[len(state.Project.Applications)-1]
			state.Project.Applications = state.Project.Applications[:len(state.Project.Applications)-1]
			state.CurrentApp = &lastApp
			return state, models.AppConfigScreen, screens.NewAppConfigModel(lastApp.Type)
		}
		// No apps, go back to app selection
		return state, models.AddAppsScreen, nil
	}
	
	// Handle special case: navigating back from AppConfigScreen
	if current == models.AppConfigScreen {
		// Clear the current app being configured
		state.CurrentApp = nil
	}
	
	// Get the previous screen
	previousScreen := nf.GetPreviousScreen(current, state)
	
	return state, previousScreen, nil
}

// GetScreenFactory returns a factory function for creating new screen models.
// This centralizes screen creation logic and ensures consistent initialization.
func (nf *NavigationFlow) GetScreenFactory(screen models.Screen) func(...interface{}) interface{} {
	switch screen {
	case models.WelcomeScreen:
		return func(...interface{}) interface{} { return screens.NewWelcomeModel() }
	case models.ProjectSetupScreen:
		return func(...interface{}) interface{} { return screens.NewProjectSetupModel() }
	case models.ArchitectureScreen:
		return func(...interface{}) interface{} { return screens.NewArchitectureModel() }
	case models.AddAppsScreen:
		return func(...interface{}) interface{} { return screens.NewAddAppsModel() }
	case models.AppConfigScreen:
		return func(args ...interface{}) interface{} {
			if len(args) > 0 {
				if appType, ok := args[0].(models.AppType); ok {
					return screens.NewAppConfigModel(appType)
				}
			}
			return screens.NewAppConfigModel(models.AppTypeReact) // Default fallback
		}
	case models.AddAnotherAppScreen:
		return func(args ...interface{}) interface{} {
			appCount := 0
			architecture := models.ArchitectureTurborepo
			if len(args) > 0 {
				if count, ok := args[0].(int); ok {
					appCount = count
				}
			}
			if len(args) > 1 {
				if arch, ok := args[1].(models.ArchitectureType); ok {
					architecture = arch
				}
			}
			return screens.NewAddAnotherAppModel(appCount, architecture)
		}
	case models.DevToolsScreen:
		return func(...interface{}) interface{} { return screens.NewDevToolsModel() }
	case models.InfrastructureScreen:
		return func(...interface{}) interface{} { return screens.NewInfrastructureModel() }
	case models.CIPipelineScreen:
		return func(...interface{}) interface{} { return screens.NewCIPipelineModel() }
	case models.AIToolsScreen:
		return func(...interface{}) interface{} { return screens.NewAIToolsModel() }
	case models.GeneratingScreen:
		return func(...interface{}) interface{} { return screens.NewGeneratingModel() }
	case models.CompleteScreen:
		return func(args ...interface{}) interface{} {
			projectName := "project"
			if len(args) > 0 {
				if name, ok := args[0].(string); ok {
					projectName = name
				}
			}
			return screens.NewCompleteModel(projectName)
		}
	default:
		return func(...interface{}) interface{} { return screens.NewWelcomeModel() }
	}
}

// GetNextScreen determines the next screen based on the current screen and a completion message.
// This centralizes the forward navigation logic.
func (nf *NavigationFlow) GetNextScreen(current models.Screen, msgType string) models.Screen {
	switch current {
	case models.WelcomeScreen:
		if msgType == "WelcomeComplete" {
			return models.ProjectSetupScreen
		}
	case models.ProjectSetupScreen:
		if msgType == "ProjectSetupComplete" {
			return models.ArchitectureScreen
		}
	case models.ArchitectureScreen:
		if msgType == "ArchitectureSelected" {
			return models.AddAppsScreen
		}
	case models.AddAppsScreen:
		if msgType == "AppTypeSelected" {
			return models.AppConfigScreen
		}
	case models.AppConfigScreen:
		if msgType == "AppConfigComplete" {
			return models.AddAnotherAppScreen
		}
	case models.AddAnotherAppScreen:
		if msgType == "AddAnotherAppSelected" {
			// This is handled in the main logic with the action parameter
			return models.AddAppsScreen // Default, but will be overridden
		}
	case models.DevToolsScreen:
		if msgType == "DevToolsSelected" {
			return models.InfrastructureScreen
		}
	case models.InfrastructureScreen:
		if msgType == "InfrastructureSelected" {
			return models.CIPipelineScreen
		}
	case models.CIPipelineScreen:
		if msgType == "CIPipelineSelected" {
			return models.AIToolsScreen
		}
	case models.AIToolsScreen:
		if msgType == "AIToolsSelected" {
			return models.GeneratingScreen
		}
	case models.GeneratingScreen:
		if msgType == "GenerationComplete" {
			return models.CompleteScreen
		}
	}
	
	// No transition defined
	return current
}