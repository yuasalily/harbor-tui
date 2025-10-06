package dockeradapter

import (
	"context"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/docker"
)

type dockerClient interface {
	Info(ctx context.Context) (docker.DaemonInfo, error)
	ImagesList(ctx context.Context, opts docker.ImagesListOptions) ([]docker.ImageSummary, error)
}

type Adapter struct {
	c dockerClient
}

func New(c dockerClient) *Adapter { return &Adapter{c: c} }

func NewFromEnv() (*Adapter, error) {
	cli, err := docker.New()
	if err != nil {
		return nil, err
	}
	return &Adapter{c: cli}, nil
}

func (a *Adapter) Info(ctx context.Context) (ports.DaemonInfo, error) {
	info, err := a.c.Info(ctx)
	if err != nil {
		return ports.DaemonInfo{}, err
	}
	return ports.DaemonInfo{Version: info.Version, OS: info.OS, APIVersion: info.APIVersion}, nil
}

func (a *Adapter) ImagesList(ctx context.Context, opts ports.ImagesListOptions) ([]ports.ImageSummary, error) {
	imgs, err := a.c.ImagesList(ctx, docker.ImagesListOptions{
		All:       opts.All,
		Dangling:  opts.Dangling,
		Reference: opts.Reference,
	})
	if err != nil {
		return nil, err
	}
	out := make([]ports.ImageSummary, 0, len(imgs))
	for _, it := range imgs {
		out = append(out, ports.ImageSummary{
			ID:        it.ID,
			RepoTags:  it.RepoTags,
			Size:      it.Size,
			CreatedAt: it.CreatedAt,
		})
	}
	return out, nil
}
