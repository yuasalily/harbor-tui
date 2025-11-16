package logs

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/ui/pages/containers/details"
)

type Model struct {
	deps   details.Dependencies
	target string
	W, H   int
	focus  bool

	buf  []string
	ch   <-chan string
	done <-chan error
}

func New(deps details.Dependencies) *Model { return &Model{deps: deps} }

func (m *Model) SetTarget(id string) { m.target = id }
func (m *Model) SetSize(w, h int)    { m.W, m.H = w, h }
func (m *Model) SetFocused(f bool)   { m.focus = f }

func (m *Model) ShortKeys() []key.Binding {
	// escはグローバルに近い扱いでルーター側が処理する
	return nil
}

func (m *Model) Init() tea.Cmd {
	if m.target == "" {
		return nil
	}
	return cmds.StartContainerLogsCmd(m.deps.API, ports.ContainerLogsOptions{
		ContainerID: m.target,
		Tail:        200,
		Stdout:      true,
		Stderr:      true,
		Follow:      true,
	})
}

func (m *Model) Update(msg tea.Msg) (details.Detail, tea.Cmd) {
	switch v := msg.(type) {
	case cmds.ContainerLogsStartedMsg:
		if v.ContainerID != m.target {
			return m, nil
		}
		m.ch = v.Ch
		m.done = v.Done
		return m, cmds.NextContainerLogLine(m.target, m.ch, m.done)
	case cmds.ContainerLogLineMsg:
		if v.ContainerID != m.target {
			return m, nil
		}
		m.buf = append(m.buf, v.Line)
		return m, cmds.NextContainerLogLine(m.target, m.ch, m.done)
	case cmds.ContainerLogsEndedMsg:
		if v.ContainerID != m.target {
			return m, nil
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) Close() tea.Cmd {
	// StartContainerLogsCmdの内部で io.ReadCloserをCloseしている。
	// 追加のキャンセルが必要になればDependenciesにコンテキストを持たせる設計へ拡張
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString("  Container Logs\n")
	for _, ln := range m.buf {
		b.WriteString("  ")
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return b.String()
}
