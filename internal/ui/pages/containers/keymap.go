package containers

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up, Down key.Binding
	Refresh  key.Binding
	Logs     key.Binding
	Remove   key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
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
		Logs: key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
		Remove: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove"),
		),
	}
}

func (k KeyMap) Short() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Refresh, k.Logs, k.Remove}
}


