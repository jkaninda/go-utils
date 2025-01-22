package goutils

import (
	"fmt"
	"time"
)

// FormatDuration formats the duration to either "X.Xms", "X.Xs", "X.Xm", or "X.Xh"
// with a customizable number of decimal places.
func FormatDuration(d time.Duration, decimalCount int) string {
	// Create the format string dynamically based on the decimalCount
	format := fmt.Sprintf("%%.%df%%s", decimalCount)

	switch {
	case d < time.Millisecond:
		// Format as nanoseconds
		return fmt.Sprintf("%dns", d.Nanoseconds())
	case d < time.Second:
		// Format as milliseconds
		return fmt.Sprintf(format, float64(d.Milliseconds()), "ms")
	case d < time.Minute:
		// Format as seconds
		return fmt.Sprintf(format, d.Seconds(), "s")
	case d < time.Hour:
		// Format as minutes
		return fmt.Sprintf(format, d.Minutes(), "m")
	default:
		// Format as hours
		return fmt.Sprintf(format, d.Hours(), "h")
	}
}

// ParseDuration parses the duration string and returns the duration
func ParseDuration(durationStr string) (time.Duration, error) {
	if durationStr == "" {
		return 0, nil
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
