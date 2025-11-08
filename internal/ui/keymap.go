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
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "focus"),
		),
	}
}

func (g GlobalKeyMap) ShortForNav() []key.Binding {
	return []key.Binding{g.Quit, g.Tab}
}

func (g GlobalKeyMap) ShortForPage() []key.Binding {
	return []key.Binding{g.Quit, g.Tab}
}

func (g GlobalKeyMap) ShortForDialog() []key.Binding {
	// 確認ダイアログ中はQuitのみ表示(コマンドはダイアログ内に表示する)
	return []key.Binding{g.Quit}
}
