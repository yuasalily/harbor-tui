package views

import (
	"fmt"
	"strings"
)

func RenderInfo(status, info string) string {
	var b strings.Builder
	b.WriteString("  Overview\n\n")
	b.WriteString(fmt.Sprintf("  Status: %s\n", status))
	if info != "" {
		b.WriteString(fmt.Sprintf("  %s\n", info))
	}
	return b.String()
}

func RenderOverview(status, info string) string {
	return RenderInfo(status, info)
}
