package ui

import "github.com/charmbracelet/bubbles/key"

type GlobalKeyMap struct {
	Quit   key.Binding
	Select key.Binding
	Tab    key.Binding
}

func NewGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "focus"),
		),
	}
}

func (g GlobalKeyMap) ShortForNav() []key.Binding {
	return []key.Binding{g.Quit, g.Select, g.Tab}
}

func (g GlobalKeyMap) ShortForPage() []key.Binding {
	return []key.Binding{g.Quit, g.Tab}
}

