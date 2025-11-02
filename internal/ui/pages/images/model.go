package images

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/core"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	uidialog "github.com/yuasalily/harbor-tui/internal/ui/dialog"
	"github.com/yuasalily/harbor-tui/internal/ui/pages"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type Model struct {
	core *core.Core
	W, H int

	Tbl     table.Model
	Keys    KeyMap
	focused bool

	selectedIDs    map[string]struct{}
	resetCursorTop bool                // 削除後に先頭へカーソルを移動するためのフラグ
	inUseIDs       map[string]struct{} // コンテナが使用中のImageID
}

func (m *Model) SetFocused(f bool) {
	m.focused = f
	if f {
		m.Tbl.Focus()
	} else {
		m.Tbl.Blur()
	}
}

func New(core *core.Core) *Model {
	cols := views.ImageColumns(80)
	return &Model{
		core:        core,
		Tbl:         components.NewTable(cols, 12, true),
		Keys:        NewKeyMap(),
		focused:     false,
		selectedIDs: map[string]struct{}{},
		inUseIDs:    map[string]struct{}{},
	}
}

func (m *Model) SetSize(w, h int) {
	m.W, m.H = w, h
	inner := max(w-2, 20)
	m.Tbl.SetColumns(views.ImageColumns(inner))

	tableH := max(m.H-6, 5)
	m.Tbl.SetHeight(tableH)
}

func (m *Model) Init() tea.Cmd {
	// 使用中判定のため、イメージとコンテナを取得
	return tea.Batch(
		cmds.FetchImagesCmd(m.core.API, m.core.Images.Filter, 0),
		cmds.FetchContainersCmd(m.core.API, m.core.Containers.Filter, 0),
	)
}

func (m *Model) toggleSelectAtCursor() {
	idx := m.Tbl.Cursor()
	if idx < 0 || idx >= len(m.core.Images.List) {
		return
	}
	id := m.core.Images.List[idx].ID
	// 使用中は選択不可
	if _, locked := m.inUseIDs[id]; locked {
		return
	}
	if _, ok := m.selectedIDs[id]; ok {
		delete(m.selectedIDs, id)
	} else {
		m.selectedIDs[id] = struct{}{}
	}
	m.decorateSelectionOnRows()
}

func (m *Model) decorateSelectionOnRows() {
	// SEL列(index=0)に [*] / [ ] / [L] をいれる
	rows := m.Tbl.Rows()
	for i := range rows {
		if i >= len(m.core.Images.List) {
			continue
		}
		id := m.core.Images.List[i].ID
		mark := "[ ]"
		if _, locked := m.inUseIDs[id]; locked {
			mark = "[L]"
		}
		if _, ok := m.selectedIDs[id]; ok {
			mark = "[*]"
		}
		if len(rows[i]) > 0 {
			rows[i][0] = mark
		}
	}
	m.Tbl.SetRows(rows)
}

func (m *Model) rebuildInUseSet(cts []ports.ContainerSummary) {
	set := make(map[string]struct{}, len(cts))
	for _, c := range cts {
		if c.ImageID != "" {
			set[c.ImageID] = struct{}{}
		}
	}
	m.inUseIDs = set
	m.pruneSelectionToCurrentList()
}

func (m *Model) pruneSelectionToCurrentList() {
	if len(m.selectedIDs) == 0 {
		return
	}
	next := make(map[string]struct{}, len(m.selectedIDs))
	for _, it := range m.core.Images.List {
		if _, ok := m.selectedIDs[it.ID]; ok {
			next[it.ID] = struct{}{}
		}
	}
	m.selectedIDs = next
}

func (m *Model) Update(msg tea.Msg) (pages.Page, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		// フォーカスがないときは無視
		if !m.focused {
			return m, nil
		}
		switch {
		case key.Matches(x, m.Keys.Down):
			m.Tbl.MoveDown(1)
			return m, nil
		case key.Matches(x, m.Keys.Up):
			m.Tbl.MoveUp(1)
			return m, nil
		case key.Matches(x, m.Keys.Toggle):
			m.toggleSelectAtCursor()
			return m, nil
		case key.Matches(x, m.Keys.Refresh):
			return m, m.Init()
		case key.Matches(x, m.Keys.Delete):
			refs := make([]string, 0, len(m.selectedIDs))
			for id := range m.selectedIDs {
				if _, locked := m.inUseIDs[id]; !locked {
					refs = append(refs, id)
				}
			}
			if len(refs) > 0 {
				return m, func() tea.Msg {
					return uidialog.OpenConfirmDialogMsg{
						Title:   "Delete images",
						Body:    "This will remove the selected iamge(s)",
						Hint:    "[y] Delete    [n] Cancel",
						Payload: refs,
					}
				}
			}
			return m, nil
		}
	default:
		switch v := msg.(type) {
		case uidialog.DialogResultMsg:
			if !v.Confirmed {
				return m, nil
			}
			refs := v.Payload.([]string)
			return m, cmds.DeleteImagesCmd(m.core.API, refs, ports.ImageRemoveOptions{Force: true, PruneChildlen: true}, 0)
		}
		m.core.Reduce(msg)
		switch msg.(type) {
		case cmds.ImagesListedMsg:
			views.ApplyImages(&m.Tbl, m.core.Images.List)
			m.pruneSelectionToCurrentList()
			m.decorateSelectionOnRows()
			if m.resetCursorTop {
				m.Tbl.SetCursor(0)
				m.resetCursorTop = false
			}
		case cmds.ContainersListMsg:
			m.rebuildInUseSet(m.core.Containers.List)
			m.decorateSelectionOnRows()
		case cmds.ImagesDeletedMsg:
			m.selectedIDs = map[string]struct{}{}
			m.resetCursorTop = true
			return m, m.Init()
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString("  Images\n\n")
	b.WriteString(views.RenderImages(m.Tbl, "  "))
	return b.String()
}
