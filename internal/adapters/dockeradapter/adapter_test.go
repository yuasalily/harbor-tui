package dockeradapter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
	"github.com/yuasalily/harbor-tui/internal/docker"
)

type fakeClient struct {
	info docker.DaemonInfo
	imgs []docker.ImageSummary
	err  error
}

func (f *fakeClient) Info(ctx context.Context) (docker.DaemonInfo, error) {
	if f.err != nil {
		return docker.DaemonInfo{}, f.err
	}
	return f.info, nil
}

func (f *fakeClient) ImagesList(ctx context.Context, opts docker.ImagesListOptions) ([]docker.ImageSummary, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.imgs, nil
}

func TestAdapter_Info_OK(t *testing.T) {
	fa := &fakeClient{info: docker.DaemonInfo{Version: "27.2.0", OS: "linux", APIVersion: "1.47"}}
	a := New(fa)
	out, err := a.Info(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out.Version != "27.2.0" || out.OS != "linux" || out.APIVersion != "1.47" {
		t.Fatalf("unexpected info: %#v", out)
	}
}

func TestAdapter_Info_Error(t *testing.T) {
	a := New(&fakeClient{err: errors.New("boom")})
	if _, err := a.Info(context.Background()); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestAdapter_ImagesList_OK(t *testing.T) {
	fa := &fakeClient{imgs: []docker.ImageSummary{
		{ID: "sha256:a", RepoTags: []string{"alpine:latest"}, Size: 10, CreatedAt: time.Unix(1000, 0)},
		{ID: "sha256:b", RepoTags: []string{"busybox:1.36"}, Size: 20, CreatedAt: time.Unix(2000, 0)},
	}}
	a := New(fa)
	out, err := a.ImagesList(context.Background(), ports.ImagesListOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(out) != 2 || out[0].ID != "sha256:a" || out[1].RepoTags[0] != "busybox:1.36" {
		t.Fatalf("unexpected out: %#v", out)
	}
}

func TestAdapter_ImagesList_Error(t *testing.T) {
	a := New(&fakeClient{err: errors.New("boom")})
	if _, err := a.ImagesList(context.Background(), ports.ImagesListOptions{}); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
