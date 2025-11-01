package cmds

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type ImagesDeletedMsg struct {
	Refs []string
	Err  error
}

const defaultDeleteTimeout = 10 * time.Second

func DeleteImagesCmd(api ports.DockerAPI, refs []string, opts ports.ImageRemoveOptions, timeout time.Duration) tea.Cmd {
	if timeout <= 0 {
		timeout = defaultDeleteTimeout
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		var err error
		for _, r := range refs {
			if e := api.ImageRemove(ctx, r, opts); e != nil {
				err = e
				break
			}
		}
		return ImagesDeletedMsg{Refs: refs, Err: err}
	}
}
