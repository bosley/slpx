package tui

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bosley/slpx/pkg/object"
	"github.com/bosley/slpx/pkg/repl"
	"github.com/bosley/slpx/pkg/slp"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type viewMode uint

const (
	normalMode viewMode = iota
	fullscreenInputMode
	scrollableOutputMode
)

type focusPane uint

const (
	editorFocus focusPane = iota
	historyFocus
)

var (
	promptStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
	resultStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	focusedStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	blurredStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	historyItemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dirtyPromptStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
)

type capturedIO struct {
	mu     sync.Mutex
	buffer strings.Builder
}

func (c *capturedIO) ReadLine() (string, error) {
	return "", fmt.Errorf("stdin not available in TUI mode")
}

func (c *capturedIO) ReadAll() ([]byte, error) {
	return nil, fmt.Errorf("stdin not available in TUI mode")
}

func (c *capturedIO) Write(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buffer.Write(data)
}

func (c *capturedIO) WriteString(s string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buffer.WriteString(s)
}

func (c *capturedIO) WriteError(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buffer.Write(data)
}

func (c *capturedIO) WriteErrorString(s string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buffer.WriteString(s)
}

func (c *capturedIO) Flush() error {
	return nil
}

func (c *capturedIO) SetStdin(r io.Reader) {
}

func (c *capturedIO) SetStdout(w io.Writer) {
}

func (c *capturedIO) SetStderr(w io.Writer) {
}

func (c *capturedIO) GetAndClear() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := c.buffer.String()
	c.buffer.Reset()
	return result
}

type model struct {
	logger           *slog.Logger
	session          *repl.Session
	capturedIO       *capturedIO
	viewport         viewport.Model
	commandHistory   []string
	outputHistory    []string
	textarea         textarea.Model
	mode             viewMode
	focus            focusPane
	historySelection int
	width            int
	height           int
	ready            bool
	showDirtyPrompt  bool
	pendingHistory   string
}

func initialModel(logger *slog.Logger) model {
	ta := textarea.New()
	ta.Placeholder = "Enter SLPX code..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 0
	ta.SetWidth(30)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to SLPX REPL!
Type SLPX code and press Enter to evaluate.
ctrl+e: editor/history • ctrl+o: scroll output`)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	sessionPath := filepath.Join(cwd, ".repl.slpx")

	capturedIO := &capturedIO{}
	session := repl.NewSessionBuilder(logger).
		WithIO(capturedIO).
		Build(sessionPath)

	return model{
		logger:           logger,
		session:          session,
		capturedIO:       capturedIO,
		textarea:         ta,
		commandHistory:   []string{},
		outputHistory:    []string{},
		viewport:         vp,
		mode:             normalMode,
		focus:            editorFocus,
		historySelection: 0,
		ready:            false,
		showDirtyPrompt:  false,
		pendingHistory:   "",
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) evaluateInput(input string) string {
	m.capturedIO.GetAndClear()

	result, err := m.session.Evaluate(input)

	capturedOutput := m.capturedIO.GetAndClear()

	var output strings.Builder

	if capturedOutput != "" {
		output.WriteString(capturedOutput)
		if !strings.HasSuffix(capturedOutput, "\n") {
			output.WriteString("\n")
		}
	}

	if err != nil {
		if parseErr, ok := err.(*slp.ParseError); ok {
			output.WriteString(errorStyle.Render(fmt.Sprintf("Parse Error: %s", parseErr.Message)))
		} else {
			output.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", err)))
		}
		return output.String()
	}

	if result.Type == object.OBJ_TYPE_ERROR {
		errObj := result.D.(object.Error)
		output.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", errObj.Message)))
		return output.String()
	}

	if result.Type != object.OBJ_TYPE_NONE {
		output.WriteString(resultStyle.Render(result.Encode()))
	}

	return output.String()
}

func (m *model) renderHistoryForViewport() string {
	return strings.Join(m.outputHistory, "\n")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		if m.mode == fullscreenInputMode {
			editorWidth := (msg.Width / 2) - 2
			m.textarea.SetWidth(editorWidth)
			m.textarea.SetHeight(msg.Height - 5)
		} else {
			m.viewport.Width = msg.Width
			m.textarea.SetWidth(msg.Width)
			m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap) - 2
		}

		content := m.renderHistoryForViewport()
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))

		if m.mode == normalMode {
			m.viewport.GotoBottom()
		}

	case tea.KeyMsg:
		switch m.mode {
		case normalMode:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "ctrl+e":
				m.mode = fullscreenInputMode
				m.focus = editorFocus
				m.historySelection = len(m.commandHistory) - 1
				if m.historySelection < 0 {
					m.historySelection = 0
				}
				m.showDirtyPrompt = false
				m.textarea.ShowLineNumbers = true
				m.textarea.KeyMap.InsertNewline.SetEnabled(true)
				editorWidth := (m.width / 2) - 2
				m.textarea.SetWidth(editorWidth)
				m.textarea.SetHeight(m.height - 5)
				return m, nil
			case "ctrl+o":
				m.mode = scrollableOutputMode
				m.textarea.Blur()
				m.viewport.Width = m.width - 2
				content := m.renderHistoryForViewport()
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
				return m, nil
			case "enter":
				value := m.textarea.Value()
				if value != "" {
					if value == "clear" || value == "cls" {
						m.outputHistory = []string{}
						m.viewport.SetContent("")
						m.textarea.Reset()
						return m, tea.WindowSize()
					}
					m.commandHistory = append(m.commandHistory, value)
					output := m.evaluateInput(value)
					m.outputHistory = append(m.outputHistory, promptStyle.Render("> ")+value)
					if output != "" {
						m.outputHistory = append(m.outputHistory, output)
					}
					content := m.renderHistoryForViewport()
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
					m.textarea.Reset()
					m.viewport.GotoBottom()
				}
				return m, nil
			}

		case fullscreenInputMode:
			if m.showDirtyPrompt {
				switch msg.String() {
				case "a":
					currentValue := m.textarea.Value()
					if currentValue != "" {
						m.textarea.SetValue(currentValue + "\n" + m.pendingHistory)
					} else {
						m.textarea.SetValue(m.pendingHistory)
					}
					m.showDirtyPrompt = false
					m.pendingHistory = ""
					m.focus = editorFocus
					cmd := m.textarea.Focus()
					return m, cmd
				case "o":
					m.textarea.SetValue(m.pendingHistory)
					m.showDirtyPrompt = false
					m.pendingHistory = ""
					m.focus = editorFocus
					cmd := m.textarea.Focus()
					return m, cmd
				case "c", "esc":
					m.showDirtyPrompt = false
					m.pendingHistory = ""
					m.focus = editorFocus
					cmd := m.textarea.Focus()
					return m, cmd
				}
				return m, nil
			}

			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl+e":
				m.mode = normalMode
				m.textarea.ShowLineNumbers = false
				m.textarea.KeyMap.InsertNewline.SetEnabled(false)
				m.textarea.SetWidth(m.width)
				m.textarea.SetHeight(3)
				m.viewport.Height = m.height - m.textarea.Height() - lipgloss.Height(gap) - 2
				cmd := m.textarea.Focus()
				return m, cmd
			case "esc":
				value := m.textarea.Value()
				if value != "" {
					if value == "clear" || value == "cls" {
						m.outputHistory = []string{}
						m.viewport.SetContent("")
						m.textarea.Reset()
						m.mode = normalMode
						m.textarea.ShowLineNumbers = false
						m.textarea.KeyMap.InsertNewline.SetEnabled(false)
						m.textarea.SetWidth(m.width)
						m.textarea.SetHeight(3)
						m.viewport.Height = m.height - m.textarea.Height() - lipgloss.Height(gap) - 2
						return m, tea.Batch(m.textarea.Focus(), tea.WindowSize())
					}
					m.commandHistory = append(m.commandHistory, value)
					output := m.evaluateInput(value)
					m.outputHistory = append(m.outputHistory, promptStyle.Render("> ")+value)
					if output != "" {
						m.outputHistory = append(m.outputHistory, output)
					}
					content := m.renderHistoryForViewport()
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
					m.textarea.Reset()
					m.viewport.GotoBottom()
				}
				m.mode = normalMode
				m.textarea.ShowLineNumbers = false
				m.textarea.KeyMap.InsertNewline.SetEnabled(false)
				m.textarea.SetWidth(m.width)
				m.textarea.SetHeight(3)
				m.viewport.Height = m.height - m.textarea.Height() - lipgloss.Height(gap) - 2
				cmd := m.textarea.Focus()
				return m, cmd
			case "tab":
				if m.focus == editorFocus {
					m.focus = historyFocus
					m.textarea.Blur()
					return m, nil
				} else {
					m.focus = editorFocus
					cmd := m.textarea.Focus()
					return m, cmd
				}
			case "up":
				if m.focus == historyFocus && len(m.commandHistory) > 0 {
					if m.historySelection > 0 {
						m.historySelection--
					}
					return m, nil
				}
			case "down":
				if m.focus == historyFocus && len(m.commandHistory) > 0 {
					if m.historySelection < len(m.commandHistory)-1 {
						m.historySelection++
					}
					return m, nil
				}
			case "enter":
				if m.focus == historyFocus && len(m.commandHistory) > 0 {
					selectedInput := m.commandHistory[m.historySelection]

					if m.textarea.Value() != "" {
						m.showDirtyPrompt = true
						m.pendingHistory = selectedInput
						return m, nil
					} else {
						m.textarea.SetValue(selectedInput)
						m.focus = editorFocus
						cmd := m.textarea.Focus()
						return m, cmd
					}
				}
			}

		case scrollableOutputMode:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc", "ctrl+o":
				m.mode = normalMode
				m.viewport.Width = m.width
				content := m.renderHistoryForViewport()
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
				cmd := m.textarea.Focus()
				m.viewport.GotoBottom()
				return m, cmd
			case "q":
				m.mode = normalMode
				m.viewport.Width = m.width
				content := m.renderHistoryForViewport()
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
				cmd := m.textarea.Focus()
				m.viewport.GotoBottom()
				return m, cmd
			}
		}
	}

	if m.mode == normalMode || m.mode == fullscreenInputMode {
		m.textarea, tiCmd = m.textarea.Update(msg)
	}

	if m.mode == normalMode {
		m.viewport, vpCmd = m.viewport.Update(msg)
	} else if m.mode == scrollableOutputMode {
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var helpText string
	switch m.mode {
	case normalMode:
		ctrlE := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("ctrl+e")
		ctrlO := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("ctrl+o")
		ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")
		esc := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("esc")

		helpText = helpStyle.Render(fmt.Sprintf("%s: editor/history • %s: scroll output • %s/%s: quit",
			ctrlE, ctrlO, ctrlC, esc))
		return fmt.Sprintf("%s%s%s\n%s", m.viewport.View(), gap, m.textarea.View(), helpText)
	case fullscreenInputMode:
		return m.renderFullscreenSplit()
	case scrollableOutputMode:
		arrows := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("↑/↓")
		escKey := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("esc")
		ctrlO := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("ctrl+o")
		qKey := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("q")
		ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")

		helpText = helpStyle.Render(fmt.Sprintf("%s: scroll • %s/%s/%s: back to REPL • %s: quit",
			arrows, escKey, ctrlO, qKey, ctrlC))
		return fmt.Sprintf("%s%s%s\n%s", focusedStyle.Render(m.viewport.View()), gap, m.textarea.View(), helpText)
	}

	return ""
}

func (m model) renderFullscreenSplit() string {
	editorWidth := (m.width / 2) - 2
	historyWidth := m.width - editorWidth - 4

	editorStyle := blurredStyle
	historyStyle := blurredStyle
	if m.focus == editorFocus {
		editorStyle = focusedStyle
	} else {
		historyStyle = focusedStyle
	}

	contentHeight := m.height - 3

	editorView := editorStyle.Width(editorWidth).Height(contentHeight).Render(m.textarea.View())
	historyView := m.renderHistory(historyWidth, contentHeight)
	historyView = historyStyle.Width(historyWidth).Height(contentHeight).Render(historyView)

	splitView := lipgloss.JoinHorizontal(lipgloss.Top, editorView, historyView)

	var helpText string
	if m.showDirtyPrompt {
		aKey := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("[a]")
		oKey := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("[o]")
		cKey := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("[c]")
		helpText = dirtyPromptStyle.Render(fmt.Sprintf("Editor has content! %s ppend • %s verwrite • %s ancel", aKey, oKey, cKey))
	} else if m.focus == editorFocus {
		tab := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("tab")
		ctrlE := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("ctrl+e")
		esc := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render("esc")
		ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")
		helpText = helpStyle.Render(fmt.Sprintf("%s: history • %s: exit • %s: eval & exit • %s: quit",
			tab, ctrlE, esc, ctrlC))
	} else {
		arrows := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("↑/↓")
		enter := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render("enter")
		tab := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("tab")
		ctrlE := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("ctrl+e")
		ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")
		helpText = helpStyle.Render(fmt.Sprintf("%s: navigate • %s: select • %s: editor • %s: exit • %s: quit",
			arrows, enter, tab, ctrlE, ctrlC))
	}

	return fmt.Sprintf("%s\n%s", splitView, helpText)
}

func (m model) renderHistory(width, height int) string {
	if len(m.commandHistory) == 0 {
		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No history yet")
	}

	var historyItems []string
	historyItems = append(historyItems, lipgloss.NewStyle().Bold(true).Render("Command History:"))
	historyItems = append(historyItems, "")

	startIdx := 0
	endIdx := len(m.commandHistory)
	maxVisible := height - 3

	if len(m.commandHistory) > maxVisible {
		if m.historySelection < maxVisible/2 {
			endIdx = maxVisible
		} else if m.historySelection > len(m.commandHistory)-maxVisible/2 {
			startIdx = len(m.commandHistory) - maxVisible
		} else {
			startIdx = m.historySelection - maxVisible/2
			endIdx = m.historySelection + maxVisible/2
		}
	}

	for i := startIdx; i < endIdx && i < len(m.commandHistory); i++ {
		input := m.commandHistory[i]
		displayInput := strings.ReplaceAll(input, "\n", " ")
		if len(displayInput) > width-4 {
			displayInput = displayInput[:width-7] + "..."
		}

		if i == m.historySelection {
			historyItems = append(historyItems, selectedItemStyle.Render("► "+displayInput))
		} else {
			historyItems = append(historyItems, historyItemStyle.Render("  "+displayInput))
		}
	}

	return strings.Join(historyItems, "\n")
}

func Launch(logger *slog.Logger) {
	p := tea.NewProgram(initialModel(logger), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
