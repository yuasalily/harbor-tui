package images

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Refresh key.Binding
	Delete  key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		Up:  key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete image"),
		),
	}
}
