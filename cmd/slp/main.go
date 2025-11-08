/*
This is the simplified SLP CLI. It does not offer a tui or engage in any of the specialized
runtime activity. It can take in a single slpx file and execute vanilla slp commands

Use this with "tests/primitive/main.slpx" to run core tests on the SLP language implementation
without the larger runtime overhead

bosley
*/

package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bosley/slpx/pkg/slp/object"
	"github.com/bosley/slpx/pkg/slp/repl"
	"github.com/bosley/slpx/pkg/slp/slp"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		absFilePath = filePath
	}

	session := repl.NewSessionBuilder(logger).Build(absFilePath)

	result, err := session.Evaluate(string(content))
	session.GetIO().Flush()
	if err != nil {
		if parseErr, ok := err.(*slp.ParseError); ok {
			line, col, lineStart, lineEnd := positionToLineCol(string(content), parseErr.Position)
			fmt.Fprintf(os.Stderr, "Parse error in %s at line %d, column %d:\n", absFilePath, line, col)

			if lineStart < len(content) && lineEnd <= len(content) {
				lineContent := string(content[lineStart:lineEnd])
				fmt.Fprintf(os.Stderr, "  %d | %s\n", line, lineContent)

				fmt.Fprintf(os.Stderr, "      ")
				for i := 1; i < col; i++ {
					fmt.Fprintf(os.Stderr, " ")
				}
				fmt.Fprintf(os.Stderr, "^\n")
			}

			fmt.Fprintf(os.Stderr, "%s\n", parseErr.Message)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	if result.Type == object.OBJ_TYPE_ERROR {
		errObj := result.D.(object.Error)
		fmt.Fprintf(os.Stderr, "%s\n", formatError(errObj, string(content)))
		os.Exit(1)
	}

	fmt.Printf("Result: %s\n", result.Encode())
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
