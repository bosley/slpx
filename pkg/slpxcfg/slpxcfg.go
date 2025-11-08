package slpxcfg

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bosley/slpx/pkg/slp/env"
	"github.com/bosley/slpx/pkg/slp/object"
	"github.com/bosley/slpx/pkg/slp/repl"
	"github.com/bosley/slpx/pkg/slp/slp"
)

/*
Package slpxcfg provides functionality to load and extract typed configuration variables
from SLPX script files.

LoadConfig reads an SLPX file, evaluates it in a full execution environment, and extracts
specified variables from the resulting memory state with type validation.

Usage:
	logger := slog.Default()

	variables := []slpxcfg.Variable{
		{Identifier: "app_name", Type: object.OBJ_TYPE_STRING, Required: true},
		{Identifier: "port", Type: object.OBJ_TYPE_INTEGER, Required: true},
		{Identifier: "debug_mode", Type: object.OBJ_TYPE_INTEGER, Required: false},
	}

	config, err := slpxcfg.LoadConfig(logger, "config.slpx", variables, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	appName := config["app_name"].D.(string)
	port := config["port"].D.(object.Integer)

Validation:
- Required variables must exist in the file or an error is returned
- Optional variables (Required: false) are skipped if not found
- Type checking is enforced: actual type must match expected type
- Use object.OBJ_TYPE_ANY to accept any type without validation

Example config.slpx:
	(set app_name "MyApp")
	(set port 8080)
	(set debug_mode 1)

The file is evaluated with full SLPX capabilities including all standard library functions.
Only the specified variables are extracted and returned in the result map.
*/

var (
	ErrTimeout = errors.New("config evaluation timed out")
)

type Variable struct {
	Identifier object.Identifier
	Type       object.ObjType
	Required   bool
}

type evalResult struct {
	result object.Obj
	err    error
}

type Loader interface {
	Load(file string, variables []Variable) (map[object.Identifier]object.Obj, error)
}

type loaderImpl struct {
	logger     *slog.Logger
	maxTimeout time.Duration
	fs         env.FS
	io         env.IO
}

func New(logger *slog.Logger, maxTimeout time.Duration, fs env.FS, io env.IO) Loader {
	return &loaderImpl{
		logger:     logger,
		maxTimeout: maxTimeout,
		fs:         fs,
		io:         io,
	}
}

func (l *loaderImpl) Load(file string, variables []Variable) (map[object.Identifier]object.Obj, error) {
	return loadFile(l.logger, file, l.maxTimeout, variables, l.fs, l.io)
}

func Load(logger *slog.Logger, file string, variables []Variable, timeout time.Duration) (map[object.Identifier]object.Obj, error) {
	fs := env.DefaultFS()
	io := env.DefaultIO()
	return loadFile(logger, file, timeout, variables, fs, io)
}

func loadFile(logger *slog.Logger, file string, timeout time.Duration, variables []Variable, fs env.FS, io env.IO) (map[object.Identifier]object.Obj, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// We let it create its own mem
	session := repl.NewSessionBuilder(logger).
		WithFS(fs).
		WithIO(io).
		Build(file)

	resultChan := make(chan evalResult, 1)
	go func() {
		result, err := session.Evaluate(string(content))
		resultChan <- evalResult{result, err}
	}()

	var result object.Obj
	select {
	case res := <-resultChan:
		result = res.result
		err = res.err
	case <-time.After(timeout):
		return nil, ErrTimeout
	}

	if err != nil {
		if parseErr, ok := err.(*slp.ParseError); ok {
			line, col, lineStart, lineEnd := positionToLineCol(string(content), parseErr.Position)
			var errMsg strings.Builder
			errMsg.WriteString(fmt.Sprintf("Parse error in %s at line %d, column %d:\n", file, line, col))

			if lineStart < len(content) && lineEnd <= len(content) {
				lineContent := string(content[lineStart:lineEnd])
				errMsg.WriteString(fmt.Sprintf("  %d | %s\n", line, lineContent))
				errMsg.WriteString("      ")
				for i := 1; i < col; i++ {
					errMsg.WriteString(" ")
				}
				errMsg.WriteString("^\n")
			}
			errMsg.WriteString(parseErr.Message)
			return nil, fmt.Errorf("%s", errMsg.String())
		}
		return nil, err
	}

	if result.Type == object.OBJ_TYPE_ERROR {
		errObj := result.D.(object.Error)
		formatted := formatError(errObj, string(content))
		return nil, fmt.Errorf("evaluation error:\n%s", formatted)
	}

	mem := session.GetMEM()

	resultMap := make(map[object.Identifier]object.Obj)

	for _, variable := range variables {
		obj, err := mem.Get(variable.Identifier, true)
		if err != nil {
			if variable.Required {
				return nil, fmt.Errorf("required variable '%s' not found in config", variable.Identifier)
			}
			continue
		}

		if variable.Type != object.OBJ_TYPE_ANY && obj.Type != variable.Type {
			return nil, fmt.Errorf("type mismatch for variable '%s': expected %s, got %s", variable.Identifier, variable.Type, obj.Type)
		}

		resultMap[variable.Identifier] = obj
	}

	return resultMap, nil
}

func positionToLineCol(content string, position int) (line int, col int, lineStart int, lineEnd int) {
	line = 1
	col = 1
	lineStart = 0

	for i := 0; i < len(content) && i < position; i++ {
		if content[i] == '\n' {
			line++
			col = 1
			lineStart = i + 1
		} else {
			col++
		}
	}

	lineEnd = lineStart
	for lineEnd < len(content) && content[lineEnd] != '\n' {
		lineEnd++
	}

	return line, col, lineStart, lineEnd
}

func formatError(err object.Error, sourceContent string) string {
	var output strings.Builder

	if err.File != "" {
		if err.Position == 0 {
			output.WriteString(fmt.Sprintf("Error in %s:\n", err.File))
			output.WriteString(err.Message)
		} else {
			line, col, lineStart, lineEnd := positionToLineCol(sourceContent, err.Position)

			output.WriteString(fmt.Sprintf("Error in %s at line %d, column %d:\n", err.File, line, col))

			if lineStart < len(sourceContent) && lineEnd <= len(sourceContent) {
				lineContent := sourceContent[lineStart:lineEnd]
				output.WriteString(fmt.Sprintf("  %d | %s\n", line, lineContent))

				output.WriteString("      ")
				for i := 1; i < col; i++ {
					output.WriteString(" ")
				}
				output.WriteString("^\n")
			}

			output.WriteString(err.Message)
		}
	} else {
		output.WriteString(fmt.Sprintf("Error: %s", err.Message))
	}

	return output.String()
}
