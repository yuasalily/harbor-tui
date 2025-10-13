package overview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/pages"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type Model struct {
	core *core.Core
	W, H int
}

func New(core *core.Core) *Model { return &Model{core: core} }

func (m *Model) Name() string     { return "Overview" }
func (m *Model) SetSize(w, h int) { m.W, m.H = w, h }
func (m *Model) Init() tea.Cmd    { return cmds.FetchDaemonInfoCmd(m.core.API, 0) }
func (m *Model) Update(msg tea.Msg) (pages.Page, tea.Cmd) {
	m.core.Reduce(msg)
	return m, nil
}

func (m *Model) View() string {
	status := "NOT CONNECTED"
	info := ""
	if m.core.Overview.Conn.OK {
		status = "CONNECTED"
		info = fmt.Sprintf("Docker %s (%s)", m.core.Overview.Daemon.Version, m.core.Overview.Daemon.OS)
	}
	if m.core.Overview.Conn.Err != "" {
		info = "error: " + m.core.Overview.Conn.Err
	}
	return views.RenderOverview(status, info)
}
