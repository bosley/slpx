package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OutputScreen struct {
	viewport viewport.Model
}

func NewOutputScreen() *OutputScreen {
	vp := viewport.New(30, 5)
	return &OutputScreen{
		viewport: vp,
	}
}

func (s *OutputScreen) OnEnter(shared *SharedState) tea.Cmd {
	s.viewport.Width = shared.Width - 2
	s.viewport.Height = shared.Height - 4
	content := shared.RenderOutput()
	s.viewport.SetContent(lipgloss.NewStyle().Width(s.viewport.Width).Render(content))
	return nil
}

func (s *OutputScreen) OnExit(shared *SharedState) {
}

func (s *OutputScreen) Update(shared *SharedState, msg tea.Msg) (Screen, tea.Cmd) {
	var vpCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		shared.Width = msg.Width
		shared.Height = msg.Height
		s.viewport.Width = msg.Width - 2
		s.viewport.Height = msg.Height - 4
		content := shared.RenderOutput()
		s.viewport.SetContent(lipgloss.NewStyle().Width(s.viewport.Width).Render(content))

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "esc", shared.TuiConfig.CmdToggleOutput, "q":
			return NewREPLScreen(), nil
		}
	}

	s.viewport, vpCmd = s.viewport.Update(msg)
	return s, vpCmd
}

func (s *OutputScreen) View(shared *SharedState) string {
	arrows := shared.PromptStyle().Render("↑/↓")
	escKey := shared.SecondaryActionStyle().Render("esc")
	ctrlO := shared.SecondaryActionStyle().Render(shared.TuiConfig.CmdToggleOutput)
	qKey := shared.SecondaryActionStyle().Render("q")
	ctrlC := shared.ErrorStyle().Render("ctrl+c")

	helpText := shared.HelpStyle().Render(fmt.Sprintf("%s: scroll • %s/%s/%s: back to REPL • %s: quit",
		arrows, escKey, ctrlO, qKey, ctrlC))

	return fmt.Sprintf("%s\n\n%s", shared.FocusedStyle().Render(s.viewport.View()), helpText)
}

