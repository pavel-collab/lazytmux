package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle dialog input first
		if m.activeDialog == CreateSessionDialog || m.activeDialog == CreateWindowDialog {
			return m.handleDialogInput(msg)
		}
		if m.activeDialog == ConfirmDeleteDialog {
			return m.handleConfirmDialog(msg)
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case TmuxStateMsg:
		if msg.Err != nil {
			m.lastError = msg.Err
		} else {
			m.tmuxState = msg.State
			m.lastError = nil
			// Reset cursors if out of bounds
			if m.sessionCursor >= len(m.tmuxState.Sessions) {
				m.sessionCursor = max(0, len(m.tmuxState.Sessions)-1)
			}
			if m.windowCursor >= len(m.tmuxState.Windows) {
				m.windowCursor = max(0, len(m.tmuxState.Windows)-1)
			}
		}
		return m, nil

	case SessionCreatedMsg:
		cmds = append(cmds, setStatusCmd("Session created: "+msg.Name))
		cmds = append(cmds, RefreshCmd(m.client))
		return m, tea.Batch(cmds...)

	case SessionDeletedMsg:
		cmds = append(cmds, setStatusCmd("Session deleted: "+msg.Name))
		cmds = append(cmds, RefreshCmd(m.client))
		return m, tea.Batch(cmds...)

	case WindowCreatedMsg:
		cmds = append(cmds, setStatusCmd("Window created: "+msg.WindowName))
		cmds = append(cmds, RefreshCmd(m.client))
		return m, tea.Batch(cmds...)

	case WindowDeletedMsg:
		cmds = append(cmds, setStatusCmd("Window deleted"))
		cmds = append(cmds, RefreshCmd(m.client))
		return m, tea.Batch(cmds...)

	case SessionSwitchedMsg:
		cmds = append(cmds, setStatusCmd("Switched to session: "+msg.Name))
		return m, tea.Batch(cmds...)

	case WindowSwitchedMsg:
		cmds = append(cmds, setStatusCmd("Switched to window: "+msg.WindowName))
		return m, tea.Batch(cmds...)

	case ErrorMsg:
		m.lastError = msg.Err
		return m, nil

	case StatusMsg:
		m.statusMessage = msg.Message
		return m, nil

	case ClearStatusMsg:
		m.statusMessage = ""
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keyMap.Help):
		m.showHelp = !m.showHelp
		return m, nil

	case key.Matches(msg, m.keyMap.Refresh):
		return m, RefreshCmd(m.client)

	case key.Matches(msg, m.keyMap.Tab), key.Matches(msg, m.keyMap.Right):
		m.focusedPanel = (m.focusedPanel + 1) % 2
		// Load windows when switching to windows panel
		if m.focusedPanel == WindowsPanel && len(m.tmuxState.Sessions) > 0 {
			return m, m.loadWindowsCmd()
		}
		return m, nil

	case key.Matches(msg, m.keyMap.ShiftTab), key.Matches(msg, m.keyMap.Left):
		if m.focusedPanel == 0 {
			m.focusedPanel = 1
		} else {
			m.focusedPanel--
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Up):
		m.moveCursorUp()
		if m.focusedPanel == SessionsPanel {
			return m, m.loadWindowsCmd()
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Down):
		m.moveCursorDown()
		if m.focusedPanel == SessionsPanel {
			return m, m.loadWindowsCmd()
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Create):
		return m.openCreateDialog()

	case key.Matches(msg, m.keyMap.Delete):
		return m.openDeleteDialog()

	case key.Matches(msg, m.keyMap.Attach):
		if m.focusedPanel == SessionsPanel && len(m.tmuxState.Sessions) > 0 {
			session := m.tmuxState.Sessions[m.sessionCursor]
			m.attachCmd = m.client.AttachSession(session.Name)
			return m, tea.Quit
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Select):
		if m.focusedPanel == SessionsPanel && len(m.tmuxState.Sessions) > 0 {
			session := m.tmuxState.Sessions[m.sessionCursor]
			return m, SwitchSessionCmd(m.client, session.Name)
		} else if m.focusedPanel == WindowsPanel && len(m.tmuxState.Windows) > 0 && m.tmuxState.CurrentSession != nil {
			window := m.tmuxState.Windows[m.windowCursor]
			return m, SwitchWindowCmd(m.client, m.tmuxState.CurrentSession.Name, window.Index, window.Name)
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Enter):
		if m.focusedPanel == SessionsPanel && len(m.tmuxState.Sessions) > 0 {
			m.focusedPanel = WindowsPanel
			return m, m.loadWindowsCmd()
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleDialogInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.activeDialog = NoDialog
		m.dialogInput.Reset()
		return m, nil

	case tea.KeyEnter:
		name := m.dialogInput.Value()
		if name == "" {
			return m, nil
		}

		var cmd tea.Cmd
		if m.activeDialog == CreateSessionDialog {
			cmd = CreateSessionCmd(m.client, name)
		} else if m.activeDialog == CreateWindowDialog && m.tmuxState.CurrentSession != nil {
			cmd = CreateWindowCmd(m.client, m.tmuxState.CurrentSession.Name, name)
		}

		m.activeDialog = NoDialog
		m.dialogInput.Reset()
		return m, cmd
	}

	var cmd tea.Cmd
	m.dialogInput, cmd = m.dialogInput.Update(msg)
	return m, cmd
}

func (m Model) handleConfirmDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Confirm):
		var cmd tea.Cmd
		if m.focusedPanel == SessionsPanel && m.deleteTarget != "" {
			cmd = DeleteSessionCmd(m.client, m.deleteTarget)
		} else if m.focusedPanel == WindowsPanel && m.tmuxState.CurrentSession != nil {
			cmd = DeleteWindowCmd(m.client, m.tmuxState.CurrentSession.Name, m.deleteWindowIdx)
		}
		m.activeDialog = NoDialog
		m.deleteTarget = ""
		return m, cmd

	case key.Matches(msg, m.keyMap.Cancel), msg.Type == tea.KeyEsc:
		m.activeDialog = NoDialog
		m.deleteTarget = ""
		return m, nil
	}
	return m, nil
}

func (m *Model) moveCursorUp() {
	if m.focusedPanel == SessionsPanel {
		if m.sessionCursor > 0 {
			m.sessionCursor--
		}
	} else {
		if m.windowCursor > 0 {
			m.windowCursor--
		}
	}
}

func (m *Model) moveCursorDown() {
	if m.focusedPanel == SessionsPanel {
		if m.sessionCursor < len(m.tmuxState.Sessions)-1 {
			m.sessionCursor++
		}
	} else {
		if m.windowCursor < len(m.tmuxState.Windows)-1 {
			m.windowCursor++
		}
	}
}

func (m Model) openCreateDialog() (tea.Model, tea.Cmd) {
	if m.focusedPanel == SessionsPanel {
		m.activeDialog = CreateSessionDialog
		m.dialogInput.Placeholder = "session name"
	} else {
		if m.tmuxState.CurrentSession == nil {
			return m, nil
		}
		m.activeDialog = CreateWindowDialog
		m.dialogInput.Placeholder = "window name"
	}
	m.dialogInput.Focus()
	return m, textinput.Blink
}

func (m Model) openDeleteDialog() (tea.Model, tea.Cmd) {
	if m.focusedPanel == SessionsPanel {
		if len(m.tmuxState.Sessions) == 0 {
			return m, nil
		}
		m.activeDialog = ConfirmDeleteDialog
		m.deleteTarget = m.tmuxState.Sessions[m.sessionCursor].Name
	} else {
		if len(m.tmuxState.Windows) == 0 {
			return m, nil
		}
		m.activeDialog = ConfirmDeleteDialog
		m.deleteWindowIdx = m.tmuxState.Windows[m.windowCursor].Index
	}
	return m, nil
}

func (m Model) loadWindowsCmd() tea.Cmd {
	if len(m.tmuxState.Sessions) == 0 {
		return nil
	}
	session := m.tmuxState.Sessions[m.sessionCursor]
	return func() tea.Msg {
		windows, err := m.client.ListWindows(session.Name)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		state := m.tmuxState
		state.CurrentSession = &session
		state.Windows = windows
		if len(windows) > 0 {
			state.CurrentWindow = &windows[0]
		} else {
			state.CurrentWindow = nil
		}
		return TmuxStateMsg{State: state}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
