package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/yuasalily/harbor-tui/internal/app/ports"
)

func ImageColumns(totalWidth int) []table.Column {
	if totalWidth < 60 {
		totalWidth = 60
	}
	idWidth := 14
	sizeWidth := 8
	createdWidth := 16
	mergin := 3

	fixed := idWidth + sizeWidth + createdWidth + mergin
	tagWidth := max(totalWidth - fixed, 20)
	return []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "TAG", Width: tagWidth},
		{Title: "SIZE", Width: sizeWidth},
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
			shortID(it.ID),
			tag,
			humanBytes(it.Size),
			it.CreatedAt.Local().Format("2006-01-02 15:04"),
		})
	}
	t.SetRows(rows)
}

func RenderIndented(t table.Model, indent string) string {
	var b strings.Builder
	for ln := range strings.SplitSeq(t.View(), "\n") {
		b.WriteString(indent)
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for n >= unit && exp < 4 {
		n /= unit
		div *= unit
		exp++
	}
	suffix := []string{"KB", "MB", "GB", "TB"}[exp-1]
	return fmt.Sprintf("%d%s", n, suffix)
}
