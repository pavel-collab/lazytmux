package ui

import "github.com/charmbracelet/lipgloss"

// Styles holds all the application styles
type Styles struct {
	// Panels
	FocusedPanel   lipgloss.Style
	UnfocusedPanel lipgloss.Style

	// List items
	SelectedItem lipgloss.Style
	NormalItem   lipgloss.Style
	ActiveItem   lipgloss.Style

	// Text
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	StatusBar  lipgloss.Style
	ErrorText  lipgloss.Style
	HelpText   lipgloss.Style
	DimText    lipgloss.Style

	// Dialogs
	Dialog      lipgloss.Style
	DialogTitle lipgloss.Style
	Input       lipgloss.Style
}

// DefaultStyles returns the default color scheme
func DefaultStyles() Styles {
	return Styles{
		FocusedPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1),

		UnfocusedPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true).
			Padding(0, 1),

		NormalItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),

		ActiveItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true).
			Padding(0, 1),

		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),

		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),

		ErrorText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),

		HelpText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),

		DimText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),

		Dialog: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2),

		DialogTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			MarginBottom(1),

		Input: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1),
	}
}
