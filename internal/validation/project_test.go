package validation

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		expectedError bool
		errorContains string
	}{
		// Valid names
		{"valid simple name", "my-project", false, ""},
		{"valid with underscore", "my_project", false, ""},
		{"valid mixed case", "MyProject", false, ""},
		{"valid with numbers", "project123", false, ""},
		{"valid mixed", "my-project_v2", false, ""},
		
		// Invalid names - empty/length
		{"empty name", "", true, "cannot be empty"},
		{"too short", "a", true, "too short"},
		{"too long", "this-is-a-very-long-project-name-that-exceeds-the-maximum-allowed-length-of-fifty-characters", true, "too long"},
		
		// Invalid names - path traversal
		{"path traversal dots", "../project", true, "path separators"},
		{"path traversal slash", "project/test", true, "path separators"},
		{"path traversal backslash", "project\\test", true, "path separators"},
		
		// Invalid names - reserved
		{"reserved con", "con", true, "reserved system name"},
		{"reserved Con", "Con", true, "reserved system name"},
		{"reserved node_modules", "node_modules", true, "reserved system name"},
		{"reserved .git", ".git", true, "reserved system name"},
		
		// Invalid names - invalid characters
		{"invalid space", "my project", true, "invalid character"},
		{"invalid dot", "my.project", true, "invalid character"},
		{"invalid special", "my@project", true, "invalid character"},
		
		// Invalid names - start/end with special chars
		{"starts with hyphen", "-project", true, "cannot start"},
		{"ends with hyphen", "project-", true, "cannot start or end"},
		{"starts with underscore", "_project", true, "cannot start"},
		{"ends with underscore", "project_", true, "cannot start or end"},
		
		// Invalid names - consecutive special chars
		{"consecutive hyphens", "my--project", true, "consecutive special"},
		{"consecutive underscores", "my__project", true, "consecutive special"},
		{"mixed consecutive", "my-_project", true, "consecutive special"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.projectName)
			
			if tt.expectedError && err == nil {
				t.Errorf("Expected error for project name '%s', but got nil", tt.projectName)
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error for project name '%s', but got: %v", tt.projectName, err)
			}
			
			if tt.expectedError && err != nil && tt.errorContains != "" {
				if !containsSubstring(err.Error(), tt.errorContains) {
					t.Errorf("Expected error message to contain '%s', but got: %v", tt.errorContains, err)
				}
			}
		})
	}
}

func TestValidateProjectDescription(t *testing.T) {
	tests := []struct {
		name          string
		description   string
		expectedError bool
		errorContains string
	}{
		// Valid descriptions
		{"empty description", "", false, ""},
		{"short description", "A simple project", false, ""},
		{"long description", "This is a very long description that contains many words and should still be valid as long as it doesn't exceed the maximum character limit of two hundred characters.", false, ""},
		{"description with newlines", "Line 1\nLine 2", false, ""},
		{"description with tabs", "Text\twith\ttabs", false, ""},
		
		// Invalid descriptions
		{"too long", "This is an extremely long description that exceeds the maximum allowed length of two hundred characters. It just keeps going and going and going until it becomes way too long for a reasonable project description and should be rejected by the validation function.", true, "too long"},
		{"control character", "Text with \x00 control", true, "control character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectDescription(tt.description)
			
			if tt.expectedError && err == nil {
				t.Errorf("Expected error for description '%s', but got nil", tt.description)
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error for description '%s', but got: %v", tt.description, err)
			}
			
			if tt.expectedError && err != nil && tt.errorContains != "" {
				if !containsSubstring(err.Error(), tt.errorContains) {
					t.Errorf("Expected error message to contain '%s', but got: %v", tt.errorContains, err)
				}
			}
		})
	}
}

func TestSanitizeProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "MyProject", "myproject"},
		{"with spaces", "My Project", "my-project"},
		{"with invalid chars", "My@Project#2024", "myproject2024"},
		{"consecutive hyphens", "my--project", "my-project"},
		{"consecutive underscores", "my__project", "my_project"},
		{"leading/trailing special", "-_project_-", "project"},
		{"too short result", "a@", "my-project"},
		{"too long result", "this-is-a-very-long-project-name-that-exceeds-the-maximum-allowed-length-of-fifty-characters", "this-is-a-very-long-project-name-that-exceeds-the"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', but got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsValidProjectNameChar(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		{"lowercase letter", 'a', true},
		{"uppercase letter", 'A', true},
		{"digit", '5', true},
		{"hyphen", '-', true},
		{"underscore", '_', true},
		{"space", ' ', false},
		{"dot", '.', false},
		{"special char", '@', false},
		{"unicode", 'é', true}, // Unicode letters should be valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidProjectNameChar(tt.char)
			if result != tt.expected {
				t.Errorf("Expected %v for char '%c', but got %v", tt.expected, tt.char, result)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}