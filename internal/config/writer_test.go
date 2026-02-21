package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveConfigNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetValue("mouse", "on")
	cfg.SetValue("history-limit", "10000")

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	// Read the file and verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Check that our values are present
	if !strings.Contains(contentStr, "mouse") || !strings.Contains(contentStr, "on") {
		t.Error("saved config should contain mouse setting")
	}
	if !strings.Contains(contentStr, "history-limit") || !strings.Contains(contentStr, "10000") {
		t.Error("saved config should contain history-limit setting")
	}
}

func TestSaveConfigPreservesComments(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	// Create initial config with comments
	initialContent := `# My tmux configuration
# Created by me

set -g mouse off

# Status bar settings
set -g status on
`
	err := os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("failed to create initial config: %v", err)
	}

	// Load and modify
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	cfg.SetValue("mouse", "on")

	// Save
	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	// Verify comments are preserved
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "# My tmux configuration") {
		t.Error("comment should be preserved")
	}
	if !strings.Contains(contentStr, "# Status bar settings") {
		t.Error("comment should be preserved")
	}
}

func TestSaveConfigUpdatesExistingValue(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	initialContent := `set -g mouse off
set -g history-limit 2000
`
	err := os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("failed to create initial config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	cfg.SetValue("mouse", "on")

	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Should have updated value
	if !strings.Contains(contentStr, "set -g mouse on") {
		t.Errorf("mouse value should be updated, got:\n%s", contentStr)
	}

	// Original value should not be present
	if strings.Contains(contentStr, "mouse off") {
		t.Error("old mouse value should be replaced")
	}
}

func TestSaveConfigAddsNewPlugin(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)
	cfg.SetPluginEnabled("tmux-plugins/tmux-resurrect", true)

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Check plugins are present
	if !strings.Contains(contentStr, "@plugin 'tmux-plugins/tpm'") {
		t.Error("tpm plugin should be in config")
	}
	if !strings.Contains(contentStr, "@plugin 'tmux-plugins/tmux-resurrect'") {
		t.Error("resurrect plugin should be in config")
	}

	// Check TPM initialization line is present
	if !strings.Contains(contentStr, "tpm/tpm") {
		t.Error("TPM run line should be present")
	}
}

func TestSaveConfigRemovesDisabledPlugin(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	initialContent := `set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-resurrect'
run '~/.tmux/plugins/tpm/tpm'
`
	err := os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("failed to create initial config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Disable resurrect
	cfg.SetPluginEnabled("tmux-plugins/tmux-resurrect", false)

	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// TPM should still be present
	if !strings.Contains(contentStr, "@plugin 'tmux-plugins/tpm'") {
		t.Error("tpm plugin should still be in config")
	}

	// Resurrect should be removed
	if strings.Contains(contentStr, "tmux-resurrect") {
		t.Error("disabled plugin should be removed from config")
	}
}

func TestSaveConfigWithPluginSettings(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetPluginEnabled("tmux-plugins/tmux-resurrect", true)
	cfg.SetPluginSetting("tmux-plugins/tmux-resurrect", "@resurrect-capture-pane-contents", "on")

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Plugin should be present
	if !strings.Contains(contentStr, "tmux-resurrect") {
		t.Error("plugin should be in config")
	}

	// Setting should be present (only non-default settings are written)
	if !strings.Contains(contentStr, "@resurrect-capture-pane-contents") {
		t.Error("plugin setting should be in config")
	}
}

func TestFormatPluginValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple value",
			input:    "on",
			expected: "on",
		},
		{
			name:     "value with space",
			input:    "hello world",
			expected: "'hello world'",
		},
		{
			name:     "value with hash",
			input:    "#ffffff",
			expected: "'#ffffff'",
		},
		{
			name:     "empty value",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPluginValue(tt.input)
			if result != tt.expected {
				t.Errorf("formatPluginValue(%q) = %q, expected %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateConfigLine(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "session scope option",
			key:      "mouse",
			value:    "on",
			expected: "set -g mouse on",
		},
		{
			name:     "window scope option",
			key:      "automatic-rename",
			value:    "off",
			expected: "setw -g automatic-rename off",
		},
		{
			name:     "server scope option",
			key:      "escape-time",
			value:    "0",
			expected: "set -s escape-time 0",
		},
		{
			name:     "unknown option defaults to session",
			key:      "unknown-option",
			value:    "value",
			expected: "set -g unknown-option value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateConfigLine(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("GenerateConfigLine(%q, %q) = %q, expected %q",
					tt.key, tt.value, result, tt.expected)
			}
		})
	}
}

func TestFormatValueForFile(t *testing.T) {
	tests := []struct {
		name     string
		optType  OptionType
		value    string
		expected string
	}{
		{
			name:     "bool value",
			optType:  TypeBool,
			value:    "on",
			expected: "on",
		},
		{
			name:     "number value",
			optType:  TypeNumber,
			value:    "500",
			expected: "500",
		},
		{
			name:     "string without spaces",
			optType:  TypeString,
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces",
			optType:  TypeString,
			value:    "hello world",
			expected: `"hello world"`,
		},
		{
			name:     "style with spaces",
			optType:  TypeStyle,
			value:    "fg=red bg=blue",
			expected: `"fg=red bg=blue"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Option{Type: tt.optType}
			result := FormatValueForFile(opt, tt.value)
			if result != tt.expected {
				t.Errorf("FormatValueForFile(..., %q) = %q, expected %q",
					tt.value, result, tt.expected)
			}
		})
	}
}

func TestSaveConfigTPMRunLineAtEnd(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetValue("mouse", "on")
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// TPM run line should be at the end
	lastNonEmptyLine := ""
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			lastNonEmptyLine = lines[i]
			break
		}
	}

	if !strings.Contains(lastNonEmptyLine, "tpm/tpm") {
		t.Errorf("TPM run line should be at the end, last line: %q", lastNonEmptyLine)
	}
}

func TestSaveConfigPreservesTPMRunLine(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	initialContent := `set -g @plugin 'tmux-plugins/tpm'
set -g mouse off

run '~/.tmux/plugins/tpm/tpm'
`
	err := os.WriteFile(configPath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("failed to create initial config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	cfg.SetValue("mouse", "on")

	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Original TPM line format should be preserved
	if !strings.Contains(contentStr, "run '~/.tmux/plugins/tpm/tpm'") {
		t.Error("original TPM run line should be preserved")
	}
}

func TestSaveConfigNoPluginsNoTPMLine(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetValue("mouse", "on")
	// No plugins enabled

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Should not have TPM line when no plugins
	if strings.Contains(contentStr, "tpm/tpm") {
		t.Error("TPM run line should not be present when no plugins enabled")
	}
}

func TestSaveConfigWindowScopeOptions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetValue("automatic-rename", "off")
	cfg.SetValue("pane-base-index", "1")

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Window scope options should use setw
	if !strings.Contains(contentStr, "setw -g automatic-rename off") {
		t.Error("window scope option should use setw")
	}
}

func TestSaveConfigServerScopeOptions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	cfg := NewConfig(configPath)
	cfg.SetValue("escape-time", "0")

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	contentStr := string(content)

	// Server scope options should use set -s
	if !strings.Contains(contentStr, "set -s escape-time 0") {
		t.Error("server scope option should use set -s")
	}
}

func TestSaveConfigReloadAndVerify(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	// Create and save config
	cfg := NewConfig(configPath)
	cfg.SetValue("mouse", "on")
	cfg.SetValue("history-limit", "50000")
	cfg.SetPluginEnabled("tmux-plugins/tpm", true)

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error: %v", err)
	}

	// Reload and verify
	cfg2, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	if cfg2.GetValue("mouse") != "on" {
		t.Errorf("reloaded mouse = %q, expected 'on'", cfg2.GetValue("mouse"))
	}

	if cfg2.GetValue("history-limit") != "50000" {
		t.Errorf("reloaded history-limit = %q, expected '50000'", cfg2.GetValue("history-limit"))
	}

	if !cfg2.IsPluginEnabled("tmux-plugins/tpm") {
		t.Error("reloaded tpm plugin should be enabled")
	}
}
