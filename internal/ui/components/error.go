package components

import (
	"fmt"
	"strings"
	"time"

	"teapot/internal/errors"
	"teapot/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ErrorDisplayMode determines how errors are displayed
type ErrorDisplayMode int

const (
	// ErrorDisplayInline shows errors inline with other content
	ErrorDisplayInline ErrorDisplayMode = iota
	// ErrorDisplayModal shows errors in a modal overlay
	ErrorDisplayModal
	// ErrorDisplayBanner shows errors as a banner at the top
	ErrorDisplayBanner
)

// ErrorDisplay manages the display of error messages in the UI
type ErrorDisplay struct {
	// currentError is the error currently being displayed
	currentError *errors.TeapotError
	// displayMode determines how errors are shown
	displayMode ErrorDisplayMode
	// autoHide determines if errors should auto-hide after a timeout
	autoHide bool
	// hideTimeout specifies how long to show errors before hiding
	hideTimeout time.Duration
	// showTime tracks when the current error was first shown
	showTime time.Time
}

// NewErrorDisplay creates a new error display component
func NewErrorDisplay(mode ErrorDisplayMode, autoHide bool, timeout time.Duration) *ErrorDisplay {
	return &ErrorDisplay{
		displayMode: mode,
		autoHide:    autoHide,
		hideTimeout: timeout,
	}
}

// ShowError displays an error message
func (ed *ErrorDisplay) ShowError(err *errors.TeapotError) {
	ed.currentError = err
	ed.showTime = time.Now()
}

// ClearError clears the current error
func (ed *ErrorDisplay) ClearError() {
	ed.currentError = nil
}

// HasError returns true if there's an error to display
func (ed *ErrorDisplay) HasError() bool {
	if ed.currentError == nil {
		return false
	}
	
	// Check if auto-hide timeout has passed
	if ed.autoHide && time.Since(ed.showTime) > ed.hideTimeout {
		ed.ClearError()
		return false
	}
	
	return true
}

// Render renders the error display
func (ed *ErrorDisplay) Render(width, height int) string {
	if !ed.HasError() {
		return ""
	}
	
	switch ed.displayMode {
	case ErrorDisplayModal:
		return ed.renderModal(width, height)
	case ErrorDisplayBanner:
		return ed.renderBanner(width)
	default:
		return ed.renderInline(width)
	}
}

// renderInline renders the error inline with content
func (ed *ErrorDisplay) renderInline(width int) string {
	if ed.currentError == nil {
		return ""
	}
	
	// Choose style based on error type
	var style lipgloss.Style
	var icon string
	
	switch ed.currentError.Type {
	case errors.ErrorTypeValidation:
		style = styles.CardStyle.Copy().
			BorderForeground(styles.ColorWarning).
			Foreground(styles.ColorWarning)
		icon = "⚠️"
	case errors.ErrorTypePanic:
		style = styles.CardStyle.Copy().
			BorderForeground(styles.ColorError).
			Foreground(styles.ColorError)
		icon = "💥"
	case errors.ErrorTypeSystem:
		style = styles.CardStyle.Copy().
			BorderForeground(styles.ColorDanger).
			Foreground(styles.ColorDanger)
		icon = "🚨"
	default:
		style = styles.CardStyle.Copy().
			BorderForeground(styles.ColorError).
			Foreground(styles.ColorError)
		icon = "❌"
	}
	
	// Format error message
	title := fmt.Sprintf("%s %s Error", icon, ed.currentError.Type)
	message := ed.currentError.Message
	
	// Add recovery action if available
	var recovery string
	if ed.currentError.Recoverable && ed.currentError.RecoveryAction != "" {
		recovery = fmt.Sprintf("\n💡 %s", ed.currentError.RecoveryAction)
	}
	
	content := fmt.Sprintf("%s\n%s%s", title, message, recovery)
	
	// Apply styling and width constraints
	maxWidth := width - 4 // Account for padding
	if maxWidth < 20 {
		maxWidth = 20
	}
	
	return style.Width(maxWidth).Render(content)
}

// renderBanner renders the error as a banner at the top
func (ed *ErrorDisplay) renderBanner(width int) string {
	if ed.currentError == nil {
		return ""
	}
	
	// Choose color based on error type
	var bgColor, fgColor lipgloss.Color
	var icon string
	
	switch ed.currentError.Type {
	case errors.ErrorTypeValidation:
		bgColor = styles.ColorWarning
		fgColor = styles.ColorBgPrimary
		icon = "⚠️"
	case errors.ErrorTypePanic:
		bgColor = styles.ColorError
		fgColor = styles.ColorTextPrimary
		icon = "💥"
	case errors.ErrorTypeSystem:
		bgColor = styles.ColorDanger
		fgColor = styles.ColorTextPrimary
		icon = "🚨"
	default:
		bgColor = styles.ColorError
		fgColor = styles.ColorTextPrimary
		icon = "❌"
	}
	
	// Format banner message
	message := fmt.Sprintf("%s %s", icon, ed.currentError.Message)
	
	// Create banner style
	bannerStyle := lipgloss.NewStyle().
		Background(bgColor).
		Foreground(fgColor).
		Bold(true).
		Padding(0, 1).
		Width(width).
		Align(lipgloss.Center)
	
	return bannerStyle.Render(message)
}

// renderModal renders the error as a modal overlay
func (ed *ErrorDisplay) renderModal(width, height int) string {
	if ed.currentError == nil {
		return ""
	}
	
	// Calculate modal dimensions
	modalWidth := width / 2
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 60 {
		modalWidth = 60
	}
	
	modalHeight := height / 3
	if modalHeight < 8 {
		modalHeight = 8
	}
	if modalHeight > 15 {
		modalHeight = 15
	}
	
	// Choose style based on error type
	var borderColor lipgloss.Color
	var icon string
	
	switch ed.currentError.Type {
	case errors.ErrorTypeValidation:
		borderColor = styles.ColorWarning
		icon = "⚠️"
	case errors.ErrorTypePanic:
		borderColor = styles.ColorError
		icon = "💥"
	case errors.ErrorTypeSystem:
		borderColor = styles.ColorDanger
		icon = "🚨"
	default:
		borderColor = styles.ColorError
		icon = "❌"
	}
	
	// Format modal content
	title := fmt.Sprintf("%s %s Error", icon, ed.currentError.Type)
	message := ed.currentError.Message
	
	// Add recovery information
	var recovery string
	if ed.currentError.Recoverable {
		recovery = fmt.Sprintf("\n💡 %s", ed.currentError.RecoveryAction)
		if ed.autoHide {
			recovery += "\n⏱️ This dialog will auto-close"
		}
	}
	
	// Add dismiss instructions
	dismiss := "\n\nPress ESC to dismiss"
	
	content := fmt.Sprintf("%s\n\n%s%s%s", title, message, recovery, dismiss)
	
	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(borderColor).
		Background(styles.ColorBgSecondary).
		Foreground(styles.ColorTextPrimary).
		Padding(2, 3).
		Width(modalWidth).
		Height(modalHeight)
	
	// Create modal
	modal := modalStyle.Render(content)
	
	// Center the modal
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}

// RenderErrorSummary renders a summary of recent errors
func RenderErrorSummary(errorStats errors.ErrorStats, width int) string {
	if errorStats.TotalErrors == 0 {
		return ""
	}
	
	// Create summary header
	header := fmt.Sprintf("🔍 Error Summary (%d total)", errorStats.TotalErrors)
	
	// Create breakdown
	var breakdown []string
	for errType, count := range errorStats.TypeBreakdown {
		if count > 0 {
			breakdown = append(breakdown, fmt.Sprintf("  %s: %d", errType, count))
		}
	}
	
	// Add recovery rate
	recoveryRate := ""
	if errorStats.TotalErrors > 0 {
		rate := float64(errorStats.RecoverableCount) / float64(errorStats.TotalErrors) * 100
		recoveryRate = fmt.Sprintf("  Recovery Rate: %.1f%%", rate)
	}
	
	content := header + "\n" + strings.Join(breakdown, "\n")
	if recoveryRate != "" {
		content += "\n" + recoveryRate
	}
	
	// Style the summary
	summaryStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorTextMuted).
		Foreground(styles.ColorTextSecondary).
		Padding(1, 2).
		Width(width).
		Margin(1, 0)
	
	return summaryStyle.Render(content)
}

// RenderErrorList renders a list of recent errors
func RenderErrorList(errorLog []errors.TeapotError, maxItems int, width int) string {
	if len(errorLog) == 0 {
		return ""
	}
	
	// Limit the number of items shown
	if maxItems > len(errorLog) {
		maxItems = len(errorLog)
	}
	
	// Get recent errors (last N)
	recentErrors := errorLog[len(errorLog)-maxItems:]
	
	// Create header
	header := fmt.Sprintf("📋 Recent Errors (last %d)", maxItems)
	
	// Create error list
	var items []string
	for i, err := range recentErrors {
		timestamp := err.Timestamp.Format("15:04:05")
		status := "❌"
		if err.Recoverable {
			status = "⚠️"
		}
		
		item := fmt.Sprintf("%d. %s [%s] %s: %s", 
			i+1, status, timestamp, err.Type, err.Message)
		items = append(items, item)
	}
	
	content := header + "\n" + strings.Join(items, "\n")
	
	// Style the list
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorTextMuted).
		Foreground(styles.ColorTextSecondary).
		Padding(1, 2).
		Width(width).
		Margin(1, 0)
	
	return listStyle.Render(content)
}