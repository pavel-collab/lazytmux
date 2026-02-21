package config

// OptionType defines the type of an option's value
type OptionType int

const (
	TypeBool   OptionType = iota // on/off
	TypeChoice                   // select from list
	TypeNumber                   // numeric value
	TypeString                   // arbitrary string
	TypeStyle                    // style (fg=..., bg=...)
)

// Scope defines where the option applies
type Scope string

const (
	ScopeServer  Scope = "server"
	ScopeSession Scope = "session"
	ScopeWindow  Scope = "window"
)

// Category groups related options
type Category struct {
	ID      string
	NameEN  string
	NameRU  string
	Options []Option
}

// Option describes a single tmux setting
type Option struct {
	Key     string     // tmux option name (e.g. "mouse")
	DescEN  string     // description in English
	DescRU  string     // description in Russian
	Type    OptionType // value type
	Default string     // default value
	Choices []string   // possible values for TypeChoice
	Min     int        // minimum for TypeNumber
	Max     int        // maximum for TypeNumber
	Scope   Scope      // server, session, or window
}

// ConfigValue stores the current value of an option
type ConfigValue struct {
	Key      string
	Value    string
	Modified bool   // modified relative to file
	Source   string // "default", "file", or "user"
}

// Config holds all configuration
type Config struct {
	Values       map[string]ConfigValue
	FilePath     string
	RawLines     []string // original file lines (to preserve comments)
	ParseErrors  []string
	Plugins      map[string]*PluginState // plugin states by repo
	PluginLines  []string                // raw plugin lines from file
}

// NewConfig creates a new Config with defaults
func NewConfig(path string) *Config {
	cfg := &Config{
		Values:      make(map[string]ConfigValue),
		FilePath:    path,
		RawLines:    []string{},
		Plugins:     make(map[string]*PluginState),
		PluginLines: []string{},
	}

	// Initialize with default values
	for key, opt := range GetAllOptions() {
		cfg.Values[key] = ConfigValue{
			Key:    key,
			Value:  opt.Default,
			Source: "default",
		}
	}

	// Initialize plugins (all disabled by default)
	for _, p := range GetPlugins() {
		cfg.Plugins[p.Repo] = &PluginState{
			Repo:     p.Repo,
			Enabled:  false,
			Settings: make(map[string]string),
		}
		// Set default settings
		for _, s := range p.Settings {
			cfg.Plugins[p.Repo].Settings[s.Key] = s.Default
		}
	}

	return cfg
}

// GetValue returns the current value of an option
func (c *Config) GetValue(key string) string {
	if val, ok := c.Values[key]; ok {
		return val.Value
	}
	return ""
}

// SetValue sets a value and marks it as modified
func (c *Config) SetValue(key, value string) {
	c.Values[key] = ConfigValue{
		Key:      key,
		Value:    value,
		Modified: true,
		Source:   "user",
	}
}

// HasChanges returns true if any value was modified
func (c *Config) HasChanges() bool {
	for _, v := range c.Values {
		if v.Modified {
			return true
		}
	}
	return false
}

// ResetToDefaults resets all options to their default values
func (c *Config) ResetToDefaults() {
	for key, opt := range GetAllOptions() {
		c.Values[key] = ConfigValue{
			Key:      key,
			Value:    opt.Default,
			Modified: true,
			Source:   "user",
		}
	}
}

// IsPluginEnabled returns whether a plugin is enabled
func (c *Config) IsPluginEnabled(repo string) bool {
	if p, ok := c.Plugins[repo]; ok {
		return p.Enabled
	}
	return false
}

// SetPluginEnabled enables or disables a plugin
func (c *Config) SetPluginEnabled(repo string, enabled bool) {
	if p, ok := c.Plugins[repo]; ok {
		p.Enabled = enabled
	} else {
		c.Plugins[repo] = &PluginState{
			Repo:     repo,
			Enabled:  enabled,
			Settings: make(map[string]string),
		}
	}
}

// GetPluginSetting returns a plugin setting value
func (c *Config) GetPluginSetting(repo, key string) string {
	if p, ok := c.Plugins[repo]; ok {
		if v, ok := p.Settings[key]; ok {
			return v
		}
	}
	// Return default from plugin definition
	if plugin, ok := GetPlugin(repo); ok {
		for _, s := range plugin.Settings {
			if s.Key == key {
				return s.Default
			}
		}
	}
	return ""
}

// SetPluginSetting sets a plugin setting value
func (c *Config) SetPluginSetting(repo, key, value string) {
	if p, ok := c.Plugins[repo]; ok {
		p.Settings[key] = value
	}
}

// HasPluginChanges returns true if any plugin state was modified
func (c *Config) HasPluginChanges() bool {
	// Compare current state with what was loaded from file
	for repo, state := range c.Plugins {
		// Check if this plugin line exists in PluginLines
		found := false
		for _, line := range c.PluginLines {
			if containsPluginRepo(line, repo) {
				found = true
				break
			}
		}
		if state.Enabled != found {
			return true
		}
	}
	return false
}

// containsPluginRepo checks if a line contains the plugin repo
func containsPluginRepo(line, repo string) bool {
	return len(line) > 0 && (len(repo) > 0 && (line == repo ||
		(len(line) > len(repo) && line[len(line)-len(repo):] == repo)))
}
