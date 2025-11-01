package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuasalily/harbor-tui/internal/ui/components"
)

type Model struct {
	Visible bool
	W, H    int

	Title string
	Body  string
	Hint  string

	BorderColor lipgloss.TerminalColor
	Payload     any
}

func (m *Model) SetSize(w, h int) { m.W, m.H = w, h }

func (m *Model) OpenConfirm(title, body, hint string, border lipgloss.TerminalColor, payload any) {
	m.Visible = true
	m.Title = title
	m.Body = body
	m.Hint = hint
	m.BorderColor = border
	m.Payload = payload
}

func (m *Model) Close() { m.Visible = false }

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "y", "Y":
			payload := m.Payload
			m.Close()
			return *m, func() tea.Msg { return DialogResultMsg{Confirmed: true, Payload: payload} }
		case "n", "N", "esc":
			payload := m.Payload
			m.Close()
			return *m, func() tea.Msg { return DialogResultMsg{Confirmed: false, Payload: payload} }
		}
	}
	return *m, nil
}

func (m Model) View() string {
	if !m.Visible {
		return ""
	}

	avail := max(m.W-6, 20)
	w := max(36, min(64, avail)) // 36 - 64の範囲でクランプ
	return components.RenderDialog(w, m.Title, m.Body, m.Hint, "  ", m.BorderColor)
}
