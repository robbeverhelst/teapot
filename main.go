package main

import (
	"fmt"
	"os"

	"teapot/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Enable debug logging in development
	if os.Getenv("DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Printf("Could not open debug log: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen(),
	)
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running Teapot: %v\n", err)
		os.Exit(1)
	}
}