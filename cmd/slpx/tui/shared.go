package tui

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/bosley/slpx/pkg/rt"
	"github.com/bosley/slpx/pkg/slp/object"
	"github.com/bosley/slpx/pkg/slp/repl"
	"github.com/bosley/slpx/pkg/slp/slp"
	"github.com/charmbracelet/lipgloss"
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

type SharedState struct {
	Logger         *slog.Logger
	Session        *repl.Session
	CapturedIO     *capturedIO
	CommandHistory []string
	OutputHistory  []string
	Width          int
	Height         int
	PendingInput   string
	Runtime        rt.Runtime
	ActiveContext  rt.ActiveContext
	TuiConfig      rt.TuiConfig
}

func (s *SharedState) PromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.PromptColor)).Bold(true)
}

func (s *SharedState) ResultStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.ResultColor))
}

func (s *SharedState) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.ErrorColor)).Bold(true)
}

func (s *SharedState) HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.HelpColor))
}

func (s *SharedState) FocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(s.TuiConfig.FocusedBorderColor))
}

func (s *SharedState) BlurredStyle() lipgloss.Style {
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(s.TuiConfig.BlurredBorderColor))
}

func (s *SharedState) SelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.SelectedItemColor)).Bold(true)
}

func (s *SharedState) HistoryItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.HistoryItemColor))
}

func (s *SharedState) DirtyPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.DirtyPromptColor)).Bold(true)
}

func (s *SharedState) SecondaryActionStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.TuiConfig.SecondaryActionColor)).Bold(true)
}

func (s *SharedState) EvaluateInput(input string) string {
	s.CapturedIO.GetAndClear()

	result, err := s.Session.Evaluate(input)

	capturedOutput := s.CapturedIO.GetAndClear()

	var output strings.Builder

	if capturedOutput != "" {
		output.WriteString(capturedOutput)
		if !strings.HasSuffix(capturedOutput, "\n") {
			output.WriteString("\n")
		}
	}

	if err != nil || result.Type == object.OBJ_TYPE_ERROR {
		routedResult, routeErr := s.tryCommandRoute(input)
		if routeErr == nil && routedResult.Type != object.OBJ_TYPE_ERROR && routedResult.Type != "" {
			routedOutput := s.CapturedIO.GetAndClear()
			if routedOutput != "" {
				output.WriteString(routedOutput)
				if !strings.HasSuffix(routedOutput, "\n") {
					output.WriteString("\n")
				}
			}
			if routedResult.Type != object.OBJ_TYPE_NONE {
				output.WriteString(s.ResultStyle().Render(routedResult.Encode()))
			}
			return output.String()
		}

		if err != nil {
			if parseErr, ok := err.(*slp.ParseError); ok {
				output.WriteString(s.ErrorStyle().Render(fmt.Sprintf("Parse Error: %s", parseErr.Message)))
			} else {
				output.WriteString(s.ErrorStyle().Render(fmt.Sprintf("Error: %v", err)))
			}
			return output.String()
		}

		if result.Type == object.OBJ_TYPE_ERROR {
			errObj := result.D.(object.Error)
			output.WriteString(s.ErrorStyle().Render(fmt.Sprintf("Error: %s", errObj.Message)))
			return output.String()
		}
	}

	if result.Type != object.OBJ_TYPE_NONE {
		output.WriteString(s.ResultStyle().Render(result.Encode()))
	}

	return output.String()
}

func (s *SharedState) tryCommandRoute(input string) (object.Obj, error) {
	if s.TuiConfig.CommandRouter.Body == nil {
		return object.Obj{}, nil // no command router, so we don't try to route
	}

	routerIdent := object.Identifier("command_router")
	_, err := s.Session.GetMEM().Get(routerIdent, true)
	if err != nil {
		return object.Obj{}, fmt.Errorf("command_router not in memory")
	}

	callExpr := fmt.Sprintf("(command_router %q)", input)
	result, evalErr := s.Session.Evaluate(callExpr)
	if evalErr != nil {
		return object.Obj{}, evalErr
	}

	return result, nil
}

func (s *SharedState) AddCommand(input string, output string) {
	s.CommandHistory = append(s.CommandHistory, input)
	s.OutputHistory = append(s.OutputHistory, s.PromptStyle().Render("> ")+input)
	if output != "" {
		s.OutputHistory = append(s.OutputHistory, output)
	}
}

func (s *SharedState) ClearOutput() {
	s.OutputHistory = []string{}
}

func (s *SharedState) RenderOutput() string {
	return strings.Join(s.OutputHistory, "\n")
}
