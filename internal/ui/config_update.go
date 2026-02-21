package ui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"lazytmux/internal/config"
)

// ConfigLoadedMsg is sent when config is loaded
type ConfigLoadedMsg struct {
	Config *config.Config
	Err    error
}

// ConfigSavedMsg is sent when config is saved
type ConfigSavedMsg struct {
	Err error
}

// ExitConfigEditorMsg signals exit from config editor
type ExitConfigEditorMsg struct{}

// LoadConfigCmd creates a command to load config
func LoadConfigCmd() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.LoadConfig(config.ConfigPath())
		return ConfigLoadedMsg{Config: cfg, Err: err}
	}
}

// SaveConfigCmd creates a command to save config
func SaveConfigCmd(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		err := config.SaveConfig(cfg)
		return ConfigSavedMsg{Err: err}
	}
}

// ExitConfigEditorCmd creates a command to exit config editor
func ExitConfigEditorCmd() tea.Cmd {
	return func() tea.Msg {
		return ExitConfigEditorMsg{}
	}
}

// Update handles messages for the config editor
func (m ConfigEditorModel) Update(msg tea.Msg) (ConfigEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case ConfigLoadedMsg:
		if msg.Err == nil && msg.Config != nil {
			m.config = msg.Config
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m ConfigEditorModel) handleKeyPress(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	// Handle text input mode
	if m.editing {
		return m.handleTextInput(msg)
	}

	// Handle choice selection mode
	if m.choosing {
		return m.handleChoiceSelection(msg)
	}

	// Handle reset confirmation
	if m.confirmReset {
		return m.handleResetConfirm(msg)
	}

	// Handle save confirmation
	if m.confirmSave {
		return m.handleSaveConfirm(msg)
	}

	// Handle TPM install dialog
	if m.showTPMInstall {
		return m.handleTPMInstallDialog(msg)
	}

	// Route to appropriate tab handler
	switch m.activeTab {
	case PluginsTab:
		return m.handlePluginsTabKeyPress(msg)
	default:
		return m.handleOptionsTabKeyPress(msg)
	}
}

func (m ConfigEditorModel) handleOptionsTabKeyPress(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		if m.HasChanges() {
			m.confirmSave = true
			return m, nil
		}
		return m, ExitConfigEditorCmd()

	case "j", "down":
		m.moveDown()

	case "k", "up":
		m.moveUp()

	case "h", "left":
		if m.focusOnOptions {
			m.focusOnOptions = false
		}

	case "l", "right":
		if !m.focusOnOptions {
			m.focusOnOptions = true
		}

	case "enter":
		if !m.focusOnOptions {
			m.focusOnOptions = true
		} else {
			return m.startEditing()
		}

	case "tab":
		m.focusOnOptions = !m.focusOnOptions
		if !m.focusOnOptions {
			m.optionCursor = 0
		}

	case " ":
		// Quick toggle for bool options
		return m.toggleBool()

	case "s", "ctrl+s":
		// Save
		if m.config != nil && m.HasChanges() {
			return m, SaveConfigCmd(m.config)
		}

	case "r":
		// Reset to defaults
		m.confirmReset = true

	case "L":
		// Toggle language
		m.ToggleLanguage()

	case "1":
		// Stay on Options tab
		m.activeTab = OptionsTab

	case "2":
		// Switch to Plugins tab
		m.activeTab = PluginsTab

	case "p":
		// Switch to Plugins tab
		m.activeTab = PluginsTab
	}

	return m, nil
}

func (m ConfigEditorModel) handlePluginsTabKeyPress(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		if m.HasChanges() {
			m.confirmSave = true
			return m, nil
		}
		return m, ExitConfigEditorCmd()

	case "j", "down":
		m.moveDownPlugins()

	case "k", "up":
		m.moveUpPlugins()

	case "h", "left":
		if m.focusOnPluginSettings {
			m.focusOnPluginSettings = false
		}

	case "l", "right":
		p := m.CurrentPlugin()
		if p != nil && len(p.Settings) > 0 && m.IsPluginEnabled(p.Repo) {
			m.focusOnPluginSettings = true
		}

	case "tab":
		p := m.CurrentPlugin()
		if p != nil && len(p.Settings) > 0 && m.IsPluginEnabled(p.Repo) {
			m.focusOnPluginSettings = !m.focusOnPluginSettings
		}

	case " ", "enter":
		if m.focusOnPluginSettings {
			return m.startPluginSettingEditing()
		} else {
			return m.toggleCurrentPlugin()
		}

	case "s", "ctrl+s":
		// Save
		if m.config != nil && m.HasChanges() {
			return m, SaveConfigCmd(m.config)
		}

	case "r":
		// Reset to defaults
		m.confirmReset = true

	case "L":
		// Toggle language
		m.ToggleLanguage()

	case "1":
		// Switch to Options tab
		m.activeTab = OptionsTab

	case "2":
		// Stay on Plugins tab
		m.activeTab = PluginsTab

	case "o":
		// Switch to Options tab
		m.activeTab = OptionsTab

	case "?":
		// Show TPM install instructions
		m.showTPMInstall = true
	}

	return m, nil
}

func (m ConfigEditorModel) toggleCurrentPlugin() (ConfigEditorModel, tea.Cmd) {
	p := m.CurrentPlugin()
	if p == nil {
		return m, nil
	}

	// TPM must be enabled first if enabling other plugins
	if p.Repo != "tmux-plugins/tpm" && !m.IsPluginEnabled("tmux-plugins/tpm") {
		// Auto-enable TPM when enabling any plugin
		m.TogglePlugin("tmux-plugins/tpm")
	}

	m.TogglePlugin(p.Repo)

	// When disabling a plugin, also disable plugins that require it
	if !m.IsPluginEnabled(p.Repo) {
		for _, plugin := range m.plugins {
			for _, req := range plugin.Requires {
				if req == p.Repo && m.IsPluginEnabled(plugin.Repo) {
					m.TogglePlugin(plugin.Repo)
				}
			}
		}
	}

	// When enabling a plugin, also enable required plugins
	if m.IsPluginEnabled(p.Repo) {
		for _, req := range p.Requires {
			if !m.IsPluginEnabled(req) {
				m.TogglePlugin(req)
			}
		}
	}

	return m, nil
}

func (m ConfigEditorModel) startPluginSettingEditing() (ConfigEditorModel, tea.Cmd) {
	p := m.CurrentPlugin()
	if p == nil || !m.IsPluginEnabled(p.Repo) {
		return m, nil
	}

	setting := m.CurrentPluginSetting()
	if setting == nil {
		return m, nil
	}

	switch setting.Type {
	case config.TypeBool:
		// Toggle boolean setting
		current := m.GetPluginSettingValue(p.Repo, setting.Key)
		newValue := "on"
		if current == "on" {
			newValue = "off"
		}
		m.SetPluginSettingValue(p.Repo, setting.Key, newValue)

	case config.TypeChoice:
		m.choosing = true
		m.editingPluginSetting = setting
		m.editingPluginRepo = p.Repo
		currentVal := m.GetPluginSettingValue(p.Repo, setting.Key)
		m.choiceCursor = 0
		for i, choice := range setting.Choices {
			if choice == currentVal {
				m.choiceCursor = i
				break
			}
		}

	default:
		m.editing = true
		m.editingPluginSetting = setting
		m.editingPluginRepo = p.Repo
		m.editInput.SetValue(m.GetPluginSettingValue(p.Repo, setting.Key))
		m.editInput.Focus()
		return m, textinput.Blink
	}

	return m, nil
}

func (m ConfigEditorModel) handleTPMInstallDialog(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "enter", " ":
		m.showTPMInstall = false
	}
	return m, nil
}

func (m ConfigEditorModel) startEditing() (ConfigEditorModel, tea.Cmd) {
	opt := m.CurrentOption()
	if opt == nil {
		return m, nil
	}

	switch opt.Type {
	case config.TypeBool:
		return m.toggleBool()

	case config.TypeChoice:
		m.choosing = true
		// Find current selection
		currentVal := m.GetOptionValue(opt.Key)
		m.choiceCursor = 0
		for i, choice := range opt.Choices {
			if choice == currentVal {
				m.choiceCursor = i
				break
			}
		}

	case config.TypeNumber, config.TypeString, config.TypeStyle:
		m.editing = true
		m.editingOption = opt
		m.editInput.SetValue(m.GetOptionValue(opt.Key))
		m.editInput.Focus()
		return m, textinput.Blink
	}

	return m, nil
}

func (m ConfigEditorModel) toggleBool() (ConfigEditorModel, tea.Cmd) {
	opt := m.CurrentOption()
	if opt == nil || opt.Type != config.TypeBool {
		return m, nil
	}

	current := m.GetOptionValue(opt.Key)
	newValue := "on"
	if current == "on" {
		newValue = "off"
	}
	m.SetOptionValue(opt.Key, newValue)

	return m, nil
}

func (m ConfigEditorModel) handleChoiceSelection(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	// Handle plugin setting choice
	if m.editingPluginSetting != nil {
		choices := m.editingPluginSetting.Choices

		switch msg.String() {
		case "j", "down":
			if m.choiceCursor < len(choices)-1 {
				m.choiceCursor++
			}
		case "k", "up":
			if m.choiceCursor > 0 {
				m.choiceCursor--
			}
		case "enter", " ":
			if m.choiceCursor < len(choices) {
				m.SetPluginSettingValue(m.editingPluginRepo, m.editingPluginSetting.Key, choices[m.choiceCursor])
			}
			m.choosing = false
			m.editingPluginSetting = nil
			m.editingPluginRepo = ""
		case "esc", "q":
			m.choosing = false
			m.editingPluginSetting = nil
			m.editingPluginRepo = ""
		}
		return m, nil
	}

	// Handle regular option choice
	opt := m.CurrentOption()
	if opt == nil {
		m.choosing = false
		return m, nil
	}

	switch msg.String() {
	case "j", "down":
		if m.choiceCursor < len(opt.Choices)-1 {
			m.choiceCursor++
		}
	case "k", "up":
		if m.choiceCursor > 0 {
			m.choiceCursor--
		}
	case "enter", " ":
		if m.choiceCursor < len(opt.Choices) {
			m.SetOptionValue(opt.Key, opt.Choices[m.choiceCursor])
		}
		m.choosing = false
	case "esc", "q":
		m.choosing = false
	}

	return m, nil
}

func (m ConfigEditorModel) handleTextInput(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		value := m.editInput.Value()

		// Handle plugin setting
		if m.editingPluginSetting != nil {
			// Validate numbers for plugin settings
			if m.editingPluginSetting.Type == config.TypeNumber {
				_, err := strconv.Atoi(value)
				if err != nil {
					// Invalid number, keep editing
					return m, nil
				}
			}
			m.SetPluginSettingValue(m.editingPluginRepo, m.editingPluginSetting.Key, value)
			m.editing = false
			m.editingPluginSetting = nil
			m.editingPluginRepo = ""
			m.editInput.Blur()
			return m, nil
		}

		// Handle regular option
		if m.editingOption != nil {
			// Validate numbers
			if m.editingOption.Type == config.TypeNumber {
				num, err := strconv.Atoi(value)
				if err != nil {
					// Invalid number, keep editing
					return m, nil
				}
				if num < m.editingOption.Min {
					value = strconv.Itoa(m.editingOption.Min)
				} else if num > m.editingOption.Max {
					value = strconv.Itoa(m.editingOption.Max)
				}
			}

			m.SetOptionValue(m.editingOption.Key, value)
		}
		m.editing = false
		m.editingOption = nil
		m.editInput.Blur()
		return m, nil

	case tea.KeyEsc:
		m.editing = false
		m.editingOption = nil
		m.editingPluginSetting = nil
		m.editingPluginRepo = ""
		m.editInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.editInput, cmd = m.editInput.Update(msg)
	return m, cmd
}

func (m ConfigEditorModel) handleResetConfirm(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if m.config != nil {
			m.config.ResetToDefaults()
		}
		m.confirmReset = false
	case "n", "N", "esc", "q":
		m.confirmReset = false
	}
	return m, nil
}

func (m ConfigEditorModel) handleSaveConfirm(msg tea.KeyMsg) (ConfigEditorModel, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Save and exit
		m.confirmSave = false
		if m.config != nil {
			return m, tea.Batch(SaveConfigCmd(m.config), ExitConfigEditorCmd())
		}
		return m, ExitConfigEditorCmd()
	case "n", "N":
		// Exit without saving
		m.confirmSave = false
		return m, ExitConfigEditorCmd()
	case "esc", "q", "c":
		// Cancel, stay in editor
		m.confirmSave = false
	}
	return m, nil
}
