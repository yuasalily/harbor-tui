package core

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type Page int

const (
	PageInfo Page = iota
	PageImages
)

type Model struct {
	API ports.DockerAPI

	Page      Page
	
	Conn struct {OK bool; Err string}
	Daemon ports.DaemonInfo
	Images    []ports.ImageSummary
}

func New(api ports.DockerAPI) Model { return Model{API: api, Page: PageInfo} }

func (m Model) Init() tea.Cmd { return cmds.FetchDaemonInfoCmd(m.API, 3*time.Second) }

func (m *Model) CmdFetchImages() tea.Cmd {
	return cmds.FetchImagesCmd(m.API, ports.ImagesListOptions{All: true}, 5*time.Second)
}

func (m *Model) Reduce(msg tea.Msg) {
	switch x := msg.(type) {
	case cmds.DaemonInfoMsg:
		if x.Err != nil {
			m.Conn.OK = false
			m.Conn.Err = x.Err.Error()
			return
		}
		m.Conn.OK = true
		m.Conn.Err = ""
		m.Daemon = x.Info
	case cmds.ImagesListedMsg:
		if x.Err != nil {
			m.Conn.Err = x.Err.Error()
			return
		}
		m.Conn.Err = ""
		m.Images = x.Items
	}
}
