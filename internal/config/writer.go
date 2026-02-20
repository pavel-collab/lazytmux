package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// SaveConfig saves the configuration to file
func SaveConfig(cfg *Config) error {
	allOpts := GetAllOptions()

	// Track which keys we've updated in existing lines
	processedKeys := make(map[string]bool)

	// Patterns to match existing set commands
	setOptionRe := regexp.MustCompile(`^(set(?:-option)?\s+(?:-[sg]\s+)?(?:-[sg]\s+)?)(\S+)(\s+)(.+)$`)
	setWindowRe := regexp.MustCompile(`^(setw(?:-window-option)?\s+(?:-g\s+)?)(\S+)(\s+)(.+)$`)

	// Update existing lines
	updatedLines := make([]string, len(cfg.RawLines))
	copy(updatedLines, cfg.RawLines)

	for i, line := range updatedLines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines, comments, and plugin lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || isPluginLine(trimmed) {
			continue
		}

		var key string
		var newLine string

		if match := setOptionRe.FindStringSubmatch(trimmed); match != nil {
			key = match[2]
			if val, ok := cfg.Values[key]; ok && val.Modified {
				if _, known := allOpts[key]; known {
					newLine = match[1] + key + match[3] + val.Value
				}
			}
		} else if match := setWindowRe.FindStringSubmatch(trimmed); match != nil {
			key = match[2]
			if val, ok := cfg.Values[key]; ok && val.Modified {
				if _, known := allOpts[key]; known {
					newLine = match[1] + key + match[3] + val.Value
				}
			}
		}

		if newLine != "" {
			// Preserve leading whitespace from original line
			leadingSpace := ""
			for _, ch := range line {
				if ch == ' ' || ch == '\t' {
					leadingSpace += string(ch)
				} else {
					break
				}
			}
			updatedLines[i] = leadingSpace + newLine
			processedKeys[key] = true
		}
	}

	// Collect new lines that weren't in the original file
	var newLines []string
	for key, val := range cfg.Values {
		if val.Modified && !processedKeys[key] {
			if opt, known := allOpts[key]; known {
				var line string
				switch opt.Scope {
				case ScopeWindow:
					line = fmt.Sprintf("setw -g %s %s", key, val.Value)
				case ScopeServer:
					line = fmt.Sprintf("set -s %s %s", key, val.Value)
				default:
					line = fmt.Sprintf("set -g %s %s", key, val.Value)
				}
				newLines = append(newLines, line)
			}
		}
	}

	// Build final content
	var content strings.Builder

	for _, line := range updatedLines {
		content.WriteString(line)
		content.WriteString("\n")
	}

	// Add new options at the end
	if len(newLines) > 0 {
		// Add a marker comment if the file had content
		if len(updatedLines) > 0 {
			content.WriteString("\n# Added by lazytmux\n")
		}
		for _, line := range newLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	// Handle empty file case
	if len(updatedLines) == 0 && len(newLines) > 0 {
		content.Reset()
		content.WriteString("# tmux configuration - edited by lazytmux\n\n")
		for _, line := range newLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	// Write to file
	return os.WriteFile(cfg.FilePath, []byte(content.String()), 0644)
}

// GenerateConfigLine generates a config line for an option
func GenerateConfigLine(key, value string) string {
	opt, ok := GetOption(key)
	if !ok {
		return fmt.Sprintf("set -g %s %s", key, value)
	}

	switch opt.Scope {
	case ScopeWindow:
		return fmt.Sprintf("setw -g %s %s", key, value)
	case ScopeServer:
		return fmt.Sprintf("set -s %s %s", key, value)
	default:
		return fmt.Sprintf("set -g %s %s", key, value)
	}
}

// FormatValueForFile formats a value for writing to config file
func FormatValueForFile(opt Option, value string) string {
	switch opt.Type {
	case TypeString, TypeStyle:
		// Quote strings if they contain spaces
		if strings.Contains(value, " ") {
			return fmt.Sprintf("\"%s\"", value)
		}
	}
	return value
}
