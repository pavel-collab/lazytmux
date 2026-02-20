package tmux

// ListSessions returns all tmux sessions
func (c *Client) ListSessions() ([]Session, error) {
	format := "#{session_name}:#{session_windows}:#{session_attached}:#{session_created}"
	output, err := c.Execute("list-sessions", "-F", format)
	if err != nil {
		if err == ErrNoServer {
			return []Session{}, nil
		}
		return nil, err
	}
	return ParseSessions(output)
}

// CreateSession creates a new tmux session
func (c *Client) CreateSession(name string) error {
	_, err := c.Execute("new-session", "-d", "-s", name)
	return err
}

// KillSession terminates a session
func (c *Client) KillSession(name string) error {
	_, err := c.Execute("kill-session", "-t", name)
	return err
}

// RenameSession renames a session
func (c *Client) RenameSession(oldName, newName string) error {
	_, err := c.Execute("rename-session", "-t", oldName, newName)
	return err
}

// AttachSession returns the command to attach (for exec)
func (c *Client) AttachSession(name string) []string {
	return []string{c.tmuxPath, "attach-session", "-t", name}
}

// SwitchClient switches the current client to another session
func (c *Client) SwitchClient(sessionName string) error {
	_, err := c.Execute("switch-client", "-t", sessionName)
	return err
}

// DetachClient detaches the current client from its session
func (c *Client) DetachClient() error {
	_, err := c.Execute("detach-client")
	return err
}
