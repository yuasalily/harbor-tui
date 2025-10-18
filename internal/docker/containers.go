package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// DTO
type ContainerSummary struct {
	ID        string
	Names     []string
	Image     string
	State     string
	Status    string
	CreatedAt time.Time
}

type ContainersListOptions struct {
	All    bool
	Name   string
	Status string
}

func (c *Client) ContainersList(ctx context.Context, opts ContainersListOptions) ([]ContainerSummary, error) {
	args := filters.NewArgs()
	if opts.Name != "" {
		args.Add("name", opts.Name)
	}
	if opts.Status != "" {
		args.Add("status", opts.Status)
	}

	listOpts := container.ListOptions{All: opts.All, Filters: args}
	items, err := c.cli.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	out := make([]ContainerSummary, 0, len(items))
	for _, it := range items {
		created := time.Unix(it.Created, 0)
		out = append(out, ContainerSummary{
			ID:        it.ID,
			Names:     it.Names,
			Image:     it.Image,
			State:     it.State,
			Status:    it.Status,
			CreatedAt: created,
		})
	}
	return out, nil
}
