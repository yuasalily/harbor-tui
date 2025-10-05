package docker

import "context"

type DaemonInfo struct {
	Version    string
	OS         string
	APIVersion string
}

func Info(ctx context.Context) (DaemonInfo, error) {
	cli, err := New()
	if err != nil {
		return DaemonInfo{}, err
	}
	defer cli.Close()
	return cli.Info(ctx)
}

func (c *Client) Info(ctx context.Context) (DaemonInfo, error) {
	if _, err := c.cli.Ping(ctx); err != nil {
		return DaemonInfo{}, err
	}

	ver, err := c.cli.ServerVersion(ctx)
	if err != nil {
		return DaemonInfo{}, err
	}

	return DaemonInfo{
		Version:    ver.Version,
		OS:         ver.Os,
		APIVersion: ver.APIVersion,
	}, nil
}
