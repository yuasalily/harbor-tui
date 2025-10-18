package containers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	"github.com/yuasalily/harbor-tui/internal/ui/pages"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type Model struct {
	core *core.Core
	W, H int

	Tbl  table.Model
	Keys KeyMap
}

func New(core *core.Core) *Model {
	cols := views.ContainerColumns(80)
	return &Model{
		core: core,
		Tbl:  components.NewTable(cols, 12, true),
		Keys: NewKeyMap(),
	}
}

func (m *Model) SetSize(w, h int) {
	m.W, m.H = w, h
	inner := max(w-2, 20)
	m.Tbl.SetColumns(views.ContainerColumns(inner))

	tableH := max(m.H-6, 5)
	m.Tbl.SetHeight(tableH)
}

func (m *Model) Init() tea.Cmd {
	return cmds.FetchContainersCmd(m.core.API, m.core.Containers.Filter, 0)
}

func (m *Model) Update(msg tea.Msg) (pages.Page, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(x, m.Keys.Down):
			m.Tbl.MoveDown(1)
			return m, nil
		case key.Matches(x, m.Keys.Up):
			m.Tbl.MoveUp(1)
			return m, nil
		case key.Matches(x, m.Keys.Refresh):
			return m, m.Init()
		case key.Matches(x, m.Keys.Logs):
			// TODO: ログ表示
			return m, nil
		case key.Matches(x, m.Keys.Remove):
			// TODO: コンテナ削除
			return m, nil
		}
	default:
		m.core.Reduce(msg)
		switch msg.(type) {
		case cmds.ContainersListMsg:
			views.ApplyContainers(&m.Tbl, m.core.Containers.List)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Containers:\n  - total: %d\n\n", len(m.core.Containers.List)))
	b.WriteString(views.RenderContainers(m.Tbl, "  "))
	return b.String()
}
