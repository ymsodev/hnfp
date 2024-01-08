package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ymsodev/hnfp/hackernews"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	numStories int
	choices    []*hackernews.Item
	cursor     int
}

func newModel(numStories int) model {
	ids, err := hackernews.GetTopStories()
	if err != nil {
		log.Fatal(err)
	}

	items := make([]*hackernews.Item, 0, numStories)
	for _, id := range ids {
		item, err := hackernews.GetItem(id)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return model{
		choices: items,
		cursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			openUrl(m.choices[m.cursor].Url)
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
		log.Printf("cannot open %s, unsupported OS: %s\n", url, runtime.GOOS)
		return
	}

	cmd := exec.Command(cmdName, url)
	if err := cmd.Run(); err != nil {
		log.Printf("failed to open '%s': %s\n", url, err.Error())
	}
}

func (m model) View() string {
	var sb strings.Builder
	sb.WriteString("Hacker News Front Page\n")
	sb.WriteString("----------------------\n")
	for i, choice := range m.choices {
		if i == m.cursor {
			sb.WriteString("> ")
		} else {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("[%d] %s\n", choice.Score, choice.Title))
	}
	sb.WriteString("\npress enter to open, esc to quit\n")
	return sb.String()
}

func main() {
	countPtr := flag.Int("count", 10, "number of top stories")
	flag.Parse()

	if _, err := tea.NewProgram(newModel(*countPtr)).Run(); err != nil {
		log.Fatal(err)
	}
}
