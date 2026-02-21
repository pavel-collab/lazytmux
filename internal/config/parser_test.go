package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath(t *testing.T) {
	path := ConfigPath()

	// Should end with .tmux.conf
	if filepath.Base(path) != ".tmux.conf" {
		t.Errorf("ConfigPath() = %q, expected to end with .tmux.conf", path)
	}

	// Should be an absolute path (unless home lookup fails)
	home, err := os.UserHomeDir()
	if err == nil {
		expected := filepath.Join(home, ".tmux.conf")
		if path != expected {
			t.Errorf("ConfigPath() = %q, expected %q", path, expected)
		}
	}
}

func TestCleanValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain value",
			input:    "on",
			expected: "on",
		},
		{
			name:     "value with spaces",
			input:    "  on  ",
			expected: "on",
		},
		{
			name:     "double quoted value",
			input:    `"hello world"`,
			expected: "hello world",
		},
		{
			name:     "single quoted value",
			input:    `'hello world'`,
			expected: "hello world",
		},
		{
			name:     "value with trailing comment",
			input:    "on # this is a comment",
			expected: "on",
		},
		{
			name:     "value with hash but no space before",
			input:    "color#ffffff",
			expected: "color#ffffff",
		},
		{
			name:     "quoted value with hash inside",
			input:    `"#ffffff"`,
			expected: "#ffffff",
		},
		{
			name:     "empty value",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "mismatched quotes - left only",
			input:    `"hello`,
			expected: `"hello`,
		},
		{
			name:     "number value",
			input:    "500",
			expected: "500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanValue(tt.input)
			if result != tt.expected {
				t.Errorf("cleanValue(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPluginLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "run-shell command",
			input:    "run-shell ~/.tmux/plugins/tpm/tpm",
			expected: true,
		},
		{
			name:     "run command",
			input:    "run '~/.tmux/plugins/tpm/tpm'",
			expected: true,
		},
		{
			name:     "plugin declaration",
			input:    "set -g @plugin 'tmux-plugins/tpm'",
			expected: true,
		},
		{
			name:     "user option with @",
			input:    "set -g @my-option 'value'",
			expected: true,
		},
		{
			name:     "regular set option",
			input:    "set -g mouse on",
			expected: false,
		},
		{
			name:     "setw option",
			input:    "setw -g mode-keys vi",
			expected: false,
		},
		{
			name:     "empty line",
			input:    "",
			expected: false,
		},
		{
			name:     "comment",
			input:    "# this is a comment",
			expected: false,
		},
		{
			name:     "bind command",
			input:    "bind r source-file ~/.tmux.conf",
			expected: false,
		},
		{
			name:     "RUN-SHELL uppercase",
			input:    "RUN-SHELL command",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPluginLine(tt.input)
			if result != tt.expected {
				t.Errorf("isPluginLine(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// Test loading a config from a non-existent file
	cfg, err := LoadConfig("/nonexistent/path/.tmux.conf")
	if err != nil {
		t.Errorf("LoadConfig() returned unexpected error: %v", err)
	}
	if cfg == nil {
		t.Error("LoadConfig() returned nil config")
		return
	}

	// Should have default values
	if len(cfg.Values) == 0 {
		t.Error("expected default values to be set")
	}

	// Check a specific default
	mouseVal := cfg.GetValue("mouse")
	if mouseVal != "off" {
		t.Errorf("default mouse value = %q, expected 'off'", mouseVal)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	content := `# Test tmux configuration
set -g mouse on
set -g history-limit 5000
setw -g mode-keys vi
set -g status-position top

# Plugin
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-resurrect'

# Bindings (should be skipped)
bind r source-file ~/.tmux.conf
unbind C-b
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Check parsed values
	tests := []struct {
		key      string
		expected string
	}{
		{"mouse", "on"},
		{"history-limit", "5000"},
		{"mode-keys", "vi"},
		{"status-position", "top"},
	}

	for _, tt := range tests {
		t.Run("value_"+tt.key, func(t *testing.T) {
			val := cfg.GetValue(tt.key)
			if val != tt.expected {
				t.Errorf("GetValue(%q) = %q, expected %q", tt.key, val, tt.expected)
			}
		})
	}

	// Check plugins
	if !cfg.IsPluginEnabled("tmux-plugins/tpm") {
		t.Error("expected tpm plugin to be enabled")
	}
	if !cfg.IsPluginEnabled("tmux-plugins/tmux-resurrect") {
		t.Error("expected resurrect plugin to be enabled")
	}

	// Check raw lines preserved
	if len(cfg.RawLines) == 0 {
		t.Error("expected RawLines to be preserved")
	}
}

func TestLoadConfigWithQuotedValues(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	content := `set -g status-left "Session: #S"
set -g status-right '#H %Y-%m-%d'
set -g default-terminal 'tmux-256color'
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Check that quotes were removed from known options
	val := cfg.GetValue("default-terminal")
	if val != "tmux-256color" {
		t.Errorf("GetValue('default-terminal') = %q, expected 'tmux-256color'", val)
	}
}

func TestLoadConfigWithComments(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	content := `# This is a comment
set -g mouse on  # Enable mouse

# Another comment
set -g base-index 1
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Value should not include comment
	val := cfg.GetValue("mouse")
	if val != "on" {
		t.Errorf("GetValue('mouse') = %q, expected 'on'", val)
	}

	val = cfg.GetValue("base-index")
	if val != "1" {
		t.Errorf("GetValue('base-index') = %q, expected '1'", val)
	}
}

func TestLoadConfigWithPluginSettings(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	content := `set -g @plugin 'tmux-plugins/tmux-resurrect'
set -g @resurrect-capture-pane-contents 'on'
set -g @resurrect-strategy-vim 'session'
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Check plugin is enabled
	if !cfg.IsPluginEnabled("tmux-plugins/tmux-resurrect") {
		t.Error("expected resurrect plugin to be enabled")
	}

	// Check plugin settings
	setting := cfg.GetPluginSetting("tmux-plugins/tmux-resurrect", "@resurrect-capture-pane-contents")
	if setting != "on" {
		t.Errorf("plugin setting = %q, expected 'on'", setting)
	}
}

func TestLoadConfigEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	err := os.WriteFile(configPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Should still have defaults
	if len(cfg.Values) == 0 {
		t.Error("expected default values")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with existing file
	existingPath := filepath.Join(tmpDir, "exists.txt")
	err := os.WriteFile(existingPath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if !FileExists(existingPath) {
		t.Error("FileExists() returned false for existing file")
	}

	// Test with non-existent file
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist.txt")
	if FileExists(nonExistentPath) {
		t.Error("FileExists() returned true for non-existent file")
	}
}

func TestLoadConfigVariousSetFormats(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".tmux.conf")

	// Test various set command formats
	content := `set -g mouse on
set-option -g history-limit 10000
set -s escape-time 0
setw -g mode-keys vi
set-window-option -g automatic-rename on
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"mouse", "on"},
		{"history-limit", "10000"},
		{"escape-time", "0"},
		{"mode-keys", "vi"},
		{"automatic-rename", "on"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			val := cfg.GetValue(tt.key)
			if val != tt.expected {
				t.Errorf("GetValue(%q) = %q, expected %q", tt.key, val, tt.expected)
			}
		})
	}
}
