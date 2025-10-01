package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/client"
)

type model struct {
	width  int
	height int

	// Docker
	dockerOk       bool
	dockerErr      string
	serverVersion  string
	daemonPlatform string
}

type dockerStatusMsg struct {
	ok       bool
	version  string
	platform string
	err      error
}

func initialModel() model { return model{} }

func checkDockerCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cli, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			return dockerStatusMsg{ok: false, err: err}
		}
		defer cli.Close()

		if _, err := cli.Ping(ctx); err != nil {
			return dockerStatusMsg{ok: false, err: err}
		}

		ver, err := cli.ServerVersion(ctx)
		if err != nil {
			return dockerStatusMsg{ok: false, err: err}
		}
		return dockerStatusMsg{ok: true, version: ver.Version, platform: ver.Os}
	}
}

func (m model) Init() tea.Cmd { return checkDockerCmd() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "crtl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case dockerStatusMsg:
		m.dockerOk = msg.ok
		m.serverVersion = msg.version
		m.daemonPlatform = msg.platform
		if msg.err != nil {
			m.dockerErr = msg.err.Error()
		}
	}
	return m, nil
}

func (m model) View() string {
	status := "NOT CONNECTED"
	info := ""
	if m.dockerOk {
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

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
