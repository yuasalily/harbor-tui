package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

func ImageColumns(totalWidth int) []table.Column {
	if totalWidth < 60 {
		totalWidth = 60
	}
	selWidth := 5
	idWidth := 14
	createdWidth := 16
	mergin := 3

	fixed := selWidth + idWidth + createdWidth + mergin
	tagWidth := max(totalWidth-fixed, 20)
	return []table.Column{
		{Title: "SEL", Width: selWidth},
		{Title: "ID", Width: idWidth},
		{Title: "TAG", Width: tagWidth},
		{Title: "CREATED", Width: createdWidth},
	}
}

func ApplyImages(t *table.Model, images []ports.ImageSummary) {
	rows := make([]table.Row, 0, len(images))
	for _, it := range images {
		tag := "<none>"
		if len(it.RepoTags) > 0 {
			tag = it.RepoTags[0]
		}
		rows = append(rows, table.Row{
			"", // 別で埋める
			shortID(it.ID),
			tag,
			it.CreatedAt.Local().Format("2006-01-02 15:04"),
		})
	}
	t.SetRows(rows)
}

func RenderImages(t table.Model, indent string) string {
	var b strings.Builder
	for ln := range strings.SplitSeq(t.View(), "\n") {
		b.WriteString(indent)
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}
