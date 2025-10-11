package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type Model struct {
	Core *core.Model
	W, H int
	Tbl  table.Model
}

func New(core *core.Model) Model {
	cols := views.ImageColumns(80)
	return Model{Core: core, Tbl: components.NewTable(cols, 12, true)}
}

func (m Model) Init() tea.Cmd { return m.Core.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.Core.Page = core.PageInfo
			return m, nil
		case "2":
			if m.Core.Page == core.PageInfo {
				m.Core.Page = core.PageImages
			} else {
				m.Core.Page = core.PageInfo
			}
			return m, nil
		case "i":
			return m, m.Core.CmdFetchImages()
		case "down":
			if m.Core.Page == core.PageImages {
				m.Tbl.MoveDown(1)
			}
			return m, nil
		case "up":
			if m.Core.Page == core.PageImages {
				m.Tbl.MoveUp(1)
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.W, m.H = x.Width, x.Height
		m.Tbl.SetColumns(views.ImageColumns(m.W - 6))
		h := max(m.H - 10, 5)
		m.Tbl.SetHeight(h)
	default:
		m.Core.Reduce(msg)
		if m.Core.Page == core.PageImages {
			views.ApplyImages(&m.Tbl, m.Core.Images)
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("\n  Harbor-TUI: Bubble Tea + Docker SDK\n\n")
	if m.Core.Page == core.PageInfo {
		b.WriteString("  [Info]    Images\n\n")
	} else {
		b.WriteString("   Info    [Images]\n\n")
	}

	switch m.Core.Page {
	case core.PageInfo:
		status := "NOT CONNECTED"
		info := ""
		if m.Core.Conn.OK {
			status = "CONNECTED"
			info = fmt.Sprintf("Docker %s (%s)", m.Core.Daemon.Version, m.Core.Daemon.OS)
		}
		if m.Core.Conn.Err != "" {
			info = "error: " + m.Core.Conn.Err
		}
		b.WriteString(views.RenderInfo(status, info))
		b.WriteString("\n  [Keys] q: quit  1: info  2: images\n")
	case core.PageImages:
		b.WriteString(fmt.Sprintf("  Images:\n  - total: %d\n\n", len(m.Core.Images)))
		b.WriteString(views.RenderIndented(m.Tbl, "  "))
		b.WriteString("\n\n  [Keys] q: quit  1: info  2: images  i: fetch\n")
	}
	return b.String()
}
