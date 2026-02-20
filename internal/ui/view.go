package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lazytmux/internal/tmux"
)

// View renders the entire UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Handle dialogs
	if m.activeDialog != NoDialog {
		return m.renderWithDialog()
	}

	return m.renderMainView()
}

func (m Model) renderMainView() string {
	// Calculate panel dimensions
	// Layout: Left column (1/3) | Right column (2/3)
	availableWidth := m.width - 4 // borders
	leftColumnWidth := availableWidth / 3
	rightColumnWidth := availableWidth - leftColumnWidth

	panelHeight := m.height - 4 // status bar + help

	// Left column: Sessions (1/2) and Windows (1/2) stacked vertically
	sessionsHeight := panelHeight / 2
	windowsHeight := panelHeight - sessionsHeight

	// Right column: Info (3/4) and Logs (1/4) stacked vertically
	infoHeight := panelHeight * 3 / 4
	logsHeight := panelHeight - infoHeight

	// Render left column panels
	sessionsPanel := m.renderSessionsPanel(leftColumnWidth, sessionsHeight)
	windowsPanel := m.renderWindowsPanel(leftColumnWidth, windowsHeight)

	// Stack left column vertically
	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		sessionsPanel,
		windowsPanel,
	)

	// Render right column panels
	infoPanel := m.renderInfoPanel(rightColumnWidth, infoHeight)
	logsPanel := m.renderLogsPanel(rightColumnWidth, logsHeight)

	// Stack right column vertically
	rightColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		infoPanel,
		logsPanel,
	)

	// Join columns horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		rightColumn,
	)

	// Render status bar
	statusBar := m.renderStatusBar()

	// Render help
	helpView := m.renderHelp()

	// Stack vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		statusBar,
		helpView,
	)
}

func (m Model) renderSessionsPanel(width, height int) string {
	style := m.styles.UnfocusedPanel
	if m.focusedPanel == SessionsPanel {
		style = m.styles.FocusedPanel
	}

	title := m.styles.Title.Render("Sessions")

	var content strings.Builder
	if len(m.tmuxState.Sessions) == 0 {
		content.WriteString(m.styles.DimText.Render("No sessions\nPress 'n' to create"))
	} else {
		for i, session := range m.tmuxState.Sessions {
			line := m.formatSessionLine(session)
			if i == m.sessionCursor && m.focusedPanel == SessionsPanel {
				content.WriteString(m.styles.SelectedItem.Render(line))
			} else if session.Attached {
				content.WriteString(m.styles.ActiveItem.Render(line))
			} else {
				content.WriteString(m.styles.NormalItem.Render(line))
			}
			content.WriteString("\n")
		}
	}

	innerContent := lipgloss.JoinVertical(lipgloss.Left, title, content.String())

	return style.
		Width(width).
		Height(height).
		Render(innerContent)
}

func (m Model) formatSessionLine(s tmux.Session) string {
	indicator := " "
	if s.Attached {
		indicator = "*"
	}
	return fmt.Sprintf("%s %s (%d)", indicator, s.Name, s.Windows)
}

func (m Model) renderWindowsPanel(width, height int) string {
	style := m.styles.UnfocusedPanel
	if m.focusedPanel == WindowsPanel {
		style = m.styles.FocusedPanel
	}

	title := m.styles.Title.Render("Windows")

	var content strings.Builder
	if m.tmuxState.CurrentSession == nil {
		content.WriteString(m.styles.DimText.Render("Select a session"))
	} else if len(m.tmuxState.Windows) == 0 {
		content.WriteString(m.styles.DimText.Render("No windows\nPress 'n' to create"))
	} else {
		for i, window := range m.tmuxState.Windows {
			line := m.formatWindowLine(window)
			if i == m.windowCursor && m.focusedPanel == WindowsPanel {
				content.WriteString(m.styles.SelectedItem.Render(line))
			} else if window.Active {
				content.WriteString(m.styles.ActiveItem.Render(line))
			} else {
				content.WriteString(m.styles.NormalItem.Render(line))
			}
			content.WriteString("\n")
		}
	}

	innerContent := lipgloss.JoinVertical(lipgloss.Left, title, content.String())

	return style.
		Width(width).
		Height(height).
		Render(innerContent)
}

func (m Model) formatWindowLine(w tmux.Window) string {
	indicator := " "
	if w.Active {
		indicator = "*"
	}
	return fmt.Sprintf("%s %d: %s (%d panes)", indicator, w.Index, w.Name, w.Panes)
}

func (m Model) renderInfoPanel(width, height int) string {
	style := m.styles.UnfocusedPanel

	title := m.styles.Title.Render("Info")

	var content strings.Builder

	if m.focusedPanel == SessionsPanel && len(m.tmuxState.Sessions) > 0 {
		s := m.tmuxState.Sessions[m.sessionCursor]
		content.WriteString(fmt.Sprintf("Session: %s\n", s.Name))
		content.WriteString(fmt.Sprintf("Windows: %d\n", s.Windows))
		content.WriteString(fmt.Sprintf("Attached: %v\n", s.Attached))
		content.WriteString(fmt.Sprintf("Created: %s", s.Created.Format("2006-01-02 15:04")))
	} else if m.focusedPanel == WindowsPanel && len(m.tmuxState.Windows) > 0 {
		w := m.tmuxState.Windows[m.windowCursor]
		content.WriteString(fmt.Sprintf("Window: %s\n", w.Name))
		content.WriteString(fmt.Sprintf("Index: %d\n", w.Index))
		content.WriteString(fmt.Sprintf("Panes: %d\n", w.Panes))
		content.WriteString(fmt.Sprintf("Active: %v\n", w.Active))
		if m.tmuxState.CurrentSession != nil {
			content.WriteString(fmt.Sprintf("Session: %s", m.tmuxState.CurrentSession.Name))
		}
	} else {
		content.WriteString(m.styles.DimText.Render("Select an item"))
	}

	innerContent := lipgloss.JoinVertical(lipgloss.Left, title, content.String())

	return style.
		Width(width).
		Height(height).
		Render(innerContent)
}

func (m Model) renderLogsPanel(width, height int) string {
	style := m.styles.UnfocusedPanel

	title := m.styles.Title.Render("Logs")

	content := m.styles.DimText.Render("Logs will be displayed here...")

	innerContent := lipgloss.JoinVertical(lipgloss.Left, title, content)

	return style.
		Width(width).
		Height(height).
		Render(innerContent)
}

func (m Model) renderStatusBar() string {
	var msg string
	if m.lastError != nil {
		msg = m.styles.ErrorText.Render("Error: " + m.lastError.Error())
	} else if m.statusMessage != "" {
		msg = m.statusMessage
	} else {
		if m.tmuxState.ServerRunning {
			msg = "tmux server running"
		} else {
			msg = "tmux server not running"
		}
	}

	return m.styles.StatusBar.Width(m.width).Render(msg)
}

func (m Model) renderHelp() string {
	if m.showHelp {
		return m.help.View(m.keyMap)
	}
	return m.styles.HelpText.Render("j/k: navigate • Tab: switch panel • n: new • d: delete • a: attach • q: quit • ?: help")
}

func (m Model) renderWithDialog() string {
	// Render main view dimmed
	mainView := m.renderMainView()

	// Render dialog
	var dialogContent string

	switch m.activeDialog {
	case CreateSessionDialog:
		dialogContent = m.styles.DialogTitle.Render("Create Session") + "\n\n" +
			m.styles.Input.Render(m.dialogInput.View()) + "\n\n" +
			m.styles.DimText.Render("Enter: confirm • Esc: cancel")

	case CreateWindowDialog:
		dialogContent = m.styles.DialogTitle.Render("Create Window") + "\n\n" +
			m.styles.Input.Render(m.dialogInput.View()) + "\n\n" +
			m.styles.DimText.Render("Enter: confirm • Esc: cancel")

	case ConfirmDeleteDialog:
		var target string
		if m.focusedPanel == SessionsPanel {
			target = fmt.Sprintf("session '%s'", m.deleteTarget)
		} else {
			target = fmt.Sprintf("window %d", m.deleteWindowIdx)
		}
		dialogContent = m.styles.DialogTitle.Render("Confirm Delete") + "\n\n" +
			fmt.Sprintf("Delete %s?", target) + "\n\n" +
			m.styles.DimText.Render("y: confirm • Esc: cancel")
	}

	dialog := m.styles.Dialog.Render(dialogContent)

	// Center the dialog
	dialogWidth := lipgloss.Width(dialog)
	dialogHeight := lipgloss.Height(dialog)

	x := (m.width - dialogWidth) / 2
	y := (m.height - dialogHeight) / 2

	// Overlay dialog on main view
	return placeOverlay(x, y, dialog, mainView)
}

// placeOverlay places an overlay on top of a background at the given position
func placeOverlay(x, y int, overlay, background string) string {
	bgLines := strings.Split(background, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, line := range overlayLines {
		bgY := y + i
		if bgY < 0 || bgY >= len(bgLines) {
			continue
		}

		bgLine := bgLines[bgY]
		bgRunes := []rune(bgLine)

		// Pad background line if needed
		for len(bgRunes) < x+len([]rune(line)) {
			bgRunes = append(bgRunes, ' ')
		}

		// Replace portion of background with overlay
		overlayRunes := []rune(line)
		for j, r := range overlayRunes {
			if x+j >= 0 && x+j < len(bgRunes) {
				bgRunes[x+j] = r
			}
		}

		bgLines[bgY] = string(bgRunes)
	}

	return strings.Join(bgLines, "\n")
}
