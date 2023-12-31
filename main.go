package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor  int
}

func newModel() model {
	return model{
		choices: []string{
			"https://google.com",
			"https://bing.com",
			"https://duckduckgo.com",
		},
		cursor: 1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			openUrl(m.choices[m.cursor])
		}
	}

	return m, nil
}

func openUrl(url string) {
	var cmdName string
	switch runtime.GOOS {
	case "linux":
		cmdName = "xdg-open"
	case "windows":
		cmdName = "start"
	case "darwin":
		cmdName = "open"
	default:
		log.Printf("unsupported OS: %s\n", runtime.GOOS)
		return
	}

	cmd := exec.Command(cmdName, url)
	if err := cmd.Run(); err != nil {
		log.Printf("failed to open '%s': %s\n", url, err.Error())
	}
}

func (m model) View() string {
	var sb strings.Builder
	for i, choice := range m.choices {
		if i == m.cursor {
			sb.WriteString("> ")
		}
		sb.WriteString(fmt.Sprintf("[%2d] %s\n", i+1, choice))
	}
	return sb.String()
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
