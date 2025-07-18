// Package generator provides functionality for generating project configuration files
// and scaffolding project structures based on user selections.
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"teapot/internal/models"
	"gopkg.in/yaml.v3"
)

// TeapotConfig represents the complete configuration for a Teapot project
type TeapotConfig struct {
	Version     string                `yaml:"version"`
	Project     ProjectConfig         `yaml:"project"`
	Architecture string               `yaml:"architecture"`
	Applications []ApplicationConfig  `yaml:"applications"`
	DevTools    DevToolsConfig       `yaml:"devTools"`
	Infrastructure InfrastructureConfig `yaml:"infrastructure"`
	CIPipeline  CIPipelineConfig     `yaml:"ciPipeline"`
	AITools     AIToolsConfig        `yaml:"aiTools"`
}

type ProjectConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type ApplicationConfig struct {
	ID      string                 `yaml:"id"`
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}

type DevToolsConfig struct {
	Linting     string `yaml:"linting"`
	TypeScript  bool   `yaml:"typescript"`
	Husky       bool   `yaml:"husky"`
	LintStaged  bool   `yaml:"lintStaged"`
}

type InfrastructureConfig struct {
	Docker         bool `yaml:"docker"`
	DockerCompose  bool `yaml:"dockerCompose"`
	Pulumi         bool `yaml:"pulumi"`
	Terraform      bool `yaml:"terraform"`
}

type CIPipelineConfig struct {
	Provider string   `yaml:"provider"`
	Features []string `yaml:"features"`
}

type AIToolsConfig struct {
	Editor     string   `yaml:"editor"`
	Extensions []string `yaml:"extensions"`
}

// GenerateTeapotYAML generates a teapot.yml file from the project configuration
func GenerateTeapotYAML(project models.ProjectConfig) (string, error) {
	// Convert models.ProjectConfig to TeapotConfig
	config := TeapotConfig{
		Version: "1.0",
		Project: ProjectConfig{
			Name:        project.Name,
			Description: project.Description,
		},
		Architecture: string(project.Architecture),
		Applications: make([]ApplicationConfig, len(project.Applications)),
		DevTools: DevToolsConfig{
			Linting:    project.DevTools.Linting,
			TypeScript: project.DevTools.TypeScript,
			Husky:      project.DevTools.Husky,
			LintStaged: project.DevTools.LintStaged,
		},
		Infrastructure: InfrastructureConfig{
			Docker:        project.Infrastructure.Docker,
			DockerCompose: project.Infrastructure.DockerCompose,
			Pulumi:        project.Infrastructure.Pulumi,
			Terraform:     project.Infrastructure.Terraform,
		},
		CIPipeline: CIPipelineConfig{
			Provider: project.CIPipeline.Provider,
			Features: project.CIPipeline.Features,
		},
		AITools: AIToolsConfig{
			Editor:     project.AITools.Editor,
			Extensions: project.AITools.Extensions,
		},
	}

	// Convert applications
	for i, app := range project.Applications {
		config.Applications[i] = ApplicationConfig{
			ID:      app.ID,
			Name:    app.Name,
			Type:    string(app.Type),
			Options: app.Options,
		}
	}

	// Generate YAML content
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(yamlData), nil
}

// SaveTeapotYAML saves the teapot.yml file to the specified directory
func SaveTeapotYAML(project models.ProjectConfig, outputDir string) error {
	yamlContent, err := GenerateTeapotYAML(project)
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write YAML file
	filePath := filepath.Join(outputDir, "teapot.yml")
	if err := os.WriteFile(filePath, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// FormatYAMLForDisplay formats the YAML content for display in the terminal
func FormatYAMLForDisplay(yamlContent string) string {
	lines := strings.Split(yamlContent, "\n")
	var formatted []string
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		formatted = append(formatted, line)
	}
	
	return strings.Join(formatted, "\n")
}