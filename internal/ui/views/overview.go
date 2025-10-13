package views

import (
	"fmt"
	"strings"
)

func RenderOverview(status, info string) string {
	var b strings.Builder
	b.WriteString("  Overview\n\n")
	b.WriteString(fmt.Sprintf("  Status: %s\n", status))
	if info != "" {
		b.WriteString(fmt.Sprintf("  %s\n", info))
	}
	return b.String()
}
