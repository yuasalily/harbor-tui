package ports

import (
	"context"
	"time"
)

type DockerAPI interface {
	Info(ctx context.Context) (DaemonInfo, error)
	ImagesList(ctx context.Context, opts ImagesListOptions) ([]ImageSummary, error)
	ContainersList(ctx context.Context, opts ContainersListOptions) ([]ContainerSummary, error)
	ImageRemove(ctx context.Context, ref string, opts ImageRemoveOptions) error
}

type DaemonInfo struct {
	Version    string
	OS         string
	APIVersion string
}

type ImageSummary struct {
	ID        string
	RepoTags  []string
	Size      int64
	CreatedAt time.Time
}

type ImagesListOptions struct {
	All       bool  //中間イメージを含む
	Dangling  *bool // nil=未指定, true/falseをfiltersに渡す
	Reference string
}

type ImageRemoveOptions struct {
	Force bool
	PruneChildlen bool
}

type ContainerSummary struct {
	ID        string
	Names     []string
	Image     string
	State     string // created, running, ...
	Status    string // "UP 10 minutes", ...
	CreatedAt time.Time
}

type ContainersListOptions struct {
	All    bool
	Name   string
	Status string
}
