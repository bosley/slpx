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

var (
	PromptStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
	ResultStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	ErrorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	HelpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	FocusedStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	BlurredStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	SelectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	HistoryItemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	DirtyPromptStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
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

	if err != nil {
		if parseErr, ok := err.(*slp.ParseError); ok {
			output.WriteString(ErrorStyle.Render(fmt.Sprintf("Parse Error: %s", parseErr.Message)))
		} else {
			output.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %v", err)))
		}
		return output.String()
	}

	if result.Type == object.OBJ_TYPE_ERROR {
		errObj := result.D.(object.Error)
		output.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %s", errObj.Message)))
		return output.String()
	}

	if result.Type != object.OBJ_TYPE_NONE {
		output.WriteString(ResultStyle.Render(result.Encode()))
	}

	return output.String()
}

func (s *SharedState) AddCommand(input string, output string) {
	s.CommandHistory = append(s.CommandHistory, input)
	s.OutputHistory = append(s.OutputHistory, PromptStyle.Render("> ")+input)
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
