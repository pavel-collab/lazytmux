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

// Tab styles
var (
	tabActiveStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
			MarginBottom(1)

	statusSuccessStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("42")).
				Foreground(lipgloss.Color("230")).
				Bold(true).
				Padding(0, 2)

	statusErrorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("196")).
				Foreground(lipgloss.Color("230")).
				Bold(true).
				Padding(0, 2)
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
	if m.showTPMInstall {
		return m.renderTPMInstallDialog()
	}

	// Title and tabs
	title := configTitleStyle.Render("tmux Configuration Editor")

	// Add status message next to title if present
	if m.statusMessage != "" {
		var statusStyle lipgloss.Style
		if m.statusIsError {
			statusStyle = statusErrorStyle
		} else {
			statusStyle = statusSuccessStyle
		}
		statusMsg := statusStyle.Render(m.statusMessage)
		title = lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", statusMsg)
	}

	tabs := m.renderTabs()

	// Calculate dimensions
	contentHeight := m.height - 9 // title + tabs + description + help + margins

	var mainContent string
	var description string

	switch m.activeTab {
	case PluginsTab:
		mainContent = m.renderPluginsContent(contentHeight)
		description = m.renderPluginDescription()
	default:
		categoriesWidth := m.width / 4
		if categoriesWidth < 20 {
			categoriesWidth = 20
		}
		optionsWidth := m.width - categoriesWidth - 4

		categoriesPanel := m.renderCategories(categoriesWidth, contentHeight)
		optionsPanel := m.renderOptions(optionsWidth, contentHeight)

		mainContent = lipgloss.JoinHorizontal(
			lipgloss.Top,
			categoriesPanel,
			optionsPanel,
		)
		description = m.renderDescription()
	}

	// Help bar
	helpBar := m.renderHelpBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		tabs,
		mainContent,
		description,
		helpBar,
	)
}

func (m ConfigEditorModel) renderTabs() string {
	optionsLabel := "1: Options"
	pluginsLabel := "2: Plugins"
	if m.language == "ru" {
		optionsLabel = "1: Настройки"
		pluginsLabel = "2: Плагины"
	}

	var optionsTab, pluginsTab string
	if m.activeTab == OptionsTab {
		optionsTab = tabActiveStyle.Render(optionsLabel)
		pluginsTab = tabInactiveStyle.Render(pluginsLabel)
	} else {
		optionsTab = tabInactiveStyle.Render(optionsLabel)
		pluginsTab = tabActiveStyle.Render(pluginsLabel)
	}

	return tabBarStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, optionsTab, "  ", pluginsTab))
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

	if m.activeTab == PluginsTab {
		if m.language == "ru" {
			help = "j/k: навигация | Tab: панели | Space: вкл/выкл | 1: настройки | ?: инструкция TPM | s: сохранить | L: язык | Esc: выход"
		} else {
			help = "j/k: navigate | Tab: panels | Space: toggle | 1: options | ?: TPM help | s: save | L: lang | Esc: exit"
		}
	} else {
		if m.language == "ru" {
			help = "j/k: навигация | Tab: панели | Space: вкл/выкл | Enter: изменить | 2/p: плагины | s: сохранить | r: сброс | L: язык | Esc: выход"
		} else {
			help = "j/k: navigate | Tab: panels | Space: toggle | Enter: edit | 2/p: plugins | s: save | r: reset | L: lang | Esc: exit"
		}
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

// Plugin rendering styles
var (
	pluginEnabledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true)

	pluginDisabledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

	pluginRequiredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

	pluginSettingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				MarginLeft(4)
)

func (m ConfigEditorModel) renderPluginsContent(height int) string {
	pluginListWidth := m.width / 3
	if pluginListWidth < 25 {
		pluginListWidth = 25
	}
	detailsWidth := m.width - pluginListWidth - 4

	pluginList := m.renderPluginList(pluginListWidth, height)
	pluginDetails := m.renderPluginDetails(detailsWidth, height)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		pluginList,
		pluginDetails,
	)
}

func (m ConfigEditorModel) renderPluginList(width, height int) string {
	style := configCategoryUnfocusedStyle
	if !m.focusOnPluginSettings {
		style = configCategoryFocusedStyle
	}

	var content strings.Builder
	headerText := "Plugins"
	if m.language == "ru" {
		headerText = "Плагины"
	}
	content.WriteString(configTitleStyle.Render(headerText))
	content.WriteString("\n\n")

	for i, plugin := range m.plugins {
		enabled := m.IsPluginEnabled(plugin.Repo)

		// Status indicator
		var status string
		if enabled {
			status = "[ON] "
		} else {
			status = "[OFF]"
		}

		// Plugin name
		line := fmt.Sprintf("%s %s", status, plugin.Name)

		// Styling
		if i == m.pluginCursor {
			if !m.focusOnPluginSettings {
				line = configSelectedStyle.Render("> " + line)
			} else {
				line = lipgloss.NewStyle().Bold(true).Render("> " + line)
			}
		} else if enabled {
			line = "  " + pluginEnabledStyle.Render(line)
		} else {
			line = "  " + pluginDisabledStyle.Render(line)
		}

		content.WriteString(line)
		content.WriteString("\n")
	}

	return style.
		Width(width).
		Height(height).
		Render(content.String())
}

func (m ConfigEditorModel) renderPluginDetails(width, height int) string {
	style := configCategoryUnfocusedStyle
	if m.focusOnPluginSettings {
		style = configCategoryFocusedStyle
	}

	p := m.CurrentPlugin()
	if p == nil {
		return style.Width(width).Height(height).Render("")
	}

	var content strings.Builder
	content.WriteString(configTitleStyle.Render(p.Name))
	content.WriteString("\n\n")

	// Description
	desc := p.DescEN
	if m.language == "ru" {
		desc = p.DescRU
	}
	content.WriteString(desc)
	content.WriteString("\n\n")

	// Repository
	repoLabel := "Repository: "
	if m.language == "ru" {
		repoLabel = "Репозиторий: "
	}
	content.WriteString(configDescriptionStyle.Render(repoLabel + p.Repo))
	content.WriteString("\n\n")

	// Key bindings
	keys := p.KeysEN
	if m.language == "ru" {
		keys = p.KeysRU
	}
	keysLabel := "Keys: "
	if m.language == "ru" {
		keysLabel = "Клавиши: "
	}
	content.WriteString(configDescriptionStyle.Render(keysLabel + keys))
	content.WriteString("\n")

	// Requirements
	if len(p.Requires) > 0 {
		reqLabel := "\nRequires: "
		if m.language == "ru" {
			reqLabel = "\nТребует: "
		}
		content.WriteString(pluginRequiredStyle.Render(reqLabel + strings.Join(p.Requires, ", ")))
		content.WriteString("\n")
	}

	// Settings (if plugin is enabled and has settings)
	if m.IsPluginEnabled(p.Repo) && len(p.Settings) > 0 {
		settingsLabel := "\nSettings:"
		if m.language == "ru" {
			settingsLabel = "\nНастройки:"
		}
		content.WriteString(configTitleStyle.Render(settingsLabel))
		content.WriteString("\n")

		for i, setting := range p.Settings {
			settingDesc := setting.DescEN
			if m.language == "ru" {
				settingDesc = setting.DescRU
			}

			value := m.GetPluginSettingValue(p.Repo, setting.Key)

			// Format value
			var valueStr string
			switch setting.Type {
			case config.TypeBool:
				if value == "on" {
					valueStr = configValueOnStyle.Render("[ON]")
				} else {
					valueStr = configValueOffStyle.Render("[OFF]")
				}
			default:
				valueStr = value
			}

			line := fmt.Sprintf("%-30s %s", settingDesc, valueStr)

			if i == m.pluginSettingCursor && m.focusOnPluginSettings {
				line = configSelectedStyle.Render("> " + line)
			} else {
				line = "  " + line
			}

			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	return style.
		Width(width).
		Height(height).
		Render(content.String())
}

func (m ConfigEditorModel) renderPluginDescription() string {
	p := m.CurrentPlugin()
	if p == nil {
		return configDescriptionStyle.Width(m.width).Render("")
	}

	var desc strings.Builder
	desc.WriteString("Plugin: ")
	desc.WriteString(p.Repo)

	if p.RequiresTPM {
		reqText := " | Requires TPM"
		if m.language == "ru" {
			reqText = " | Требует TPM"
		}
		desc.WriteString(reqText)
	}

	return configDescriptionStyle.Width(m.width).Render(desc.String())
}

func (m ConfigEditorModel) renderTPMInstallDialog() string {
	var content strings.Builder

	titleText := "TPM Installation"
	if m.language == "ru" {
		titleText = "Установка TPM"
	}
	content.WriteString(configTitleStyle.Render(titleText))
	content.WriteString("\n\n")

	if m.language == "ru" {
		content.WriteString("Tmux Plugin Manager (TPM) необходим для работы плагинов.\n\n")
		content.WriteString("Для установки выполните:\n\n")
		content.WriteString("  git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm\n\n")
		content.WriteString("После сохранения конфигурации:\n")
		content.WriteString("  1. Перезагрузите tmux: tmux source ~/.tmux.conf\n")
		content.WriteString("  2. Установите плагины: prefix + I\n\n")
		content.WriteString("Клавиши TPM:\n")
		content.WriteString("  prefix + I     - установить плагины\n")
		content.WriteString("  prefix + U     - обновить плагины\n")
		content.WriteString("  prefix + alt+u - удалить плагины\n")
	} else {
		content.WriteString("Tmux Plugin Manager (TPM) is required for plugins.\n\n")
		content.WriteString("To install, run:\n\n")
		content.WriteString("  git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm\n\n")
		content.WriteString("After saving configuration:\n")
		content.WriteString("  1. Reload tmux: tmux source ~/.tmux.conf\n")
		content.WriteString("  2. Install plugins: prefix + I\n\n")
		content.WriteString("TPM key bindings:\n")
		content.WriteString("  prefix + I     - install plugins\n")
		content.WriteString("  prefix + U     - update plugins\n")
		content.WriteString("  prefix + alt+u - uninstall plugins\n")
	}

	content.WriteString("\n")
	helpText := "Press any key to close"
	if m.language == "ru" {
		helpText = "Нажмите любую клавишу для закрытия"
	}
	content.WriteString(configDescriptionStyle.Render(helpText))

	dialog := configDialogStyle.Width(60).Render(content.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}
