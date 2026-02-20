package config

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ConfigPath returns the default path to ~/.tmux.conf
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".tmux.conf"
	}
	return filepath.Join(home, ".tmux.conf")
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	cfg := NewConfig(path)

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		// File doesn't exist, use defaults
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Regex patterns for parsing
	// Match: set [-option] [-s|-g] option value
	// Match: setw [-window-option] [-g] option value
	setOptionRe := regexp.MustCompile(`^set(?:-option)?\s+(?:-[sg]\s+)?(-[sg]\s+)?(\S+)\s+(.+)$`)
	setWindowRe := regexp.MustCompile(`^setw(?:-window-option)?\s+(?:-g\s+)?(\S+)\s+(.+)$`)

	allOpts := GetAllOptions()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cfg.RawLines = append(cfg.RawLines, line)

		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Skip plugin lines (run-shell, @plugin)
		if isPluginLine(trimmed) {
			continue
		}

		// Skip bind and unbind commands
		if strings.HasPrefix(trimmed, "bind") || strings.HasPrefix(trimmed, "unbind") {
			continue
		}

		var key, value string

		if match := setOptionRe.FindStringSubmatch(trimmed); match != nil {
			key = match[2]
			value = cleanValue(match[3])
		} else if match := setWindowRe.FindStringSubmatch(trimmed); match != nil {
			key = match[1]
			value = cleanValue(match[2])
		}

		if key != "" {
			// Only track known options
			if _, known := allOpts[key]; known {
				cfg.Values[key] = ConfigValue{
					Key:    key,
					Value:  value,
					Source: "file",
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// isPluginLine checks if the line is related to tmux plugins
func isPluginLine(line string) bool {
	lower := strings.ToLower(line)

	// Plugin manager commands
	if strings.HasPrefix(lower, "run-shell") || strings.HasPrefix(lower, "run ") {
		return true
	}

	// TPM plugin declarations
	if strings.Contains(line, "@plugin") {
		return true
	}

	// User options (prefixed with @)
	if strings.HasPrefix(strings.TrimSpace(line), "set") && strings.Contains(line, " @") {
		return true
	}

	return false
}

// cleanValue removes quotes and trailing comments from a value
func cleanValue(value string) string {
	value = strings.TrimSpace(value)

	// Remove surrounding quotes
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	// Remove trailing comments (but be careful with # in values)
	// Only remove if there's a space before #
	if idx := strings.Index(value, " #"); idx != -1 {
		value = strings.TrimSpace(value[:idx])
	}

	return value
}

// FileExists checks if the config file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
