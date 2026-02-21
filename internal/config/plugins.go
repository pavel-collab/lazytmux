package config

// Plugin represents a tmux plugin
type Plugin struct {
	Name        string            // Short name for display
	Repo        string            // GitHub repo (e.g., "tmux-plugins/tmux-resurrect")
	DescEN      string            // Description in English
	DescRU      string            // Description in Russian
	KeysEN      string            // Key bindings description (EN)
	KeysRU      string            // Key bindings description (RU)
	Settings    []PluginSetting   // Plugin-specific settings
	RequiresTPM bool              // Whether this plugin requires TPM
	Requires    []string          // Other required plugins
}

// PluginSetting represents a plugin's configurable option
type PluginSetting struct {
	Key      string
	DescEN   string
	DescRU   string
	Type     OptionType
	Default  string
	Choices  []string
}

// PluginState represents the current state of a plugin
type PluginState struct {
	Repo      string
	Enabled   bool
	Installed bool // Whether the plugin directory exists
	Settings  map[string]string
}

// GetPlugins returns all available plugins
func GetPlugins() []Plugin {
	return []Plugin{
		// TPM - Plugin Manager (required first)
		{
			Name:   "TPM",
			Repo:   "tmux-plugins/tpm",
			DescEN: "Tmux Plugin Manager - required for other plugins",
			DescRU: "Менеджер плагинов tmux - необходим для других плагинов",
			KeysEN: "prefix + I: install plugins | prefix + U: update | prefix + alt + u: uninstall",
			KeysRU: "prefix + I: установить | prefix + U: обновить | prefix + alt + u: удалить",
			RequiresTPM: false,
		},

		// tmux-sensible - Sensible defaults
		{
			Name:   "Sensible",
			Repo:   "tmux-plugins/tmux-sensible",
			DescEN: "Sensible tmux defaults everyone can agree on",
			DescRU: "Разумные настройки tmux по умолчанию",
			KeysEN: "No keybindings, just settings",
			KeysRU: "Нет клавиш, только настройки",
			RequiresTPM: true,
		},

		// tmux-resurrect - Session save/restore
		{
			Name:   "Resurrect",
			Repo:   "tmux-plugins/tmux-resurrect",
			DescEN: "Save and restore tmux sessions",
			DescRU: "Сохранение и восстановление сессий tmux",
			KeysEN: "prefix + Ctrl-s: save | prefix + Ctrl-r: restore",
			KeysRU: "prefix + Ctrl-s: сохранить | prefix + Ctrl-r: восстановить",
			RequiresTPM: true,
			Settings: []PluginSetting{
				{
					Key:     "@resurrect-capture-pane-contents",
					DescEN:  "Capture pane contents",
					DescRU:  "Захватывать содержимое панелей",
					Type:    TypeBool,
					Default: "off",
				},
				{
					Key:     "@resurrect-strategy-vim",
					DescEN:  "Vim restore strategy",
					DescRU:  "Стратегия восстановления Vim",
					Type:    TypeChoice,
					Default: "session",
					Choices: []string{"session", "none"},
				},
			},
		},

		// tmux-continuum - Continuous auto-save
		{
			Name:   "Continuum",
			Repo:   "tmux-plugins/tmux-continuum",
			DescEN: "Automatic session save/restore every 15 min",
			DescRU: "Авто-сохранение/восстановление каждые 15 мин",
			KeysEN: "No keybindings, automatic save",
			KeysRU: "Нет клавиш, автоматическое сохранение",
			RequiresTPM: true,
			Requires: []string{"tmux-plugins/tmux-resurrect"},
			Settings: []PluginSetting{
				{
					Key:     "@continuum-restore",
					DescEN:  "Auto restore on tmux start",
					DescRU:  "Авто-восстановление при запуске",
					Type:    TypeBool,
					Default: "off",
				},
				{
					Key:     "@continuum-save-interval",
					DescEN:  "Auto-save interval (minutes, 0=disable)",
					DescRU:  "Интервал авто-сохранения (мин, 0=выкл)",
					Type:    TypeNumber,
					Default: "15",
				},
			},
		},

		// tmux-yank - System clipboard integration
		{
			Name:   "Yank",
			Repo:   "tmux-plugins/tmux-yank",
			DescEN: "Copy to system clipboard",
			DescRU: "Копирование в системный буфер обмена",
			KeysEN: "prefix + y: copy line | prefix + Y: copy pwd | in copy mode: y to copy",
			KeysRU: "prefix + y: копировать строку | prefix + Y: копировать путь | в режиме копирования: y",
			RequiresTPM: true,
		},

		// tmux-logging - Pane logging
		{
			Name:   "Logging",
			Repo:   "tmux-plugins/tmux-logging",
			DescEN: "Easy logging and screen capture for tmux",
			DescRU: "Логирование и захват экрана tmux",
			KeysEN: "prefix + Shift-p: toggle logging | prefix + Alt-p: screen capture | prefix + Alt-Shift-p: save history",
			KeysRU: "prefix + Shift-p: логирование | prefix + Alt-p: захват экрана | prefix + Alt-Shift-p: сохранить историю",
			RequiresTPM: true,
			Settings: []PluginSetting{
				{
					Key:     "@logging-path",
					DescEN:  "Log files directory",
					DescRU:  "Директория для логов",
					Type:    TypeString,
					Default: "#{pane_current_path}",
				},
				{
					Key:     "@screen-capture-path",
					DescEN:  "Screen capture directory",
					DescRU:  "Директория для захвата",
					Type:    TypeString,
					Default: "#{pane_current_path}",
				},
			},
		},

		// tmux-copycat - Regex search
		{
			Name:   "Copycat",
			Repo:   "tmux-plugins/tmux-copycat",
			DescEN: "Regex search and copy in tmux",
			DescRU: "Поиск по регулярным выражениям",
			KeysEN: "prefix + /: regex search | prefix + Ctrl-f: files | prefix + Ctrl-u: URLs",
			KeysRU: "prefix + /: поиск regex | prefix + Ctrl-f: файлы | prefix + Ctrl-u: URL",
			RequiresTPM: true,
		},

		// tmux-open - Open files and URLs
		{
			Name:   "Open",
			Repo:   "tmux-plugins/tmux-open",
			DescEN: "Open highlighted selection in browser/editor",
			DescRU: "Открыть выделенное в браузере/редакторе",
			KeysEN: "o: open | Ctrl-o: open in editor | Shift-s: search",
			KeysRU: "o: открыть | Ctrl-o: в редакторе | Shift-s: поиск",
			RequiresTPM: true,
		},

		// tmux-pain-control - Pane controls
		{
			Name:   "Pain Control",
			Repo:   "tmux-plugins/tmux-pain-control",
			DescEN: "Standard pane key-bindings for tmux",
			DescRU: "Стандартные клавиши для управления панелями",
			KeysEN: "prefix + h/j/k/l: navigate | prefix + H/J/K/L: resize | prefix + |/-: split",
			KeysRU: "prefix + h/j/k/l: навигация | prefix + H/J/K/L: размер | prefix + |/-: разделить",
			RequiresTPM: true,
		},

		// tmux-sessionist - Session management
		{
			Name:   "Sessionist",
			Repo:   "tmux-plugins/tmux-sessionist",
			DescEN: "Lightweight session management",
			DescRU: "Легкое управление сессиями",
			KeysEN: "prefix + g: switch | prefix + C: create | prefix + X: kill | prefix + @: join pane",
			KeysRU: "prefix + g: переключить | prefix + C: создать | prefix + X: удалить | prefix + @: присоединить",
			RequiresTPM: true,
		},

		// tmux-prefix-highlight - Prefix indicator
		{
			Name:   "Prefix Highlight",
			Repo:   "tmux-plugins/tmux-prefix-highlight",
			DescEN: "Highlight when prefix key is pressed",
			DescRU: "Подсветка при нажатии prefix",
			KeysEN: "No keybindings, status bar indicator",
			KeysRU: "Нет клавиш, индикатор в статус-баре",
			RequiresTPM: true,
		},

		// tmux-cpu - CPU/RAM info
		{
			Name:   "CPU",
			Repo:   "tmux-plugins/tmux-cpu",
			DescEN: "CPU and RAM info in status bar",
			DescRU: "Информация о CPU и RAM в статус-баре",
			KeysEN: "No keybindings, status bar integration",
			KeysRU: "Нет клавиш, интеграция в статус-бар",
			RequiresTPM: true,
		},

		// tmux-battery - Battery indicator
		{
			Name:   "Battery",
			Repo:   "tmux-plugins/tmux-battery",
			DescEN: "Battery percentage and icon in status bar",
			DescRU: "Индикатор заряда батареи в статус-баре",
			KeysEN: "No keybindings, status bar integration",
			KeysRU: "Нет клавиш, интеграция в статус-бар",
			RequiresTPM: true,
		},
	}
}

// GetPlugin returns a plugin by repo name
func GetPlugin(repo string) (Plugin, bool) {
	for _, p := range GetPlugins() {
		if p.Repo == repo {
			return p, true
		}
	}
	return Plugin{}, false
}

// TPMInstallPath returns the path where TPM should be installed
func TPMInstallPath() string {
	return "~/.tmux/plugins/tpm"
}

// PluginsDir returns the plugins directory
func PluginsDir() string {
	return "~/.tmux/plugins"
}
