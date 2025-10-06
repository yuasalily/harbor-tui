package cmds

import (
	"errors"
	"testing"
	"time"

	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

func TestFetchImagesCmd_Ok(t *testing.T) {
	api := NewFakeAPI(WithImages([]ports.ImageSummary{{ID: "sha256:a"}, {ID: "sha256:b"}}))
	cmd := FetchImagesCmd(api, ports.ImagesListOptions{All: true}, time.Second)
	msg := cmd()
	m := msg.(ImagesListedMsg)
	if m.Err != nil {
		t.Fatalf("unexpeted err: %v", m.Err)
	}
	if len(m.Items) != 2 || m.Items[0].ID != "sha256:a" {
		t.Fatalf("unexpected items: %#v", m.Items)
	}
}

func TestFetchImagesCmd_Error(t *testing.T) {
	api := NewFakeAPI(WithError(errors.New("boom")))
	cmd := FetchImagesCmd(api, ports.ImagesListOptions{All: true}, time.Second)
	msg := cmd().(ImagesListedMsg)
	if msg.Err == nil {
		t.Fatalf("expected error, got nil")
	}
}
