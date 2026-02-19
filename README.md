# LazyTmux

TUI-приложение для удобного управления tmux-сессиями, написанное на Go с использованием [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Возможности

- Просмотр всех tmux-сессий и окон
- Создание и удаление сессий
- Создание и удаление окон
- Подключение к сессиям
- Навигация в стиле vim (hjkl)

## Требования

- Go 1.21+
- tmux

## Установка и запуск

```bash
# Клонировать репозиторий
git clone <repo-url>
cd lazytmux

# Скачать зависимости
go mod tidy

# Собрать
go build -o lazytmux .

# Запустить
./lazytmux
```

### Установка в систему

```bash
go install .
```

После этого команда `lazytmux` будет доступна глобально (если `$GOPATH/bin` в PATH).

## Управление

| Клавиша | Действие |
|---------|----------|
| `j` / `↓` | Вниз |
| `k` / `↑` | Вверх |
| `h` / `←` | Предыдущая панель |
| `l` / `→` | Следующая панель |
| `Tab` | Переключить панель |
| `n` | Создать сессию/окно |
| `d` | Удалить сессию/окно |
| `a` | Подключиться к сессии |
| `R` | Обновить |
| `?` | Справка |
| `q` | Выход |

---

## For Developers

### Структура проекта

```
lazytmux/
├── main.go                 # Точка входа
├── go.mod                  # Зависимости Go
├── go.sum                  # Контрольные суммы зависимостей
└── internal/
    ├── tmux/               # Работа с tmux
    │   ├── client.go       # Клиент для выполнения команд tmux
    │   ├── types.go        # Структуры данных (Session, Window, TmuxState)
    │   ├── session.go      # Операции с сессиями
    │   ├── window.go       # Операции с окнами
    │   └── parser.go       # Парсинг вывода tmux
    └── ui/                  # Пользовательский интерфейс (Bubble Tea)
        ├── model.go        # Главная модель приложения
        ├── update.go       # Обработка событий и сообщений
        ├── view.go         # Рендеринг UI
        ├── keys.go         # Определение горячих клавиш
        ├── styles.go       # Стили и цвета (lipgloss)
        └── messages.go     # Типы сообщений и команды
```

### Описание файлов

#### `main.go`
Точка входа в приложение. Инициализирует tmux-клиент, создаёт Bubble Tea программу и запускает её. После завершения проверяет, нужно ли подключиться к сессии (через `syscall.Exec`).

#### `internal/tmux/`

| Файл | Описание |
|------|----------|
| `client.go` | Обёртка над командой `tmux`. Содержит `Client` с методом `Execute()` для выполнения произвольных tmux-команд. Обрабатывает ошибки (например, `ErrNoServer`). |
| `types.go` | Определяет структуры данных: `Session` (имя, кол-во окон, attached), `Window` (индекс, имя, panes), `TmuxState` (полное состояние). |
| `session.go` | Методы для работы с сессиями: `ListSessions()`, `CreateSession()`, `KillSession()`, `RenameSession()`, `AttachSession()`. |
| `window.go` | Методы для работы с окнами: `ListWindows()`, `CreateWindow()`, `KillWindow()`, `SelectWindow()`, `RenameWindow()`. |
| `parser.go` | Парсинг текстового вывода команд `tmux list-sessions` и `tmux list-windows` в структуры Go. |

#### `internal/ui/`

| Файл | Описание |
|------|----------|
| `model.go` | Главная модель Bubble Tea (`Model`). Содержит состояние UI: текущий курсор, активная панель, состояние диалогов, tmux-клиент. Реализует `Init()` для первоначальной загрузки данных. |
| `update.go` | Метод `Update()` — обработка всех событий: нажатия клавиш, изменение размера окна, сообщения от tmux-операций. Содержит логику навигации, открытия диалогов, выполнения действий. |
| `view.go` | Метод `View()` — рендеринг интерфейса. Отрисовка трёх панелей (Sessions, Windows, Info), статусной строки, справки и модальных диалогов. Использует `lipgloss` для стилизации. |
| `keys.go` | Определение горячих клавиш (`KeyMap`). Использует `bubbles/key` для биндингов. Реализует `ShortHelp()` и `FullHelp()` для отображения справки. |
| `styles.go` | Стили UI (`Styles`): цвета панелей, выделенных элементов, текста ошибок, диалогов. Использует `lipgloss`. |
| `messages.go` | Типы сообщений для Bubble Tea: `TmuxStateMsg`, `SessionCreatedMsg`, `ErrorMsg` и т.д. Также содержит команды (`Cmd`) для асинхронных операций: `RefreshCmd`, `CreateSessionCmd`, `DeleteSessionCmd`. |

### Архитектура

Приложение построено на [Bubble Tea](https://github.com/charmbracelet/bubbletea) — фреймворке для TUI на Go, использующем паттерн Elm Architecture:

1. **Model** — состояние приложения
2. **Update** — обработка событий, возврат нового состояния
3. **View** — рендеринг состояния в строку

Взаимодействие с tmux происходит через пакет `internal/tmux`, который выполняет команды через `os/exec` и парсит их вывод.

### Зависимости

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI фреймворк
- [bubbles](https://github.com/charmbracelet/bubbles) — готовые компоненты (help, textinput, key)
- [lipgloss](https://github.com/charmbracelet/lipgloss) — стилизация терминального вывода
