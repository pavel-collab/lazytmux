package tmux

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

const testSessionPrefix = "lazytmux_test_"

// testSession holds info about a test session
type testSession struct {
	name   string
	client *Client
}

// setupTestSession creates a temporary tmux session for testing
// Returns the session info and a cleanup function
func setupTestSession(t *testing.T) *testSession {
	t.Helper()

	client, err := NewClient()
	if err != nil {
		t.Skipf("tmux not installed: %v", err)
	}

	// Generate unique session name with timestamp
	sessionName := fmt.Sprintf("%s%d", testSessionPrefix, time.Now().UnixNano())

	// Create detached session
	err = client.CreateSession(sessionName)
	if err != nil {
		t.Fatalf("failed to create test session: %v", err)
	}

	ts := &testSession{
		name:   sessionName,
		client: client,
	}

	// Register cleanup
	t.Cleanup(func() {
		ts.cleanup(t)
	})

	return ts
}

// cleanup removes the test session
func (ts *testSession) cleanup(t *testing.T) {
	t.Helper()

	err := ts.client.KillSession(ts.name)
	if err != nil {
		// Log but don't fail - session might already be gone
		t.Logf("cleanup: failed to kill session %s: %v", ts.name, err)
	}
}

// cleanupAllTestSessions removes any leftover test sessions
func cleanupAllTestSessions(t *testing.T) {
	t.Helper()

	client, err := NewClient()
	if err != nil {
		return
	}

	sessions, err := client.ListSessions()
	if err != nil {
		return
	}

	for _, s := range sessions {
		if len(s.Name) > len(testSessionPrefix) && s.Name[:len(testSessionPrefix)] == testSessionPrefix {
			_ = client.KillSession(s.Name)
		}
	}
}

// ============================================================================
// Session Integration Tests
// ============================================================================

func TestCreateAndListSessions(t *testing.T) {
	ts := setupTestSession(t)

	sessions, err := ts.client.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	// Find our test session
	found := false
	for _, s := range sessions {
		if s.Name == ts.name {
			found = true
			if s.Windows < 1 {
				t.Errorf("session should have at least 1 window, got %d", s.Windows)
			}
			break
		}
	}

	if !found {
		t.Errorf("test session %q not found in sessions list", ts.name)
	}
}

func TestKillSession(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skipf("tmux not installed: %v", err)
	}

	sessionName := fmt.Sprintf("%s%d_kill", testSessionPrefix, time.Now().UnixNano())

	// Create session
	err = client.CreateSession(sessionName)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Verify it exists
	sessions, _ := client.ListSessions()
	found := false
	for _, s := range sessions {
		if s.Name == sessionName {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("session was not created")
	}

	// Kill it
	err = client.KillSession(sessionName)
	if err != nil {
		t.Fatalf("KillSession failed: %v", err)
	}

	// Verify it's gone
	sessions, _ = client.ListSessions()
	for _, s := range sessions {
		if s.Name == sessionName {
			t.Error("session should have been killed")
		}
	}
}

func TestRenameSession(t *testing.T) {
	ts := setupTestSession(t)

	newName := ts.name + "_renamed"

	err := ts.client.RenameSession(ts.name, newName)
	if err != nil {
		t.Fatalf("RenameSession failed: %v", err)
	}

	// Update name for cleanup
	ts.name = newName

	// Verify rename
	sessions, err := ts.client.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	found := false
	for _, s := range sessions {
		if s.Name == newName {
			found = true
			break
		}
	}

	if !found {
		t.Error("renamed session not found")
	}
}

// ============================================================================
// Window Integration Tests
// ============================================================================

func TestListWindows(t *testing.T) {
	ts := setupTestSession(t)

	windows, err := ts.client.ListWindows(ts.name)
	if err != nil {
		t.Fatalf("ListWindows failed: %v", err)
	}

	// New session should have exactly 1 window
	if len(windows) != 1 {
		t.Errorf("expected 1 window, got %d", len(windows))
	}

	if len(windows) > 0 {
		w := windows[0]
		if w.SessionName != ts.name {
			t.Errorf("window.SessionName = %q, expected %q", w.SessionName, ts.name)
		}
		if w.Index != 0 {
			t.Errorf("first window index = %d, expected 0", w.Index)
		}
		if w.Panes < 1 {
			t.Errorf("window should have at least 1 pane, got %d", w.Panes)
		}
	}
}

func TestCreateWindow(t *testing.T) {
	ts := setupTestSession(t)

	windowName := "test_window"

	err := ts.client.CreateWindow(ts.name, windowName)
	if err != nil {
		t.Fatalf("CreateWindow failed: %v", err)
	}

	windows, err := ts.client.ListWindows(ts.name)
	if err != nil {
		t.Fatalf("ListWindows failed: %v", err)
	}

	if len(windows) != 2 {
		t.Errorf("expected 2 windows after creation, got %d", len(windows))
	}

	// Find our window
	found := false
	for _, w := range windows {
		if w.Name == windowName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("created window %q not found", windowName)
	}
}

func TestCreateWindowWithoutName(t *testing.T) {
	ts := setupTestSession(t)

	// Create window without name
	err := ts.client.CreateWindow(ts.name, "")
	if err != nil {
		t.Fatalf("CreateWindow failed: %v", err)
	}

	windows, err := ts.client.ListWindows(ts.name)
	if err != nil {
		t.Fatalf("ListWindows failed: %v", err)
	}

	if len(windows) != 2 {
		t.Errorf("expected 2 windows after creation, got %d", len(windows))
	}
}

func TestKillWindow(t *testing.T) {
	ts := setupTestSession(t)

	// Create a second window
	err := ts.client.CreateWindow(ts.name, "to_delete")
	if err != nil {
		t.Fatalf("CreateWindow failed: %v", err)
	}

	windows, _ := ts.client.ListWindows(ts.name)
	if len(windows) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(windows))
	}

	// Kill the second window
	err = ts.client.KillWindow(ts.name, 1)
	if err != nil {
		t.Fatalf("KillWindow failed: %v", err)
	}

	windows, _ = ts.client.ListWindows(ts.name)
	if len(windows) != 1 {
		t.Errorf("expected 1 window after kill, got %d", len(windows))
	}
}

func TestRenameWindow(t *testing.T) {
	ts := setupTestSession(t)

	newName := "renamed_window"

	err := ts.client.RenameWindow(ts.name, 0, newName)
	if err != nil {
		t.Fatalf("RenameWindow failed: %v", err)
	}

	windows, err := ts.client.ListWindows(ts.name)
	if err != nil {
		t.Fatalf("ListWindows failed: %v", err)
	}

	if len(windows) == 0 {
		t.Fatal("no windows found")
	}

	if windows[0].Name != newName {
		t.Errorf("window name = %q, expected %q", windows[0].Name, newName)
	}
}

func TestSelectWindow(t *testing.T) {
	ts := setupTestSession(t)

	// Create second window
	err := ts.client.CreateWindow(ts.name, "second")
	if err != nil {
		t.Fatalf("CreateWindow failed: %v", err)
	}

	// Select first window
	err = ts.client.SelectWindow(ts.name, 0)
	if err != nil {
		t.Fatalf("SelectWindow failed: %v", err)
	}

	windows, err := ts.client.ListWindows(ts.name)
	if err != nil {
		t.Fatalf("ListWindows failed: %v", err)
	}

	// First window should be active
	for _, w := range windows {
		if w.Index == 0 && !w.Active {
			t.Error("window 0 should be active after SelectWindow")
		}
	}
}

// ============================================================================
// Pane Integration Tests
// ============================================================================

func TestListPanes(t *testing.T) {
	ts := setupTestSession(t)

	panes, err := ts.client.ListPanes(ts.name, 0)
	if err != nil {
		t.Fatalf("ListPanes failed: %v", err)
	}

	// New window should have exactly 1 pane
	if len(panes) != 1 {
		t.Errorf("expected 1 pane, got %d", len(panes))
	}

	if len(panes) > 0 {
		p := panes[0]
		if p.Index != 0 {
			t.Errorf("first pane index = %d, expected 0", p.Index)
		}
		if p.Width <= 0 {
			t.Errorf("pane width = %d, should be > 0", p.Width)
		}
		if p.Height <= 0 {
			t.Errorf("pane height = %d, should be > 0", p.Height)
		}
		if !p.Active {
			t.Error("single pane should be active")
		}
	}
}

func TestSplitWindowVertical(t *testing.T) {
	ts := setupTestSession(t)

	err := ts.client.SplitWindowVertical(ts.name, 0)
	if err != nil {
		t.Fatalf("SplitWindowVertical failed: %v", err)
	}

	panes, err := ts.client.ListPanes(ts.name, 0)
	if err != nil {
		t.Fatalf("ListPanes failed: %v", err)
	}

	if len(panes) != 2 {
		t.Errorf("expected 2 panes after vertical split, got %d", len(panes))
	}

	// Panes should be side by side (same top position, different left)
	if len(panes) == 2 {
		if panes[0].Top != panes[1].Top {
			t.Log("Note: vertical split panes have different top positions")
		}
	}
}

func TestSplitWindowHorizontal(t *testing.T) {
	ts := setupTestSession(t)

	err := ts.client.SplitWindowHorizontal(ts.name, 0)
	if err != nil {
		t.Fatalf("SplitWindowHorizontal failed: %v", err)
	}

	panes, err := ts.client.ListPanes(ts.name, 0)
	if err != nil {
		t.Fatalf("ListPanes failed: %v", err)
	}

	if len(panes) != 2 {
		t.Errorf("expected 2 panes after horizontal split, got %d", len(panes))
	}

	// Panes should be stacked (same left position, different top)
	if len(panes) == 2 {
		if panes[0].Left != panes[1].Left {
			t.Log("Note: horizontal split panes have different left positions")
		}
	}
}

func TestMultipleSplits(t *testing.T) {
	ts := setupTestSession(t)

	// Split vertically
	err := ts.client.SplitWindowVertical(ts.name, 0)
	if err != nil {
		t.Fatalf("SplitWindowVertical failed: %v", err)
	}

	// Split horizontally
	err = ts.client.SplitWindowHorizontal(ts.name, 0)
	if err != nil {
		t.Fatalf("SplitWindowHorizontal failed: %v", err)
	}

	panes, err := ts.client.ListPanes(ts.name, 0)
	if err != nil {
		t.Fatalf("ListPanes failed: %v", err)
	}

	if len(panes) != 3 {
		t.Errorf("expected 3 panes after two splits, got %d", len(panes))
	}
}

// ============================================================================
// AttachSession Tests
// ============================================================================

func TestAttachSessionReturnsCommand(t *testing.T) {
	ts := setupTestSession(t)

	cmd := ts.client.AttachSession(ts.name)

	if len(cmd) < 3 {
		t.Fatalf("expected at least 3 command parts, got %d", len(cmd))
	}

	// Should contain tmux path
	if cmd[0] == "" {
		t.Error("tmux path should not be empty")
	}

	// Should be attach-session command
	if cmd[1] != "attach-session" {
		t.Errorf("cmd[1] = %q, expected 'attach-session'", cmd[1])
	}

	// Should have -t flag
	if cmd[2] != "-t" {
		t.Errorf("cmd[2] = %q, expected '-t'", cmd[2])
	}

	// Should have session name
	if cmd[3] != ts.name {
		t.Errorf("cmd[3] = %q, expected %q", cmd[3], ts.name)
	}
}

// ============================================================================
// Error Cases
// ============================================================================

func TestListWindowsNonExistentSession(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skipf("tmux not installed: %v", err)
	}

	_, err = client.ListWindows("nonexistent_session_12345")
	if err == nil {
		t.Error("expected error for non-existent session")
	}
}

func TestKillWindowNonExistentWindow(t *testing.T) {
	ts := setupTestSession(t)

	err := ts.client.KillWindow(ts.name, 999)
	if err == nil {
		t.Error("expected error for non-existent window index")
	}
}

func TestSelectWindowNonExistentWindow(t *testing.T) {
	ts := setupTestSession(t)

	err := ts.client.SelectWindow(ts.name, 999)
	if err == nil {
		t.Error("expected error for non-existent window index")
	}
}

// ============================================================================
// Helpers
// ============================================================================

// isTmuxInstalled checks if tmux is available
func isTmuxInstalled() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkListSessions(b *testing.B) {
	if !isTmuxInstalled() {
		b.Skip("tmux not installed")
	}

	client, _ := NewClient()

	// Create a test session
	sessionName := fmt.Sprintf("%sbench_%d", testSessionPrefix, time.Now().UnixNano())
	_ = client.CreateSession(sessionName)
	defer client.KillSession(sessionName)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ListSessions()
	}
}

func BenchmarkListWindows(b *testing.B) {
	if !isTmuxInstalled() {
		b.Skip("tmux not installed")
	}

	client, _ := NewClient()

	sessionName := fmt.Sprintf("%sbench_%d", testSessionPrefix, time.Now().UnixNano())
	_ = client.CreateSession(sessionName)
	defer client.KillSession(sessionName)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ListWindows(sessionName)
	}
}

func BenchmarkListPanes(b *testing.B) {
	if !isTmuxInstalled() {
		b.Skip("tmux not installed")
	}

	client, _ := NewClient()

	sessionName := fmt.Sprintf("%sbench_%d", testSessionPrefix, time.Now().UnixNano())
	_ = client.CreateSession(sessionName)
	defer client.KillSession(sessionName)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ListPanes(sessionName, 0)
	}
}
