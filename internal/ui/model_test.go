package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"lazytmux/internal/tmux"
)

// Mock client for testing - we can't create real client without tmux
func createTestModel() Model {
	// Create model without client for basic state testing
	m := Model{
		focusedPanel: SessionsPanel,
		help:         help.New(),
		keyMap:       DefaultKeyMap(),
		styles:       DefaultStyles(),
		configEditor: NewConfigEditorModel(),
	}
	return m
}

func TestFocusedPanelConstants(t *testing.T) {
	if SessionsPanel != 0 {
		t.Error("SessionsPanel should be 0")
	}
	if WindowsPanel != 1 {
		t.Error("WindowsPanel should be 1")
	}
}

func TestDialogTypeConstants(t *testing.T) {
	if NoDialog != 0 {
		t.Error("NoDialog should be 0")
	}
	if ConfirmDeleteDialog != 1 {
		t.Error("ConfirmDeleteDialog should be 1")
	}
	if CreateSessionDialog != 2 {
		t.Error("CreateSessionDialog should be 2")
	}
	if CreateWindowDialog != 3 {
		t.Error("CreateWindowDialog should be 3")
	}
}

func TestModelInitialState(t *testing.T) {
	m := createTestModel()

	// Check default panel focus
	if m.focusedPanel != SessionsPanel {
		t.Errorf("focusedPanel = %v, expected SessionsPanel", m.focusedPanel)
	}

	// Check no dialog active
	if m.activeDialog != NoDialog {
		t.Errorf("activeDialog = %v, expected NoDialog", m.activeDialog)
	}

	// Check cursors at start
	if m.sessionCursor != 0 {
		t.Errorf("sessionCursor = %d, expected 0", m.sessionCursor)
	}
	if m.windowCursor != 0 {
		t.Errorf("windowCursor = %d, expected 0", m.windowCursor)
	}

	// Check config editor not active
	if m.configEditorActive {
		t.Error("configEditorActive should be false initially")
	}
}

func TestModelGetAttachCmd(t *testing.T) {
	m := createTestModel()

	// Initially should be nil
	if m.GetAttachCmd() != nil {
		t.Error("GetAttachCmd() should be nil initially")
	}

	// Set attach cmd
	m.attachCmd = []string{"tmux", "attach-session", "-t", "test"}
	cmd := m.GetAttachCmd()

	if len(cmd) != 4 {
		t.Errorf("GetAttachCmd() length = %d, expected 4", len(cmd))
	}
	if cmd[0] != "tmux" {
		t.Errorf("cmd[0] = %q, expected 'tmux'", cmd[0])
	}
}

func TestMoveCursorUp(t *testing.T) {
	m := createTestModel()
	m.tmuxState = tmux.TmuxState{
		Sessions: []tmux.Session{
			{Name: "session1"},
			{Name: "session2"},
			{Name: "session3"},
		},
	}
	m.sessionCursor = 2

	// Move cursor up in sessions panel
	m.focusedPanel = SessionsPanel
	m.moveCursorUp()
	if m.sessionCursor != 1 {
		t.Errorf("sessionCursor = %d, expected 1", m.sessionCursor)
	}

	m.moveCursorUp()
	if m.sessionCursor != 0 {
		t.Errorf("sessionCursor = %d, expected 0", m.sessionCursor)
	}

	// Should not go below 0
	m.moveCursorUp()
	if m.sessionCursor != 0 {
		t.Errorf("sessionCursor = %d, expected 0 (no underflow)", m.sessionCursor)
	}
}

func TestMoveCursorDown(t *testing.T) {
	m := createTestModel()
	m.tmuxState = tmux.TmuxState{
		Sessions: []tmux.Session{
			{Name: "session1"},
			{Name: "session2"},
			{Name: "session3"},
		},
	}
	m.sessionCursor = 0

	// Move cursor down in sessions panel
	m.focusedPanel = SessionsPanel
	m.moveCursorDown()
	if m.sessionCursor != 1 {
		t.Errorf("sessionCursor = %d, expected 1", m.sessionCursor)
	}

	m.moveCursorDown()
	if m.sessionCursor != 2 {
		t.Errorf("sessionCursor = %d, expected 2", m.sessionCursor)
	}

	// Should not go beyond last item
	m.moveCursorDown()
	if m.sessionCursor != 2 {
		t.Errorf("sessionCursor = %d, expected 2 (no overflow)", m.sessionCursor)
	}
}

func TestMoveCursorInWindowsPanel(t *testing.T) {
	m := createTestModel()
	m.tmuxState = tmux.TmuxState{
		Windows: []tmux.Window{
			{Name: "window1"},
			{Name: "window2"},
		},
	}
	m.focusedPanel = WindowsPanel
	m.windowCursor = 0

	m.moveCursorDown()
	if m.windowCursor != 1 {
		t.Errorf("windowCursor = %d, expected 1", m.windowCursor)
	}

	m.moveCursorUp()
	if m.windowCursor != 0 {
		t.Errorf("windowCursor = %d, expected 0", m.windowCursor)
	}
}

func TestUpdateWithWindowSizeMsg(t *testing.T) {
	m := createTestModel()

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.width != 120 {
		t.Errorf("width = %d, expected 120", updated.width)
	}
	if updated.height != 40 {
		t.Errorf("height = %d, expected 40", updated.height)
	}
}

func TestUpdateWithTmuxStateMsg(t *testing.T) {
	m := createTestModel()

	state := tmux.TmuxState{
		ServerRunning: true,
		Sessions: []tmux.Session{
			{Name: "test-session", Windows: 2, Attached: true},
		},
	}

	msg := TmuxStateMsg{State: state, Err: nil}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if len(updated.tmuxState.Sessions) != 1 {
		t.Errorf("sessions count = %d, expected 1", len(updated.tmuxState.Sessions))
	}
	if updated.tmuxState.Sessions[0].Name != "test-session" {
		t.Errorf("session name = %q, expected 'test-session'", updated.tmuxState.Sessions[0].Name)
	}
	if updated.lastError != nil {
		t.Errorf("lastError should be nil, got %v", updated.lastError)
	}
}

func TestUpdateWithTmuxStateMsgError(t *testing.T) {
	m := createTestModel()

	msg := TmuxStateMsg{Err: tmux.ErrNoServer}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.lastError != tmux.ErrNoServer {
		t.Errorf("lastError = %v, expected ErrNoServer", updated.lastError)
	}
}

func TestUpdateWithStatusMsg(t *testing.T) {
	m := createTestModel()

	msg := StatusMsg{Message: "Test status message"}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.statusMessage != "Test status message" {
		t.Errorf("statusMessage = %q, expected 'Test status message'", updated.statusMessage)
	}
}

func TestUpdateWithClearStatusMsg(t *testing.T) {
	m := createTestModel()
	m.statusMessage = "Some message"

	msg := ClearStatusMsg{}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.statusMessage != "" {
		t.Errorf("statusMessage = %q, expected empty", updated.statusMessage)
	}
}

func TestUpdateWithErrorMsg(t *testing.T) {
	m := createTestModel()

	testErr := tmux.ErrNoServer
	msg := ErrorMsg{Err: testErr}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.lastError != testErr {
		t.Errorf("lastError = %v, expected %v", updated.lastError, testErr)
	}
}

func TestUpdateWithPanesLoadedMsg(t *testing.T) {
	m := createTestModel()

	panes := []tmux.Pane{
		{ID: "%1", Index: 0, Width: 80, Height: 24},
		{ID: "%2", Index: 1, Width: 80, Height: 24},
	}

	msg := PanesLoadedMsg{
		SessionName: "test",
		WindowIndex: 0,
		Panes:       panes,
	}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if len(updated.currentPanes) != 2 {
		t.Errorf("currentPanes length = %d, expected 2", len(updated.currentPanes))
	}
}

func TestCursorBoundsAfterStateChange(t *testing.T) {
	m := createTestModel()
	m.sessionCursor = 5 // Set cursor beyond what will be available

	state := tmux.TmuxState{
		Sessions: []tmux.Session{
			{Name: "session1"},
			{Name: "session2"},
		},
	}

	msg := TmuxStateMsg{State: state}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	// Cursor should be clamped to valid range
	if updated.sessionCursor >= len(state.Sessions) {
		t.Errorf("sessionCursor = %d, should be < %d", updated.sessionCursor, len(state.Sessions))
	}
}

func TestMaxHelper(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{1, 2, 2},
		{2, 1, 2},
		{0, 0, 0},
		{-1, 0, 0},
		{-5, -3, -3},
	}

	for _, tt := range tests {
		result := max(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("max(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

// Message type tests

func TestTmuxStateMsgStruct(t *testing.T) {
	msg := TmuxStateMsg{
		State: tmux.TmuxState{ServerRunning: true},
		Err:   nil,
	}

	if !msg.State.ServerRunning {
		t.Error("ServerRunning should be true")
	}
	if msg.Err != nil {
		t.Error("Err should be nil")
	}
}

func TestSessionCreatedMsgStruct(t *testing.T) {
	msg := SessionCreatedMsg{Name: "test-session"}
	if msg.Name != "test-session" {
		t.Errorf("Name = %q, expected 'test-session'", msg.Name)
	}
}

func TestSessionDeletedMsgStruct(t *testing.T) {
	msg := SessionDeletedMsg{Name: "deleted-session"}
	if msg.Name != "deleted-session" {
		t.Errorf("Name = %q, expected 'deleted-session'", msg.Name)
	}
}

func TestWindowCreatedMsgStruct(t *testing.T) {
	msg := WindowCreatedMsg{
		SessionName: "session1",
		WindowName:  "window1",
	}
	if msg.SessionName != "session1" {
		t.Errorf("SessionName = %q, expected 'session1'", msg.SessionName)
	}
	if msg.WindowName != "window1" {
		t.Errorf("WindowName = %q, expected 'window1'", msg.WindowName)
	}
}

func TestWindowDeletedMsgStruct(t *testing.T) {
	msg := WindowDeletedMsg{
		SessionName: "session1",
		WindowIndex: 2,
	}
	if msg.SessionName != "session1" {
		t.Errorf("SessionName = %q, expected 'session1'", msg.SessionName)
	}
	if msg.WindowIndex != 2 {
		t.Errorf("WindowIndex = %d, expected 2", msg.WindowIndex)
	}
}

func TestSessionSwitchedMsgStruct(t *testing.T) {
	msg := SessionSwitchedMsg{Name: "new-session"}
	if msg.Name != "new-session" {
		t.Errorf("Name = %q, expected 'new-session'", msg.Name)
	}
}

func TestWindowSwitchedMsgStruct(t *testing.T) {
	msg := WindowSwitchedMsg{
		SessionName: "session1",
		WindowName:  "window1",
	}
	if msg.SessionName != "session1" {
		t.Errorf("SessionName = %q, expected 'session1'", msg.SessionName)
	}
	if msg.WindowName != "window1" {
		t.Errorf("WindowName = %q, expected 'window1'", msg.WindowName)
	}
}

func TestDetachedMsgStruct(t *testing.T) {
	// Just verify the type exists and can be created
	msg := DetachedMsg{}
	_ = msg
}

func TestErrorMsgStruct(t *testing.T) {
	testErr := tmux.ErrNoServer
	msg := ErrorMsg{Err: testErr}
	if msg.Err != testErr {
		t.Errorf("Err = %v, expected %v", msg.Err, testErr)
	}
}

func TestStatusMsgStruct(t *testing.T) {
	msg := StatusMsg{Message: "Status message"}
	if msg.Message != "Status message" {
		t.Errorf("Message = %q, expected 'Status message'", msg.Message)
	}
}

func TestClearStatusMsgStruct(t *testing.T) {
	msg := ClearStatusMsg{}
	_ = msg
}

func TestAttachMsgStruct(t *testing.T) {
	msg := AttachMsg{SessionName: "attach-session"}
	if msg.SessionName != "attach-session" {
		t.Errorf("SessionName = %q, expected 'attach-session'", msg.SessionName)
	}
}

func TestPaneSplitMsgStruct(t *testing.T) {
	msg := PaneSplitMsg{
		SessionName: "session1",
		WindowIndex: 0,
		Vertical:    true,
	}
	if msg.SessionName != "session1" {
		t.Errorf("SessionName = %q, expected 'session1'", msg.SessionName)
	}
	if msg.WindowIndex != 0 {
		t.Errorf("WindowIndex = %d, expected 0", msg.WindowIndex)
	}
	if !msg.Vertical {
		t.Error("Vertical should be true")
	}
}

func TestPanesLoadedMsgStruct(t *testing.T) {
	panes := []tmux.Pane{{ID: "%1"}}
	msg := PanesLoadedMsg{
		SessionName: "session1",
		WindowIndex: 0,
		Panes:       panes,
		Err:         nil,
	}
	if msg.SessionName != "session1" {
		t.Errorf("SessionName = %q, expected 'session1'", msg.SessionName)
	}
	if len(msg.Panes) != 1 {
		t.Errorf("Panes length = %d, expected 1", len(msg.Panes))
	}
}

func TestTmuxStateStruct(t *testing.T) {
	session := tmux.Session{Name: "test", Windows: 2}
	window := tmux.Window{ID: "@1", Name: "main"}

	state := tmux.TmuxState{
		Sessions:       []tmux.Session{session},
		CurrentSession: &session,
		Windows:        []tmux.Window{window},
		CurrentWindow:  &window,
		ServerRunning:  true,
	}

	if len(state.Sessions) != 1 {
		t.Error("Sessions should have 1 element")
	}
	if state.CurrentSession == nil {
		t.Error("CurrentSession should not be nil")
	}
	if len(state.Windows) != 1 {
		t.Error("Windows should have 1 element")
	}
	if state.CurrentWindow == nil {
		t.Error("CurrentWindow should not be nil")
	}
	if !state.ServerRunning {
		t.Error("ServerRunning should be true")
	}
}

func TestSessionStruct(t *testing.T) {
	created := time.Now()
	session := tmux.Session{
		Name:     "test-session",
		ID:       "session-id",
		Windows:  5,
		Created:  created,
		Attached: true,
	}

	if session.Name != "test-session" {
		t.Errorf("Name = %q, expected 'test-session'", session.Name)
	}
	if session.Windows != 5 {
		t.Errorf("Windows = %d, expected 5", session.Windows)
	}
	if !session.Attached {
		t.Error("Attached should be true")
	}
}

func TestWindowStruct(t *testing.T) {
	window := tmux.Window{
		ID:          "@1",
		Index:       0,
		Name:        "editor",
		SessionName: "main",
		Active:      true,
		Panes:       3,
		PaneList:    []tmux.Pane{{ID: "%1"}},
	}

	if window.ID != "@1" {
		t.Errorf("ID = %q, expected '@1'", window.ID)
	}
	if window.Name != "editor" {
		t.Errorf("Name = %q, expected 'editor'", window.Name)
	}
	if !window.Active {
		t.Error("Active should be true")
	}
	if window.Panes != 3 {
		t.Errorf("Panes = %d, expected 3", window.Panes)
	}
}

func TestPaneStruct(t *testing.T) {
	pane := tmux.Pane{
		ID:     "%1",
		Index:  0,
		Width:  80,
		Height: 24,
		Left:   0,
		Top:    0,
		Active: true,
	}

	if pane.ID != "%1" {
		t.Errorf("ID = %q, expected '%%1'", pane.ID)
	}
	if pane.Width != 80 {
		t.Errorf("Width = %d, expected 80", pane.Width)
	}
	if pane.Height != 24 {
		t.Errorf("Height = %d, expected 24", pane.Height)
	}
	if !pane.Active {
		t.Error("Active should be true")
	}
}
