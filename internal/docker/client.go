package docker

import (
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
