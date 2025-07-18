package components

import (
	"strings"
	"sync"
	"time"

	"teapot/internal/cache"
	"teapot/internal/models"
	"teapot/internal/ui/styles"
	
	"github.com/charmbracelet/lipgloss"
)

var (
	// structureCache is a singleton cache for project structure rendering
	structureCache *cache.StructureCache
	// cacheOnce ensures the cache is initialized only once
	cacheOnce sync.Once
)

// getStructureCache returns the singleton structure cache instance
func getStructureCache() *cache.StructureCache {
	cacheOnce.Do(func() {
		// Initialize cache with reasonable defaults:
		// - maxSize: 50 entries (should be plenty for terminal resizing)
		// - ttl: 5 minutes (structures change infrequently)
		structureCache = cache.NewStructureCache(50, 5*time.Minute)
	})
	return structureCache
}

func RenderProjectStructure(project models.ProjectConfig, terminalWidth, terminalHeight int) string {
	// Check cache first
	cache := getStructureCache()
	if cached, found := cache.GetStructure(project, terminalWidth, terminalHeight); found {
		return cached
	}
	
	if project.Name == "" {
		// Enhanced header with neon styling
		header := lipgloss.NewStyle().
			Foreground(styles.ColorAccent).
			Bold(true).
			Align(lipgloss.Center).
			Width(44).
			Padding(0, 1).
			Background(styles.ColorBgSecondary).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(styles.ColorBorderNeon).
			Margin(0, 0, 1, 0).
			Render("📁 Project Structure")
		
		// Interactive preview message
		emptyMessage := lipgloss.NewStyle().
			Foreground(styles.ColorTextPrimary).
			Align(lipgloss.Center).
			Bold(true).
			Padding(1, 2).
			Background(styles.ColorBgTertiary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorderSecondary).
			Render("⚡ Make selections to see\nyour project structure")
		
		emptyResult := styles.GetRightPanelStyle(terminalWidth, terminalHeight).Render(
			header + "\n\n" + emptyMessage,
		)
		
		// Cache the empty result
		cache.SetStructure(project, terminalWidth, terminalHeight, emptyResult)
		return emptyResult
	}

	var structure strings.Builder
	
	// Root folder with neon accent
	rootFolder := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Bold(true).
		Background(styles.ColorBgSecondary).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderNeon).
		Render("📁 " + project.Name + "/")
	structure.WriteString(rootFolder + "\n")

	if len(project.Applications) > 0 {
		appsFolder := lipgloss.NewStyle().
			Foreground(styles.ColorAccent).
			Bold(true).
			Padding(0, 1).
			Background(styles.ColorBgTertiary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorderSecondary).
			Render("  📁 apps/")
		structure.WriteString(appsFolder + "\n")
		
		for _, app := range project.Applications {
			appName := getAppFolderName(string(app.Type))
			appFolder := lipgloss.NewStyle().
				Foreground(styles.ColorSuccess).
				Bold(true).
				Render("    📁 " + appName + "/")
			structure.WriteString(appFolder + "\n")
			
			// Enhanced file styling with better contrast
			packageJson := lipgloss.NewStyle().
				Foreground(styles.ColorWarning).
				Bold(true).
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
				Foreground(styles.ColorTextNeon).
				Bold(true).
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

	// Root files with neon highlights
	rootPackageJson := lipgloss.NewStyle().
		Foreground(styles.ColorWarning).
		Bold(true).
		Background(styles.ColorBgSecondary).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorWarning).
		Render("  📄 package.json")
	structure.WriteString(rootPackageJson + "\n")
	
	gitignore := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Render("  📄 .gitignore")
	structure.WriteString(gitignore + "\n")
	
	readme := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Bold(true).
		Render("  📖 README.md")
	structure.WriteString(readme + "\n")

	// Enhanced header with neon styling
	header := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Align(lipgloss.Center).
		Width(44).
		Padding(0, 1).
		Background(styles.ColorBgSecondary).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.ColorBorderNeon).
		Margin(0, 0, 1, 0).
		Render("📁 Project Structure")
	
	content := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(structure.String())

	result := styles.GetRightPanelStyle(terminalWidth, terminalHeight).Render(
		header + "\n\n" + content,
	)
	
	// Cache the result
	cache.SetStructure(project, terminalWidth, terminalHeight, result)
	return result
}

// ClearStructureCache clears the entire structure cache.
// This is useful when the project configuration changes significantly.
func ClearStructureCache() {
	cache := getStructureCache()
	cache.Clear()
}

// GetStructureCacheStats returns statistics about the structure cache.
// This is useful for monitoring cache performance and debugging.
func GetStructureCacheStats() cache.CacheStats {
	cache := getStructureCache()
	return cache.GetStats()
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