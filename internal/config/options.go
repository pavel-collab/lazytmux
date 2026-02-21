package config

// GetCategories returns all option categories
func GetCategories() []Category {
	return []Category{
		{
			ID:     "general",
			NameEN: "General",
			NameRU: "Общие",
			Options: []Option{
				{
					Key:     "default-terminal",
					DescEN:  "Default terminal type",
					DescRU:  "Тип терминала по умолчанию",
					Type:    TypeChoice,
					Default: "screen-256color",
					Choices: []string{"screen", "screen-256color", "tmux", "tmux-256color", "xterm-256color"},
					Scope:   ScopeServer,
				},
				{
					Key:     "escape-time",
					DescEN:  "Escape key delay (ms)",
					DescRU:  "Задержка клавиши Escape (мс)",
					Type:    TypeNumber,
					Default: "500",
					Min:     0,
					Max:     2000,
					Scope:   ScopeServer,
				},
				{
					Key:     "focus-events",
					DescEN:  "Pass focus events to applications",
					DescRU:  "Передавать события фокуса приложениям",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeServer,
				},
				{
					Key:     "history-limit",
					DescEN:  "Scrollback buffer size",
					DescRU:  "Размер буфера прокрутки",
					Type:    TypeNumber,
					Default: "2000",
					Min:     0,
					Max:     100000,
					Scope:   ScopeSession,
				},
				{
					Key:     "base-index",
					DescEN:  "Starting index for windows",
					DescRU:  "Начальный индекс окон",
					Type:    TypeNumber,
					Default: "0",
					Min:     0,
					Max:     99,
					Scope:   ScopeSession,
				},
				{
					Key:     "pane-base-index",
					DescEN:  "Starting index for panes",
					DescRU:  "Начальный индекс панелей",
					Type:    TypeNumber,
					Default: "0",
					Min:     0,
					Max:     99,
					Scope:   ScopeWindow,
				},
				{
					Key:     "renumber-windows",
					DescEN:  "Renumber windows on close",
					DescRU:  "Перенумеровывать окна при закрытии",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeSession,
				},
			},
		},
		{
			ID:     "statusbar",
			NameEN: "Status Bar",
			NameRU: "Статус-бар",
			Options: []Option{
				{
					Key:     "status",
					DescEN:  "Show status bar",
					DescRU:  "Показывать статус-бар",
					Type:    TypeBool,
					Default: "on",
					Scope:   ScopeSession,
				},
				{
					Key:     "status-position",
					DescEN:  "Status bar position",
					DescRU:  "Позиция статус-бара",
					Type:    TypeChoice,
					Default: "bottom",
					Choices: []string{"top", "bottom"},
					Scope:   ScopeSession,
				},
				{
					Key:     "status-interval",
					DescEN:  "Status refresh interval (sec)",
					DescRU:  "Интервал обновления (сек)",
					Type:    TypeNumber,
					Default: "15",
					Min:     0,
					Max:     600,
					Scope:   ScopeSession,
				},
				{
					Key:     "status-justify",
					DescEN:  "Window list alignment",
					DescRU:  "Выравнивание списка окон",
					Type:    TypeChoice,
					Default: "left",
					Choices: []string{"left", "centre", "right"},
					Scope:   ScopeSession,
				},
				{
					Key:     "status-left-length",
					DescEN:  "Max length of left status",
					DescRU:  "Макс. длина левой части",
					Type:    TypeNumber,
					Default: "10",
					Min:     0,
					Max:     200,
					Scope:   ScopeSession,
				},
				{
					Key:     "status-right-length",
					DescEN:  "Max length of right status",
					DescRU:  "Макс. длина правой части",
					Type:    TypeNumber,
					Default: "40",
					Min:     0,
					Max:     200,
					Scope:   ScopeSession,
				},
			},
		},
		{
			ID:     "mouse",
			NameEN: "Mouse",
			NameRU: "Мышь",
			Options: []Option{
				{
					Key:     "mouse",
					DescEN:  "Enable mouse support",
					DescRU:  "Включить поддержку мыши",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeSession,
				},
			},
		},
		{
			ID:     "visual",
			NameEN: "Visual",
			NameRU: "Визуальные",
			Options: []Option{
				{
					Key:     "monitor-activity",
					DescEN:  "Monitor window activity",
					DescRU:  "Отслеживать активность окон",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeWindow,
				},
				{
					Key:     "visual-activity",
					DescEN:  "Visual activity indicator",
					DescRU:  "Визуальный индикатор активности",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeSession,
				},
				{
					Key:     "monitor-bell",
					DescEN:  "Monitor bell signals",
					DescRU:  "Отслеживать звуковые сигналы",
					Type:    TypeBool,
					Default: "on",
					Scope:   ScopeWindow,
				},
				{
					Key:     "visual-bell",
					DescEN:  "Visual bell indicator",
					DescRU:  "Визуальный индикатор звонка",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeSession,
				},
				{
					Key:     "bell-action",
					DescEN:  "Bell notification mode",
					DescRU:  "Режим уведомления звуком",
					Type:    TypeChoice,
					Default: "any",
					Choices: []string{"none", "any", "current", "other"},
					Scope:   ScopeSession,
				},
				{
					Key:     "display-time",
					DescEN:  "Message display time (ms)",
					DescRU:  "Время показа сообщений (мс)",
					Type:    TypeNumber,
					Default: "750",
					Min:     0,
					Max:     10000,
					Scope:   ScopeSession,
				},
			},
		},
		{
			ID:     "keys",
			NameEN: "Keys",
			NameRU: "Клавиши",
			Options: []Option{
				{
					Key:     "mode-keys",
					DescEN:  "Copy mode key bindings",
					DescRU:  "Клавиши в режиме копирования",
					Type:    TypeChoice,
					Default: "emacs",
					Choices: []string{"vi", "emacs"},
					Scope:   ScopeSession,
				},
				{
					Key:     "status-keys",
					DescEN:  "Command line key bindings",
					DescRU:  "Клавиши командной строки",
					Type:    TypeChoice,
					Default: "emacs",
					Choices: []string{"vi", "emacs"},
					Scope:   ScopeSession,
				},
				{
					Key:     "repeat-time",
					DescEN:  "Key repeat timeout (ms)",
					DescRU:  "Таймаут повтора клавиш (мс)",
					Type:    TypeNumber,
					Default: "500",
					Min:     0,
					Max:     5000,
					Scope:   ScopeServer,
				},
			},
		},
		{
			ID:     "windows",
			NameEN: "Windows",
			NameRU: "Окна",
			Options: []Option{
				{
					Key:     "allow-rename",
					DescEN:  "Allow programs to rename windows",
					DescRU:  "Разрешить программам переименовывать окна",
					Type:    TypeBool,
					Default: "on",
					Scope:   ScopeWindow,
				},
				{
					Key:     "automatic-rename",
					DescEN:  "Automatically rename windows",
					DescRU:  "Автоматически переименовывать окна",
					Type:    TypeBool,
					Default: "on",
					Scope:   ScopeWindow,
				},
				{
					Key:     "aggressive-resize",
					DescEN:  "Resize windows aggressively",
					DescRU:  "Агрессивное изменение размера",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeWindow,
				},
				{
					Key:     "remain-on-exit",
					DescEN:  "Keep pane after command exits",
					DescRU:  "Сохранять панель после выхода",
					Type:    TypeBool,
					Default: "off",
					Scope:   ScopeWindow,
				},
			},
		},
	}
}

// GetAllOptions returns a flat map of all options by key
func GetAllOptions() map[string]Option {
	result := make(map[string]Option)
	for _, cat := range GetCategories() {
		for _, opt := range cat.Options {
			result[opt.Key] = opt
		}
	}
	return result
}

// GetOption returns a single option by key
func GetOption(key string) (Option, bool) {
	opts := GetAllOptions()
	opt, ok := opts[key]
	return opt, ok
}
