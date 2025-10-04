package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/client"
)

type DockerChecker func(ctx context.Context) (version string, platform string, err error)

type Model struct {
	witdh, height int

	// Docker
	dockerOK       bool
	dockerErr      string
	serverVersion  string
	daemonPlatform string

	dockerChecker DockerChecker
}

// 関数型オプションパターン
type Option func(*Model)

func WithDockerChecker(c DockerChecker) Option {
	return func(m *Model) { m.dockerChecker = c }
}

func New(opts ...Option) Model {
	m := Model{dockerChecker: defaultDockerChecker}
	for _, o := range opts {
		o(&m)
	}
	return m
}

// Bubble Tea ライフサイクル
func (m Model) Init() tea.Cmd { return m.checkDockerCmd() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.witdh, m.height = msg.Width, msg.Height
	case dockerStatusMsg:
		m.dockerOK = msg.ok
		m.serverVersion = msg.version
		m.daemonPlatform = msg.platform
		if msg.err != nil {
			m.dockerErr = msg.err.Error()
		}
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

type dockerStatusMsg struct {
	ok       bool
	version  string
	platform string
	err      error
}

func (m Model) checkDockerCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		ver, plat, err := m.dockerChecker(ctx)
		if err != nil {
			return dockerStatusMsg{ok: false, err: err}
		}
		return dockerStatusMsg{ok: true, version: ver, platform: plat}
	}
}

func defaultDockerChecker(ctx context.Context) (string, string, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return "", "", err
	}
	defer cli.Close()

	if _, err := cli.Ping(ctx); err != nil {
		return "", "", err
	}
	ver, err := cli.ServerVersion(ctx)
	if err != nil {
		return "", "", err
	}
	return ver.Version, ver.Os, nil
}
