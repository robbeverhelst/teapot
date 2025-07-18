// Package validation provides input validation functions for the Teapot CLI application.
// It includes validation for project names, descriptions, and other user inputs with
// comprehensive security checks and user-friendly error messages.
package validation

import (
	"fmt"
	"strings"
	"unicode"
	
	"teapot/internal/errors"
)

// ValidateProjectName validates a project name for security and formatting requirements.
// It checks for:
// - Length constraints (1-50 characters)
// - Path traversal attempts
// - Reserved system names
// - Valid character set (alphanumeric, hyphens, underscores)
func ValidateProjectName(name string) error {
	if len(name) == 0 {
		return errors.NewValidationError("project name cannot be empty", nil)
	}
	
	if len(name) > 50 {
		return errors.NewValidationError("project name too long (max 50 characters)", nil)
	}
	
	if len(name) < 2 {
		return errors.NewValidationError("project name too short (min 2 characters)", nil)
	}
	
	// Check for path traversal attempts
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return errors.NewValidationError("project name cannot contain path separators", nil)
	}
	
	// Check for reserved names (case-insensitive)
	reserved := []string{
		"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", 
		"com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", 
		"lpt5", "lpt6", "lpt7", "lpt8", "lpt9", "node_modules", "dist",
		"build", "tmp", "temp", "cache", ".git", ".svn", ".hg",
	}
	
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return errors.NewValidationError(fmt.Sprintf("project name cannot be a reserved system name: %s", r), nil)
		}
	}
	
	// Check for invalid starting/ending characters
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return errors.NewValidationError("project name cannot start or end with hyphens", nil)
	}
	
	if strings.HasPrefix(name, "_") || strings.HasSuffix(name, "_") {
		return errors.NewValidationError("project name cannot start or end with underscores", nil)
	}
	
	// Validate character set
	for i, char := range name {
		if !isValidProjectNameChar(char) {
			return errors.NewValidationError(fmt.Sprintf("invalid character '%c' at position %d", char, i+1), nil)
		}
	}
	
	// Check for consecutive special characters
	if strings.Contains(name, "--") || strings.Contains(name, "__") || strings.Contains(name, "-_") || strings.Contains(name, "_-") {
		return errors.NewValidationError("project name cannot contain consecutive special characters", nil)
	}
	
	return nil
}

// isValidProjectNameChar checks if a character is valid for project names.
// Valid characters are: letters (a-z, A-Z), digits (0-9), hyphens (-), underscores (_)
func isValidProjectNameChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-' || char == '_'
}

// ValidateProjectDescription validates a project description.
// It checks for:
// - Length constraints (max 200 characters)
// - No control characters
func ValidateProjectDescription(description string) error {
	if len(description) > 200 {
		return errors.NewValidationError("project description too long (max 200 characters)", nil)
	}
	
	// Check for control characters (except newlines and tabs)
	for i, char := range description {
		if unicode.IsControl(char) && char != '\n' && char != '\t' {
			return errors.NewValidationError(fmt.Sprintf("invalid control character at position %d", i+1), nil)
		}
	}
	
	return nil
}

// SanitizeProjectName sanitizes a project name by removing invalid characters
// and replacing them with valid alternatives.
func SanitizeProjectName(name string) string {
	// Convert to lowercase for consistency
	name = strings.ToLower(name)
	
	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	
	// Remove invalid characters
	var result strings.Builder
	for _, char := range name {
		if isValidProjectNameChar(char) {
			result.WriteRune(char)
		}
	}
	
	sanitized := result.String()
	
	// Remove consecutive hyphens/underscores
	for strings.Contains(sanitized, "--") {
		sanitized = strings.ReplaceAll(sanitized, "--", "-")
	}
	for strings.Contains(sanitized, "__") {
		sanitized = strings.ReplaceAll(sanitized, "__", "_")
	}
	
	// Trim leading/trailing special characters
	sanitized = strings.Trim(sanitized, "-_")
	
	// Ensure minimum length
	if len(sanitized) < 2 {
		sanitized = "my-project"
	}
	
	// Ensure maximum length
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
		sanitized = strings.TrimRight(sanitized, "-_")
	}
	
	return sanitized
}