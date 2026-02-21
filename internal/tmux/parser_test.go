package tmux

import (
	"testing"
	"time"
)

func TestParseSessions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Session
		wantErr  bool
	}{
		{
			name:     "empty input",
			input:    "",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "whitespace only",
			input:    "   \n\t\n   ",
			expected: nil,
			wantErr:  false,
		},
		{
			name:  "single session not attached",
			input: "main:3:0:1700000000",
			expected: []Session{
				{
					Name:     "main",
					Windows:  3,
					Attached: false,
					Created:  time.Unix(1700000000, 0),
				},
			},
			wantErr: false,
		},
		{
			name:  "single session attached",
			input: "dev:5:1:1700000000",
			expected: []Session{
				{
					Name:     "dev",
					Windows:  5,
					Attached: true,
					Created:  time.Unix(1700000000, 0),
				},
			},
			wantErr: false,
		},
		{
			name: "multiple sessions",
			input: `main:3:0:1700000000
dev:5:1:1700001000
work:2:0:1700002000`,
			expected: []Session{
				{Name: "main", Windows: 3, Attached: false, Created: time.Unix(1700000000, 0)},
				{Name: "dev", Windows: 5, Attached: true, Created: time.Unix(1700001000, 0)},
				{Name: "work", Windows: 2, Attached: false, Created: time.Unix(1700002000, 0)},
			},
			wantErr: false,
		},
		{
			name:  "session name with special characters",
			input: "my-session_123:1:0:1700000000",
			expected: []Session{
				{Name: "my-session_123", Windows: 1, Attached: false, Created: time.Unix(1700000000, 0)},
			},
			wantErr: false,
		},
		{
			name:     "malformed line - too few parts",
			input:    "main:3:0",
			expected: nil,
			wantErr:  false,
		},
		{
			name: "mixed valid and invalid lines",
			input: `main:3:0:1700000000
invalid:line
dev:2:1:1700001000`,
			expected: []Session{
				{Name: "main", Windows: 3, Attached: false, Created: time.Unix(1700000000, 0)},
				{Name: "dev", Windows: 2, Attached: true, Created: time.Unix(1700001000, 0)},
			},
			wantErr: false,
		},
		{
			name: "empty lines between sessions",
			input: `main:3:0:1700000000

dev:2:1:1700001000`,
			expected: []Session{
				{Name: "main", Windows: 3, Attached: false, Created: time.Unix(1700000000, 0)},
				{Name: "dev", Windows: 2, Attached: true, Created: time.Unix(1700001000, 0)},
			},
			wantErr: false,
		},
		{
			name:  "zero windows",
			input: "empty:0:0:1700000000",
			expected: []Session{
				{Name: "empty", Windows: 0, Attached: false, Created: time.Unix(1700000000, 0)},
			},
			wantErr: false,
		},
		{
			name:  "invalid window count - defaults to 0",
			input: "test:abc:0:1700000000",
			expected: []Session{
				{Name: "test", Windows: 0, Attached: false, Created: time.Unix(1700000000, 0)},
			},
			wantErr: false,
		},
		{
			name:  "extra colons in data - parses first 4 parts only",
			input: "name:with:colons:5:0:1700000000",
			// Parser takes first 4 parts: name=name, windows=with (fails->0), attached=colons (!=1->false), created=5
			expected: []Session{
				{Name: "name", Windows: 0, Attached: false, Created: time.Unix(5, 0)},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions, err := ParseSessions(tt.input)

			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(sessions) != len(tt.expected) {
				t.Errorf("got %d sessions, expected %d", len(sessions), len(tt.expected))
				return
			}

			for i, s := range sessions {
				exp := tt.expected[i]
				if s.Name != exp.Name {
					t.Errorf("session[%d].Name = %q, expected %q", i, s.Name, exp.Name)
				}
				if s.Windows != exp.Windows {
					t.Errorf("session[%d].Windows = %d, expected %d", i, s.Windows, exp.Windows)
				}
				if s.Attached != exp.Attached {
					t.Errorf("session[%d].Attached = %v, expected %v", i, s.Attached, exp.Attached)
				}
				if !s.Created.Equal(exp.Created) {
					t.Errorf("session[%d].Created = %v, expected %v", i, s.Created, exp.Created)
				}
			}
		})
	}
}

func TestParseWindows(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		input       string
		expected    []Window
		wantErr     bool
	}{
		{
			name:        "empty input",
			sessionName: "main",
			input:       "",
			expected:    nil,
			wantErr:     false,
		},
		{
			name:        "single window not active",
			sessionName: "main",
			input:       "@1:0:bash:0:1",
			expected: []Window{
				{ID: "@1", Index: 0, Name: "bash", SessionName: "main", Active: false, Panes: 1},
			},
			wantErr: false,
		},
		{
			name:        "single window active",
			sessionName: "dev",
			input:       "@5:2:vim:1:3",
			expected: []Window{
				{ID: "@5", Index: 2, Name: "vim", SessionName: "dev", Active: true, Panes: 3},
			},
			wantErr: false,
		},
		{
			name:        "multiple windows",
			sessionName: "work",
			input: `@1:0:editor:1:1
@2:1:terminal:0:2
@3:2:logs:0:1`,
			expected: []Window{
				{ID: "@1", Index: 0, Name: "editor", SessionName: "work", Active: true, Panes: 1},
				{ID: "@2", Index: 1, Name: "terminal", SessionName: "work", Active: false, Panes: 2},
				{ID: "@3", Index: 2, Name: "logs", SessionName: "work", Active: false, Panes: 1},
			},
			wantErr: false,
		},
		{
			name:        "window name with special chars",
			sessionName: "main",
			input:       "@1:0:my-window_123:0:1",
			expected: []Window{
				{ID: "@1", Index: 0, Name: "my-window_123", SessionName: "main", Active: false, Panes: 1},
			},
			wantErr: false,
		},
		{
			name:        "malformed line - too few parts",
			sessionName: "main",
			input:       "@1:0:bash:0",
			expected:    nil,
			wantErr:     false,
		},
		{
			name:        "empty lines skipped",
			sessionName: "main",
			input: `@1:0:bash:1:1

@2:1:vim:0:1`,
			expected: []Window{
				{ID: "@1", Index: 0, Name: "bash", SessionName: "main", Active: true, Panes: 1},
				{ID: "@2", Index: 1, Name: "vim", SessionName: "main", Active: false, Panes: 1},
			},
			wantErr: false,
		},
		{
			name:        "high index numbers",
			sessionName: "main",
			input:       "@100:99:test:0:50",
			expected: []Window{
				{ID: "@100", Index: 99, Name: "test", SessionName: "main", Active: false, Panes: 50},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			windows, err := ParseWindows(tt.sessionName, tt.input)

			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(windows) != len(tt.expected) {
				t.Errorf("got %d windows, expected %d", len(windows), len(tt.expected))
				return
			}

			for i, w := range windows {
				exp := tt.expected[i]
				if w.ID != exp.ID {
					t.Errorf("window[%d].ID = %q, expected %q", i, w.ID, exp.ID)
				}
				if w.Index != exp.Index {
					t.Errorf("window[%d].Index = %d, expected %d", i, w.Index, exp.Index)
				}
				if w.Name != exp.Name {
					t.Errorf("window[%d].Name = %q, expected %q", i, w.Name, exp.Name)
				}
				if w.SessionName != exp.SessionName {
					t.Errorf("window[%d].SessionName = %q, expected %q", i, w.SessionName, exp.SessionName)
				}
				if w.Active != exp.Active {
					t.Errorf("window[%d].Active = %v, expected %v", i, w.Active, exp.Active)
				}
				if w.Panes != exp.Panes {
					t.Errorf("window[%d].Panes = %d, expected %d", i, w.Panes, exp.Panes)
				}
			}
		})
	}
}

func TestParsePanes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Pane
		wantErr  bool
	}{
		{
			name:     "empty input",
			input:    "",
			expected: nil,
			wantErr:  false,
		},
		{
			name:  "single pane not active",
			input: "%1:0:80:24:0:0:0",
			expected: []Pane{
				{ID: "%1", Index: 0, Width: 80, Height: 24, Left: 0, Top: 0, Active: false},
			},
			wantErr: false,
		},
		{
			name:  "single pane active",
			input: "%5:2:120:40:10:5:1",
			expected: []Pane{
				{ID: "%5", Index: 2, Width: 120, Height: 40, Left: 10, Top: 5, Active: true},
			},
			wantErr: false,
		},
		{
			name: "multiple panes - typical split layout",
			input: `%1:0:80:24:0:0:1
%2:1:80:24:81:0:0`,
			expected: []Pane{
				{ID: "%1", Index: 0, Width: 80, Height: 24, Left: 0, Top: 0, Active: true},
				{ID: "%2", Index: 1, Width: 80, Height: 24, Left: 81, Top: 0, Active: false},
			},
			wantErr: false,
		},
		{
			name: "horizontal split - panes stacked",
			input: `%1:0:160:12:0:0:1
%2:1:160:12:0:13:0`,
			expected: []Pane{
				{ID: "%1", Index: 0, Width: 160, Height: 12, Left: 0, Top: 0, Active: true},
				{ID: "%2", Index: 1, Width: 160, Height: 12, Left: 0, Top: 13, Active: false},
			},
			wantErr: false,
		},
		{
			name: "four pane grid",
			input: `%1:0:80:12:0:0:1
%2:1:80:12:81:0:0
%3:2:80:12:0:13:0
%4:3:80:12:81:13:0`,
			expected: []Pane{
				{ID: "%1", Index: 0, Width: 80, Height: 12, Left: 0, Top: 0, Active: true},
				{ID: "%2", Index: 1, Width: 80, Height: 12, Left: 81, Top: 0, Active: false},
				{ID: "%3", Index: 2, Width: 80, Height: 12, Left: 0, Top: 13, Active: false},
				{ID: "%4", Index: 3, Width: 80, Height: 12, Left: 81, Top: 13, Active: false},
			},
			wantErr: false,
		},
		{
			name:     "malformed line - too few parts",
			input:    "%1:0:80:24:0:0",
			expected: nil,
			wantErr:  false,
		},
		{
			name: "empty lines skipped",
			input: `%1:0:80:24:0:0:1

%2:1:80:24:81:0:0`,
			expected: []Pane{
				{ID: "%1", Index: 0, Width: 80, Height: 24, Left: 0, Top: 0, Active: true},
				{ID: "%2", Index: 1, Width: 80, Height: 24, Left: 81, Top: 0, Active: false},
			},
			wantErr: false,
		},
		{
			name:  "large dimensions",
			input: "%1:0:1920:1080:0:0:1",
			expected: []Pane{
				{ID: "%1", Index: 0, Width: 1920, Height: 1080, Left: 0, Top: 0, Active: true},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panes, err := ParsePanes(tt.input)

			if tt.wantErr && err == nil {
				t.Error("expected error but got nil")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(panes) != len(tt.expected) {
				t.Errorf("got %d panes, expected %d", len(panes), len(tt.expected))
				return
			}

			for i, p := range panes {
				exp := tt.expected[i]
				if p.ID != exp.ID {
					t.Errorf("pane[%d].ID = %q, expected %q", i, p.ID, exp.ID)
				}
				if p.Index != exp.Index {
					t.Errorf("pane[%d].Index = %d, expected %d", i, p.Index, exp.Index)
				}
				if p.Width != exp.Width {
					t.Errorf("pane[%d].Width = %d, expected %d", i, p.Width, exp.Width)
				}
				if p.Height != exp.Height {
					t.Errorf("pane[%d].Height = %d, expected %d", i, p.Height, exp.Height)
				}
				if p.Left != exp.Left {
					t.Errorf("pane[%d].Left = %d, expected %d", i, p.Left, exp.Left)
				}
				if p.Top != exp.Top {
					t.Errorf("pane[%d].Top = %d, expected %d", i, p.Top, exp.Top)
				}
				if p.Active != exp.Active {
					t.Errorf("pane[%d].Active = %v, expected %v", i, p.Active, exp.Active)
				}
			}
		})
	}
}

// Benchmark tests for parser performance
func BenchmarkParseSessions(b *testing.B) {
	input := `main:10:1:1700000000
dev:5:0:1700001000
work:3:0:1700002000
test:2:0:1700003000
staging:1:0:1700004000`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseSessions(input)
	}
}

func BenchmarkParseWindows(b *testing.B) {
	input := `@1:0:editor:1:1
@2:1:terminal:0:2
@3:2:logs:0:1
@4:3:build:0:1
@5:4:debug:0:3`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseWindows("main", input)
	}
}

func BenchmarkParsePanes(b *testing.B) {
	input := `%1:0:80:24:0:0:1
%2:1:80:24:81:0:0
%3:2:80:12:0:25:0
%4:3:80:12:81:25:0`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParsePanes(input)
	}
}
