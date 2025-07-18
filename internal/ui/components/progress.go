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
	steps := []string{"Setup", "Architecture", "Apps", "Config", "Tools", "Infrastructure", "CI/CD", "AI Tools", "Preview"}
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
		models.YAMLPreviewScreen:    8,
		models.GeneratingScreen:     9,
		models.CompleteScreen:       9,
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
			// Completed step - green checkmark
			icon = "✓"
			style = lipgloss.NewStyle().
				Foreground(styles.ColorSuccess).
				Bold(true)
		} else if i == currentStep {
			// Current step - neon cyan dot
			icon = "●"
			style = lipgloss.NewStyle().
				Foreground(styles.ColorPrimary).
				Bold(true)
		} else {
			// Future step - muted circle
			icon = "○"
			style = lipgloss.NewStyle().
				Foreground(styles.ColorTextMuted)
		}

		if i == currentStep {
			// Show step name for current step
			label = style.Render(icon + " " + steps[i])
		} else {
			// Show just icon for other steps
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
	// Enhanced title with neon styling and vertical gradient effect
	title := "🫖  T E A P O T"
	subtitle := "Modern Monorepo Builder"
	
	// Main title with neon glow effect and larger font simulation
	titleStyle := styles.LargeHeaderStyle.Copy().
		Foreground(styles.ColorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(44).
		Padding(1, 2).
		Margin(0, 0, 1, 0).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.ColorBorderNeon)
	
	// Subtitle with improved contrast and modern styling
	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Align(lipgloss.Center).
		Width(44).
		Bold(true).
		Background(styles.ColorBgGradient).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderSecondary).
		Padding(0, 2).
		Margin(0, 0, 1, 0)
	
	return titleStyle.Render(title) + "\n" + subtitleStyle.Render(subtitle)
}

func RenderSubtitle(text string) string {
	return lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Underline(true).
		Align(lipgloss.Center).
		Width(44).
		Padding(0, 1).
		Background(styles.ColorBgSecondary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderNeon).
		Margin(1, 0, 1, 0).
		Render(text)
}