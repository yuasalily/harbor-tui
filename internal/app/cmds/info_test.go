package cmds

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type fakeAPI struct {
	info ports.DaemonInfo
	err  error
}

func (f *fakeAPI) Info(ctx context.Context) (ports.DaemonInfo, error) {
	if f.err != nil {
		return ports.DaemonInfo{}, f.err
	}
	return f.info, nil
}

func (f *fakeAPI) ImagesList(ctx context.Context, opts ports.ImagesListOptions) ([]ports.ImageSummary, error) {
	return nil, errors.New("not used in this test")
}

func TestFetchDaemonInfoCmd_OK(t *testing.T) {
	api := &fakeAPI{info: ports.DaemonInfo{Version: "27.2.0", OS: "linux"}}
	cmd := FetchDaemonInfoCmd(api, time.Second)
	msg := cmd()
	m := msg.(DaemonInfoMsg)
	if m.Err != nil || m.Info.Version != "27.2.0" || m.Info.OS != "linux" {
		t.Fatalf("unexpected msg: %#v", m)
	}
}

func TestFetchDaemonInfoCmd_Error(t *testing.T) {
	api := &fakeAPI{err: errors.New("boom")}
	cmd := FetchDaemonInfoCmd(api, time.Second)
	msg := cmd()
	m := msg.(DaemonInfoMsg)
	if m.Err == nil {
		t.Fatalf("expected error, got nil")
	}
}
