package ui

import (
	"lazytmux/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ConfigTab represents the active tab in config editor
type ConfigTab int

const (
	OptionsTab ConfigTab = iota
	PluginsTab
)

// ConfigEditorModel is the model for the configuration editor
type ConfigEditorModel struct {
	config     *config.Config
	categories []config.Category
	plugins    []config.Plugin

	// Tab navigation
	activeTab ConfigTab

	// Navigation (Options tab)
	categoryCursor int  // current category index
	optionCursor   int  // current option index in category
	focusOnOptions bool // true = focus on options list, false = on categories

	// Navigation (Plugins tab)
	pluginCursor        int  // current plugin index
	pluginSettingCursor int  // current setting index
	focusOnPluginSettings bool // true = focus on settings, false = on plugin list

	// Editing
	editing       bool
	editInput     textinput.Model
	editingOption *config.Option
	editingPluginSetting *config.PluginSetting
	editingPluginRepo    string

	// Choice selection mode (for TypeChoice)
	choosing     bool
	choiceCursor int

	// Confirmation dialogs
	confirmReset bool
	confirmSave  bool

	// TPM install dialog
	showTPMInstall bool

	// Dimensions
	width  int
	height int

	// Language (ru/en)
	language string

	// Status message (shown after save)
	statusMessage string
	statusIsError bool
}

// NewConfigEditorModel creates a new config editor model
func NewConfigEditorModel() ConfigEditorModel {
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Width = 40

	return ConfigEditorModel{
		categories: config.GetCategories(),
		plugins:    config.GetPlugins(),
		editInput:  ti,
		language:   "en",
		activeTab:  OptionsTab,
	}
}

// Init initializes the config editor
func (m ConfigEditorModel) Init() tea.Cmd {
	return LoadConfigCmd()
}

// SetSize sets the dimensions of the editor
func (m *ConfigEditorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetConfig sets the loaded configuration
func (m *ConfigEditorModel) SetConfig(cfg *config.Config) {
	m.config = cfg
}

// CurrentCategory returns the current category
func (m ConfigEditorModel) CurrentCategory() config.Category {
	if m.categoryCursor < len(m.categories) {
		return m.categories[m.categoryCursor]
	}
	return config.Category{}
}

// CurrentOption returns the current option
func (m ConfigEditorModel) CurrentOption() *config.Option {
	cat := m.CurrentCategory()
	if m.optionCursor < len(cat.Options) {
		opt := cat.Options[m.optionCursor]
		return &opt
	}
	return nil
}

// GetOptionValue returns the current value of an option
func (m ConfigEditorModel) GetOptionValue(key string) string {
	if m.config != nil {
		return m.config.GetValue(key)
	}
	return ""
}

// SetOptionValue sets a value for an option
func (m *ConfigEditorModel) SetOptionValue(key, value string) {
	if m.config != nil {
		m.config.SetValue(key, value)
	}
}

// HasChanges returns true if there are unsaved changes
func (m ConfigEditorModel) HasChanges() bool {
	if m.config != nil {
		return m.config.HasChanges()
	}
	return false
}

// IsModified returns true if a specific option was modified
func (m ConfigEditorModel) IsModified(key string) bool {
	if m.config != nil {
		if val, ok := m.config.Values[key]; ok {
			return val.Modified
		}
	}
	return false
}

// ToggleLanguage switches between English and Russian
func (m *ConfigEditorModel) ToggleLanguage() {
	if m.language == "en" {
		m.language = "ru"
	} else {
		m.language = "en"
	}
}

// moveUp moves the cursor up
func (m *ConfigEditorModel) moveUp() {
	if m.focusOnOptions {
		if m.optionCursor > 0 {
			m.optionCursor--
		}
	} else {
		if m.categoryCursor > 0 {
			m.categoryCursor--
			m.optionCursor = 0
		}
	}
}

// moveDown moves the cursor down
func (m *ConfigEditorModel) moveDown() {
	if m.focusOnOptions {
		cat := m.CurrentCategory()
		if m.optionCursor < len(cat.Options)-1 {
			m.optionCursor++
		}
	} else {
		if m.categoryCursor < len(m.categories)-1 {
			m.categoryCursor++
			m.optionCursor = 0
		}
	}
}

// resetState resets editing/selection state
func (m *ConfigEditorModel) resetState() {
	m.editing = false
	m.choosing = false
	m.confirmReset = false
	m.confirmSave = false
	m.editingOption = nil
	m.editingPluginSetting = nil
	m.editingPluginRepo = ""
	m.showTPMInstall = false
}

// CurrentPlugin returns the current plugin
func (m ConfigEditorModel) CurrentPlugin() *config.Plugin {
	if m.pluginCursor < len(m.plugins) {
		p := m.plugins[m.pluginCursor]
		return &p
	}
	return nil
}

// CurrentPluginSetting returns the current plugin setting
func (m ConfigEditorModel) CurrentPluginSetting() *config.PluginSetting {
	p := m.CurrentPlugin()
	if p == nil || m.pluginSettingCursor >= len(p.Settings) {
		return nil
	}
	s := p.Settings[m.pluginSettingCursor]
	return &s
}

// IsPluginEnabled returns whether the current plugin is enabled
func (m ConfigEditorModel) IsPluginEnabled(repo string) bool {
	if m.config != nil {
		return m.config.IsPluginEnabled(repo)
	}
	return false
}

// TogglePlugin toggles a plugin on/off
func (m *ConfigEditorModel) TogglePlugin(repo string) {
	if m.config != nil {
		current := m.config.IsPluginEnabled(repo)
		m.config.SetPluginEnabled(repo, !current)
	}
}

// GetPluginSettingValue returns a plugin setting value
func (m ConfigEditorModel) GetPluginSettingValue(repo, key string) string {
	if m.config != nil {
		return m.config.GetPluginSetting(repo, key)
	}
	return ""
}

// SetPluginSettingValue sets a plugin setting value
func (m *ConfigEditorModel) SetPluginSettingValue(repo, key, value string) {
	if m.config != nil {
		m.config.SetPluginSetting(repo, key, value)
	}
}

// moveUpPlugins moves cursor up in plugins tab
func (m *ConfigEditorModel) moveUpPlugins() {
	if m.focusOnPluginSettings {
		if m.pluginSettingCursor > 0 {
			m.pluginSettingCursor--
		}
	} else {
		if m.pluginCursor > 0 {
			m.pluginCursor--
			m.pluginSettingCursor = 0
		}
	}
}

// moveDownPlugins moves cursor down in plugins tab
func (m *ConfigEditorModel) moveDownPlugins() {
	if m.focusOnPluginSettings {
		p := m.CurrentPlugin()
		if p != nil && m.pluginSettingCursor < len(p.Settings)-1 {
			m.pluginSettingCursor++
		}
	} else {
		if m.pluginCursor < len(m.plugins)-1 {
			m.pluginCursor++
			m.pluginSettingCursor = 0
		}
	}
}

// SwitchTab switches between Options and Plugins tabs
func (m *ConfigEditorModel) SwitchTab() {
	if m.activeTab == OptionsTab {
		m.activeTab = PluginsTab
	} else {
		m.activeTab = OptionsTab
	}
}
