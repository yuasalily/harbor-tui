package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type HelpBar struct {
	h help.Model
}

func NewHelpBar() HelpBar {
	h := help.New()
	h.ShowAll = false
	return HelpBar{h: h}
}

func (hb HelpBar) Render(bindings ...[]key.Binding) string {
	var flat []key.Binding
	for _, bs := range bindings {
		flat = append(flat, bs...)
	}
	if len(flat) == 0 {
		return ""
	}
	view := hb.h.ShortHelpView(flat)
	view = strings.ReplaceAll(view, "\n", " ")
	return lipgloss.NewStyle().Bold(true).Render(view)
}
