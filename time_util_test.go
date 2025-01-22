package goutils

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	now := time.Now()
	time.Sleep(2 * time.Second)
	duration := time.Since(now)
	fmt.Println(FormatDuration(duration, 2))

}

func TestParseDuration(t *testing.T) {
	durationStr := "2s"
	duration, err := ParseDuration(durationStr)
	if err != nil {
		t.Errorf("Error parsing duration: %v", err)
	}
	fmt.Println(duration)
}
