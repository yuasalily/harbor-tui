package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/cmds"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type Model struct {
	witdh, height int

	// Docker
	dockerOK       bool
	dockerErr      string
	serverVersion  string
	daemonPlatform string

	api ports.DockerAPI
}

func New(api ports.DockerAPI) Model {
	return Model{api: api}
}

// Bubble Tea ライフサイクル
func (m Model) Init() tea.Cmd { return cmds.FetchDaemonInfoCmd(m.api, 3*time.Second) }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.witdh, m.height = msg.Width, msg.Height
	case cmds.DaemonInfoMsg:
		if msg.Err != nil {
			m.dockerErr = msg.Err.Error()
			m.dockerOK = false
			return m, nil
		}
		m.dockerOK = true
		m.serverVersion = msg.Info.Version
		m.daemonPlatform = msg.Info.OS
	}
	return m, nil
}

func (m Model) View() string {
	status := "NOT CONNECTED"
	info := ""
	if m.dockerOK {
		status = "CONNECTED"
		info = fmt.Sprintf("Docker %s (%s)", m.serverVersion, m.daemonPlatform)
	} else if m.dockerErr != "" {
		info = fmt.Sprintf("error: %s", m.dockerErr)
	}
	return fmt.Sprintf(`
	Harbor-TUI: Bubble Tea + Docker SDK

	Status: %s
	%s

	Press 'q' to quit.
	`, status, info)
}
