package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	"github.com/yuasalily/harbor-tui/internal/ui/pages"
	pimages "github.com/yuasalily/harbor-tui/internal/ui/pages/images"
	poverview "github.com/yuasalily/harbor-tui/internal/ui/pages/overview"
)

type Model struct {
	Core *core.Core
	W, H int
	Nav  list.Model
	Keys GlobalKeyMap

	pages   map[pages.ID]pages.Page
	current pages.ID // 現在ページ
}

func New(core *core.Core) Model {
	ov := poverview.New(core)
	im := pimages.New(core)

	pageMap := map[pages.ID]pages.Page{
		pages.PageOverview: ov,
		pages.PageImages:   im,
	}

	var items []string
	for _, meta := range pages.Metas() {
		items = append(items, meta.Title)
	}

	return Model{
		Core:    core,
		Nav:     components.NewSidebar(items, 20, 12),
		Keys:    NewGlobalKeyMap(),
		pages:   pageMap,
		current: pages.PageOverview,
	}
}

func (m Model) Init() tea.Cmd { return m.currentPage().Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(x, m.Keys.Quit):
			return m, tea.Quit
		case key.Matches(x, m.Keys.Select):
			m.applySidebarSelection()
			return m, m.currentPage().Init()
		}
		var cmd tea.Cmd
		m.Nav, cmd = m.Nav.Update(x)
		return m, cmd

	case tea.WindowSizeMsg:
		m.W, m.H = x.Width, x.Height
		for _, p := range m.pages {
			p.SetSize(m.W, m.H)
		}
	}

	p, cmd := m.currentPage().Update(msg)
	m.pages[m.current] = p
	return m, cmd
}

func (m *Model) applySidebarSelection() {
	if len(m.Nav.Items()) == 0 {
		return
	}
	title := m.Nav.SelectedItem().FilterValue()
	m.current = pages.FromTitle(title)
}

func (m Model) currentPage() pages.Page {
	p, ok := m.pages[m.current]
	if !ok {
		panic(fmt.Sprintf("no page instance for id=%d", m.current))
	}
	return p
}

var (
	leftStyle  = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	rightStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
)

func (m Model) View() string {
	left := leftStyle.Render(m.Nav.View())
	right := rightStyle.Render(m.currentPage().View())
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
