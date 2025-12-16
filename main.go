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
	ready      bool
}

func initialModel() model {
	return model{
		panes: []pane{
			{idx: 0, content: "Pane 0\nThis is the first pane\nLine 3\nLine 4\nLine 5"},
			{idx: 1, content: "Pane 1\nThis is the second pane\nLine 3\nLine 4\nLine 5"},
		},
		activePane: 0,
		ready:      false,
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
		m.ready = true
		if len(m.panes) == 0 {
			return m, nil
		}
		// 縦に分割: 各ペインの高さ = 全体の高さ / ペイン数
		h := msg.Height / len(m.panes)
		for i := range m.panes {
			// ボーダー分を引く（上下で2行）
			m.panes[i].viewport = viewport.New(msg.Width-4, h-3)
			m.panes[i].viewport.SetContent(m.panes[i].content)
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	views := []string{}
	for i, p := range m.panes {
		// アクティブなペインは二重線
		border := lipgloss.NormalBorder()
		borderColor := lipgloss.Color("240")
		if i == m.activePane {
			border = lipgloss.DoubleBorder()
			borderColor = lipgloss.Color("63")
		}

		style := lipgloss.NewStyle().
			Border(border).
			BorderForeground(borderColor).
			Width(p.viewport.Width + 2).
			Height(p.viewport.Height + 2)

		views = append(views, style.Render(p.viewport.View()))
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
