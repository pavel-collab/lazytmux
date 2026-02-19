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

// Window represents a tmux window within a session
type Window struct {
	ID          string
	Index       int
	Name        string
	SessionName string
	Active      bool
	Panes       int
}

// TmuxState holds the complete current state
type TmuxState struct {
	Sessions       []Session
	CurrentSession *Session
	Windows        []Window
	CurrentWindow  *Window
	ServerRunning  bool
}
