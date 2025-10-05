package docker

import (
	"context"

	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func New() (*Client, error) {
	c, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, err
	}
	return &Client{cli: c}, nil
}

func (c *Client) Close() error { return c.cli.Close() }

func (c *Client) PingAndVersion(ctx context.Context) (version string, platform string, err error) {
	if _, err = c.cli.Ping(ctx); err != nil {
		return "", "", err
	}
	ver, err := c.cli.ServerVersion(ctx)
	if err != nil {
		return "", "", err
	}
	return ver.Version, ver.Os, nil
}