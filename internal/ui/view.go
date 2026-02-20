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
			content.WriteString(fmt.Sprintf("Session: %s\n", m.tmuxState.CurrentSession.Name))
		}
		// Draw pane layout sketch
		if len(m.currentPanes) > 0 {
			content.WriteString("\nLayout:\n")
			sketch := m.renderPaneSketch(width-6, 10) // account for padding/borders
			content.WriteString(sketch)
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

// renderPaneSketch draws an ASCII art representation of the pane layout
func (m Model) renderPaneSketch(sketchWidth, sketchHeight int) string {
	if len(m.currentPanes) == 0 {
		return ""
	}

	// Find total window dimensions from panes
	var totalWidth, totalHeight int
	for _, p := range m.currentPanes {
		right := p.Left + p.Width
		bottom := p.Top + p.Height
		if right > totalWidth {
			totalWidth = right
		}
		if bottom > totalHeight {
			totalHeight = bottom
		}
	}

	if totalWidth == 0 || totalHeight == 0 {
		return ""
	}

	// Ensure minimum sketch size
	if sketchWidth < 10 {
		sketchWidth = 10
	}
	if sketchHeight < 5 {
		sketchHeight = 5
	}

	// Create grid for the sketch
	grid := make([][]rune, sketchHeight)
	for i := range grid {
		grid[i] = make([]rune, sketchWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Draw each pane
	for _, pane := range m.currentPanes {
		// Normalize coordinates to sketch size
		x1 := pane.Left * (sketchWidth - 1) / totalWidth
		y1 := pane.Top * (sketchHeight - 1) / totalHeight
		x2 := (pane.Left + pane.Width) * (sketchWidth - 1) / totalWidth
		y2 := (pane.Top + pane.Height) * (sketchHeight - 1) / totalHeight

		// Ensure minimum size
		if x2 <= x1 {
			x2 = x1 + 1
		}
		if y2 <= y1 {
			y2 = y1 + 1
		}

		// Clamp to grid bounds
		if x2 >= sketchWidth {
			x2 = sketchWidth - 1
		}
		if y2 >= sketchHeight {
			y2 = sketchHeight - 1
		}

		// Draw borders
		for x := x1; x <= x2; x++ {
			if y1 < sketchHeight {
				grid[y1][x] = '─'
			}
			if y2 < sketchHeight {
				grid[y2][x] = '─'
			}
		}
		for y := y1; y <= y2; y++ {
			if x1 < sketchWidth {
				grid[y][x1] = '│'
			}
			if x2 < sketchWidth {
				grid[y][x2] = '│'
			}
		}

		// Draw corners
		if y1 < sketchHeight && x1 < sketchWidth {
			grid[y1][x1] = '┌'
		}
		if y1 < sketchHeight && x2 < sketchWidth {
			grid[y1][x2] = '┐'
		}
		if y2 < sketchHeight && x1 < sketchWidth {
			grid[y2][x1] = '└'
		}
		if y2 < sketchHeight && x2 < sketchWidth {
			grid[y2][x2] = '┘'
		}

		// Draw pane index in the center
		centerX := (x1 + x2) / 2
		centerY := (y1 + y2) / 2
		if centerY > y1 && centerY < y2 && centerX > x1 && centerX < x2 {
			label := fmt.Sprintf("%d", pane.Index)
			if pane.Active {
				label = "*" + label
			}
			for i, ch := range label {
				if centerX+i < x2 && centerX+i < sketchWidth {
					grid[centerY][centerX+i] = ch
				}
			}
		}
	}

	// Fix intersections (where pane borders meet)
	for y := 0; y < sketchHeight; y++ {
		for x := 0; x < sketchWidth; x++ {
			// Check for intersection points and fix them
			hasTop := y > 0 && (grid[y-1][x] == '│' || grid[y-1][x] == '┌' || grid[y-1][x] == '┐' || grid[y-1][x] == '├' || grid[y-1][x] == '┤' || grid[y-1][x] == '┬' || grid[y-1][x] == '┼')
			hasBottom := y < sketchHeight-1 && (grid[y+1][x] == '│' || grid[y+1][x] == '└' || grid[y+1][x] == '┘' || grid[y+1][x] == '├' || grid[y+1][x] == '┤' || grid[y+1][x] == '┴' || grid[y+1][x] == '┼')
			hasLeft := x > 0 && (grid[y][x-1] == '─' || grid[y][x-1] == '┌' || grid[y][x-1] == '└' || grid[y][x-1] == '┬' || grid[y][x-1] == '┴' || grid[y][x-1] == '├' || grid[y][x-1] == '┼')
			hasRight := x < sketchWidth-1 && (grid[y][x+1] == '─' || grid[y][x+1] == '┐' || grid[y][x+1] == '┘' || grid[y][x+1] == '┬' || grid[y][x+1] == '┴' || grid[y][x+1] == '┤' || grid[y][x+1] == '┼')

			current := grid[y][x]
			if current == '┌' || current == '┐' || current == '└' || current == '┘' || current == '─' || current == '│' {
				if hasTop && hasBottom && hasLeft && hasRight {
					grid[y][x] = '┼'
				} else if hasTop && hasBottom && hasRight {
					grid[y][x] = '├'
				} else if hasTop && hasBottom && hasLeft {
					grid[y][x] = '┤'
				} else if hasLeft && hasRight && hasBottom {
					grid[y][x] = '┬'
				} else if hasLeft && hasRight && hasTop {
					grid[y][x] = '┴'
				}
			}
		}
	}

	// Convert grid to string
	var result strings.Builder
	for _, row := range grid {
		result.WriteString(string(row))
		result.WriteString("\n")
	}

	return result.String()
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
	return m.styles.HelpText.Render("j/k: navigate • Tab: switch panel • n: new • d: delete • v/s: split • a: attach • q: quit • ?: help")
}

func (m Model) renderWithDialog() string {
	// Render dialog content
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

	// Use lipgloss.Place to properly center the dialog
	// This avoids ANSI escape code issues that cause visual artifacts
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
		lipgloss.WithWhitespaceChars(" "),
	)
}

