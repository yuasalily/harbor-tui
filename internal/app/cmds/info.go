package cmds

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type DaemonInfoMsg struct {
	Info ports.DaemonInfo
	Err  error
}

const defaultInfoTimeout = 3 * time.Second

func FetchDaemonInfoCmd(api ports.DockerAPI, timeout time.Duration) tea.Cmd {
	if timeout <= 0 {
		timeout = defaultInfoTimeout
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		info, err := api.Info(ctx)
		return DaemonInfoMsg{Info: info, Err: err}
	}
}
