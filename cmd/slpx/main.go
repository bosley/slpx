package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bosley/slpx/cmd/slpx/tui"
	"github.com/bosley/slpx/pkg/rt"
	"github.com/bosley/slpx/pkg/slp/object"
	"github.com/bosley/slpx/pkg/slp/slp"
	"github.com/fatih/color"
)

func main() {

	slpxHome := setupSLPXHome()
	setupContent := loadSetupFile(slpxHome)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if len(os.Args) < 2 {
		tui.Launch(logger, slpxHome, setupContent)
		return
	}

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]

	switch filePath {
	case "uninstall":
		uninstall(logger)
		return
	case "install":
		install(logger)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		absFilePath = filePath
	}

	runtime, err := rt.New(rt.Config{
		Logger:          logger,
		SLPXHome:        slpxHome,
		LaunchDirectory: absFilePath,
		SetupContent:    setupContent,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating runtime: %v\n", err)
		os.Exit(1)
	}

	defer runtime.Stop()

	ac, err := runtime.NewActiveContext("main")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating active context: %v\n", err)
		os.Exit(1)
	}

	session := ac.GetRepl()

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

func loadSetupFile(slpxHome string) string {
	setupFile := filepath.Join(slpxHome, "init.slp")
	if _, err := os.Stat(setupFile); os.IsNotExist(err) {
		return ""
	}
	content, err := os.ReadFile(setupFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading setup file: %v\n", err)
		os.Exit(1)
	}
	return string(content)
}

func setupSLPXHome() string {

	slpxHome := os.Getenv("SLPX_HOME")
	if slpxHome == "" {
		home, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
			os.Exit(1)
		}
		slpxHome = filepath.Join(home, "slpx")

		_, err = os.Stat(slpxHome)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
				os.Exit(1)
			}
			color.HiYellow("Config directory %s does not exist, creating it...", slpxHome)
			if err := os.MkdirAll(slpxHome, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
				os.Exit(1)
			}
			time.Sleep(250 * time.Millisecond)
		}
		os.Setenv("SLPX_HOME", slpxHome)
	}
	return slpxHome
}

func install(logger *slog.Logger) {
	slpxHome := setupSLPXHome()
	writeDefaultSetupFile(slpxHome)
	time.Sleep(250 * time.Millisecond)
	color.HiGreen("SLPX installed successfully")
	os.Exit(0)
}

func writeDefaultSetupFile(slpxHome string) {

	content := `
(set text_foreground "#000000")
(set text_background "#FFFF00")
(set cmd_toggle_editor "ctrl+e")
(set cmd_toggle_output "ctrl+o")
(set cmd_clear "clear")
	`

	setupFile := filepath.Join(slpxHome, "init.slp")
	if _, err := os.Stat(setupFile); os.IsNotExist(err) {
		color.HiYellow("Creating default setup file %s...", setupFile)
		if err := os.WriteFile(setupFile, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default setup file: %v\n", err)
			os.Exit(1)
		}
	}
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

func uninstall(logger *slog.Logger) {
	slpxHome := os.Getenv("SLPX_HOME")
	if slpxHome == "" {
		home, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
			os.Exit(1)
		}
		slpxHome = filepath.Join(home, "slpx")

		_, err = os.Stat(slpxHome)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
				os.Exit(1)
			}
			color.HiYellow("Config directory %s does not exist, skipping uninstall", slpxHome)
			return
		}
	}

	color.HiRed("Uninstalling SLPX...")
	os.RemoveAll(slpxHome)
	os.Unsetenv("SLPX_HOME")
	color.HiGreen("SLPX uninstalled successfully")
}
