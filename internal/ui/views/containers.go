package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

func ContainerColumns(totalWidth int) []table.Column {
	if totalWidth < 60 {
		totalWidth = 60
	}
	idW, nameW, imageW, stateusW, createdW := 14, 18, 18, 20, 16
	fixed := idW + nameW + imageW + stateusW + createdW + 4
	if totalWidth > fixed {
		nameW += totalWidth - fixed
	}
	return []table.Column{
		{Title: "ID", Width: idW},
		{Title: "NAME", Width: nameW},
		{Title: "IMAGE", Width: imageW},
		{Title: "STATUS", Width: stateusW},
		{Title: "CREATED", Width: createdW},
	}
}

func ApplyContainers(t *table.Model, items []ports.ContainerSummary) {
	rows := make([]table.Row, 0, len(items))
	for _, it := range items {
		name := ""
		if len(it.Names) > 0 {
			name = strings.TrimPrefix(it.Names[0], "/")
		}
		rows = append(rows, table.Row{
			shortID(it.ID),
			name,
			it.Image,
			it.Status,
			it.CreatedAt.Local().Format("2006-01-02 15:04"),
		})
	}
	t.SetRows(rows)
}

func RenderContainers(t table.Model, indent string) string {
	var b strings.Builder
	for ln := range strings.SplitSeq(t.View(), "\n") {
		b.WriteString(indent)
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}
