package tmux

import (
	"strconv"
	"strings"
	"time"
)

// ParseSessions parses `list-sessions -F` output
func ParseSessions(output string) ([]Session, error) {
	var sessions []Session
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		windows, _ := strconv.Atoi(parts[1])
		attached := parts[2] == "1"
		created, _ := strconv.ParseInt(parts[3], 10, 64)

		sessions = append(sessions, Session{
			Name:     parts[0],
			Windows:  windows,
			Attached: attached,
			Created:  time.Unix(created, 0),
		})
	}
	return sessions, nil
}

// ParseWindows parses `list-windows -F` output
func ParseWindows(sessionName, output string) ([]Window, error) {
	var windows []Window
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 5 {
			continue
		}

		index, _ := strconv.Atoi(parts[1])
		active := parts[3] == "1"
		panes, _ := strconv.Atoi(parts[4])

		windows = append(windows, Window{
			ID:          parts[0],
			Index:       index,
			Name:        parts[2],
			SessionName: sessionName,
			Active:      active,
			Panes:       panes,
		})
	}
	return windows, nil
}
