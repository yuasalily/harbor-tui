package app

import (
	"context"
	"time"
)

type DockerAPI interface {
	Info(ctx context.Context) (DaemonInfo, error)
	ImagesList(ctx context.Context, opts ImagesListOptions) ([]ImageSummary, error)
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
