package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane struct {
	idx      int
	viewport viewport.Model
	content  string
}

type model struct {
	panes      []pane
	activePane int
}

func initialModel() model {
	return model{
		panes: []pane{
			{idx: 0, content: "Pane 0"},
			{idx: 1, content: "Pane 1"},
		},
		activePane: 0,
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
		case "tab":
			if len(m.panes) > 0 {
				m.activePane = (m.activePane + 1) % len(m.panes)
			}
		}
	case tea.WindowSizeMsg:
		if len(m.panes) == 0 {
			return m, nil
		}
		// 縦に分割: 各ペインの高さ = 全体の高さ / ペイン数
		h := msg.Height / len(m.panes)
		for i := range m.panes {
			// ボーダー分引く
			m.panes[i].viewport = viewport.New(msg.Width-2, h-2)
			m.panes[i].viewport.SetContent(m.panes[i].content)
		}
	}
	return m, nil
}

func (m model) View() string {
	views := []string{}
	for i, p := range m.panes {
		// アクティブなペインは二重線
		border := lipgloss.NormalBorder()
		if i == m.activePane {
			border = lipgloss.DoubleBorder()
		}

		style := lipgloss.NewStyle().
			Border(border).
			Width(p.viewport.Width).
			Height(p.viewport.Height)

		views = append(views, style.Render(p.viewport.View()))
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
