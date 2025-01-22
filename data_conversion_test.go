package goutils

import (
	"fmt"
	"testing"
)

func TestConvertBytes(t *testing.T) {
	byteSizes := []uint64{
		512,
		1024,
		2048,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024 * 1024 * 1024,
	}

	for _, size := range byteSizes {
		fmt.Printf("%d bytes -> %s\n", size, ConvertBytes(size))
	}
}

func TestConvertToBytes(t *testing.T) {
	byteSizes := []string{
		"1Ki",
		"1Mi",
		"1Gi",
		"1Ti",
		"1Pi",
		"1Ei",
	}

	for _, size := range byteSizes {
		bytes, err := ConvertToBytes(size)
		if err != nil {
			t.Errorf("Error converting %s to bytes: %v", size, err)
		}
		fmt.Printf("%s -> %d bytes\n", size, bytes)
	}
}
