package cmds

import (
	"bufio"
	"context"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type ContainerLogsStartedMsg struct {
	ContainerID string
	Ch          <-chan string
	Done        <-chan error
}

type ContainerLogLineMsg struct {
	ContainerID string
	Line        string
}

type ContainerLogsEndedMsg struct {
	ContainerID string
	Err         error
}

const defaultLogsTimeout = 0

func StartContainerLogsCmd(api ports.DockerAPI, opts ports.ContainerLogsOptions) tea.Cmd {
	if !opts.Stdout && !opts.Stderr {
		opts.Stdout = true
	}
	return func() tea.Msg {
		ctx := context.Background()
		rc, err := api.ContainerLogs(ctx, opts)
		if err != nil {
			return ContainerLogsEndedMsg{ContainerID: opts.ContainerID, Err: err}
		}
		ch := make(chan string, 64)
		done := make(chan error, 1)
		go func() {
			defer close(ch)
			defer rc.Close()
			reader := bufio.NewReader(rc)
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					// 開業は表示側で統一するので削る
					if line[len(line)-1] == '\n' {
						line = line[:len(line)-1]
					}
					ch <- line
				}
				if err != nil {
					if err == io.EOF {
						done <- nil
					} else {
						done <- err
					}
					close(done)
					return
				}
			}
		}()
		return ContainerLogsStartedMsg{ContainerID: opts.ContainerID, Ch: ch, Done: done}
	}
}

func NextContainerLogLine(containerID string, ch <-chan string, done <-chan error) tea.Cmd {
	return func() tea.Msg {
		select {
		case ln, ok := <-ch:
			if !ok {
				return ContainerLogsEndedMsg{ContainerID: containerID, Err: nil}
			}
			return ContainerLogLineMsg{ContainerID: containerID, Line: ln}
		case err, ok := <-done:
			if !ok {
				return ContainerLogsEndedMsg{ContainerID: containerID, Err: nil}
			}
			return ContainerLogsEndedMsg{ContainerID: containerID, Err: err}
		case <-time.After(200 * time.Microsecond):
			return nil
		}
	}
}
