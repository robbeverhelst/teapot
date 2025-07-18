// Package styles provides a comprehensive design system for the Teapot CLI application.
// It includes color palettes, typography styles, and responsive layout components
// that create a consistent and professional user interface.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Modern neon-accented color palette with high contrast
	ColorPrimary     = lipgloss.Color("#00FFCC")    // Neon cyan - main brand
	ColorSecondary   = lipgloss.Color("#00E5FF")    // Bright cyan - secondary actions
	ColorAccent      = lipgloss.Color("#00FFCC")    // Neon cyan - highlights
	ColorSuccess     = lipgloss.Color("#00FF88")    // Neon green - success states
	ColorWarning     = lipgloss.Color("#FFD700")    // Gold - warnings
	ColorError       = lipgloss.Color("#FF6B6B")    // Coral red - errors
	ColorDanger      = lipgloss.Color("#FF4757")    // Bright red - critical
	
	// Text colors with maximum contrast for readability
	ColorTextPrimary   = lipgloss.Color("#FFFFFF")  // Pure white
	ColorTextSecondary = lipgloss.Color("#F0F0F0")  // Off white
	ColorTextMuted     = lipgloss.Color("#B0B0B0")  // Light gray
	ColorTextDisabled  = lipgloss.Color("#808080")  // Medium gray
	ColorTextNeon      = lipgloss.Color("#00FFCC")  // Neon cyan text
	
	// Background colors - darker for better contrast
	ColorBgPrimary   = lipgloss.Color("#1E1E2E")    // Dark purple-gray
	ColorBgSecondary = lipgloss.Color("#2A2A3A")    // Lighter purple-gray
	ColorBgTertiary  = lipgloss.Color("#3A3A4A")    // Medium gray
	ColorBgSurface   = lipgloss.Color("#4A4A5A")    // Light surface
	ColorBgGradient  = lipgloss.Color("#2A2A3A")    // For gradient effects
	
	// Border colors with neon accents
	ColorBorderPrimary   = lipgloss.Color("#4A4A5A") // Subtle border
	ColorBorderSecondary = lipgloss.Color("#6A6A7A") // Visible border
	ColorBorderAccent    = lipgloss.Color("#00FFCC") // Neon cyan border
	ColorBorderNeon      = lipgloss.Color("#00FFCC") // Neon glow border
	
	// Legacy aliases for backward compatibility
	ColorText       = ColorTextPrimary
	ColorMuted      = ColorTextMuted
	ColorBorder     = ColorBorderPrimary
	ColorBackground = ColorBgPrimary
)

var (
	// Typography hierarchy with neon styling and better visual weight
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Margin(1, 0).
			Align(lipgloss.Center)
	
	// Large header with neon glow effect
	LargeHeaderStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Margin(1, 0).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNeon).
			Padding(1, 2)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Italic(true).
			Margin(0, 0, 1, 0)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Bold(true).
			Margin(1, 0, 0, 2)

	SectionHeaderStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Underline(true).
			Margin(1, 0, 0, 0)

	BodyTextStyle = lipgloss.NewStyle().
			Foreground(ColorTextSecondary).
			Margin(0, 0, 1, 0)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Bold(true).
			Margin(1, 0, 0, 0)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorBorderSecondary).
			Padding(2, 3)

	ProgressIndicatorStyle = lipgloss.NewStyle().
				Margin(1, 0, 0, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Padding(0, 1).
			Margin(1, 0, 0, 2)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1).
			Background(ColorBgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNeon)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorTextSecondary)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			Padding(0, 1).
			Background(ColorBgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess)
	
	// Neon glow effect for checkmarks
	NeonCheckStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1).
			Background(ColorBgSecondary).
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorBorderNeon)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	// Enhanced interactive styles with neon hover effects
	HoverStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1).
			Background(ColorBgTertiary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderNeon)
	
	// Subtle hover effect
	SubtleHoverStyle = lipgloss.NewStyle().
			Foreground(ColorTextPrimary).
			Background(ColorBgTertiary).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderSecondary)

	FocusedStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1).
			Background(ColorBgTertiary).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorBorderNeon)

	// Modern card-like styles
	CardStyle = lipgloss.NewStyle().
			Background(ColorBgSecondary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderSecondary).
			Padding(1, 2).
			Margin(1, 0)

	ElevatedCardStyle = lipgloss.NewStyle().
			Background(ColorBgSecondary).
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorBorderAccent).
			Padding(2, 3).
			Margin(1, 0)

	// Status indicator styles
	StatusActiveStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			Padding(0, 1)

	StatusInactiveStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Padding(0, 1)

	// Fullscreen specific styles
	FullscreenContainerStyle = lipgloss.NewStyle().
								Align(lipgloss.Center, lipgloss.Center).
								Background(ColorBgPrimary)
)

// GetLeftPanelStyle returns a dynamically sized left panel style with better proportions
func GetLeftPanelStyle(terminalWidth, terminalHeight int) lipgloss.Style {
	// Better proportions - 45% for left panel with more breathing room
	panelWidth := int(float64(terminalWidth) * 0.45) - 4
	maxPanelHeight := terminalHeight - 8

	// Ensure minimum dimensions for usability
	if panelWidth < 45 {
		panelWidth = 45
	}
	if maxPanelHeight < 25 {
		maxPanelHeight = 25
	}

	return BorderStyle.Copy().
		Width(panelWidth).
		Height(maxPanelHeight).
		Margin(2, 2).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorBorderNeon)
}

// GetRightPanelStyle returns a dynamically sized right panel style with better proportions
func GetRightPanelStyle(terminalWidth, terminalHeight int) lipgloss.Style {
	// Better proportions - 50% for right panel (project structure gets more space)
	panelWidth := int(float64(terminalWidth) * 0.50) - 4
	panelHeight := terminalHeight - 8

	// Ensure minimum dimensions for usability
	if panelWidth < 45 {
		panelWidth = 45
	}
	if panelHeight < 25 {
		panelHeight = 25
	}

	return BorderStyle.Copy().
		Width(panelWidth).
		Height(panelHeight).
		Margin(2, 2).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorBorderNeon)
}

// GetFullScreenStyle returns a style that fills the entire terminal
func GetFullScreenStyle(terminalWidth, terminalHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(terminalWidth).
		Height(terminalHeight)
}