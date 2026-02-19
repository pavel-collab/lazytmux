package main

import (
	"fmt"
	"os"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"lazytmux/internal/tmux"
	"lazytmux/internal/ui"
)

func main() {
	// Initialize tmux client
	client, err := tmux.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure tmux is installed and available in PATH")
		os.Exit(1)
	}

	// Create and run the Bubble Tea program
	model := ui.NewModel(client)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	// Check if we need to attach to a session
	if m, ok := finalModel.(ui.Model); ok {
		if attachCmd := m.GetAttachCmd(); len(attachCmd) > 0 {
			// Execute tmux attach, replacing this process
			err := syscall.Exec(attachCmd[0], attachCmd, os.Environ())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error attaching to session: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
