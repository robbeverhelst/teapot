package components

import (
	"fmt"
	"strings"

	"teapot/internal/models"
	"teapot/internal/ui/styles"
	
	"github.com/charmbracelet/lipgloss"
)

func RenderProgressIndicator(currentScreen models.Screen) string {
	// Map screens to progress steps (skip welcome screen)
	steps := []string{"Setup", "Architecture", "Apps", "Config", "Tools", "Infrastructure", "CI/CD", "AI Tools"}
	screenToStep := map[models.Screen]int{
		models.WelcomeScreen:        -1, // Not shown in progress
		models.ProjectSetupScreen:   0,
		models.ArchitectureScreen:   1,
		models.AddAppsScreen:        2,
		models.AppConfigScreen:      3,
		models.AddAnotherAppScreen:  3, // Same step as app config
		models.DevToolsScreen:       4,
		models.InfrastructureScreen: 5,
		models.CIPipelineScreen:     6,
		models.AIToolsScreen:        7,
		models.GeneratingScreen:     8,
		models.CompleteScreen:       8,
	}

	currentStep := screenToStep[currentScreen]
	
	// Don't show progress on welcome screen
	if currentStep == -1 {
		return ""
	}

	totalSteps := len(steps)
	indicators := make([]string, totalSteps)

	for i := 0; i < totalSteps; i++ {
		var icon, label string
		var style lipgloss.Style

		if i < currentStep {
			icon = "✓"
			style = lipgloss.NewStyle().
				Foreground(styles.ColorSuccess).
				Bold(true)
		} else if i == currentStep {
			icon = "●"
			style = lipgloss.NewStyle().
				Foreground(styles.ColorPrimary).
				Bold(true)
		} else {
			icon = "○"
			style = lipgloss.NewStyle().
				Foreground(styles.ColorTextMuted)
		}

		if i == currentStep {
			label = style.Render(icon + " " + steps[i])
		} else {
			label = style.Render(icon)
		}

		indicators[i] = label
	}

	progressText := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Step %d of %d", currentStep+1, totalSteps))

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Margin(1, 0).
		Render(strings.Join(indicators, "  ") + "\n" + progressText)
}

func RenderHelp(helpText string) string {
	helpBox := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Padding(0, 2).
		Align(lipgloss.Center).
		Render("💡 " + helpText)
	
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(helpBox)
}

func RenderTitle() string {
	// Enhanced title with better typography and ASCII art
	title := "🫖  T E A P O T"
	subtitle := "Modern Monorepo Builder"
	
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(40).
		Padding(0, 1)
	
	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Align(lipgloss.Center).
		Width(40).
		Italic(true).
		Margin(1, 0, 0, 0)
	
	return titleStyle.Render(title) + "\n" + subtitleStyle.Render(subtitle)
}

func RenderSubtitle(text string) string {
	return lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Underline(true).
		Align(lipgloss.Center).
		Width(40).
		Margin(1, 0, 1, 0).
		Render(text)
}