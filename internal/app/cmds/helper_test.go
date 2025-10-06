package cmds

import (
	"context"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type fakeAPI struct {
	info ports.DaemonInfo
	imgs []ports.ImageSummary
	err  error
}

// Functional Options Pattern
type Option func(*fakeAPI)

func WithInfo(info ports.DaemonInfo) Option {
	return func(f *fakeAPI) { f.info = info }
}

func WithImages(imgs []ports.ImageSummary) Option {
	return func(f *fakeAPI) { f.imgs = imgs }
}

func WithError(err error) Option {
	return func(f *fakeAPI) { f.err = err }
}

func NewFakeAPI(opts ...Option) *fakeAPI {
	f := &fakeAPI{}
	for _, o := range opts {
		o(f)
	}
	return f
}

func (f *fakeAPI) Info(ctx context.Context) (ports.DaemonInfo, error) {
	if f.err != nil {
		return ports.DaemonInfo{}, f.err
	}
	return f.info, nil
}

func (f *fakeAPI) ImagesList(ctx context.Context, opts ports.ImagesListOptions) ([]ports.ImageSummary, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.imgs, nil
}
