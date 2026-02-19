package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrNoServer = errors.New("no tmux server running")

// Client wraps tmux command execution
type Client struct {
	tmuxPath string
}

// NewClient creates a new tmux client
func NewClient() (*Client, error) {
	path, err := exec.LookPath("tmux")
	if err != nil {
		return nil, fmt.Errorf("tmux not found in PATH: %w", err)
	}
	return &Client{tmuxPath: path}, nil
}

// Execute runs a tmux command and returns output
func (c *Client) Execute(args ...string) (string, error) {
	cmd := exec.Command(c.tmuxPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := stderr.String()
		if strings.Contains(errStr, "no server running") ||
			strings.Contains(errStr, "no current client") {
			return "", ErrNoServer
		}
		return "", fmt.Errorf("tmux error: %s", strings.TrimSpace(errStr))
	}
	return stdout.String(), nil
}

// IsServerRunning checks if tmux server is running
func (c *Client) IsServerRunning() bool {
	_, err := c.Execute("list-sessions")
	return err != ErrNoServer && err == nil
}
