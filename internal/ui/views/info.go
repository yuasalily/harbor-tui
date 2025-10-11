package views

import (
	"fmt"
	"strings"
)

func RenderInfo(status, info string) string {
	var b strings.Builder
	b.WriteString("  Info\n\n")
	b.WriteString(fmt.Sprintf("  Status: %s\n", status))
	if info != "" {
		b.WriteString(fmt.Sprintf("  %s\n", info))
	}
	return b.String()
}