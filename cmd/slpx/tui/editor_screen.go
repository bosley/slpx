package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusPane uint

const (
	editorFocus focusPane = iota
	historyFocus
)

type EditorScreen struct {
	textarea         textarea.Model
	focus            focusPane
	historySelection int
	showDirtyPrompt  bool
	pendingHistory   string
}

func NewEditorScreen() *EditorScreen {
	ta := textarea.New()
	ta.Placeholder = "Enter SLPX code..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 0
	ta.ShowLineNumbers = true
	ta.KeyMap.InsertNewline.SetEnabled(true)

	return &EditorScreen{
		textarea:         ta,
		focus:            editorFocus,
		historySelection: 0,
		showDirtyPrompt:  false,
		pendingHistory:   "",
	}
}

func (s *EditorScreen) OnEnter(shared *SharedState) tea.Cmd {
	editorWidth := (shared.Width / 2) - 2
	s.textarea.SetWidth(editorWidth)
	s.textarea.SetHeight(shared.Height - 5)
	s.historySelection = len(shared.CommandHistory) - 1
	if s.historySelection < 0 {
		s.historySelection = 0
	}
	return s.textarea.Focus()
}

func (s *EditorScreen) OnExit(shared *SharedState) {
}

func (s *EditorScreen) Update(shared *SharedState, msg tea.Msg) (Screen, tea.Cmd) {
	var tiCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		shared.Width = msg.Width
		shared.Height = msg.Height
		editorWidth := (msg.Width / 2) - 2
		s.textarea.SetWidth(editorWidth)
		s.textarea.SetHeight(msg.Height - 5)

	case tea.KeyMsg:
		if s.showDirtyPrompt {
			switch msg.String() {
			case "a":
				currentValue := s.textarea.Value()
				if currentValue != "" {
					s.textarea.SetValue(currentValue + "\n" + s.pendingHistory)
				} else {
					s.textarea.SetValue(s.pendingHistory)
				}
				s.showDirtyPrompt = false
				s.pendingHistory = ""
				s.focus = editorFocus
				cmd := s.textarea.Focus()
				return s, cmd
			case "o":
				s.textarea.SetValue(s.pendingHistory)
				s.showDirtyPrompt = false
				s.pendingHistory = ""
				s.focus = editorFocus
				cmd := s.textarea.Focus()
				return s, cmd
			case "c", "esc":
				s.showDirtyPrompt = false
				s.pendingHistory = ""
				s.focus = editorFocus
				cmd := s.textarea.Focus()
				return s, cmd
			}
			return s, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "ctrl+e":
			value := s.textarea.Value()
			if value != "" {
				shared.PendingInput = value
			}
			return NewREPLScreen(), nil
		case "esc":
			value := s.textarea.Value()
			if value != "" {
				if value == "clear" || value == "cls" {
					shared.ClearOutput()
					s.textarea.Reset()
					return NewREPLScreen(), tea.WindowSize()
				}
				output := shared.EvaluateInput(value)
				shared.AddCommand(value, output)
				s.textarea.Reset()
			}
			return NewREPLScreen(), nil
		case "tab":
			if s.focus == editorFocus {
				s.focus = historyFocus
				s.textarea.Blur()
				return s, nil
			} else {
				s.focus = editorFocus
				cmd := s.textarea.Focus()
				return s, cmd
			}
		case "up":
			if s.focus == historyFocus && len(shared.CommandHistory) > 0 {
				if s.historySelection > 0 {
					s.historySelection--
				}
				return s, nil
			}
		case "down":
			if s.focus == historyFocus && len(shared.CommandHistory) > 0 {
				if s.historySelection < len(shared.CommandHistory)-1 {
					s.historySelection++
				}
				return s, nil
			}
		case "enter":
			if s.focus == historyFocus && len(shared.CommandHistory) > 0 {
				selectedInput := shared.CommandHistory[s.historySelection]

				if s.textarea.Value() != "" {
					s.showDirtyPrompt = true
					s.pendingHistory = selectedInput
					return s, nil
				} else {
					s.textarea.SetValue(selectedInput)
					s.focus = editorFocus
					cmd := s.textarea.Focus()
					return s, cmd
				}
			}
		}
	}

	s.textarea, tiCmd = s.textarea.Update(msg)
	return s, tiCmd
}

func (s *EditorScreen) View(shared *SharedState) string {
	editorWidth := (shared.Width / 2) - 2
	historyWidth := shared.Width - editorWidth - 4

	editorStyle := BlurredStyle
	historyStyle := BlurredStyle
	if s.focus == editorFocus {
		editorStyle = FocusedStyle
	} else {
		historyStyle = FocusedStyle
	}

	contentHeight := shared.Height - 3

	editorView := editorStyle.Width(editorWidth).Height(contentHeight).Render(s.textarea.View())
	historyView := s.renderHistory(shared, historyWidth, contentHeight)
	historyView = historyStyle.Width(historyWidth).Height(contentHeight).Render(historyView)

	splitView := lipgloss.JoinHorizontal(lipgloss.Top, editorView, historyView)

	var helpText string
	if s.showDirtyPrompt {
		aKey := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("[a]")
		oKey := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("[o]")
		cKey := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("[c]")
		helpText = DirtyPromptStyle.Render(fmt.Sprintf("Editor has content! %s ppend • %s verwrite • %s ancel", aKey, oKey, cKey))
	} else if s.focus == editorFocus {
		tab := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("tab")
		ctrlE := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("ctrl+e")
		esc := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render("esc")
		ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")
		helpText = HelpStyle.Render(fmt.Sprintf("%s: history • %s: exit • %s: eval & exit • %s: quit",
			tab, ctrlE, esc, ctrlC))
	} else {
		arrows := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("↑/↓")
		enter := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render("enter")
		tab := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("tab")
		ctrlE := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("ctrl+e")
		ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")
		helpText = HelpStyle.Render(fmt.Sprintf("%s: navigate • %s: select • %s: editor • %s: exit • %s: quit",
			arrows, enter, tab, ctrlE, ctrlC))
	}

	return fmt.Sprintf("%s\n%s", splitView, helpText)
}

func (s *EditorScreen) renderHistory(shared *SharedState, width, height int) string {
	if len(shared.CommandHistory) == 0 {
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
	endIdx := len(shared.CommandHistory)
	maxVisible := height - 3

	if len(shared.CommandHistory) > maxVisible {
		if s.historySelection < maxVisible/2 {
			endIdx = maxVisible
		} else if s.historySelection > len(shared.CommandHistory)-maxVisible/2 {
			startIdx = len(shared.CommandHistory) - maxVisible
		} else {
			startIdx = s.historySelection - maxVisible/2
			endIdx = s.historySelection + maxVisible/2
		}
	}

	for i := startIdx; i < endIdx && i < len(shared.CommandHistory); i++ {
		input := shared.CommandHistory[i]
		displayInput := strings.ReplaceAll(input, "\n", " ")
		if len(displayInput) > width-4 {
			displayInput = displayInput[:width-7] + "..."
		}

		if i == s.historySelection {
			historyItems = append(historyItems, SelectedItemStyle.Render("► "+displayInput))
		} else {
			historyItems = append(historyItems, HistoryItemStyle.Render("  "+displayInput))
		}
	}

	return strings.Join(historyItems, "\n")
}

