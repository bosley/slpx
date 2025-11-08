package installer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bosley/slpx/cmd/slpx/assets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

type screen int

const (
	screenInstallOptions screen = iota
	screenUninstallOptions
	screenConfirmUninstall
	screenProgress
	screenDone
)

type action int

const (
	actionNone action = iota
	actionInstallDefault
	actionInstallAdvanced
	actionUninstall
	actionExit
)

type model struct {
	screen         screen
	selected       int
	slpxHome       string
	logger         *slog.Logger
	isInstalled    bool
	currentAction  action
	progress       float64
	progressMsg    string
	err            error
	width          int
	height         int
	filesToProcess []string
	filesContent   map[string]string
	currentFile    int
}

type progressMsg struct {
	progress float64
	message  string
}

type doneMsg struct{}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#5f87ff"))

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(lipgloss.Color("#d0d0d0"))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#d75fd7")).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)
)

func Launch(logger *slog.Logger, slpxHome string) {
	isInstalled := checkInstalled(slpxHome)

	initialScreen := screenInstallOptions
	if isInstalled {
		initialScreen = screenUninstallOptions
	}

	p := tea.NewProgram(
		model{
			screen:      initialScreen,
			selected:    0,
			slpxHome:    slpxHome,
			logger:      logger,
			isInstalled: isInstalled,
		},
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running installer: %v\n", err)
		os.Exit(1)
	}
}

func checkInstalled(slpxHome string) bool {
	initFile := filepath.Join(slpxHome, "init.slpx")
	_, err := os.Stat(initFile)
	return err == nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.screen == screenInstallOptions || m.screen == screenUninstallOptions || m.screen == screenConfirmUninstall {
				if m.selected > 0 {
					m.selected--
				}
			}

		case "down", "j":
			if m.screen == screenInstallOptions {
				if m.selected < 2 {
					m.selected++
				}
			} else if m.screen == screenUninstallOptions {
				if m.selected < 1 {
					m.selected++
				}
			} else if m.screen == screenConfirmUninstall {
				if m.selected < 1 {
					m.selected++
				}
			}

		case "enter":
			return m.handleSelection()
		}

	case installPreparedMsg:
		m.filesContent = msg.files
		m.filesToProcess = msg.fileList
		m.currentFile = 0

		if len(m.filesToProcess) > 0 {
			filename := m.filesToProcess[0]
			m.progressMsg = fmt.Sprintf("Writing %s...", filename)
			return m, processFile(m.slpxHome, filename, m.filesContent[filename])
		}
		return m, nil

	case processNextFileMsg:
		m.currentFile++

		if m.currentFile >= len(m.filesToProcess) {
			m.progress = 1.0
			m.progressMsg = "Installation complete!"
			m.screen = screenDone
			return m, nil
		}

		m.progress = float64(m.currentFile) / float64(len(m.filesToProcess))
		filename := m.filesToProcess[m.currentFile]
		m.progressMsg = fmt.Sprintf("Writing %s...", filename)
		return m, processFile(m.slpxHome, filename, m.filesContent[filename])

	case uninstallPreparedMsg:
		m.filesToProcess = msg.fileList
		m.currentFile = 0

		if len(m.filesToProcess) > 0 {
			filename := m.filesToProcess[0]
			m.progressMsg = fmt.Sprintf("Removing %s...", filename)
			return m, removeFile(m.slpxHome, filename)
		}
		return m, nil

	case removeNextFileMsg:
		m.currentFile++

		if m.currentFile >= len(m.filesToProcess) {
			m.progress = 1.0
			m.progressMsg = "Uninstallation complete!"
			m.screen = screenDone
			return m, nil
		}

		m.progress = float64(m.currentFile) / float64(len(m.filesToProcess))
		filename := m.filesToProcess[m.currentFile]
		m.progressMsg = fmt.Sprintf("Removing %s...", filename)
		return m, removeFile(m.slpxHome, filename)

	case progressMsg:
		m.progress = msg.progress
		m.progressMsg = msg.message
		if m.progress >= 1.0 {
			m.screen = screenDone
		}
		return m, nil

	case doneMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenInstallOptions:
		switch m.selected {
		case 0:
			m.currentAction = actionInstallDefault
			m.screen = screenProgress
			m.progress = 0
			m.progressMsg = "Installing default variant..."
			return m, m.performInstall(false)
		case 1:
			m.currentAction = actionInstallAdvanced
			m.screen = screenProgress
			m.progress = 0
			m.progressMsg = "Installing advanced variant..."
			return m, m.performInstall(true)
		case 2:
			return m, tea.Quit
		}

	case screenUninstallOptions:
		switch m.selected {
		case 0:
			m.screen = screenConfirmUninstall
			m.selected = 0
		case 1:
			return m, tea.Quit
		}

	case screenConfirmUninstall:
		switch m.selected {
		case 0:
			m.currentAction = actionUninstall
			m.screen = screenProgress
			m.progress = 0
			m.progressMsg = "Uninstalling SLPX..."
			return m, m.performUninstall()
		case 1:
			m.screen = screenUninstallOptions
			m.selected = 0
		}

	case screenDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) performInstall(advanced bool) tea.Cmd {
	return func() tea.Msg {
		var files map[string]string
		if advanced {
			files = assets.LoadAdvancedVariant()
		} else {
			files = assets.LoadDefaultVariant()
		}

		fileList := make([]string, 0, len(files))
		for filename := range files {
			fileList = append(fileList, filename)
		}

		return installPreparedMsg{
			files:    files,
			fileList: fileList,
		}
	}
}

type installPreparedMsg struct {
	files    map[string]string
	fileList []string
}

type processNextFileMsg struct{}

func processFile(slpxHome string, filename string, content string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond)

		filePath := filepath.Join(slpxHome, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return progressMsg{progress: 1.0, message: fmt.Sprintf("Error: %v", err)}
		}

		return processNextFileMsg{}
	}
}

func (m model) performUninstall() tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(m.slpxHome)
		if err != nil {
			return progressMsg{progress: 1.0, message: fmt.Sprintf("Error: %v", err)}
		}

		if len(entries) == 0 {
			return progressMsg{progress: 1.0, message: "Directory already empty"}
		}

		fileList := make([]string, 0, len(entries))
		for _, entry := range entries {
			fileList = append(fileList, entry.Name())
		}

		return uninstallPreparedMsg{
			fileList: fileList,
		}
	}
}

type uninstallPreparedMsg struct {
	fileList []string
}

type removeNextFileMsg struct{}

func removeFile(slpxHome string, filename string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond)

		filePath := filepath.Join(slpxHome, filename)
		info, err := os.Stat(filePath)
		if err == nil {
			if info.IsDir() {
				os.RemoveAll(filePath)
			} else {
				os.Remove(filePath)
			}
		}

		return removeNextFileMsg{}
	}
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.screen {
	case screenInstallOptions:
		return m.renderInstallOptions()
	case screenUninstallOptions:
		return m.renderUninstallOptions()
	case screenConfirmUninstall:
		return m.renderConfirmUninstall()
	case screenProgress:
		return m.renderProgress()
	case screenDone:
		return m.renderDone()
	}

	return ""
}

func (m model) renderInstallOptions() string {
	var b strings.Builder

	b.WriteString("\n")
	title := titleStyle.Render("SLPX Installation")
	b.WriteString(center(m.width, title))
	b.WriteString("\n\n\n")

	options := []string{
		"Install Default Configuration",
		"Install Advanced Configuration",
		"Exit",
	}

	for i, option := range options {
		var line string
		if i == m.selected {
			line = selectedItemStyle.Render("▸ " + option)
		} else {
			line = menuItemStyle.Render("  " + option)
		}
		b.WriteString(center(m.width, line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("↑/↓: navigate • enter: select • q: quit")
	b.WriteString(center(m.width, help))

	return centerVertical(m.height, b.String())
}

func (m model) renderUninstallOptions() string {
	var b strings.Builder

	b.WriteString("\n")
	title := titleStyle.Render("SLPX Setup")
	b.WriteString(center(m.width, title))
	b.WriteString("\n\n\n")

	msg := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render("✓ SLPX is already installed")
	b.WriteString(center(m.width, msg))
	b.WriteString("\n")
	b.WriteString(center(m.width, fmt.Sprintf("Location: %s", m.slpxHome)))
	b.WriteString("\n\n")

	options := []string{
		"Uninstall",
		"Exit",
	}

	for i, option := range options {
		var line string
		if i == m.selected {
			line = selectedItemStyle.Render("▸ " + option)
		} else {
			line = menuItemStyle.Render("  " + option)
		}
		b.WriteString(center(m.width, line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("↑/↓: navigate • enter: select • q: quit")
	b.WriteString(center(m.width, help))

	return centerVertical(m.height, b.String())
}

func (m model) renderConfirmUninstall() string {
	var b strings.Builder

	b.WriteString("\n")
	title := titleStyle.Render("Confirm Uninstallation")
	b.WriteString(center(m.width, title))
	b.WriteString("\n\n\n")

	warning := errorStyle.Render("⚠ This will delete all SLPX configuration files")
	b.WriteString(center(m.width, warning))
	b.WriteString("\n")
	b.WriteString(center(m.width, fmt.Sprintf("From: %s", m.slpxHome)))
	b.WriteString("\n\n")

	options := []string{
		"Yes, uninstall",
		"No, go back",
	}

	for i, option := range options {
		var line string
		if i == m.selected {
			line = selectedItemStyle.Render("▸ " + option)
		} else {
			line = menuItemStyle.Render("  " + option)
		}
		b.WriteString(center(m.width, line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	help := helpStyle.Render("↑/↓: navigate • enter: select • q: quit")
	b.WriteString(center(m.width, help))

	return centerVertical(m.height, b.String())
}

func (m model) renderProgress() string {
	var b strings.Builder

	b.WriteString("\n")
	title := titleStyle.Render("Please Wait")
	b.WriteString(center(m.width, title))
	b.WriteString("\n\n\n")

	b.WriteString(center(m.width, m.progressMsg))
	b.WriteString("\n\n")

	barWidth := min(60, m.width-10)
	filledWidth := int(float64(barWidth) * m.progress)
	emptyWidth := barWidth - filledWidth

	bar := lipgloss.NewStyle().Foreground(lipgloss.Color("#5f87ff")).Render(strings.Repeat("█", filledWidth)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#585858")).Render(strings.Repeat("░", emptyWidth))

	b.WriteString(center(m.width, bar))
	b.WriteString("\n")

	percentage := fmt.Sprintf("%.0f%%", m.progress*100)
	b.WriteString(center(m.width, percentage))

	return centerVertical(m.height, b.String())
}

func (m model) renderDone() string {
	var b strings.Builder

	var title, message string
	switch m.currentAction {
	case actionInstallDefault, actionInstallAdvanced:
		title = "Installation Complete"
		message = successStyle.Render("✓ SLPX has been successfully installed!")
	case actionUninstall:
		title = "Uninstallation Complete"
		message = successStyle.Render("✓ SLPX has been successfully uninstalled")
	}

	b.WriteString("\n")
	b.WriteString(center(m.width, titleStyle.Render(title)))
	b.WriteString("\n\n\n")
	b.WriteString(center(m.width, message))
	b.WriteString("\n\n")

	help := helpStyle.Render("Press any key to exit")
	b.WriteString(center(m.width, help))

	time.AfterFunc(2*time.Second, func() {
		os.Exit(0)
	})

	return centerVertical(m.height, b.String())
}

func center(width int, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth < width {
			padding := (width - lineWidth) / 2
			lines[i] = strings.Repeat(" ", padding) + line
		}
	}
	return strings.Join(lines, "\n")
}

func centerVertical(height int, s string) string {
	lines := strings.Split(s, "\n")
	contentHeight := len(lines)

	if contentHeight < height {
		padding := (height - contentHeight) / 2
		topPadding := strings.Repeat("\n", padding)
		return topPadding + s
	}

	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func InstallDefault(logger *slog.Logger, slpxHome string) {
	files := assets.LoadDefaultVariant()
	for filename, content := range files {
		filePath := filepath.Join(slpxHome, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			color.HiYellow("Creating %s...", filePath)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
				os.Exit(1)
			}
		}
	}
	color.HiGreen("SLPX installed successfully")
}

func Uninstall(logger *slog.Logger, slpxHome string) error {
	color.HiRed("Uninstalling SLPX...")
	if err := os.RemoveAll(slpxHome); err != nil {
		return err
	}
	os.Unsetenv("SLPX_HOME")
	color.HiGreen("SLPX uninstalled successfully")
	return nil
}
