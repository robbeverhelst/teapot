package models

type Screen int

const (
	WelcomeScreen Screen = iota
	ProjectSetupScreen
	ArchitectureScreen
	AddAppsScreen
	AppConfigScreen
	AddAnotherAppScreen
	DevToolsScreen
	InfrastructureScreen
	CIPipelineScreen
	AIToolsScreen
	GeneratingScreen
	CompleteScreen
)

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

type AppType string

const (
	AppTypeReact        AppType = "react"
	AppTypeNext         AppType = "next"
	AppTypeTanStack     AppType = "tanstack-start"
	AppTypeExpo         AppType = "expo"
	AppTypeNest         AppType = "nest"
	AppTypeBasicNode    AppType = "basic-node"
)

var AppTypeNames = map[AppType]string{
	AppTypeReact:     "React",
	AppTypeNext:      "Next.js",
	AppTypeTanStack:  "TanStack Start",
	AppTypeExpo:      "Expo",
	AppTypeNest:      "Nest.js",
	AppTypeBasicNode: "Basic Node.js",
}

type ArchitectureType string

const (
	ArchitectureTurborepo ArchitectureType = "turborepo"
	ArchitectureSingle    ArchitectureType = "single"
	ArchitectureNx        ArchitectureType = "nx" // Coming soon
)

var ArchitectureNames = map[ArchitectureType]string{
	ArchitectureTurborepo: "Turborepo (Recommended)",
	ArchitectureSingle:    "Single Application",
	ArchitectureNx:        "Nx (Coming Soon)",
}

type Application struct {
	ID          string
	Name        string
	Type        AppType
	Description string
	Options     map[string]interface{}
}

type DevTools struct {
	Linting     string // "prettier-eslint", "biome", "custom"
	TypeScript  bool
	Husky       bool
	LintStaged  bool
}

type Infrastructure struct {
	Docker         bool
	DockerCompose  bool
	Pulumi         bool
	Terraform      bool
	CloudProvider  string // "aws", "vercel", "railway"
}

type CIPipeline struct {
	Provider string // "github", "gitlab", "jenkins"
	Features []string // "testing", "linting", "docker", "deployment", "security"
}

type AITools struct {
	Editor     string // "claude-code", "cursor", "windsurf", "continue"
	Extensions []string
}

type ProjectConfig struct {
	Name           string
	Description    string
	Architecture   ArchitectureType
	Applications   []Application
	DevTools       DevTools
	Infrastructure Infrastructure
	CIPipeline     CIPipeline
	AITools        AITools
}

type AppState struct {
	CurrentScreen    Screen
	Project          ProjectConfig
	CurrentApp       *Application // For app configuration flow
	CurrentAppIndex  int         // For iterating through apps
	Quitting         bool
}