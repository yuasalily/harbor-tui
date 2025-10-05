package docker

import (
	"context"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

// ImageListのDTO
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

// Dockerデーモンからイメージ一覧を取得してDTOに変換
func (c *Client) ImagesList(ctx context.Context, opts ImagesListOptions) ([]ImageSummary, error) {
	args := filters.NewArgs()
	if opts.Dangling != nil {
		args.Add("dangling", strconv.FormatBool(*opts.Dangling))
	}
	if opts.Reference != "" {
		args.Add("refenrece", opts.Reference)
	}

	listOpts := image.ListOptions{All: opts.All, Filters: args}
	imgs, err := c.cli.ImageList(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	res := make([]ImageSummary, 0, len(imgs))
	for _, it := range imgs {
		res = append(res, ImageSummary{
			ID: it.ID,
			RepoTags: it.RepoTags,
			Size: it.Size,
			CreatedAt: time.Unix(it.Created, 0),
		})
	}
	return res, nil
}
