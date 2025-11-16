package components

import (
	"github.com/charmbracelet/bubbles/list"
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

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(1)

	l := list.New(it, delegate, width, height)
	
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return l
}
