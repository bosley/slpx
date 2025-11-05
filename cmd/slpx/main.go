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
