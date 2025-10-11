package components

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true)
	cellStyle     = lipgloss.NewStyle()
	selectedStyle = lipgloss.NewStyle().Bold(true)
)

func NewTable(cols []table.Column, height int, focused bool) table.Model {
	if height <= 0 {
		height = 12
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(focused),
		table.WithHeight(height),
	)
	t.SetStyles(table.Styles{
		Header:   headerStyle,
		Cell:     cellStyle,
		Selected: selectedStyle,
	})
	return t
}
