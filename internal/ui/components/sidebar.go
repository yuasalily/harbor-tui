package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type navItem string

func (n navItem) Title() string       { return string(n) }
func (n navItem) Description() string { return "" }
func (n navItem) FilterValue() string { return string(n) }

func NewSidebar(items []string, width, height int) list.Model {
	it := make([]list.Item, 0, len(items))
	for _, s := range items {
		it = append(it, navItem(s))
	}
	l := list.New(it, list.NewDefaultDelegate(), width, height)
	l.Title = "Pages"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true)
	return l
}
