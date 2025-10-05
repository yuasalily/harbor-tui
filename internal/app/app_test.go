package app

import (
	"context"
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type fakeAPI struct {
	info DaemonInfo
	err  error
}

func (f *fakeAPI) Info(ctx context.Context) (DaemonInfo, error) {
	if f.err != nil {
		return DaemonInfo{}, f.err
	}
	return f.info, nil
}

func (f *fakeAPI) ImagesList(ctx context.Context, opts ImagesListOptions) ([]ImageSummary, error) {
	return nil, errors.New("not implemented in this test")
}

func TestQuitKey(t *testing.T) {
	m := New(&fakeAPI{})
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatalf("expected a comamnd, got nil")
	}

	if msg := cmd(); msg != (tea.QuitMsg{}) {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

func TestInitFetchesDaemonInfo_OK(t *testing.T) {
	api := &fakeAPI{info: DaemonInfo{Version: "27.2.0", OS: "linux"}}
	m := New(api)
	msg := m.Init()()
	updated, _ := m.Update(msg)
	m2 := updated.(Model)
	if !m2.dockerOK {
		t.Fatalf("expected dockerOK=true")
	}

	if m2.serverVersion != "27.2.0" || m2.daemonPlatform != "linux" {
		t.Fatalf("unexpected versino/platform: %s/%s", m2.serverVersion, m2.daemonPlatform)
	}
}

func TestInitFetchesDaemonInfo_Error(t *testing.T) {
	api := &fakeAPI{err: errors.New("boom")}
	m := New(api)
	msg := m.Init()()
	updated, _ := m.Update(msg)
	m2 := updated.(Model)
	if m2.dockerOK {
		t.Fatalf("expected dockerOK=false")
	}

	if m2.dockerErr == "" {
		t.Fatalf("expected error message set")
	}
}
