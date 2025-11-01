package images

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

type Model struct {
	core *core.Core
	W, H int

	Tbl     table.Model
	Keys    KeyMap
	focused bool

	confirming    bool
	pendingDelete []string
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
		core:    core,
		Tbl:     components.NewTable(cols, 12, true),
		Keys:    NewKeyMap(),
		focused: false,
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
	return cmds.FetchImagesCmd(m.core.API, m.core.Images.Filter, 0)
}

func (m *Model) Update(msg tea.Msg) (pages.Page, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		// フォーカスがないときは無視
		if !m.focused {
			return m, nil
		}
		if m.confirming {
			switch {
			case key.Matches(x, m.Keys.Yes):
				refs := append([]string(nil), m.pendingDelete...)
				m.confirming = false
				m.pendingDelete = nil
				return m, tea.Batch(
					cmds.DeleteImagesCmd(m.core.API, refs, ports.ImageRemoveOptions{Force: true, PruneChildlen: true}, 0),
				)
			case key.Matches(x, m.Keys.No):
				m.confirming = false
				m.pendingDelete = nil
				return m, nil
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
		case key.Matches(x, m.Keys.Delete):
			idx := m.Tbl.Cursor()
			if idx >= 0 && idx < len(m.core.Images.List) {
				ref := m.core.Images.List[idx].ID
				m.pendingDelete = []string{ref}
				m.confirming = true
			}
			return m, nil
		}
	default:
		m.core.Reduce(msg)
		switch msg.(type) {
		case cmds.ImagesListedMsg:
			views.ApplyImages(&m.Tbl, m.core.Images.List)
		case cmds.ImagesDeletedMsg:
			return m, m.Init()
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  Images:\n  - total: %d\n\n", len(m.core.Images.List)))
	b.WriteString(views.RenderImages(m.Tbl, "  "))
	if m.confirming {
		n := len(m.pendingDelete)
		msg := "Delete the selected image? (y/N)"
		if n > 1 {
			msg = fmt.Sprintf("Delete %d images?", n)
		}
		b.WriteString("\n")
		b.WriteString("  ")
		b.WriteString(msg)
	}
	return b.String()
}
