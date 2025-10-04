package app

import (
	"context"
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestQuitKey(t *testing.T) {
	m := New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatalf("expected a comamnd, got nil")
	}

	if msg := cmd(); msg != (tea.QuitMsg{}) {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

func TestDockerOK(t *testing.T) {
	fake := func(_ context.Context) (string, string, error) {
		return "27.2.0", "linux", nil
	}

	m := New(WithDockerChecker(fake))

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

func TestDcokerError(t *testing.T) {
	fake := func(_ context.Context) (string, string, error) {
		return "", "", errors.New("boom")
	}
	m := New(WithDockerChecker(fake))
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
