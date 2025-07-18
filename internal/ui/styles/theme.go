package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Modern professional color palette
	ColorPrimary     = lipgloss.Color("#6366F1")    // Indigo - main brand
	ColorSecondary   = lipgloss.Color("#8B5CF6")    // Purple - secondary actions
	ColorAccent      = lipgloss.Color("#06B6D4")    // Cyan - highlights
	ColorSuccess     = lipgloss.Color("#10B981")    // Green - success states
	ColorWarning     = lipgloss.Color("#F59E0B")    // Amber - warnings
	ColorError       = lipgloss.Color("#EF4444")    // Red - errors
	ColorDanger      = lipgloss.Color("#DC2626")    // Dark red - critical
	
	// Text colors with better contrast
	ColorTextPrimary   = lipgloss.Color("#F8FAFC")  // Almost white
	ColorTextSecondary = lipgloss.Color("#E2E8F0")  // Light gray
	ColorTextMuted     = lipgloss.Color("#94A3B8")  // Muted gray
	ColorTextDisabled  = lipgloss.Color("#64748B")  // Disabled
	
	// Background colors for depth
	ColorBgPrimary   = lipgloss.Color("#0F172A")    // Dark navy
	ColorBgSecondary = lipgloss.Color("#1E293B")    // Lighter navy
	ColorBgTertiary  = lipgloss.Color("#334155")    // Medium gray
	ColorBgSurface   = lipgloss.Color("#475569")    // Light surface
	
	// Border colors
	ColorBorderPrimary   = lipgloss.Color("#475569") // Subtle border
	ColorBorderSecondary = lipgloss.Color("#64748B") // Visible border
	ColorBorderAccent    = lipgloss.Color("#6366F1") // Highlighted border
	
	// Legacy aliases for backward compatibility
	ColorText       = ColorTextPrimary
	ColorMuted      = ColorTextMuted
	ColorBorder     = ColorBorderPrimary
	ColorBackground = ColorBgPrimary
)

var (
	// Typography hierarchy with better visual weight
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Margin(1, 0)

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
			Padding(0, 1)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorTextSecondary)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			Padding(0, 1)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	// Enhanced interactive styles with modern polish
	HoverStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1)

	FocusedStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1)

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
		Margin(1, 1)
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
		Margin(1, 1)
}

// GetFullScreenStyle returns a style that fills the entire terminal
func GetFullScreenStyle(terminalWidth, terminalHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(terminalWidth).
		Height(terminalHeight)
}