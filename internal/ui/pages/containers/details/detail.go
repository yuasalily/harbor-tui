package details

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

type Dependencies struct {
	API ports.DockerAPI
}

type Detail interface {
	SetTarget(id string)

	Init() tea.Cmd
	Update(tea.Msg) (Detail, tea.Cmd)
	View() string
	SetSize(w, h int)
	SetFocused(bool)
	Close() tea.Cmd

	ShortKeys() []key.Binding
}
