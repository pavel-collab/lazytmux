package tmux

import "time"

// Session represents a tmux session
type Session struct {
	Name     string
	ID       string
	Windows  int
	Created  time.Time
	Attached bool
}

// Pane represents a tmux pane within a window
type Pane struct {
	ID     string
	Index  int
	Width  int
	Height int
	Left   int  // X position
	Top    int  // Y position
	Active bool
}

// Window represents a tmux window within a session
type Window struct {
	ID          string
	Index       int
	Name        string
	SessionName string
	Active      bool
	Panes       int
	PaneList    []Pane // detailed pane information
}

// TmuxState holds the complete current state
type TmuxState struct {
	Sessions       []Session
	CurrentSession *Session
	Windows        []Window
	CurrentWindow  *Window
	ServerRunning  bool
}
