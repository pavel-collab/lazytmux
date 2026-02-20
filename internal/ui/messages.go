package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"lazytmux/internal/tmux"
)

// TmuxStateMsg carries updated tmux state
type TmuxStateMsg struct {
	State tmux.TmuxState
	Err   error
}

// SessionCreatedMsg signals successful session creation
type SessionCreatedMsg struct {
	Name string
}

// SessionDeletedMsg signals successful session deletion
type SessionDeletedMsg struct {
	Name string
}

// WindowCreatedMsg signals successful window creation
type WindowCreatedMsg struct {
	SessionName string
	WindowName  string
}

// WindowDeletedMsg signals successful window deletion
type WindowDeletedMsg struct {
	SessionName string
	WindowIndex int
}

// SessionSwitchedMsg signals successful session switch
type SessionSwitchedMsg struct {
	Name string
}

// WindowSwitchedMsg signals successful window switch
type WindowSwitchedMsg struct {
	SessionName string
	WindowName  string
}

// DetachedMsg signals successful detach from session
type DetachedMsg struct{}

// ErrorMsg carries error information
type ErrorMsg struct {
	Err error
}

// StatusMsg sets status bar message
type StatusMsg struct {
	Message string
}

// ClearStatusMsg clears the status message
type ClearStatusMsg struct{}

// AttachMsg signals to attach to a session
type AttachMsg struct {
	SessionName string
}

// RefreshCmd returns a command to refresh tmux state
func RefreshCmd(client *tmux.Client) tea.Cmd {
	return func() tea.Msg {
		state, err := fetchTmuxState(client)
		return TmuxStateMsg{State: state, Err: err}
	}
}

func fetchTmuxState(client *tmux.Client) (tmux.TmuxState, error) {
	state := tmux.TmuxState{
		ServerRunning: client.IsServerRunning(),
	}

	sessions, err := client.ListSessions()
	if err != nil && err != tmux.ErrNoServer {
		return state, err
	}
	state.Sessions = sessions

	if len(sessions) > 0 {
		state.CurrentSession = &sessions[0]
		windows, err := client.ListWindows(sessions[0].Name)
		if err == nil {
			state.Windows = windows
			if len(windows) > 0 {
				state.CurrentWindow = &windows[0]
			}
		}
	}

	return state, nil
}

// CreateSessionCmd creates a new session
func CreateSessionCmd(client *tmux.Client, name string) tea.Cmd {
	return func() tea.Msg {
		err := client.CreateSession(name)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SessionCreatedMsg{Name: name}
	}
}

// DeleteSessionCmd deletes a session
func DeleteSessionCmd(client *tmux.Client, name string) tea.Cmd {
	return func() tea.Msg {
		err := client.KillSession(name)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SessionDeletedMsg{Name: name}
	}
}

// CreateWindowCmd creates a new window
func CreateWindowCmd(client *tmux.Client, sessionName, windowName string) tea.Cmd {
	return func() tea.Msg {
		err := client.CreateWindow(sessionName, windowName)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return WindowCreatedMsg{SessionName: sessionName, WindowName: windowName}
	}
}

// DeleteWindowCmd deletes a window
func DeleteWindowCmd(client *tmux.Client, sessionName string, windowIndex int) tea.Cmd {
	return func() tea.Msg {
		err := client.KillWindow(sessionName, windowIndex)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return WindowDeletedMsg{SessionName: sessionName, WindowIndex: windowIndex}
	}
}

// SwitchSessionCmd switches to another session
func SwitchSessionCmd(client *tmux.Client, name string) tea.Cmd {
	return func() tea.Msg {
		err := client.SwitchClient(name)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SessionSwitchedMsg{Name: name}
	}
}

// SwitchWindowCmd switches to another window
func SwitchWindowCmd(client *tmux.Client, sessionName string, windowIndex int, windowName string) tea.Cmd {
	return func() tea.Msg {
		err := client.SelectWindow(sessionName, windowIndex)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return WindowSwitchedMsg{SessionName: sessionName, WindowName: windowName}
	}
}

// DetachCmd detaches the current client from its session
func DetachCmd(client *tmux.Client) tea.Cmd {
	return func() tea.Msg {
		err := client.DetachClient()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return DetachedMsg{}
	}
}
