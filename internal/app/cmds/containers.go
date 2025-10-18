package cmds

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type ContainersListMsg struct {
	Items []ports.ContainerSummary
	Err error
}

const defaultContainersTimeout = 5 * time.Second

func FetchContainersCmd(api ports.DockerAPI, opts ports.ContainersListOptions, timeout time.Duration) tea.Cmd {
	if timeout <= 0 {timeout = defaultContainersTimeout}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		items, err := api.ContainersList(ctx, opts)
		return ContainersListMsg{Items: items, Err: err}
	}
}