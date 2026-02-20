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
	Values      map[string]ConfigValue
	FilePath    string
	RawLines    []string // original file lines (to preserve comments)
	ParseErrors []string
}

// NewConfig creates a new Config with defaults
func NewConfig(path string) *Config {
	cfg := &Config{
		Values:   make(map[string]ConfigValue),
		FilePath: path,
		RawLines: []string{},
	}

	// Initialize with default values
	for key, opt := range GetAllOptions() {
		cfg.Values[key] = ConfigValue{
			Key:    key,
			Value:  opt.Default,
			Source: "default",
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
