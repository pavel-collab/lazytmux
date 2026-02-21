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
	processedPlugins := make(map[string]bool)
	processedPluginSettings := make(map[string]bool)

	// Patterns to match existing set commands
	setOptionRe := regexp.MustCompile(`^(set(?:-option)?\s+(?:-[sg]\s+)?(?:-[sg]\s+)?)(\S+)(\s+)(.+)$`)
	setWindowRe := regexp.MustCompile(`^(setw(?:-window-option)?\s+(?:-g\s+)?)(\S+)(\s+)(.+)$`)
	pluginRe := regexp.MustCompile(`^set\s+(?:-g\s+)?@plugin\s+['"]?([^'"]+)['"]?$`)
	pluginSettingRe := regexp.MustCompile(`^(set\s+(?:-g\s+)?)(@[a-zA-Z0-9_-]+)(\s+)['"]?([^'"]+)['"]?$`)
	tpmRunRe := regexp.MustCompile(`^run\s+['"]?.*tpm/tpm['"]?$`)

	// Process lines - we'll rebuild the file
	var resultLines []string
	var tpmRunLine string
	hasTPMRun := false

	for _, line := range cfg.RawLines {
		trimmed := strings.TrimSpace(line)

		// Preserve empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			resultLines = append(resultLines, line)
			continue
		}

		// Handle @plugin lines
		if match := pluginRe.FindStringSubmatch(trimmed); match != nil {
			repo := match[1]
			processedPlugins[repo] = true
			// Only keep if still enabled
			if state, ok := cfg.Plugins[repo]; ok && state.Enabled {
				resultLines = append(resultLines, line)
			}
			continue
		}

		// Handle plugin settings (@setting)
		if match := pluginSettingRe.FindStringSubmatch(trimmed); match != nil {
			settingKey := match[2]
			processedPluginSettings[settingKey] = true
			// Find the plugin and update the value
			settingHandled := false
			for _, p := range GetPlugins() {
				if state, ok := cfg.Plugins[p.Repo]; ok && state.Enabled {
					for _, s := range p.Settings {
						if s.Key == settingKey {
							newValue := cfg.GetPluginSetting(p.Repo, settingKey)
							settingLine := match[1] + settingKey + match[3] + formatPluginValue(newValue)
							resultLines = append(resultLines, settingLine)
							settingHandled = true
							break
						}
					}
				}
				if settingHandled {
					break
				}
			}
			// If plugin is disabled or setting not found, skip this line
			continue
		}

		// Track TPM run line (should be at the end)
		if tpmRunRe.MatchString(trimmed) {
			tpmRunLine = line
			hasTPMRun = true
			continue // Don't add yet, we'll add at the end
		}

		// Handle regular options
		var key string
		var newLine string

		if match := setOptionRe.FindStringSubmatch(trimmed); match != nil {
			key = match[2]
			// Skip if it's a plugin-related line we already handled
			if strings.HasPrefix(key, "@") {
				continue
			}
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
			resultLines = append(resultLines, newLine)
			processedKeys[key] = true
		} else {
			resultLines = append(resultLines, line)
			if key != "" {
				processedKeys[key] = true
			}
		}
	}

	// Collect new options that weren't in the original file
	var newOptionLines []string
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
				newOptionLines = append(newOptionLines, line)
			}
		}
	}

	// Collect new plugins that weren't in the original file
	var newPluginLines []string
	var newPluginSettingLines []string
	for repo, state := range cfg.Plugins {
		if state.Enabled && !processedPlugins[repo] {
			newPluginLines = append(newPluginLines, fmt.Sprintf("set -g @plugin '%s'", repo))
		}
		// Add new plugin settings
		if state.Enabled {
			if plugin, ok := GetPlugin(repo); ok {
				for _, s := range plugin.Settings {
					if !processedPluginSettings[s.Key] {
						value := cfg.GetPluginSetting(repo, s.Key)
						if value != s.Default { // Only add non-default settings
							newPluginSettingLines = append(newPluginSettingLines,
								fmt.Sprintf("set -g %s '%s'", s.Key, value))
						}
					}
				}
			}
		}
	}

	// Build final content
	var content strings.Builder

	// Write existing lines
	for _, line := range resultLines {
		content.WriteString(line)
		content.WriteString("\n")
	}

	// Add new options
	if len(newOptionLines) > 0 {
		if len(resultLines) > 0 {
			content.WriteString("\n# Added by lazytmux\n")
		}
		for _, line := range newOptionLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	// Add new plugins
	if len(newPluginLines) > 0 {
		content.WriteString("\n# Plugins (managed by lazytmux)\n")
		for _, line := range newPluginLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	// Add new plugin settings
	if len(newPluginSettingLines) > 0 {
		content.WriteString("\n# Plugin settings\n")
		for _, line := range newPluginSettingLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	// Always add TPM run line at the end if any plugin is enabled
	hasEnabledPlugins := false
	for _, state := range cfg.Plugins {
		if state.Enabled {
			hasEnabledPlugins = true
			break
		}
	}

	if hasEnabledPlugins {
		if !hasTPMRun {
			content.WriteString("\n# Initialize TMUX plugin manager (keep this line at the very bottom)\n")
			content.WriteString("run '~/.tmux/plugins/tpm/tpm'\n")
		} else {
			content.WriteString("\n")
			content.WriteString(tpmRunLine)
			content.WriteString("\n")
		}
	}

	// Handle empty file case
	finalContent := content.String()
	if strings.TrimSpace(finalContent) == "" && (len(newOptionLines) > 0 || len(newPluginLines) > 0) {
		content.Reset()
		content.WriteString("# tmux configuration - edited by lazytmux\n\n")
		for _, line := range newOptionLines {
			content.WriteString(line)
			content.WriteString("\n")
		}
		if len(newPluginLines) > 0 {
			content.WriteString("\n# Plugins (managed by lazytmux)\n")
			for _, line := range newPluginLines {
				content.WriteString(line)
				content.WriteString("\n")
			}
		}
		if len(newPluginSettingLines) > 0 {
			content.WriteString("\n# Plugin settings\n")
			for _, line := range newPluginSettingLines {
				content.WriteString(line)
				content.WriteString("\n")
			}
		}
		if hasEnabledPlugins {
			content.WriteString("\n# Initialize TMUX plugin manager (keep this line at the very bottom)\n")
			content.WriteString("run '~/.tmux/plugins/tpm/tpm'\n")
		}
		finalContent = content.String()
	}

	// Write to file
	return os.WriteFile(cfg.FilePath, []byte(finalContent), 0644)
}

// formatPluginValue formats a plugin setting value
func formatPluginValue(value string) string {
	if strings.Contains(value, " ") || strings.Contains(value, "#") {
		return fmt.Sprintf("'%s'", value)
	}
	return value
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
