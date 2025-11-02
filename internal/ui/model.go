package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"

	uidialog "github.com/yuasalily/harbor-tui/internal/ui/dialog"
	"github.com/yuasalily/harbor-tui/internal/ui/pages"
	pcontainers "github.com/yuasalily/harbor-tui/internal/ui/pages/containers"
	pimages "github.com/yuasalily/harbor-tui/internal/ui/pages/images"
	poverview "github.com/yuasalily/harbor-tui/internal/ui/pages/overview"
)

type Model struct {
	Core *core.Core
	W, H int

	Nav     list.Model
	Keys    GlobalKeyMap
	focus   FocusArea
	helpbar components.HelpBar

	pages   map[pages.ID]pages.Page
	current pages.ID // 現在ページ
	navSel  int      // サイドバーの選択インデックス
	dialog  uidialog.Model
}

func New(core *core.Core) Model {
	ov := poverview.New(core)
	im := pimages.New(core)
	ct := pcontainers.New(core)

	pageMap := map[pages.ID]pages.Page{
		pages.PageOverview:   ov,
		pages.PageImages:     im,
		pages.PageContainers: ct,
	}

	var items []string
	for _, meta := range pages.Metas() {
		items = append(items, meta.Title)
	}

	return Model{
		Core:    core,
		Nav:     components.NewSidebar(items, 20, 12),
		Keys:    NewGlobalKeyMap(),
		helpbar: components.NewHelpBar(),
		focus:   FocusNav,
		pages:   pageMap,
		current: pages.PageOverview,
		navSel:  0,
		dialog:  uidialog.Model{},
	}
}

func (m Model) Init() tea.Cmd { return m.currentPage().Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch {
		// ダイアログ表示中はダイアログにキーを移譲
		case m.focus == FocusDialog:
			dlg, cmd := m.dialog.Update(x)
			m.dialog = dlg
			if !m.dialog.Visible {
				m.focus = FocusPage
				if f, ok := m.currentPage().(pages.Focusable); ok {
					f.SetFocused(true)
				}
			}
			return m, cmd
		case key.Matches(x, m.Keys.Quit):
			return m, tea.Quit
		case key.Matches(x, m.Keys.Tab):
			m.toggleFocus()
			return m, nil
		}
		if m.focus == FocusNav {
			prev := m.Nav.Index()
			var cmd tea.Cmd
			m.Nav, cmd = m.Nav.Update(x)
			if m.Nav.Index() != prev {
				m.navSel = m.Nav.Index()
				m.applySidebarSelection()
				return m, m.currentPage().Init()
			}
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.W, m.H = x.Width, x.Height
		navW := m.Nav.Width()
		const panePad = 4
		rightW := max(m.W-navW-panePad, 20)
		const helpOuter = 3
		contentH := max(m.H-helpOuter, 5)

		navInnerH := max(contentH-2, 3)
		m.Nav.SetSize(navW, navInnerH)
		pageInnerH := max(contentH-2, 3)

		for _, p := range m.pages {
			p.SetSize(rightW, pageInnerH)
		}
		m.dialog.SetSize(rightW, pageInnerH)
		return m, nil
	case uidialog.OpenConfirmDialogMsg:
		m.dialog.OpenConfirm(x.Title, x.Body, x.Hint, lipgloss.Color("204"), x.Payload)
		if f, ok := m.currentPage().(pages.Focusable); ok {
			f.SetFocused(false)
		}
		m.focus = FocusDialog
		return m, nil
	case uidialog.DialogResultMsg:
		// ページに結果メッセージを渡す
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

func (m *Model) toggleFocus() {
	switch m.focus {
	case FocusNav:
		m.focus = FocusPage
		if f, ok := m.currentPage().(pages.Focusable); ok {
			f.SetFocused(true)
		}
	case FocusPage:
		if f, ok := m.currentPage().(pages.Focusable); ok {
			f.SetFocused(false)
		}
		m.focus = FocusNav
	case FocusDialog:
		// 将来: ダイアログに移譲
	}
}

func (m Model) currentPage() pages.Page {
	p, ok := m.pages[m.current]
	if !ok {
		panic(fmt.Sprintf("no page instance for id=%d", m.current))
	}
	return p
}

var (
	borderColor      = lipgloss.Color("240") // 通常時
	focusBorderColor = lipgloss.Color("63")  // フォーカス時(ライトブルー)
	baseBoxStyle     = lipgloss.NewStyle().
				Padding(0, 1).
				Border(lipgloss.NormalBorder()).
				BorderForeground(borderColor)
)

func (m Model) View() string {
	var pageKeys []key.Binding
	if pg, ok := m.currentPage().(*pimages.Model); ok {
		pageKeys = pg.Keys.Short()
	}
	if pc, ok := m.currentPage().(*pcontainers.Model); ok {
		pageKeys = pc.Keys.Short()
	}

	var help string
	switch m.focus {
	case FocusNav:
		help = m.helpbar.Render(m.Keys.ShortForNav())
	case FocusPage:
		help = m.helpbar.Render(m.Keys.ShortForPage(), pageKeys)
	case FocusDialog:
		help = m.helpbar.Render(m.Keys.ShortForDialog())
	}

	leftStyle := baseBoxStyle
	rightStyle := baseBoxStyle
	if m.focus == FocusNav {
		leftStyle = leftStyle.BorderForeground(focusBorderColor)
	}
	if m.focus == FocusPage {
		rightStyle = rightStyle.BorderForeground(focusBorderColor)
	}
	helpView := baseBoxStyle.Render(help)
	left := leftStyle.Render(m.Nav.View())
	right := rightStyle.Render(m.currentPage().View() + "\n" + m.dialog.View())

	below := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return lipgloss.JoinVertical(lipgloss.Left, helpView, below)
}
