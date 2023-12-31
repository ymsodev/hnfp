package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type HackerNewsItem struct {
	Id          int    `json:"id"`
	Deleted     bool   `json:"deleted"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Text        string `json:"text"`
	Dead        bool   `json:"dead"`
	Parent      int    `json:"parent"`
	Poll        int    `json:"poll"`
	Kids        []int  `json:"kids"`
	Url         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Parts       []int  `json:"parts"`
	Descendants int    `json:"descendants"`
}

type model struct {
	numStories int
	choices    []*HackerNewsItem
	cursor     int
}

func newModel(numStories int) model {
	ids, err := getTopStories(numStories)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]*HackerNewsItem, 0, numStories)
	for _, id := range ids {
		item, err := getHackerNewsItem(id)
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

func getTopStories(count int) ([]int, error) {
	if count > 500 {
		return nil, errors.New("top stories must be < 500")
	}

	// TODO: handle failed request? (e.g., timeout, etc.)
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	return ids[:count], nil
}

func getHackerNewsItem(id int) (*HackerNewsItem, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item *HackerNewsItem
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}
	return item, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
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
	sb.WriteString("\nPress ENTER to open, ESC to quit\n")
	return sb.String()
}

func main() {
	countPtr := flag.Int("count", 10, "number of top stories")
	flag.Parse()

	if _, err := tea.NewProgram(newModel(*countPtr)).Run(); err != nil {
		log.Fatal(err)
	}
}
