package tui

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bosley/slpx/pkg/repl"
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Update(shared *SharedState, msg tea.Msg) (Screen, tea.Cmd)
	View(shared *SharedState) string
	OnEnter(shared *SharedState) tea.Cmd
	OnExit(shared *SharedState)
}

type model struct {
	shared        *SharedState
	currentScreen Screen
}

func initialModel(logger *slog.Logger) model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	sessionPath := filepath.Join(cwd, ".repl.slpx")

	capturedIO := &capturedIO{}
	session := repl.NewSessionBuilder(logger).
		WithIO(capturedIO).
		Build(sessionPath)

	shared := &SharedState{
		Logger:         logger,
		Session:        session,
		CapturedIO:     capturedIO,
		CommandHistory: []string{},
		OutputHistory:  []string{},
		Width:          0,
		Height:         0,
		PendingInput:   "",
	}

	initialScreen := NewREPLScreen()

	return model{
		shared:        shared,
		currentScreen: initialScreen,
	}
}

func (m model) Init() tea.Cmd {
	return m.currentScreen.OnEnter(m.shared)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nextScreen, cmd := m.currentScreen.Update(m.shared, msg)

	if nextScreen != m.currentScreen {
		m.currentScreen.OnExit(m.shared)
		m.currentScreen = nextScreen
		enterCmd := m.currentScreen.OnEnter(m.shared)
		return m, tea.Batch(cmd, enterCmd)
	}

	m.currentScreen = nextScreen
	return m, cmd
}

func (m model) View() string {
	return m.currentScreen.View(m.shared)
}

func Launch(logger *slog.Logger) {
	p := tea.NewProgram(initialModel(logger), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
