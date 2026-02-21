package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("/test/path/.tmux.conf")

	if cfg.FilePath != "/test/path/.tmux.conf" {
		t.Errorf("FilePath = %q, expected %q", cfg.FilePath, "/test/path/.tmux.conf")
	}

	if cfg.Values == nil {
		t.Error("Values map should not be nil")
	}

	if cfg.Plugins == nil {
		t.Error("Plugins map should not be nil")
	}

	if cfg.RawLines == nil {
		t.Error("RawLines should not be nil")
	}

	// Check that defaults are populated
	if len(cfg.Values) == 0 {
		t.Error("expected default values to be populated")
	}

	// Check a specific default
	if val := cfg.GetValue("mouse"); val != "off" {
		t.Errorf("default mouse value = %q, expected 'off'", val)
	}

	// Check plugins are initialized (disabled by default)
	for _, plugin := range GetPlugins() {
		if cfg.IsPluginEnabled(plugin.Repo) {
			t.Errorf("plugin %q should be disabled by default", plugin.Repo)
		}
	}
}

func TestConfigGetSetValue(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Test getting a default value
	val := cfg.GetValue("mouse")
	if val != "off" {
		t.Errorf("GetValue('mouse') = %q, expected 'off'", val)
	}

	// Test setting a value
	cfg.SetValue("mouse", "on")
	val = cfg.GetValue("mouse")
	if val != "on" {
		t.Errorf("GetValue('mouse') after set = %q, expected 'on'", val)
	}

	// Check that Modified flag is set
	if !cfg.Values["mouse"].Modified {
		t.Error("Modified flag should be true after SetValue")
	}

	// Check that Source is set to "user"
	if cfg.Values["mouse"].Source != "user" {
		t.Errorf("Source = %q, expected 'user'", cfg.Values["mouse"].Source)
	}

	// Test getting non-existent key
	val = cfg.GetValue("nonexistent-key")
	if val != "" {
		t.Errorf("GetValue('nonexistent-key') = %q, expected empty string", val)
	}
}

func TestConfigHasChanges(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Initially no changes
	if cfg.HasChanges() {
		t.Error("new config should not have changes")
	}

	// After setting a value
	cfg.SetValue("mouse", "on")
	if !cfg.HasChanges() {
		t.Error("config should have changes after SetValue")
	}
}

func TestConfigHasChangesWithPlugins(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Initially no changes
	if cfg.HasChanges() {
		t.Error("new config should not have changes")
	}

	// Enable a plugin
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)
	if !cfg.HasChanges() {
		t.Error("config should have changes after enabling plugin")
	}
}

func TestConfigResetToDefaults(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Change some values
	cfg.SetValue("mouse", "on")
	cfg.SetValue("history-limit", "10000")

	// Reset to defaults
	cfg.ResetToDefaults()

	// Check values are back to defaults
	if val := cfg.GetValue("mouse"); val != "off" {
		t.Errorf("mouse after reset = %q, expected 'off'", val)
	}

	// Check that Modified flag is still true (indicates change from file state)
	if !cfg.Values["mouse"].Modified {
		t.Error("Modified flag should be true after reset")
	}
}

func TestConfigPluginOperations(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Test enabling a plugin
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)
	if !cfg.IsPluginEnabled("tmux-plugins/tpm") {
		t.Error("plugin should be enabled after SetPluginEnabled(true)")
	}

	// Test disabling a plugin
	cfg.SetPluginEnabled("tmux-plugins/tpm", false)
	if cfg.IsPluginEnabled("tmux-plugins/tpm") {
		t.Error("plugin should be disabled after SetPluginEnabled(false)")
	}

	// Test enabling unknown plugin
	cfg.SetPluginEnabled("custom/unknown-plugin", true)
	if !cfg.IsPluginEnabled("custom/unknown-plugin") {
		t.Error("unknown plugin should be enabled")
	}
}

func TestConfigPluginSettings(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Enable a plugin with settings
	cfg.SetPluginEnabled("tmux-plugins/tmux-resurrect", true)

	// Get default setting
	defaultVal := cfg.GetPluginSetting("tmux-plugins/tmux-resurrect", "@resurrect-capture-pane-contents")
	if defaultVal != "off" {
		t.Errorf("default plugin setting = %q, expected 'off'", defaultVal)
	}

	// Set a plugin setting
	cfg.SetPluginSetting("tmux-plugins/tmux-resurrect", "@resurrect-capture-pane-contents", "on")

	// Verify it was set
	val := cfg.GetPluginSetting("tmux-plugins/tmux-resurrect", "@resurrect-capture-pane-contents")
	if val != "on" {
		t.Errorf("plugin setting after set = %q, expected 'on'", val)
	}

	// Test getting setting for non-existent plugin
	val = cfg.GetPluginSetting("nonexistent/plugin", "@some-setting")
	if val != "" {
		t.Errorf("setting for non-existent plugin = %q, expected empty", val)
	}
}

func TestConfigHasPluginChanges(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Initially no plugin changes
	if cfg.HasPluginChanges() {
		t.Error("new config should not have plugin changes")
	}

	// Enable a plugin (creates change since PluginLines is empty)
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)
	if !cfg.HasPluginChanges() {
		t.Error("should have plugin changes after enabling")
	}
}

func TestConfigClearModifiedFlags(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Make some changes
	cfg.SetValue("mouse", "on")
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)

	// Clear modified flags
	cfg.ClearModifiedFlags()

	// Check that option modified flags are cleared
	if cfg.Values["mouse"].Modified {
		t.Error("Modified flag should be false after ClearModifiedFlags")
	}

	// Check that PluginLines is updated
	found := false
	for _, line := range cfg.PluginLines {
		if line == "tmux-plugins/tpm" {
			found = true
			break
		}
	}
	if !found {
		t.Error("PluginLines should include enabled plugin after ClearModifiedFlags")
	}
}

func TestContainsPluginRepo(t *testing.T) {
	tests := []struct {
		line     string
		repo     string
		expected bool
	}{
		{
			line:     "tmux-plugins/tpm",
			repo:     "tmux-plugins/tpm",
			expected: true,
		},
		{
			line:     "tmux-plugins/tpm",
			repo:     "tmux-plugins/resurrect",
			expected: false,
		},
		{
			line:     "",
			repo:     "tmux-plugins/tpm",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.line+"_"+tt.repo, func(t *testing.T) {
			result := containsPluginRepo(tt.line, tt.repo)
			if result != tt.expected {
				t.Errorf("containsPluginRepo(%q, %q) = %v, expected %v",
					tt.line, tt.repo, result, tt.expected)
			}
		})
	}
}

func TestOptionTypes(t *testing.T) {
	// Verify option type constants
	if TypeBool != 0 {
		t.Error("TypeBool should be 0")
	}
	if TypeChoice != 1 {
		t.Error("TypeChoice should be 1")
	}
	if TypeNumber != 2 {
		t.Error("TypeNumber should be 2")
	}
	if TypeString != 3 {
		t.Error("TypeString should be 3")
	}
	if TypeStyle != 4 {
		t.Error("TypeStyle should be 4")
	}
}

func TestScopeConstants(t *testing.T) {
	if ScopeServer != "server" {
		t.Error("ScopeServer should be 'server'")
	}
	if ScopeSession != "session" {
		t.Error("ScopeSession should be 'session'")
	}
	if ScopeWindow != "window" {
		t.Error("ScopeWindow should be 'window'")
	}
}

func TestConfigValueStruct(t *testing.T) {
	cv := ConfigValue{
		Key:      "mouse",
		Value:    "on",
		Modified: true,
		Source:   "user",
	}

	if cv.Key != "mouse" {
		t.Errorf("Key = %q, expected 'mouse'", cv.Key)
	}
	if cv.Value != "on" {
		t.Errorf("Value = %q, expected 'on'", cv.Value)
	}
	if !cv.Modified {
		t.Error("Modified should be true")
	}
	if cv.Source != "user" {
		t.Errorf("Source = %q, expected 'user'", cv.Source)
	}
}

func TestConfigMultipleValuesModified(t *testing.T) {
	cfg := NewConfig("/test/.tmux.conf")

	// Modify multiple values
	cfg.SetValue("mouse", "on")
	cfg.SetValue("history-limit", "10000")
	cfg.SetValue("base-index", "1")

	// All should have HasChanges return true
	if !cfg.HasChanges() {
		t.Error("HasChanges should return true")
	}

	// Count modified values
	modifiedCount := 0
	for _, v := range cfg.Values {
		if v.Modified {
			modifiedCount++
		}
	}

	if modifiedCount != 3 {
		t.Errorf("expected 3 modified values, got %d", modifiedCount)
	}
}
