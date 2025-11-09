package containers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	"github.com/yuasalily/harbor-tui/internal/ui/pages"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type viewMode int

const (
	modeList viewMode = iota
	modeDetailLogs
)

type Model struct {
	core *core.Core
	W, H int

	Tbl     table.Model
	Keys    KeyMap
	focused bool

	containerDetailMode viewMode
	containerDetailID   string
	containerDetailBuf  []string
	containerDetailCh   <-chan string
	containerDetailDone <-chan error
}

func New(core *core.Core) *Model {
	cols := views.ContainerColumns(80)
	return &Model{
		core:                core,
		Tbl:                 components.NewTable(cols, 12, true),
		Keys:                NewKeyMap(),
		focused:             false,
		containerDetailMode: modeList,
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
		if !m.focused {
			return m, nil
		}
		if m.containerDetailMode == modeDetailLogs {
			switch x.String() {
			case "esc":
				// 詳細ビューを閉じて一覧に戻る
				m.containerDetailMode = modeList
				m.containerDetailID = ""
				m.containerDetailBuf = nil
				m.containerDetailCh = nil
				m.containerDetailDone = nil
				return m, nil
			}
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
			if m.containerDetailMode == modeList {
				id := m.selectedContainerID()
				if id == "" {
					return m, nil
				}
				m.containerDetailMode = modeDetailLogs
				m.containerDetailID = id
				m.containerDetailBuf = nil
				return m, cmds.StartContainerLogsCmd(m.core.API, ports.ContainerLogsOptions{
					ContainerID: id,
					Tail: 200,
					Stdout: true,
					Stderr: true,
				})
			}
			return m, nil
		case key.Matches(x, m.Keys.Remove):
			// TODO: コンテナ削除
			return m, nil
		}
	default:
		m.core.Reduce(msg)
		switch v := msg.(type) {
		case cmds.ContainerLogsStartedMsg:
			if m.containerDetailMode == modeDetailLogs && v.ContainerID == m.containerDetailID {
				m.containerDetailCh = v.Ch
				m.containerDetailDone = v.Done
				// 次の1行を待つ
				return m, cmds.NextContainerLogLine(m.containerDetailID, m.containerDetailCh, m.containerDetailDone)
			}
		case cmds.ContainerLogLineMsg:
			if m.containerDetailMode == modeDetailLogs && v.ContainerID == m.containerDetailID {
				m.containerDetailBuf = append(m.containerDetailBuf, v.Line)
				// 次行を待つ
				return m, cmds.NextContainerLogLine(m.containerDetailID, m.containerDetailCh, m.containerDetailDone)
			}
		case cmds.ContainerLogsEndedMsg:
			return m, nil
		case cmds.ContainersListMsg:
			views.ApplyContainers(&m.Tbl, m.core.Containers.List)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.containerDetailMode == modeDetailLogs {
		var b strings.Builder
		b.WriteString("  Container Logs\n\n")
		for _, ln := range m.containerDetailBuf {
			b.WriteString("  ")
			b.WriteString(ln)
			b.WriteByte('\n')
		}
		return b.String()
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Containers:\n  - total: %d\n\n", len(m.core.Containers.List)))
	b.WriteString(views.RenderContainers(m.Tbl, "  "))
	return b.String()
}
