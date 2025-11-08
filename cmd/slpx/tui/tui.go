package tui

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bosley/slpx/pkg/rt"
	"github.com/bosley/slpx/pkg/slp/object"
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

func initialModel(logger *slog.Logger, slpxHome string, setupContent string) (model, error) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	sessionPath := filepath.Join(cwd, ".repl.slpx")

	runtime, err := rt.New(rt.Config{
		Logger:          logger,
		SLPXHome:        slpxHome,
		LaunchDirectory: sessionPath,
		SetupContent:    setupContent,
	})
	if err != nil {
		return model{}, err
	}

	ac, err := runtime.NewActiveContext("tui")
	if err != nil {
		runtime.Stop()
		return model{}, err
	}

	capturedIO := &capturedIO{}
	session := ac.GetRepl()
	session.GetIO().SetStdout(capturedIO)
	session.GetIO().SetStderr(capturedIO)

	tuiConfig := ac.GetTuiConfig()

	if tuiConfig.CommandRouter.Body != nil {
		routerObj := object.Obj{
			Type: object.OBJ_TYPE_FUNCTION,
			D:    tuiConfig.CommandRouter,
			Pos:  0,
		}
		session.GetMEM().Set(object.Identifier("command_router"), routerObj, true)
	}

	shared := &SharedState{
		Logger:         logger,
		Session:        session,
		CapturedIO:     capturedIO,
		CommandHistory: []string{},
		OutputHistory:  []string{},
		Width:          0,
		Height:         0,
		PendingInput:   "",
		Runtime:        runtime,
		ActiveContext:  ac,
		TuiConfig:      tuiConfig,
	}

	initialScreen := NewREPLScreen()

	return model{
		shared:        shared,
		currentScreen: initialScreen,
	}, nil
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

func Launch(logger *slog.Logger, slpxHome string, setupContent string) {
	m, err := initialModel(logger, slpxHome, setupContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TUI: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if m.shared.ActiveContext != nil {
			m.shared.ActiveContext.Close()
		}
		if m.shared.Runtime != nil {
			m.shared.Runtime.Stop()
		}
	}()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
