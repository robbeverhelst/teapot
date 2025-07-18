package models

import (
	"testing"
)

func TestScreenNames(t *testing.T) {
	tests := []struct {
		screen   Screen
		expected string
	}{
		{WelcomeScreen, "Welcome"},
		{ProjectSetupScreen, "Project Setup"},
		{ArchitectureScreen, "Architecture"},
		{AddAppsScreen, "Add Applications"},
		{AppConfigScreen, "App Configuration"},
		{AddAnotherAppScreen, "Add Another App"},
		{DevToolsScreen, "Development Tools"},
		{InfrastructureScreen, "Infrastructure"},
		{CIPipelineScreen, "CI/CD Pipeline"},
		{AIToolsScreen, "AI Tools"},
		{GeneratingScreen, "Generating"},
		{CompleteScreen, "Complete"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if name, exists := ScreenNames[tt.screen]; !exists {
				t.Errorf("Screen %v not found in ScreenNames map", tt.screen)
			} else if name != tt.expected {
				t.Errorf("Expected screen name '%s', but got '%s'", tt.expected, name)
			}
		})
	}
}

func TestAppTypeNames(t *testing.T) {
	tests := []struct {
		appType  AppType
		expected string
	}{
		{AppTypeReact, "React"},
		{AppTypeNext, "Next.js"},
		{AppTypeTanStack, "TanStack Start"},
		{AppTypeExpo, "Expo"},
		{AppTypeNest, "Nest.js"},
		{AppTypeBasicNode, "Basic Node.js"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if name, exists := AppTypeNames[tt.appType]; !exists {
				t.Errorf("AppType %v not found in AppTypeNames map", tt.appType)
			} else if name != tt.expected {
				t.Errorf("Expected app type name '%s', but got '%s'", tt.expected, name)
			}
		})
	}
}

func TestArchitectureNames(t *testing.T) {
	tests := []struct {
		architecture ArchitectureType
		expected     string
	}{
		{ArchitectureTurborepo, "Turborepo (Recommended)"},
		{ArchitectureSingle, "Single Application"},
		{ArchitectureNx, "Nx (Coming Soon)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if name, exists := ArchitectureNames[tt.architecture]; !exists {
				t.Errorf("ArchitectureType %v not found in ArchitectureNames map", tt.architecture)
			} else if name != tt.expected {
				t.Errorf("Expected architecture name '%s', but got '%s'", tt.expected, name)
			}
		})
	}
}

func TestAppStateInitialization(t *testing.T) {
	state := AppState{
		CurrentScreen: WelcomeScreen,
		Project:       ProjectConfig{},
		CurrentApp:    nil,
		Quitting:      false,
	}

	if state.CurrentScreen != WelcomeScreen {
		t.Errorf("Expected CurrentScreen to be WelcomeScreen, but got %v", state.CurrentScreen)
	}

	if state.CurrentApp != nil {
		t.Errorf("Expected CurrentApp to be nil, but got %v", state.CurrentApp)
	}

	if state.Quitting {
		t.Error("Expected Quitting to be false, but got true")
	}
}

func TestProjectConfigInitialization(t *testing.T) {
	config := ProjectConfig{
		Name:           "test-project",
		Description:    "A test project",
		Architecture:   ArchitectureTurborepo,
		Applications:   []Application{},
		DevTools:       DevTools{},
		Infrastructure: Infrastructure{},
		CIPipeline:     CIPipeline{},
		AITools:        AITools{},
	}

	if config.Name != "test-project" {
		t.Errorf("Expected Name to be 'test-project', but got '%s'", config.Name)
	}

	if config.Description != "A test project" {
		t.Errorf("Expected Description to be 'A test project', but got '%s'", config.Description)
	}

	if config.Architecture != ArchitectureTurborepo {
		t.Errorf("Expected Architecture to be ArchitectureTurborepo, but got %v", config.Architecture)
	}

	if len(config.Applications) != 0 {
		t.Errorf("Expected Applications to be empty, but got %d items", len(config.Applications))
	}
}

func TestApplicationStructure(t *testing.T) {
	app := Application{
		ID:          "app-react",
		Name:        "frontend",
		Type:        AppTypeReact,
		Description: "React frontend application",
		Options:     make(map[string]interface{}),
	}

	if app.ID != "app-react" {
		t.Errorf("Expected ID to be 'app-react', but got '%s'", app.ID)
	}

	if app.Name != "frontend" {
		t.Errorf("Expected Name to be 'frontend', but got '%s'", app.Name)
	}

	if app.Type != AppTypeReact {
		t.Errorf("Expected Type to be AppTypeReact, but got %v", app.Type)
	}

	if app.Options == nil {
		t.Error("Expected Options to be initialized map, but got nil")
	}
}

func TestDevToolsDefaults(t *testing.T) {
	devTools := DevTools{
		Linting:    "prettier-eslint",
		TypeScript: true,
		Husky:      true,
		LintStaged: true,
	}

	if devTools.Linting != "prettier-eslint" {
		t.Errorf("Expected Linting to be 'prettier-eslint', but got '%s'", devTools.Linting)
	}

	if !devTools.TypeScript {
		t.Error("Expected TypeScript to be true")
	}

	if !devTools.Husky {
		t.Error("Expected Husky to be true")
	}

	if !devTools.LintStaged {
		t.Error("Expected LintStaged to be true")
	}
}

func TestInfrastructureDefaults(t *testing.T) {
	infra := Infrastructure{
		Docker:        true,
		DockerCompose: true,
		Pulumi:        false,
		Terraform:     false,
		CloudProvider: "aws",
	}

	if !infra.Docker {
		t.Error("Expected Docker to be true")
	}

	if !infra.DockerCompose {
		t.Error("Expected DockerCompose to be true")
	}

	if infra.Pulumi {
		t.Error("Expected Pulumi to be false")
	}

	if infra.Terraform {
		t.Error("Expected Terraform to be false")
	}

	if infra.CloudProvider != "aws" {
		t.Errorf("Expected CloudProvider to be 'aws', but got '%s'", infra.CloudProvider)
	}
}

func TestCIPipelineConfiguration(t *testing.T) {
	pipeline := CIPipeline{
		Provider: "github",
		Features: []string{"testing", "linting", "docker"},
	}

	if pipeline.Provider != "github" {
		t.Errorf("Expected Provider to be 'github', but got '%s'", pipeline.Provider)
	}

	expectedFeatures := []string{"testing", "linting", "docker"}
	if len(pipeline.Features) != len(expectedFeatures) {
		t.Errorf("Expected %d features, but got %d", len(expectedFeatures), len(pipeline.Features))
	}

	for i, feature := range expectedFeatures {
		if i >= len(pipeline.Features) || pipeline.Features[i] != feature {
			t.Errorf("Expected feature '%s' at index %d, but got '%s'", feature, i, pipeline.Features[i])
		}
	}
}

func TestAIToolsConfiguration(t *testing.T) {
	aiTools := AITools{
		Editor:     "claude-code",
		Extensions: []string{"prettier", "eslint"},
	}

	if aiTools.Editor != "claude-code" {
		t.Errorf("Expected Editor to be 'claude-code', but got '%s'", aiTools.Editor)
	}

	expectedExtensions := []string{"prettier", "eslint"}
	if len(aiTools.Extensions) != len(expectedExtensions) {
		t.Errorf("Expected %d extensions, but got %d", len(expectedExtensions), len(aiTools.Extensions))
	}

	for i, ext := range expectedExtensions {
		if i >= len(aiTools.Extensions) || aiTools.Extensions[i] != ext {
			t.Errorf("Expected extension '%s' at index %d, but got '%s'", ext, i, aiTools.Extensions[i])
		}
	}
}