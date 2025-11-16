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
	"github.com/yuasalily/harbor-tui/internal/ui/pages/containers/details"
	"github.com/yuasalily/harbor-tui/internal/ui/pages/containers/details/logs"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type viewMode int

const (
	modeList viewMode = iota
	modeDetail
)

type Model struct {
	core *core.Core
	W, H int

	Tbl     table.Model
	Keys    KeyMap
	focused bool

	mode     viewMode
	detail   details.Detail
	detailID string
	deps     details.Dependencies
}

func New(core *core.Core) *Model {
	cols := views.ContainerColumns(80)
	return &Model{
		core:    core,
		Tbl:     components.NewTable(cols, 12, true),
		Keys:    NewKeyMap(),
		focused: false,
		mode:    modeList,
		deps:    details.Dependencies{API: core.API},
	}
}

func (m *Model) selectedContainerID() string {
	idx := m.Tbl.Cursor()
	if idx < 0 || idx >= len(m.core.Containers.List) {
		return ""
	}
	return m.core.Containers.List[idx].ID
}

func (m *Model) SetFocused(f bool) {
	m.focused = f
	if f {
		m.Tbl.Focus()
	} else {
		m.Tbl.Blur()
	}
	if m.mode == modeDetail && m.detail != nil {
		m.detail.SetFocused(f)
	}
}

func (m *Model) SetSize(w, h int) {
	m.W, m.H = w, h
	inner := max(w-2, 20)
	m.Tbl.SetColumns(views.ContainerColumns(inner))

	tableH := max(m.H-6, 5)
	m.Tbl.SetHeight(tableH)
	if m.mode == modeDetail && m.detail != nil {
		m.detail.SetSize(w, h)
	}
}

func (m *Model) Init() tea.Cmd {
	return cmds.FetchContainersCmd(m.core.API, m.core.Containers.Filter, 0)
}

func (m *Model) Update(msg tea.Msg) (pages.Page, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
		if m.mode == modeDetail {
			switch x.String() {
			case "esc":
				if m.detail != nil {
					_ = m.detail.Close()
					m.detail = nil
					m.detailID = ""
				}
				m.mode = modeList
				return m, nil
			}
			if m.detail != nil {
				d, cmd := m.detail.Update(x)
				m.detail = d
				return m, cmd
			}
			return m, nil
		}
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
			id := m.selectedContainerID()
			if id == "" {
				return m, nil
			}
			m.mode = modeDetail
			m.detailID = id
			m.detail = logs.New(m.deps)
			m.detail.SetTarget(id)
			m.detail.SetSize(m.W, m.H)
			m.detail.SetFocused(m.focused)
			return m, m.detail.Init()

		case key.Matches(x, m.Keys.Remove):
			// TODO: コンテナ削除
			return m, nil
		}
	default:
		m.core.Reduce(msg)
		if m.mode == modeDetail && m.detail != nil {
			d, cmd := m.detail.Update(msg)
			m.detail = d
			if cmd != nil {
				return m, cmd
			}
		}
		switch msg.(type) {
		case cmds.ContainersListMsg:
			views.ApplyContainers(&m.Tbl, m.core.Containers.List)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.mode == modeDetail && m.detail != nil {
		return m.detail.View()
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Containers:\n  - total: %d\n\n", len(m.core.Containers.List)))
	b.WriteString(views.RenderContainers(m.Tbl, "  "))
	return b.String()
}
