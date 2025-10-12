package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type Model struct {
	Core *core.Model
	W, H int
	Tbl  table.Model
	Nav  list.Model

	focusLeft bool
}

func New(core *core.Model) Model {
	cols := views.ImageColumns(80)
	return Model{
		Core:      core,
		Tbl:       components.NewTable(cols, 12, true),
		Nav:       components.NewSidebar([]string{"Overview", "Images"}, 20, 12),
		focusLeft: true,
	}
}

func (m Model) Init() tea.Cmd { return m.Core.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focusLeft = !m.focusLeft
			return m, nil
		case "enter":
			if m.focusLeft {
				m.applySidebarSelection()
				if m.Core.Page == core.PageImages {
					return m, m.Core.CmdFetchImages()
				}
				return m, nil
			}
		case "down":
			if m.focusLeft {
				var cmd tea.Cmd
				m.Nav, cmd = m.Nav.Update(x)
				return m, cmd
			}
			if m.Core.Page == core.PageImages {
				m.Tbl.MoveDown(1)
			}
			return m, nil
		case "up":
			if m.focusLeft {
				var cmd tea.Cmd
				m.Nav, cmd = m.Nav.Update(x)
				return m, cmd
			}
			if m.Core.Page == core.PageImages {
				m.Tbl.MoveUp(1)
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.W, m.H = x.Width, x.Height
		m.Tbl.SetColumns(views.ImageColumns(m.W - 6))
		h := max(m.H-10, 5)
		m.Tbl.SetHeight(h)
	default:
		m.Core.Reduce(msg)
		if m.Core.Page == core.PageImages {
			views.ApplyImages(&m.Tbl, m.Core.Images)
		}
	}
	return m, nil
}

func (m *Model) applySidebarSelection() {
	if len(m.Nav.Items()) == 0 {
		return
	}
	title := m.Nav.SelectedItem().FilterValue()
	switch title {
	case "Overview":
		m.Core.Page = core.PageInfo
	case "Images":
		m.Core.Page = core.PageImages
	}
}

var (
	leftStyle  = lipgloss.NewStyle().Padding(0, 1)
	rightStyle = lipgloss.NewStyle().Padding(0, 1)
)

func (m Model) View() string {
	left := leftStyle.Render(m.Nav.View())

	var rightB strings.Builder
	if m.Core.Page == core.PageInfo {
		status := "NOT CONNECTED"
		info := ""
		if m.Core.Conn.OK {
			status = "CONNECTED"
			info = fmt.Sprintf("Docker %s (%s)", m.Core.Daemon.Version, m.Core.Daemon.OS)
		}
		if m.Core.Conn.Err != "" {
			info = "error: " + m.Core.Conn.Err
		}
		rightB.WriteString(views.RenderOverview(status, info))
	} else {
		rightB.WriteString(fmt.Sprintf("  Images:\n  - total: %d\n\n", len(m.Core.Images)))
		rightB.WriteString(views.RenderIndented(m.Tbl, "  "))
	}

	right := rightStyle.Render(rightB.String())
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
