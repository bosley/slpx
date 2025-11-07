package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type REPLScreen struct {
	viewport viewport.Model
	textarea textarea.Model
	ready    bool
}

func NewREPLScreen() *REPLScreen {
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

	return &REPLScreen{
		viewport: vp,
		textarea: ta,
		ready:    false,
	}
}

func (s *REPLScreen) OnEnter(shared *SharedState) tea.Cmd {
	if shared.Width > 0 && shared.Height > 0 {
		s.viewport.Width = shared.Width
		s.textarea.SetWidth(shared.Width)
		s.viewport.Height = shared.Height - s.textarea.Height() - lipgloss.Height(gap) - 2
		content := shared.RenderOutput()
		s.viewport.SetContent(lipgloss.NewStyle().Width(s.viewport.Width).Render(content))
		s.viewport.GotoBottom()
		s.ready = true
	}
	if shared.PendingInput != "" {
		s.textarea.SetValue(shared.PendingInput)
		shared.PendingInput = ""
	}
	s.textarea.Focus()
	return textarea.Blink
}

func (s *REPLScreen) OnExit(shared *SharedState) {
}

func (s *REPLScreen) Update(shared *SharedState, msg tea.Msg) (Screen, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		shared.Width = msg.Width
		shared.Height = msg.Height
		s.ready = true

		s.viewport.Width = msg.Width
		s.textarea.SetWidth(msg.Width)
		s.viewport.Height = msg.Height - s.textarea.Height() - lipgloss.Height(gap) - 2

		content := shared.RenderOutput()
		s.viewport.SetContent(lipgloss.NewStyle().Width(s.viewport.Width).Render(content))
		s.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return s, tea.Quit
		case "ctrl+e":
			return NewEditorScreen(), nil
		case "ctrl+o":
			return NewOutputScreen(), nil
		case "enter":
			value := s.textarea.Value()
			if value != "" {
				if value == "clear" || value == "cls" {
					shared.ClearOutput()
					s.viewport.SetContent("")
					s.textarea.Reset()
					return s, tea.WindowSize()
				}
				output := shared.EvaluateInput(value)
				shared.AddCommand(value, output)
				content := shared.RenderOutput()
				s.viewport.SetContent(lipgloss.NewStyle().Width(s.viewport.Width).Render(content))
				s.textarea.Reset()
				s.viewport.GotoBottom()
			}
			return s, nil
		}
	}

	s.textarea, tiCmd = s.textarea.Update(msg)
	s.viewport, vpCmd = s.viewport.Update(msg)

	return s, tea.Batch(tiCmd, vpCmd)
}

func (s *REPLScreen) View(shared *SharedState) string {
	if !s.ready {
		return "Initializing..."
	}

	ctrlE := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("ctrl+e")
	ctrlO := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("ctrl+o")
	ctrlC := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("ctrl+c")
	esc := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("esc")

	helpText := HelpStyle.Render(fmt.Sprintf("%s: editor/history • %s: scroll output • %s/%s: quit",
		ctrlE, ctrlO, ctrlC, esc))
	return fmt.Sprintf("%s%s%s\n%s", s.viewport.View(), gap, s.textarea.View(), helpText)
}

