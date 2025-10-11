package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
	"github.com/yuasalily/harbor-tui/internal/ui/views"
)

type Page int

const (
	PageInfo Page = iota
	PageImages
)

type Model struct {
	witdh, height int
	page          Page

	// Docker
	dockerOK       bool
	dockerErr      string
	serverVersion  string
	daemonPlatform string

	images []ports.ImageSummary
	tbl    table.Model

	api ports.DockerAPI
}

func New(api ports.DockerAPI) Model {
	m := Model{api: api}
	cols := views.ImageColumns(80)
	m.tbl = components.NewTable(cols, 12, true)
	return m
}

// Bubble Tea ライフサイクル
func (m Model) Init() tea.Cmd { return cmds.FetchDaemonInfoCmd(m.api, 3*time.Second) }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.page = PageInfo
			return m, nil
		case "2":
			m.page = PageImages
			return m, nil
		case "tab":
			if m.page == PageInfo {
				m.page = PageImages
			} else {
				m.page = PageInfo
			}
			return m, nil
		case "i":
			if m.page == PageImages {
				return m, cmds.FetchImagesCmd(m.api, ports.ImagesListOptions{All: true}, 5*time.Second)
			}
			return m, nil
		case "j", "down":
			if m.page == PageImages {
				m.tbl.MoveDown(1)
			}
			return m, nil
		case "k", "up":
			if m.page == PageImages {
				m.tbl.MoveUp(1)
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.witdh, m.height = msg.Width, msg.Height
		cols := views.ImageColumns(m.witdh - 6)
		m.tbl.SetColumns(cols)
		h := max(m.height-10, 5)
		m.tbl.SetHeight(h)
	case cmds.DaemonInfoMsg:
		if msg.Err != nil {
			m.dockerErr = msg.Err.Error()
			m.dockerOK = false
			return m, nil
		}
		m.dockerOK = true
		m.serverVersion = msg.Info.Version
		m.daemonPlatform = msg.Info.OS
	case cmds.ImagesListedMsg:
		if msg.Err != nil {
			m.dockerErr = msg.Err.Error()
			return m, nil
		}
		m.dockerErr = ""
		m.images = msg.Items
		views.ApplyImages(&m.tbl, m.images)
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
	Harbor-TUI: Bubble Tea + Docker SDK

	%s
	`, m.renderTabs()))

	switch m.page {
	case PageInfo:
		status := "NOT CONNECTED"
		info := ""
		if m.dockerOK {
			status = "CONNECTED"
			info = fmt.Sprintf("Docker %s (%s)", m.serverVersion, m.daemonPlatform)
		} else if m.dockerErr != "" {
			info = fmt.Sprintf("error: %s", m.dockerErr)
		}
		b.WriteString(views.RenderInfo(status, info))
		b.WriteString("  [Keys] q: quit   1: info   2: images   tab: switch\n")
	case PageImages:
		b.WriteString(fmt.Sprintf(`
  Images:
  -  total: %d

%s

  [Keys] q: quit   1: info   2: images   tab: switch
		`, len(m.images), views.RenderIndented(m.tbl, "  ")))
	}
	return b.String()
}

func (m Model) renderTabs() string {
	active := func(name string) string { return "[" + name + "]" }
	inactive := func(name string) string { return " " + name + " " }
	if m.page == PageInfo {
		return "  " + active("Info") + "  " + inactive("Images")
	}
	return "  " + inactive("Info") + "  " + active("Images")
}
