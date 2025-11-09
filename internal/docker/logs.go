package docker

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
)

type ContainerLogsOptions struct {
	ContainerID string
	Follow      bool
	Tail        int
	Since       time.Time
	Timestamps  bool
	Stdout      bool
	Stderr      bool
}

func (c *Client) ContainerLogs(ctx context.Context, opts ContainerLogsOptions) (io.ReadCloser, error) {
	since := ""
	if !opts.Since.IsZero() {
		since = opts.Since.Format(time.RFC3339)
	}
	lo := container.LogsOptions{
		ShowStdout: opts.Stdout || !opts.Stderr,
		ShowStderr: opts.Stderr || !opts.Stdout,
		Follow:     opts.Follow,
		Tail:       "",
		Since:      since,
		Timestamps: opts.Timestamps,
		Details:    false,
	}
	if opts.Tail > 0 {
		lo.Tail = strconv.Itoa(opts.Tail)
	}
	return c.cli.ContainerLogs(ctx, opts.ContainerID, lo)
}


