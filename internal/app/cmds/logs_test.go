package cmds

import (
	"errors"
	"testing"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

func TestStartContainerLogsCmd_OK(t *testing.T) {
	api := NewFakeAPI()
	msg := StartContainerLogsCmd(api, ports.ContainerLogsOptions{ContainerID: "cid", Stdout: true})()
	started, ok := msg.(ContainerLogsStartedMsg)
	if !ok {
		t.Fatalf("unexpected msg type: %#v", msg)
	}
	m1 := NextContainerLogLine("cid", started.Ch, started.Done)()
	if _, ok := m1.(ContainerLogLineMsg); !ok {
		t.Fatalf("expected line msg, got %#v", m1)
	}
	m2 := NextContainerLogLine("cid", started.Ch, started.Done)()
	if _, ok := m2.(ContainerLogLineMsg); !ok {
		t.Fatalf("expected line msg, got %#v", m2)
	}
	m3 := NextContainerLogLine("cid", started.Ch, started.Done)()
	if _, ok := m3.(ContainerLogsEndedMsg); !ok {
		t.Fatalf("expected ended msg, got %#v", m3)
	}
}

func TestStartContainerLogsCmd_Error(t *testing.T) {
	api := NewFakeAPI(WithError(errors.New("boom")))
	msg := StartContainerLogsCmd(api, ports.ContainerLogsOptions{ContainerID: "cid", Stdout: true})()
	if _, ok := msg.(ContainerLogsEndedMsg); !ok {
		t.Fatalf("expected ended msg on error, got %#v", msg)
	}
}
