package views

import (
	"fmt"
	"strings"
)

func shortID(id string) string {
	id = strings.TrimPrefix(id, "sha256:")
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func humanBytes(n int64) string {
	const unit = 1024.0
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := unit, 0
	for v := float64(n) / unit; v >= unit && exp < 4; v /= unit {
		div *= unit
		exp++
	}
	v := float64(n) / div
	suffix := []string{"KB", "MB", "GB", "TB"}[exp]
	if v < 10 {
		return fmt.Sprintf("%.1f%s", v, suffix)
	}
	return fmt.Sprintf("%.0f%s", v, suffix)
}
