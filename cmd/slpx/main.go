package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bosley/slpx/pkg/object"
	"github.com/bosley/slpx/pkg/repl"
	"github.com/bosley/slpx/pkg/slp"
)

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

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if len(os.Args) < 2 {
		startInteractiveREPL(logger)
		return
	}

	if len(os.Args) > 2 {
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error in %s: %v\n", filePath, err)
		os.Exit(1)
	}

	if result.Type == object.OBJ_TYPE_ERROR {
		errObj := result.D.(object.Error)
		fmt.Fprintf(os.Stderr, "%s\n", formatError(errObj, string(content)))
		os.Exit(1)
	}

	fmt.Printf("Result: %s\n", result.Encode())
}

func startInteractiveREPL(logger *slog.Logger) {
	fmt.Println("SLPX Interactive REPL")
	fmt.Println("Type expressions to evaluate. Press Ctrl+D (EOF) to exit.")
	fmt.Println()

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	sessionPath := filepath.Join(cwd, ".repl.slpx")
	session := repl.NewSessionBuilder(logger).Build(sessionPath)

	scanner := bufio.NewScanner(os.Stdin)
	lineBuffer := strings.Builder{}

	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()

		lineBuffer.WriteString(line)
		lineBuffer.WriteString("\n")

		input := lineBuffer.String()

		parser := slp.NewParser(input)
		_, err := parser.ParseAll()

		if err != nil {
			if strings.Contains(err.Error(), "unexpected end of input") ||
				strings.Contains(err.Error(), "expected") {
				fmt.Print("... ")
				continue
			}

			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			lineBuffer.Reset()
			fmt.Print("> ")
			continue
		}

		result, err := session.Evaluate(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else if result.Type != object.OBJ_TYPE_NONE {
			fmt.Printf("%s\n", result.Encode())
		}

		lineBuffer.Reset()
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
	}

	fmt.Println()
}
