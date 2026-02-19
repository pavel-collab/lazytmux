package tmux

import "fmt"

// ListWindows returns windows for a session
func (c *Client) ListWindows(sessionName string) ([]Window, error) {
	format := "#{window_id}:#{window_index}:#{window_name}:#{window_active}:#{window_panes}"
	output, err := c.Execute("list-windows", "-t", sessionName, "-F", format)
	if err != nil {
		return nil, err
	}
	return ParseWindows(sessionName, output)
}

// CreateWindow creates a new window in a session
func (c *Client) CreateWindow(sessionName, windowName string) error {
	args := []string{"new-window", "-t", sessionName}
	if windowName != "" {
		args = append(args, "-n", windowName)
	}
	_, err := c.Execute(args...)
	return err
}

// KillWindow deletes a window
func (c *Client) KillWindow(sessionName string, windowIndex int) error {
	target := fmt.Sprintf("%s:%d", sessionName, windowIndex)
	_, err := c.Execute("kill-window", "-t", target)
	return err
}

// SelectWindow switches to a window
func (c *Client) SelectWindow(sessionName string, windowIndex int) error {
	target := fmt.Sprintf("%s:%d", sessionName, windowIndex)
	_, err := c.Execute("select-window", "-t", target)
	return err
}

// RenameWindow renames a window
func (c *Client) RenameWindow(sessionName string, windowIndex int, newName string) error {
	target := fmt.Sprintf("%s:%d", sessionName, windowIndex)
	_, err := c.Execute("rename-window", "-t", target, newName)
	return err
}
