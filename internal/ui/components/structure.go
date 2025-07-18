package components

import (
	"strings"

	"teapot/internal/models"
	"teapot/internal/ui/styles"
	
	"github.com/charmbracelet/lipgloss"
)

func RenderProjectStructure(project models.ProjectConfig, terminalWidth, terminalHeight int) string {
	if project.Name == "" {
		// Enhanced header with better typography
		header := lipgloss.NewStyle().
			Foreground(styles.ColorAccent).
			Bold(true).
			Underline(true).
			Align(lipgloss.Center).
			Width(40).
			Margin(0, 0, 1, 0).
			Render("📁 Project Structure")
		
		emptyMessage := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Align(lipgloss.Center).
			Italic(true).
			Render("(empty)")
		
		return styles.GetRightPanelStyle(terminalWidth, terminalHeight).Render(
			header + "\n\n" + emptyMessage,
		)
	}

	var structure strings.Builder
	
	// Root folder with enhanced styling
	rootFolder := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Render("📁 " + project.Name + "/")
	structure.WriteString(rootFolder + "\n")

	if len(project.Applications) > 0 {
		appsFolder := lipgloss.NewStyle().
			Foreground(styles.ColorSecondary).
			Bold(true).
			Render("  📁 apps/")
		structure.WriteString(appsFolder + "\n")
		
		for _, app := range project.Applications {
			appName := getAppFolderName(string(app.Type))
			appFolder := lipgloss.NewStyle().
				Foreground(styles.ColorPrimary).
				Render("    📁 " + appName + "/")
			structure.WriteString(appFolder + "\n")
			
			// Enhanced file styling
			packageJson := lipgloss.NewStyle().
				Foreground(styles.ColorWarning).
				Render("      📄 package.json")
			structure.WriteString(packageJson + "\n")
			
			if app.Type == models.AppTypeNext {
				configFile := lipgloss.NewStyle().
					Foreground(styles.ColorTextSecondary).
					Render("      ⚙️ next.config.js")
				structure.WriteString(configFile + "\n")
			} else if app.Type == models.AppTypeExpo {
				configFile := lipgloss.NewStyle().
					Foreground(styles.ColorTextSecondary).
					Render("      ⚙️ metro.config.js")
				structure.WriteString(configFile + "\n")
			}
			
			srcFolder := lipgloss.NewStyle().
				Foreground(styles.ColorSuccess).
				Render("      📁 src/")
			structure.WriteString(srcFolder + "\n")
		}
	}

	// TODO: Add packages section back later

	// Infrastructure files based on configuration
	if project.CIPipeline.Provider == "github" {
		githubFolder := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render("  📁 .github/")
		structure.WriteString(githubFolder + "\n")
		
		workflowsFolder := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render("    📁 workflows/")
		structure.WriteString(workflowsFolder + "\n")
		
		ciFile := lipgloss.NewStyle().
			Foreground(styles.ColorSuccess).
			Render("      🔧 ci.yml")
		structure.WriteString(ciFile + "\n")
	}

	if project.Infrastructure.Docker || project.Infrastructure.DockerCompose {
		if project.Infrastructure.DockerCompose {
			dockerCompose := lipgloss.NewStyle().
				Foreground(styles.ColorPrimary).
				Render("  🐳 docker-compose.yml")
			structure.WriteString(dockerCompose + "\n")
		}
		
		if project.Infrastructure.Docker {
			dockerfile := lipgloss.NewStyle().
				Foreground(styles.ColorPrimary).
				Render("  🐳 Dockerfile")
			structure.WriteString(dockerfile + "\n")
		}
	}

	// Architecture-specific config files
	if project.Architecture == models.ArchitectureNx {
		nxFile := lipgloss.NewStyle().
			Foreground(styles.ColorSecondary).
			Render("  ⚙️ nx.json")
		structure.WriteString(nxFile + "\n")
	} else if project.Architecture == models.ArchitectureTurborepo {
		turboFile := lipgloss.NewStyle().
			Foreground(styles.ColorSecondary).
			Render("  ⚙️ turbo.json")
		structure.WriteString(turboFile + "\n")
	}

	// Root files with enhanced styling
	rootPackageJson := lipgloss.NewStyle().
		Foreground(styles.ColorWarning).
		Bold(true).
		Render("  📄 package.json")
	structure.WriteString(rootPackageJson + "\n")
	
	gitignore := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Render("  📄 .gitignore")
	structure.WriteString(gitignore + "\n")
	
	readme := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Render("  📖 README.md")
	structure.WriteString(readme + "\n")

	// Enhanced header with better typography
	header := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Underline(true).
		Align(lipgloss.Center).
		Width(40).
		Margin(0, 0, 1, 0).
		Render("📁 Project Structure")
	
	content := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(structure.String())

	return styles.GetRightPanelStyle(terminalWidth, terminalHeight).Render(
		header + "\n\n" + content,
	)
}

func getAppFolderName(appType string) string {
	switch appType {
	case "next":
		return "web"
	case "expo":
		return "mobile"
	case "nest":
		return "api"
	case "react":
		return "web"
	case "tanstack-start":
		return "web"
	case "basic-node":
		return "api"
	default:
		return "app"
	}
}

func getPackageFolderName(pkg string) string {
	switch pkg {
	case "UI Components Library":
		return "ui"
	case "TypeScript Config":
		return "tsconfig"
	case "ESLint Config":
		return "eslint-config"
	case "Prettier Config":
		return "prettier-config"
	case "Database Models":
		return "database"
	case "Auth Package":
		return "auth"
	case "Utils/Helpers":
		return "utils"
	case "API Client":
		return "api-client"
	default:
		return "package"
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}