package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	// Verify all bindings are set
	if !km.Up.Enabled() {
		t.Error("Up binding should be enabled")
	}
	if !km.Down.Enabled() {
		t.Error("Down binding should be enabled")
	}
	if !km.Left.Enabled() {
		t.Error("Left binding should be enabled")
	}
	if !km.Right.Enabled() {
		t.Error("Right binding should be enabled")
	}
	if !km.Tab.Enabled() {
		t.Error("Tab binding should be enabled")
	}
	if !km.ShiftTab.Enabled() {
		t.Error("ShiftTab binding should be enabled")
	}
	if !km.Enter.Enabled() {
		t.Error("Enter binding should be enabled")
	}
	if !km.Select.Enabled() {
		t.Error("Select binding should be enabled")
	}
	if !km.Create.Enabled() {
		t.Error("Create binding should be enabled")
	}
	if !km.Delete.Enabled() {
		t.Error("Delete binding should be enabled")
	}
	if !km.Rename.Enabled() {
		t.Error("Rename binding should be enabled")
	}
	if !km.Refresh.Enabled() {
		t.Error("Refresh binding should be enabled")
	}
	if !km.Attach.Enabled() {
		t.Error("Attach binding should be enabled")
	}
	if !km.Detach.Enabled() {
		t.Error("Detach binding should be enabled")
	}
	if !km.SplitVertical.Enabled() {
		t.Error("SplitVertical binding should be enabled")
	}
	if !km.SplitHorizontal.Enabled() {
		t.Error("SplitHorizontal binding should be enabled")
	}
	if !km.OpenConfig.Enabled() {
		t.Error("OpenConfig binding should be enabled")
	}
	if !km.Help.Enabled() {
		t.Error("Help binding should be enabled")
	}
	if !km.Quit.Enabled() {
		t.Error("Quit binding should be enabled")
	}
	if !km.Cancel.Enabled() {
		t.Error("Cancel binding should be enabled")
	}
	if !km.Confirm.Enabled() {
		t.Error("Confirm binding should be enabled")
	}
}

func TestKeyBindingHelp(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		name    string
		binding key.Binding
		help    string
	}{
		{"Up", km.Up, "up"},
		{"Down", km.Down, "down"},
		{"Left", km.Left, "prev panel"},
		{"Right", km.Right, "next panel"},
		{"Tab", km.Tab, "next panel"},
		{"Create", km.Create, "new"},
		{"Delete", km.Delete, "delete"},
		{"Rename", km.Rename, "rename"},
		{"Refresh", km.Refresh, "refresh"},
		{"Attach", km.Attach, "attach"},
		{"Detach", km.Detach, "detach"},
		{"Help", km.Help, "help"},
		{"Quit", km.Quit, "quit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := tt.binding.Help()
			if help.Desc != tt.help {
				t.Errorf("%s Help().Desc = %q, expected %q", tt.name, help.Desc, tt.help)
			}
		})
	}
}

func TestShortHelp(t *testing.T) {
	km := DefaultKeyMap()
	shortHelp := km.ShortHelp()

	if len(shortHelp) == 0 {
		t.Error("ShortHelp should return at least one binding")
	}

	// Should include essential bindings
	expectedCount := 7 // Up, Down, Tab, Select, Create, Delete, Quit
	if len(shortHelp) != expectedCount {
		t.Errorf("ShortHelp length = %d, expected %d", len(shortHelp), expectedCount)
	}
}

func TestFullHelp(t *testing.T) {
	km := DefaultKeyMap()
	fullHelp := km.FullHelp()

	if len(fullHelp) == 0 {
		t.Error("FullHelp should return at least one group")
	}

	// Should have 4 groups
	expectedGroups := 4
	if len(fullHelp) != expectedGroups {
		t.Errorf("FullHelp groups = %d, expected %d", len(fullHelp), expectedGroups)
	}

	// Each group should have bindings
	for i, group := range fullHelp {
		if len(group) == 0 {
			t.Errorf("FullHelp group %d should have bindings", i)
		}
	}
}

func TestNavigationKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Test vim-style navigation keys
	upHelp := km.Up.Help()
	if upHelp.Key != "↑/k" {
		t.Errorf("Up key help = %q, expected '↑/k'", upHelp.Key)
	}

	downHelp := km.Down.Help()
	if downHelp.Key != "↓/j" {
		t.Errorf("Down key help = %q, expected '↓/j'", downHelp.Key)
	}

	leftHelp := km.Left.Help()
	if leftHelp.Key != "←/h" {
		t.Errorf("Left key help = %q, expected '←/h'", leftHelp.Key)
	}

	rightHelp := km.Right.Help()
	if rightHelp.Key != "→/l" {
		t.Errorf("Right key help = %q, expected '→/l'", rightHelp.Key)
	}
}

func TestActionKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Create = n
	createHelp := km.Create.Help()
	if createHelp.Key != "n" {
		t.Errorf("Create key = %q, expected 'n'", createHelp.Key)
	}

	// Delete = d
	deleteHelp := km.Delete.Help()
	if deleteHelp.Key != "d" {
		t.Errorf("Delete key = %q, expected 'd'", deleteHelp.Key)
	}

	// Rename = r
	renameHelp := km.Rename.Help()
	if renameHelp.Key != "r" {
		t.Errorf("Rename key = %q, expected 'r'", renameHelp.Key)
	}

	// Config = c
	configHelp := km.OpenConfig.Help()
	if configHelp.Key != "c" {
		t.Errorf("OpenConfig key = %q, expected 'c'", configHelp.Key)
	}
}

func TestSplitKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Vertical = v
	vSplitHelp := km.SplitVertical.Help()
	if vSplitHelp.Key != "v" {
		t.Errorf("SplitVertical key = %q, expected 'v'", vSplitHelp.Key)
	}
	if vSplitHelp.Desc != "vsplit" {
		t.Errorf("SplitVertical desc = %q, expected 'vsplit'", vSplitHelp.Desc)
	}

	// Horizontal = s
	hSplitHelp := km.SplitHorizontal.Help()
	if hSplitHelp.Key != "s" {
		t.Errorf("SplitHorizontal key = %q, expected 's'", hSplitHelp.Key)
	}
	if hSplitHelp.Desc != "hsplit" {
		t.Errorf("SplitHorizontal desc = %q, expected 'hsplit'", hSplitHelp.Desc)
	}
}

func TestQuitKey(t *testing.T) {
	km := DefaultKeyMap()

	quitHelp := km.Quit.Help()
	if quitHelp.Key != "q" {
		t.Errorf("Quit key = %q, expected 'q'", quitHelp.Key)
	}
}

func TestHelpKey(t *testing.T) {
	km := DefaultKeyMap()

	helpHelp := km.Help.Help()
	if helpHelp.Key != "?" {
		t.Errorf("Help key = %q, expected '?'", helpHelp.Key)
	}
}

func TestConfirmCancelKeys(t *testing.T) {
	km := DefaultKeyMap()

	// Cancel = esc
	cancelHelp := km.Cancel.Help()
	if cancelHelp.Key != "esc" {
		t.Errorf("Cancel key = %q, expected 'esc'", cancelHelp.Key)
	}

	// Confirm = y
	confirmHelp := km.Confirm.Help()
	if confirmHelp.Key != "y" {
		t.Errorf("Confirm key = %q, expected 'y'", confirmHelp.Key)
	}
}

func TestSelectKey(t *testing.T) {
	km := DefaultKeyMap()

	selectHelp := km.Select.Help()
	if selectHelp.Key != "space" {
		t.Errorf("Select key = %q, expected 'space'", selectHelp.Key)
	}
	if selectHelp.Desc != "switch" {
		t.Errorf("Select desc = %q, expected 'switch'", selectHelp.Desc)
	}
}

func TestKeyMapStruct(t *testing.T) {
	// Verify KeyMap struct has all expected fields
	km := KeyMap{}

	// Navigation fields
	_ = km.Up
	_ = km.Down
	_ = km.Left
	_ = km.Right
	_ = km.Tab
	_ = km.ShiftTab

	// Action fields
	_ = km.Enter
	_ = km.Select
	_ = km.Create
	_ = km.Delete
	_ = km.Rename
	_ = km.Refresh
	_ = km.Attach
	_ = km.Detach
	_ = km.SplitVertical
	_ = km.SplitHorizontal
	_ = km.OpenConfig

	// General fields
	_ = km.Help
	_ = km.Quit
	_ = km.Cancel
	_ = km.Confirm
}
