// Package models defines the core data structures and types for the Teapot CLI application.
// It provides screen definitions, application types, and project configuration structures
// used throughout the Bubble Tea interface.
package models

// Screen represents the different screens available in the Teapot CLI interface.
// Each screen corresponds to a step in the project setup process.
type Screen int

const (
	// WelcomeScreen is the initial screen showing the application overview
	WelcomeScreen Screen = iota
	// ProjectSetupScreen allows users to enter project name and description
	ProjectSetupScreen
	// ArchitectureScreen lets users choose their monorepo architecture
	ArchitectureScreen
	// AddAppsScreen provides options to select application types
	AddAppsScreen
	// AppConfigScreen configures individual application settings
	AppConfigScreen
	// AddAnotherAppScreen asks if the user wants to add more applications
	AddAnotherAppScreen
	// DevToolsScreen configures development tools like linting and TypeScript
	DevToolsScreen
	// InfrastructureScreen sets up infrastructure options like Docker and cloud providers
	InfrastructureScreen
	// CIPipelineScreen configures CI/CD pipeline settings
	CIPipelineScreen
	// AIToolsScreen configures AI development tools integration
	AIToolsScreen
	// YAMLPreviewScreen shows the generated teapot.yml configuration
	YAMLPreviewScreen
	// GeneratingScreen shows the project generation progress
	GeneratingScreen
	// CompleteScreen displays completion message and next steps
	CompleteScreen
)

// ScreenNames provides human-readable names for each screen type.
// This is used for progress indicators and debugging.
var ScreenNames = map[Screen]string{
	WelcomeScreen:        "Welcome",
	ProjectSetupScreen:   "Project Setup",
	ArchitectureScreen:   "Architecture",
	AddAppsScreen:        "Add Applications",
	AppConfigScreen:      "App Configuration",
	AddAnotherAppScreen:  "Add Another App",
	DevToolsScreen:       "Development Tools",
	InfrastructureScreen: "Infrastructure",
	CIPipelineScreen:     "CI/CD Pipeline",
	AIToolsScreen:        "AI Tools",
	GeneratingScreen:     "Generating",
	CompleteScreen:       "Complete",
}

// AppType represents the different types of applications that can be created
// in a Teapot monorepo project.
type AppType string

const (
	// AppTypeReact represents a React frontend application
	AppTypeReact        AppType = "react"
	// AppTypeNext represents a Next.js full-stack application
	AppTypeNext         AppType = "next"
	// AppTypeTanStack represents a TanStack Start application
	AppTypeTanStack     AppType = "tanstack-start"
	// AppTypeExpo represents an Expo React Native mobile application
	AppTypeExpo         AppType = "expo"
	// AppTypeNest represents a Nest.js backend API application
	AppTypeNest         AppType = "nest"
	// AppTypeBasicNode represents a basic Node.js backend application
	AppTypeBasicNode    AppType = "basic-node"
)

// AppTypeNames provides human-readable names for each application type.
// This is used in the UI for display purposes.
var AppTypeNames = map[AppType]string{
	AppTypeReact:     "React",
	AppTypeNext:      "Next.js",
	AppTypeTanStack:  "TanStack Start",
	AppTypeExpo:      "Expo",
	AppTypeNest:      "Nest.js",
	AppTypeBasicNode: "Basic Node.js",
}

// ArchitectureType represents the different monorepo architecture options
// that can be selected for the project.
type ArchitectureType string

const (
	// ArchitectureTurborepo represents a Turborepo-based monorepo setup (recommended)
	ArchitectureTurborepo ArchitectureType = "turborepo"
	// ArchitectureSingle represents a single application setup
	ArchitectureSingle    ArchitectureType = "single"
	// ArchitectureNx represents an Nx-based monorepo setup (coming soon)
	ArchitectureNx        ArchitectureType = "nx"
)

// ArchitectureNames provides human-readable names for each architecture type.
// This is used in the UI for display purposes.
var ArchitectureNames = map[ArchitectureType]string{
	ArchitectureTurborepo: "Turborepo (Recommended)",
	ArchitectureSingle:    "Single Application",
	ArchitectureNx:        "Nx (Coming Soon)",
}

// Application represents a single application within the monorepo project.
// It contains the configuration and metadata for generating the application.
type Application struct {
	// ID is a unique identifier for the application
	ID          string
	// Name is the display name of the application
	Name        string
	// Type specifies the framework/technology used for this application
	Type        AppType
	// Description provides a brief description of the application's purpose
	Description string
	// Options contains framework-specific configuration options
	Options     map[string]interface{}
}

// DevTools contains configuration for development tools and workflows.
// This includes linting, TypeScript setup, and git hooks.
type DevTools struct {
	// Linting specifies the linting tool to use: "prettier-eslint", "biome", or "custom"
	Linting     string
	// TypeScript indicates whether TypeScript should be configured
	TypeScript  bool
	// Husky indicates whether Husky git hooks should be set up
	Husky       bool
	// LintStaged indicates whether lint-staged should be configured
	LintStaged  bool
}

// Infrastructure contains configuration for infrastructure and deployment options.
// This includes containerization, infrastructure-as-code, and cloud providers.
type Infrastructure struct {
	// Docker indicates whether Docker should be configured
	Docker         bool
	// DockerCompose indicates whether Docker Compose should be set up
	DockerCompose  bool
	// Pulumi indicates whether Pulumi infrastructure-as-code should be configured
	Pulumi         bool
	// Terraform indicates whether Terraform infrastructure-as-code should be configured
	Terraform      bool
	// CloudProvider specifies the cloud provider: "aws", "vercel", "railway"
	CloudProvider  string
}

// CIPipeline contains configuration for CI/CD pipeline setup.
// This includes the provider and enabled features.
type CIPipeline struct {
	// Provider specifies the CI/CD provider: "github", "gitlab", "jenkins"
	Provider string
	// Features lists the enabled pipeline features: "testing", "linting", "docker", "deployment", "security"
	Features []string
}

// AITools contains configuration for AI-powered development tools.
// This includes editor selection and extensions.
type AITools struct {
	// Editor specifies the AI-powered editor: "claude-code", "cursor", "windsurf", "continue"
	Editor     string
	// Extensions lists the extensions to be configured for the selected editor
	Extensions []string
}

// ProjectConfig holds the complete configuration for a Teapot project.
// This includes all selected options across different setup screens and is used
// to generate the final project structure.
type ProjectConfig struct {
	// Name is the project name (used for directory name and package.json)
	Name           string
	// Description is an optional project description
	Description    string
	// Architecture specifies the monorepo architecture to use
	Architecture   ArchitectureType
	// Applications contains all applications to be created in the project
	Applications   []Application
	// DevTools contains development tools configuration
	DevTools       DevTools
	// Infrastructure contains infrastructure and deployment configuration
	Infrastructure Infrastructure
	// CIPipeline contains CI/CD pipeline configuration
	CIPipeline     CIPipeline
	// AITools contains AI development tools configuration
	AITools        AITools
}

// AppState manages the current state of the Teapot CLI application.
// This includes the current screen, project configuration, and navigation state.
type AppState struct {
	// CurrentScreen indicates which screen is currently being displayed
	CurrentScreen    Screen
	// Project holds the current project configuration being built
	Project          ProjectConfig
	// CurrentApp points to the application currently being configured (nil when not configuring)
	CurrentApp       *Application
	// CurrentAppIndex tracks the current position when iterating through applications
	CurrentAppIndex  int
	// Quitting indicates whether the application is in the process of quitting
	Quitting         bool
}