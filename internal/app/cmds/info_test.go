package cmds

import (
	"errors"
	"testing"
	"time"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

func TestFetchDaemonInfoCmd_OK(t *testing.T) {
	api := NewFakeAPI(WithInfo(ports.DaemonInfo{Version: "27.2.0", OS: "linux"}))
	cmd := FetchDaemonInfoCmd(api, time.Second)
	msg := cmd()
	m := msg.(DaemonInfoMsg)
	if m.Err != nil || m.Info.Version != "27.2.0" || m.Info.OS != "linux" {
		t.Fatalf("unexpected msg: %#v", m)
	}
}

func TestFetchDaemonInfoCmd_Error(t *testing.T) {
	api := NewFakeAPI(WithError(errors.New("boom")))
	cmd := FetchDaemonInfoCmd(api, time.Second)
	msg := cmd()
	m := msg.(DaemonInfoMsg)
	if m.Err == nil {
		t.Fatalf("expected error, got nil")
	}
}
