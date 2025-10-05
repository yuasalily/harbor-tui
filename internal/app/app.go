package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	witdh, height int

	// Docker
	dockerOK       bool
	dockerErr      string
	serverVersion  string
	daemonPlatform string

	api DockerAPI
}

func New(api DockerAPI) Model {
	return Model{api: api}
}

// Bubble Tea ライフサイクル
func (m Model) Init() tea.Cmd { return m.fetchDaemonInfoCmd() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.witdh, m.height = msg.Width, msg.Height
	case daemonInfoMsg:
		if msg.err != nil {
			m.dockerErr = msg.err.Error()
			m.dockerOK = false
			return m, nil
		}
		m.dockerOK = true
		m.serverVersion = msg.version
		m.daemonPlatform = msg.platform
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

type daemonInfoMsg struct {
	version  string
	platform string
	err      error
}

func (m Model) fetchDaemonInfoCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		info, err := m.api.Info(ctx)
		if err != nil {
			return daemonInfoMsg{err: err}
		}
		return daemonInfoMsg{version: info.Version, platform: info.OS}
	}
}
