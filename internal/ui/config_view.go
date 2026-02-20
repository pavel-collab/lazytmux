package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lazytmux/internal/config"
)

// Config editor styles
var (
	configTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("62")).
				Bold(true).
				MarginBottom(1)

	configCategoryFocusedStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("62")).
					Padding(0, 1)

	configCategoryUnfocusedStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("240")).
					Padding(0, 1)

	configSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Bold(true)

	configNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	configValueOnStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true)

	configValueOffStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

	configModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214")).
				Bold(true)

	configDescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Italic(true)

	configHelpBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)

	configDialogStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2)

	configWarningDialogStyle = lipgloss.NewStyle().
					Border(lipgloss.DoubleBorder()).
					BorderForeground(lipgloss.Color("196")).
					Padding(1, 2)
)

// View renders the config editor
func (m ConfigEditorModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Handle dialogs first
	if m.choosing {
		return m.renderChoiceDialog()
	}
	if m.confirmReset {
		return m.renderResetConfirmDialog()
	}
	if m.confirmSave {
		return m.renderSaveConfirmDialog()
	}
	if m.editing {
		return m.renderEditDialog()
	}

	// Title
	title := configTitleStyle.Render("tmux Configuration Editor")

	// Calculate dimensions
	categoriesWidth := m.width / 4
	if categoriesWidth < 20 {
		categoriesWidth = 20
	}
	optionsWidth := m.width - categoriesWidth - 4
	contentHeight := m.height - 7 // title + description + help + margins

	// Render panels
	categoriesPanel := m.renderCategories(categoriesWidth, contentHeight)
	optionsPanel := m.renderOptions(optionsWidth, contentHeight)

	// Join horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		categoriesPanel,
		optionsPanel,
	)

	// Description of current option
	description := m.renderDescription()

	// Help bar
	helpBar := m.renderHelpBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		mainContent,
		description,
		helpBar,
	)
}

func (m ConfigEditorModel) renderCategories(width, height int) string {
	style := configCategoryUnfocusedStyle
	if !m.focusOnOptions {
		style = configCategoryFocusedStyle
	}

	var content strings.Builder
	headerText := "Categories"
	if m.language == "ru" {
		headerText = "Категории"
	}
	content.WriteString(configTitleStyle.Render(headerText))
	content.WriteString("\n\n")

	for i, cat := range m.categories {
		name := cat.NameEN
		if m.language == "ru" {
			name = cat.NameRU
		}

		line := "  " + name
		if i == m.categoryCursor {
			line = "> " + name
			if !m.focusOnOptions {
				line = configSelectedStyle.Render(line)
			} else {
				line = lipgloss.NewStyle().Bold(true).Render(line)
			}
		}
		content.WriteString(line)
		content.WriteString("\n")
	}

	return style.
		Width(width).
		Height(height).
		Render(content.String())
}

func (m ConfigEditorModel) renderOptions(width, height int) string {
	style := configCategoryUnfocusedStyle
	if m.focusOnOptions {
		style = configCategoryFocusedStyle
	}

	cat := m.CurrentCategory()

	var content strings.Builder
	catName := cat.NameEN
	if m.language == "ru" {
		catName = cat.NameRU
	}
	content.WriteString(configTitleStyle.Render(catName))
	content.WriteString("\n\n")

	// Calculate column widths
	descWidth := width - 25 // leave space for value
	if descWidth < 20 {
		descWidth = 20
	}

	for i, opt := range cat.Options {
		line := m.renderOptionLine(opt, i == m.optionCursor && m.focusOnOptions, descWidth)
		content.WriteString(line)
		content.WriteString("\n")
	}

	return style.
		Width(width).
		Height(height).
		Render(content.String())
}

func (m ConfigEditorModel) renderOptionLine(opt config.Option, selected bool, descWidth int) string {
	// Get description
	desc := opt.DescEN
	if m.language == "ru" {
		desc = opt.DescRU
	}

	// Truncate description if too long
	if len(desc) > descWidth {
		desc = desc[:descWidth-3] + "..."
	}

	// Get value
	value := m.GetOptionValue(opt.Key)
	isModified := m.IsModified(opt.Key)

	// Format value based on type
	var valueStr string
	var valueStyle lipgloss.Style

	switch opt.Type {
	case config.TypeBool:
		if value == "on" {
			valueStr = "[ON]"
			valueStyle = configValueOnStyle
		} else {
			valueStr = "[OFF]"
			valueStyle = configValueOffStyle
		}
	default:
		valueStr = value
		valueStyle = configNormalStyle
	}

	// Mark modified values
	if isModified {
		valueStr = "*" + valueStr
		valueStyle = configModifiedStyle
	}

	// Build line
	padding := descWidth - len(desc)
	if padding < 1 {
		padding = 1
	}
	line := fmt.Sprintf("%-*s %s", descWidth, desc, valueStyle.Render(valueStr))

	if selected {
		return configSelectedStyle.Render("> " + line)
	}
	return "  " + line
}

func (m ConfigEditorModel) renderDescription() string {
	opt := m.CurrentOption()
	if opt == nil {
		return configDescriptionStyle.Width(m.width).Render("")
	}

	var desc strings.Builder
	desc.WriteString("Option: ")
	desc.WriteString(opt.Key)

	switch opt.Type {
	case config.TypeBool:
		desc.WriteString(" | Values: on, off")
	case config.TypeChoice:
		desc.WriteString(" | Values: ")
		desc.WriteString(strings.Join(opt.Choices, ", "))
	case config.TypeNumber:
		desc.WriteString(fmt.Sprintf(" | Range: %d-%d", opt.Min, opt.Max))
	}

	desc.WriteString(" | Default: ")
	desc.WriteString(opt.Default)

	return configDescriptionStyle.Width(m.width).Render(desc.String())
}

func (m ConfigEditorModel) renderHelpBar() string {
	var help string
	if m.language == "ru" {
		help = "j/k: навигация | Tab: панели | Space: вкл/выкл | Enter: изменить | s: сохранить | r: сброс | L: язык | Esc: выход"
	} else {
		help = "j/k: navigate | Tab: panels | Space: toggle | Enter: edit | s: save | r: reset | L: lang | Esc: exit"
	}

	if m.HasChanges() {
		marker := "* UNSAVED *"
		if m.language == "ru" {
			marker = "* НЕ СОХРАНЕНО *"
		}
		help = configModifiedStyle.Render(marker) + " | " + help
	}

	return configHelpBarStyle.Width(m.width).Render(help)
}

func (m ConfigEditorModel) renderChoiceDialog() string {
	opt := m.CurrentOption()
	if opt == nil {
		return ""
	}

	var content strings.Builder
	titleText := "Select value for: " + opt.Key
	if m.language == "ru" {
		titleText = "Выберите значение для: " + opt.Key
	}
	content.WriteString(configTitleStyle.Render(titleText))
	content.WriteString("\n\n")

	for i, choice := range opt.Choices {
		line := "  " + choice
		if i == m.choiceCursor {
			line = configSelectedStyle.Render("> " + choice)
		}
		content.WriteString(line)
		content.WriteString("\n")
	}

	content.WriteString("\n")
	helpText := "Enter: select | Esc: cancel"
	if m.language == "ru" {
		helpText = "Enter: выбрать | Esc: отмена"
	}
	content.WriteString(configDescriptionStyle.Render(helpText))

	dialog := configDialogStyle.Render(content.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

func (m ConfigEditorModel) renderEditDialog() string {
	opt := m.editingOption
	if opt == nil {
		return ""
	}

	var content strings.Builder
	titleText := "Edit: " + opt.Key
	if m.language == "ru" {
		titleText = "Редактирование: " + opt.Key
	}
	content.WriteString(configTitleStyle.Render(titleText))
	content.WriteString("\n\n")

	// Show input field
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
	content.WriteString(inputStyle.Render(m.editInput.View()))
	content.WriteString("\n\n")

	// Show constraints for numbers
	if opt.Type == config.TypeNumber {
		rangeText := fmt.Sprintf("Range: %d - %d", opt.Min, opt.Max)
		if m.language == "ru" {
			rangeText = fmt.Sprintf("Диапазон: %d - %d", opt.Min, opt.Max)
		}
		content.WriteString(configDescriptionStyle.Render(rangeText))
		content.WriteString("\n")
	}

	helpText := "Enter: confirm | Esc: cancel"
	if m.language == "ru" {
		helpText = "Enter: подтвердить | Esc: отмена"
	}
	content.WriteString(configDescriptionStyle.Render(helpText))

	dialog := configDialogStyle.Render(content.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

func (m ConfigEditorModel) renderResetConfirmDialog() string {
	var content strings.Builder

	titleText := "Reset Configuration"
	if m.language == "ru" {
		titleText = "Сброс конфигурации"
	}
	content.WriteString(configTitleStyle.Render(titleText))
	content.WriteString("\n\n")

	questionText := "Reset all settings to defaults?"
	if m.language == "ru" {
		questionText = "Сбросить все настройки?"
	}
	content.WriteString(questionText)
	content.WriteString("\n\n")

	helpText := "y: yes | n: no"
	if m.language == "ru" {
		helpText = "y: да | n: нет"
	}
	content.WriteString(configDescriptionStyle.Render(helpText))

	dialog := configWarningDialogStyle.Render(content.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

func (m ConfigEditorModel) renderSaveConfirmDialog() string {
	var content strings.Builder

	titleText := "Unsaved Changes"
	if m.language == "ru" {
		titleText = "Несохранённые изменения"
	}
	content.WriteString(configTitleStyle.Render(titleText))
	content.WriteString("\n\n")

	questionText := "Save changes before exit?"
	if m.language == "ru" {
		questionText = "Сохранить изменения перед выходом?"
	}
	content.WriteString(questionText)
	content.WriteString("\n\n")

	helpText := "y: save & exit | n: exit without saving | Esc: cancel"
	if m.language == "ru" {
		helpText = "y: сохранить | n: выйти без сохранения | Esc: отмена"
	}
	content.WriteString(configDescriptionStyle.Render(helpText))

	dialog := configWarningDialogStyle.Render(content.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}
