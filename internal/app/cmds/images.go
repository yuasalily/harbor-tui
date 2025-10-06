package cmds

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type ImagesListedMsg struct {
	Items []ports.ImageSummary
	Err   error
}

const defaultImagesTimeout = 5 * time.Second

func FetchImagesCmd(api ports.DockerAPI, opts ports.ImagesListOptions, timeout time.Duration) tea.Cmd {
	if timeout <= 0 {
		timeout = defaultImagesTimeout
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		items, err := api.ImagesList(ctx, opts)
		return ImagesListedMsg{Items: items, Err: err}
	}
}
