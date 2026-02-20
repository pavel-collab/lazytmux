package ui

import (
	"lazytmux/internal/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ConfigEditorModel is the model for the configuration editor
type ConfigEditorModel struct {
	config     *config.Config
	categories []config.Category

	// Navigation
	categoryCursor int  // current category index
	optionCursor   int  // current option index in category
	focusOnOptions bool // true = focus on options list, false = on categories

	// Editing
	editing       bool
	editInput     textinput.Model
	editingOption *config.Option

	// Choice selection mode (for TypeChoice)
	choosing     bool
	choiceCursor int

	// Confirmation dialogs
	confirmReset bool
	confirmSave  bool

	// Dimensions
	width  int
	height int

	// Language (ru/en)
	language string
}

// NewConfigEditorModel creates a new config editor model
func NewConfigEditorModel() ConfigEditorModel {
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Width = 40

	return ConfigEditorModel{
		categories: config.GetCategories(),
		editInput:  ti,
		language:   "en",
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
}
