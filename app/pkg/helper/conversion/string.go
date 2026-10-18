package conversion

import (
	"fmt"
	"time"
)

func BeautifyDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	switch {
	case hours > 0 && minutes > 0:
		return fmt.Sprintf("%dh %dmin", hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh", hours)
	default:
		return fmt.Sprintf("%dmin", minutes)
	}
}
