package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	dialogTitleStyle = lipgloss.NewStyle().Bold(true)
	dialogBoxStyle   = lipgloss.NewStyle().
				Padding(1, 2).MarginTop(1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("212"))
	dialogHintStyle = lipgloss.NewStyle().Faint(true)
)

func RenderDialog(w int, title, body, hint string) string {
	if w < 20 {
		w = 20
	}
	innerW := w - 4

	var sb strings.Builder
	if title != "" {
		sb.WriteString(dialogTitleStyle.Render(title))
		sb.WriteByte('\n')
	}
	if body != "" {
		sb.WriteString(body)
		sb.WriteByte('\n')
	}
	if hint != "" {
		sb.WriteString(dialogHintStyle.Render(hint))
	}
	content := sb.String()
	lines := strings.Split(content, "\n")
	for i, ln := range lines {
		if runeCount(ln) > innerW {
			lines[i] = truncate(ln, innerW)
		}
	}
	content = strings.Join(lines, "\n")
	return dialogBoxStyle.Width(w).Render(content)
}

func runeCount(s string) int { return lipgloss.Width(s) }

func truncate(s string, w int) string {
	if runeCount(s) <= w {
		return s
	}
	ellipsis := "..."
	for runeCount(s)+runeCount(ellipsis) > w {
		s = s[:len(s)-1]
	}
	return s + ellipsis
}
