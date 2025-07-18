package navigation

import (
	"testing"

	"teapot/internal/models"
)

func TestNavigationFlow_BasicTransitions(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test simple transitions
	tests := []struct {
		current  models.Screen
		expected models.Screen
	}{
		{models.ProjectSetupScreen, models.WelcomeScreen},
		{models.ArchitectureScreen, models.ProjectSetupScreen},
		{models.AppConfigScreen, models.AddAppsScreen},
		{models.DevToolsScreen, models.AddAnotherAppScreen},
		{models.InfrastructureScreen, models.DevToolsScreen},
		{models.CIPipelineScreen, models.InfrastructureScreen},
		{models.AIToolsScreen, models.CIPipelineScreen},
	}
	
	for _, tt := range tests {
		state := &models.AppState{CurrentScreen: tt.current}
		result := nf.GetPreviousScreen(tt.current, state)
		if result != tt.expected {
			t.Errorf("Expected %v -> %v, got %v", tt.current, tt.expected, result)
		}
	}
}

func TestNavigationFlow_ConditionalTransitions(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test AddAppsScreen with no applications
	state := &models.AppState{
		CurrentScreen: models.AddAppsScreen,
		Project:       models.ProjectConfig{Applications: []models.Application{}},
	}
	result := nf.GetPreviousScreen(models.AddAppsScreen, state)
	if result != models.ArchitectureScreen {
		t.Errorf("Expected AddAppsScreen with no apps to go to ArchitectureScreen, got %v", result)
	}
	
	// Test AddAppsScreen with applications
	state.Project.Applications = []models.Application{
		{Type: models.AppTypeReact, Name: "web"},
	}
	result = nf.GetPreviousScreen(models.AddAppsScreen, state)
	if result != models.AddAnotherAppScreen {
		t.Errorf("Expected AddAppsScreen with apps to go to AddAnotherAppScreen, got %v", result)
	}
}

func TestNavigationFlow_CanNavigateBack(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test screens that allow back navigation
	allowedScreens := []models.Screen{
		models.ProjectSetupScreen,
		models.ArchitectureScreen,
		models.AddAppsScreen,
		models.AppConfigScreen,
		models.AddAnotherAppScreen,
		models.DevToolsScreen,
		models.InfrastructureScreen,
		models.CIPipelineScreen,
		models.AIToolsScreen,
	}
	
	for _, screen := range allowedScreens {
		if !nf.CanNavigateBack(screen) {
			t.Errorf("Expected %v to allow back navigation", screen)
		}
	}
	
	// Test screens that don't allow back navigation
	restrictedScreens := []models.Screen{
		models.WelcomeScreen,
		models.GeneratingScreen,
		models.CompleteScreen,
	}
	
	for _, screen := range restrictedScreens {
		if nf.CanNavigateBack(screen) {
			t.Errorf("Expected %v to not allow back navigation", screen)
		}
	}
}

func TestNavigationFlow_NavigateBack(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test simple back navigation
	state := &models.AppState{
		CurrentScreen: models.ProjectSetupScreen,
		Project:       models.ProjectConfig{Name: "test"},
	}
	
	newState, newScreen, screenModel := nf.NavigateBack(models.ProjectSetupScreen, state)
	
	if newScreen != models.WelcomeScreen {
		t.Errorf("Expected navigation to WelcomeScreen, got %v", newScreen)
	}
	if screenModel != nil {
		t.Error("Expected no specific screen model for simple transition")
	}
	if newState.CurrentScreen != models.ProjectSetupScreen {
		t.Error("Expected state to be unchanged in NavigateBack")
	}
}

func TestNavigationFlow_NavigateBackFromAddAnotherApp(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test back navigation from AddAnotherAppScreen with apps
	app := models.Application{
		Type: models.AppTypeReact,
		Name: "web-app",
	}
	
	state := &models.AppState{
		CurrentScreen: models.AddAnotherAppScreen,
		Project: models.ProjectConfig{
			Applications: []models.Application{app},
		},
	}
	
	newState, newScreen, screenModel := nf.NavigateBack(models.AddAnotherAppScreen, state)
	
	if newScreen != models.AppConfigScreen {
		t.Errorf("Expected navigation to AppConfigScreen, got %v", newScreen)
	}
	if len(newState.Project.Applications) != 0 {
		t.Errorf("Expected applications to be removed, got %d", len(newState.Project.Applications))
	}
	if newState.CurrentApp == nil {
		t.Error("Expected CurrentApp to be set")
	} else if newState.CurrentApp.Type != models.AppTypeReact {
		t.Errorf("Expected CurrentApp type to be React, got %v", newState.CurrentApp.Type)
	}
	if screenModel == nil {
		t.Error("Expected screen model to be returned")
	}
}

func TestNavigationFlow_NavigateBackFromAppConfig(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test back navigation from AppConfigScreen
	currentApp := &models.Application{
		Type: models.AppTypeReact,
		Name: "web-app",
	}
	
	state := &models.AppState{
		CurrentScreen: models.AppConfigScreen,
		CurrentApp:    currentApp,
		Project:       models.ProjectConfig{},
	}
	
	newState, newScreen, screenModel := nf.NavigateBack(models.AppConfigScreen, state)
	
	if newScreen != models.AddAppsScreen {
		t.Errorf("Expected navigation to AddAppsScreen, got %v", newScreen)
	}
	if newState.CurrentApp != nil {
		t.Error("Expected CurrentApp to be cleared")
	}
	if screenModel != nil {
		t.Error("Expected no specific screen model")
	}
}

func TestNavigationFlow_GetScreenFactory(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test basic screen factories
	basicScreens := []models.Screen{
		models.WelcomeScreen,
		models.ProjectSetupScreen,
		models.ArchitectureScreen,
		models.AddAppsScreen,
		models.DevToolsScreen,
		models.InfrastructureScreen,
		models.CIPipelineScreen,
		models.AIToolsScreen,
		models.GeneratingScreen,
	}
	
	for _, screen := range basicScreens {
		factory := nf.GetScreenFactory(screen)
		if factory == nil {
			t.Errorf("Expected factory for %v, got nil", screen)
		} else {
			model := factory()
			if model == nil {
				t.Errorf("Expected model from factory for %v, got nil", screen)
			}
		}
	}
}

func TestNavigationFlow_GetScreenFactoryWithArgs(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test AppConfigScreen factory with arguments
	factory := nf.GetScreenFactory(models.AppConfigScreen)
	model := factory(models.AppTypeReact)
	if model == nil {
		t.Error("Expected model from AppConfigScreen factory with args")
	}
	
	// Test AddAnotherAppScreen factory with arguments
	factory = nf.GetScreenFactory(models.AddAnotherAppScreen)
	model = factory(2, models.ArchitectureTurborepo)
	if model == nil {
		t.Error("Expected model from AddAnotherAppScreen factory with args")
	}
	
	// Test CompleteScreen factory with arguments
	factory = nf.GetScreenFactory(models.CompleteScreen)
	model = factory("test-project")
	if model == nil {
		t.Error("Expected model from CompleteScreen factory with args")
	}
}

func TestNavigationFlow_GetNextScreen(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test forward navigation transitions
	tests := []struct {
		current     models.Screen
		msgType     string
		expected    models.Screen
	}{
		{models.WelcomeScreen, "WelcomeComplete", models.ProjectSetupScreen},
		{models.ProjectSetupScreen, "ProjectSetupComplete", models.ArchitectureScreen},
		{models.ArchitectureScreen, "ArchitectureSelected", models.AddAppsScreen},
		{models.AddAppsScreen, "AppTypeSelected", models.AppConfigScreen},
		{models.AppConfigScreen, "AppConfigComplete", models.AddAnotherAppScreen},
		{models.DevToolsScreen, "DevToolsSelected", models.InfrastructureScreen},
		{models.InfrastructureScreen, "InfrastructureSelected", models.CIPipelineScreen},
		{models.CIPipelineScreen, "CIPipelineSelected", models.AIToolsScreen},
		{models.AIToolsScreen, "AIToolsSelected", models.GeneratingScreen},
		{models.GeneratingScreen, "GenerationComplete", models.CompleteScreen},
	}
	
	for _, tt := range tests {
		result := nf.GetNextScreen(tt.current, tt.msgType)
		if result != tt.expected {
			t.Errorf("Expected %v + %s -> %v, got %v", tt.current, tt.msgType, tt.expected, result)
		}
	}
}

func TestNavigationFlow_GetNextScreenInvalid(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test invalid message types
	result := nf.GetNextScreen(models.WelcomeScreen, "InvalidMsg")
	if result != models.WelcomeScreen {
		t.Errorf("Expected invalid message to return current screen, got %v", result)
	}
	
	// Test invalid screen
	result = nf.GetNextScreen(models.CompleteScreen, "SomeMsg")
	if result != models.CompleteScreen {
		t.Errorf("Expected invalid screen to return current screen, got %v", result)
	}
}

func TestNavigationFlow_RestrictedBackNavigation(t *testing.T) {
	nf := NewNavigationFlow()
	
	// Test that restricted screens don't allow back navigation
	restrictedScreens := []models.Screen{
		models.WelcomeScreen,
		models.GeneratingScreen,
		models.CompleteScreen,
	}
	
	for _, screen := range restrictedScreens {
		state := &models.AppState{CurrentScreen: screen}
		newState, newScreen, screenModel := nf.NavigateBack(screen, state)
		
		// Should return unchanged state and screen
		if newState != state {
			t.Errorf("Expected unchanged state for %v", screen)
		}
		if newScreen != screen {
			t.Errorf("Expected unchanged screen for %v, got %v", screen, newScreen)
		}
		if screenModel != nil {
			t.Errorf("Expected no screen model for %v", screen)
		}
	}
}