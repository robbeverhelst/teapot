package ui

import (
	"testing"

	"teapot/internal/models"
	"teapot/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// Helper function to handle model updates with type assertion
func updateModel(model Model, msg tea.Msg) Model {
	updatedModel, _ := model.Update(msg)
	return updatedModel.(Model)
}

// TestCompleteUserFlow tests the complete user journey from welcome to completion
func TestCompleteUserFlow(t *testing.T) {
	// Initialize the main model
	model := NewModel()
	
	// Test 1: Initial state
	if model.state.CurrentScreen != models.WelcomeScreen {
		t.Errorf("Expected initial screen to be WelcomeScreen, got %v", model.state.CurrentScreen)
	}
	
	// Test 2: Welcome screen completion
	model = updateModel(model, screens.WelcomeCompleteMsg{})
	if model.state.CurrentScreen != models.ProjectSetupScreen {
		t.Errorf("Expected screen to be ProjectSetupScreen after welcome completion, got %v", model.state.CurrentScreen)
	}
	
	// Test 3: Project setup completion
	model = updateModel(model, screens.ProjectSetupCompleteMsg{
		ProjectName: "test-project",
		Description: "A test project for integration testing",
	})
	if model.state.CurrentScreen != models.ArchitectureScreen {
		t.Errorf("Expected screen to be ArchitectureScreen after project setup, got %v", model.state.CurrentScreen)
	}
	if model.state.Project.Name != "test-project" {
		t.Errorf("Expected project name to be 'test-project', got '%s'", model.state.Project.Name)
	}
	
	// Test 4: Architecture selection
	model = updateModel(model, screens.ArchitectureSelectedMsg{
		Architecture: models.ArchitectureTurborepo,
	})
	if model.state.CurrentScreen != models.AddAppsScreen {
		t.Errorf("Expected screen to be AddAppsScreen after architecture selection, got %v", model.state.CurrentScreen)
	}
	if model.state.Project.Architecture != models.ArchitectureTurborepo {
		t.Errorf("Expected architecture to be Turborepo, got %v", model.state.Project.Architecture)
	}
	
	// Test 5: App type selection
	model = updateModel(model, screens.AppTypeSelectedMsg{
		AppType: models.AppTypeReact,
	})
	if model.state.CurrentScreen != models.AppConfigScreen {
		t.Errorf("Expected screen to be AppConfigScreen after app type selection, got %v", model.state.CurrentScreen)
	}
	if model.state.CurrentApp == nil {
		t.Error("Expected CurrentApp to be set after app type selection")
	} else if model.state.CurrentApp.Type != models.AppTypeReact {
		t.Errorf("Expected current app type to be React, got %v", model.state.CurrentApp.Type)
	}
	
	// Test 6: App configuration completion
	model = updateModel(model, screens.AppConfigCompleteMsg{
		AppName: "web-app",
		Options: map[string]interface{}{
			"typescript": true,
			"routing":    "react-router",
		},
	})
	if model.state.CurrentScreen != models.AddAnotherAppScreen {
		t.Errorf("Expected screen to be AddAnotherAppScreen after app config, got %v", model.state.CurrentScreen)
	}
	if len(model.state.Project.Applications) != 1 {
		t.Errorf("Expected 1 application in project, got %d", len(model.state.Project.Applications))
	}
	if model.state.Project.Applications[0].Name != "web-app" {
		t.Errorf("Expected app name to be 'web-app', got '%s'", model.state.Project.Applications[0].Name)
	}
	
	// Test 7: Continue without adding another app
	model = updateModel(model, screens.AddAnotherAppSelectedMsg{
		Action: "continue",
	})
	if model.state.CurrentScreen != models.DevToolsScreen {
		t.Errorf("Expected screen to be DevToolsScreen after choosing to continue, got %v", model.state.CurrentScreen)
	}
	
	// Test 8: Dev tools selection
	model = updateModel(model, screens.DevToolsSelectedMsg{
		LintingTool: "prettier-eslint",
	})
	if model.state.CurrentScreen != models.InfrastructureScreen {
		t.Errorf("Expected screen to be InfrastructureScreen after dev tools, got %v", model.state.CurrentScreen)
	}
	if model.state.Project.DevTools.Linting != "prettier-eslint" {
		t.Errorf("Expected linting tool to be 'prettier-eslint', got '%s'", model.state.Project.DevTools.Linting)
	}
	
	// Test 9: Infrastructure selection
	model = updateModel(model, screens.InfrastructureSelectedMsg{
		Options: map[string]bool{
			"docker":         true,
			"docker-compose": true,
			"pulumi":         false,
			"terraform":      false,
		},
	})
	if model.state.CurrentScreen != models.CIPipelineScreen {
		t.Errorf("Expected screen to be CIPipelineScreen after infrastructure, got %v", model.state.CurrentScreen)
	}
	if !model.state.Project.Infrastructure.Docker {
		t.Error("Expected Docker to be enabled")
	}
	
	// Test 10: CI/CD pipeline selection
	model = updateModel(model, screens.CIPipelineSelectedMsg{
		Provider: "github",
		Features: []string{"testing", "linting", "docker"},
	})
	if model.state.CurrentScreen != models.AIToolsScreen {
		t.Errorf("Expected screen to be AIToolsScreen after CI/CD, got %v", model.state.CurrentScreen)
	}
	if model.state.Project.CIPipeline.Provider != "github" {
		t.Errorf("Expected CI provider to be 'github', got '%s'", model.state.Project.CIPipeline.Provider)
	}
	
	// Test 11: AI tools selection
	model = updateModel(model, screens.AIToolsSelectedMsg{
		Editor:     "claude-code",
		Extensions: []string{"prettier", "eslint"},
	})
	if model.state.CurrentScreen != models.YAMLPreviewScreen {
		t.Errorf("Expected screen to be YAMLPreviewScreen after AI tools, got %v", model.state.CurrentScreen)
	}
	if model.state.Project.AITools.Editor != "claude-code" {
		t.Errorf("Expected AI editor to be 'claude-code', got '%s'", model.state.Project.AITools.Editor)
	}
	
	// Test 12: YAML preview continue
	model = updateModel(model, screens.YAMLContinueMsg{Project: model.state.Project})
	if model.state.CurrentScreen != models.GeneratingScreen {
		t.Errorf("Expected screen to be GeneratingScreen after YAML continue, got %v", model.state.CurrentScreen)
	}
	
	// Test 13: Generation completion
	model = updateModel(model, screens.GenerationCompleteMsg{})
	if model.state.CurrentScreen != models.CompleteScreen {
		t.Errorf("Expected screen to be CompleteScreen after generation, got %v", model.state.CurrentScreen)
	}
	
	// Test 14: Final project validation
	project := model.state.Project
	if project.Name != "test-project" {
		t.Errorf("Final project name incorrect: expected 'test-project', got '%s'", project.Name)
	}
	if project.Architecture != models.ArchitectureTurborepo {
		t.Errorf("Final architecture incorrect: expected Turborepo, got %v", project.Architecture)
	}
	if len(project.Applications) != 1 {
		t.Errorf("Final application count incorrect: expected 1, got %d", len(project.Applications))
	}
	if project.Applications[0].Type != models.AppTypeReact {
		t.Errorf("Final app type incorrect: expected React, got %v", project.Applications[0].Type)
	}
}

// TestBackNavigationFlow tests the back navigation throughout the user flow
func TestBackNavigationFlow(t *testing.T) {
	// Initialize and navigate to a middle screen
	model := NewModel()
	
	// Navigate to DevToolsScreen
	model = updateModel(model, screens.WelcomeCompleteMsg{})
	model = updateModel(model, screens.ProjectSetupCompleteMsg{ProjectName: "test", Description: ""})
	model = updateModel(model, screens.ArchitectureSelectedMsg{Architecture: models.ArchitectureTurborepo})
	model = updateModel(model, screens.AppTypeSelectedMsg{AppType: models.AppTypeReact})
	model = updateModel(model, screens.AppConfigCompleteMsg{AppName: "web", Options: map[string]interface{}{}})
	model = updateModel(model, screens.AddAnotherAppSelectedMsg{Action: "continue"})
	
	// Should be at DevToolsScreen
	if model.state.CurrentScreen != models.DevToolsScreen {
		t.Errorf("Expected to be at DevToolsScreen, got %v", model.state.CurrentScreen)
	}
	
	// Test back navigation
	model, _ = model.navigateBack()
	if model.state.CurrentScreen != models.AddAnotherAppScreen {
		t.Errorf("Expected back navigation to AddAnotherAppScreen, got %v", model.state.CurrentScreen)
	}
	
	// Navigate back again
	model, _ = model.navigateBack()
	if model.state.CurrentScreen != models.AppConfigScreen {
		t.Errorf("Expected back navigation to AppConfigScreen, got %v", model.state.CurrentScreen)
	}
	
	// Check that the last app was removed from the project
	if len(model.state.Project.Applications) != 0 {
		t.Errorf("Expected applications to be removed during back navigation, got %d", len(model.state.Project.Applications))
	}
	
	// Check that CurrentApp is set for reconfiguration
	if model.state.CurrentApp == nil {
		t.Error("Expected CurrentApp to be set for reconfiguration")
	} else if model.state.CurrentApp.Type != models.AppTypeReact {
		t.Errorf("Expected current app type to be React, got %v", model.state.CurrentApp.Type)
	}
}

// TestWindowResizing tests window resize handling
func TestWindowResizing(t *testing.T) {
	model := NewModel()
	
	// Test initial window size
	if model.windowWidth != 0 || model.windowHeight != 0 {
		t.Error("Expected initial window dimensions to be 0")
	}
	
	// Send window size message
	model = updateModel(model, tea.WindowSizeMsg{Width: 120, Height: 40})
	
	if model.windowWidth != 120 || model.windowHeight != 40 {
		t.Errorf("Expected window dimensions to be 120x40, got %dx%d", model.windowWidth, model.windowHeight)
	}
	
	// Test that resizing works on different screens
	model = updateModel(model, screens.WelcomeCompleteMsg{})
	model = updateModel(model, tea.WindowSizeMsg{Width: 100, Height: 30})
	
	if model.windowWidth != 100 || model.windowHeight != 30 {
		t.Errorf("Expected window dimensions to be 100x30, got %dx%d", model.windowWidth, model.windowHeight)
	}
}

// TestAddMultipleApps tests adding multiple applications
func TestAddMultipleApps(t *testing.T) {
	model := NewModel()
	
	// Navigate to app selection
	model = updateModel(model, screens.WelcomeCompleteMsg{})
	model = updateModel(model, screens.ProjectSetupCompleteMsg{ProjectName: "multi-app", Description: ""})
	model = updateModel(model, screens.ArchitectureSelectedMsg{Architecture: models.ArchitectureTurborepo})
	
	// Add first app
	model = updateModel(model, screens.AppTypeSelectedMsg{AppType: models.AppTypeReact})
	model = updateModel(model, screens.AppConfigCompleteMsg{AppName: "frontend", Options: map[string]interface{}{}})
	model = updateModel(model, screens.AddAnotherAppSelectedMsg{Action: "add"})
	
	// Should be back at app selection
	if model.state.CurrentScreen != models.AddAppsScreen {
		t.Errorf("Expected to be back at AddAppsScreen, got %v", model.state.CurrentScreen)
	}
	
	// Add second app
	model = updateModel(model, screens.AppTypeSelectedMsg{AppType: models.AppTypeNest})
	model = updateModel(model, screens.AppConfigCompleteMsg{AppName: "backend", Options: map[string]interface{}{}})
	model = updateModel(model, screens.AddAnotherAppSelectedMsg{Action: "continue"})
	
	// Check that both apps are in the project
	if len(model.state.Project.Applications) != 2 {
		t.Errorf("Expected 2 applications, got %d", len(model.state.Project.Applications))
	}
	
	// Check app details
	if model.state.Project.Applications[0].Name != "frontend" || model.state.Project.Applications[0].Type != models.AppTypeReact {
		t.Error("First app details incorrect")
	}
	if model.state.Project.Applications[1].Name != "backend" || model.state.Project.Applications[1].Type != models.AppTypeNest {
		t.Error("Second app details incorrect")
	}
}

// TestQuitKeyHandling tests global quit key handling
func TestQuitKeyHandling(t *testing.T) {
	model := NewModel()
	
	// Test Ctrl+C quit
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model = updatedModel.(Model)
	if !model.state.Quitting {
		t.Error("Expected quitting to be true after Ctrl+C")
	}
	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}
	
	// Reset and test Esc quit
	model = NewModel()
	updatedModel, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(Model)
	if !model.state.Quitting {
		t.Error("Expected quitting to be true after Esc")
	}
	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}
	
	// Test that Esc doesn't quit on generation/complete screens
	model = NewModel()
	model.state.CurrentScreen = models.GeneratingScreen
	updatedModel, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(Model)
	if model.state.Quitting {
		t.Error("Expected quitting to be false during generation")
	}
	
	model.state.CurrentScreen = models.CompleteScreen
	updatedModel, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(Model)
	if model.state.Quitting {
		t.Error("Expected quitting to be false on complete screen")
	}
}