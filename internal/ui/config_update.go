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

	// Normal navigation
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
		if m.editingOption != nil {
			value := m.editInput.Value()

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
