package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"lazytmux/internal/tmux"
)

// FocusedPanel indicates which panel has keyboard focus
type FocusedPanel int

const (
	SessionsPanel FocusedPanel = iota
	WindowsPanel
)

// DialogType represents active dialog (if any)
type DialogType int

const (
	NoDialog DialogType = iota
	ConfirmDeleteDialog
	CreateSessionDialog
	CreateWindowDialog
)

// Model is the root Bubble Tea model
type Model struct {
	// Tmux state
	tmuxState      tmux.TmuxState
	sessionCursor  int
	windowCursor   int
	focusedPanel   FocusedPanel

	// Dialog state
	activeDialog   DialogType
	dialogInput    textinput.Model
	deleteTarget   string
	deleteWindowIdx int

	// UI components
	help       help.Model
	keyMap     KeyMap
	styles     Styles
	showHelp   bool

	// Dimensions
	width  int
	height int

	// Status
	lastError     error
	statusMessage string

	// Tmux client
	client *tmux.Client

	// For attach
	attachCmd []string
}

// NewModel creates a new Model
func NewModel(client *tmux.Client) Model {
	ti := textinput.New()
	ti.Placeholder = "name"
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		client:       client,
		focusedPanel: SessionsPanel,
		help:         help.New(),
		keyMap:       DefaultKeyMap(),
		styles:       DefaultStyles(),
		dialogInput:  ti,
	}
}

// Init returns the initial command
func (m Model) Init() tea.Cmd {
	return RefreshCmd(m.client)
}

// GetAttachCmd returns attach command if set
func (m Model) GetAttachCmd() []string {
	return m.attachCmd
}

// setStatus sets a status message that clears after delay
func setStatusCmd(msg string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StatusMsg{Message: msg} },
		tea.Tick(3*time.Second, func(time.Time) tea.Msg { return ClearStatusMsg{} }),
	)
}
