package goutils

import (
	"os"
	"path/filepath"
	"testing"
)

const testFolder = "tests"

func TestParseStringRange(t *testing.T) {
	tests := []struct {
		input  string
		output []int
	}{
		{"1-3", []int{1, 2, 3}},
		{"1-10", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, test := range tests {
		output, err := ParseStringRange(test.input)
		if err != nil {
			t.Errorf("Error parsing ranges: %v", err)
		}
		if len(output) != len(test.output) {
			t.Errorf("Expected %v, got %v", test.output, output)
		}
		for i, v := range output {
			if v != test.output[i] {
				t.Errorf("Expected %v, got %v", test.output, output)
			}
		}
	}

}

func TestParseStringRanges(t *testing.T) {
	tests := []struct {
		inputs []string
		output []int
	}{
		{[]string{"1-3", "4-6"}, []int{1, 2, 3, 4, 5, 6}},
		{[]string{"1-3", "4-6", "7-9"}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	result, err := ParseStringRanges(tests[0].inputs)
	if err != nil {
		t.Errorf("Error parsing ranges: %v", err)
	}
	if len(result) != len(tests[0].output) {
		t.Errorf("Expected %v, got %v", tests[0].output, result)
	}

}

func TestChangePermission(t *testing.T) {
	// Create a folder
	err := os.MkdirAll(testFolder, 0777)
	if err != nil {
		t.Errorf("Error creating folder: %v", err)
	}
	// Create a file
	err = WriteToFile(filepath.Join(testFolder, "test.txt"), "Hello, World!")
	if err != nil {
		t.Errorf("Error writing to file: %v", err)
	}
	err = ChangePermission(filepath.Join(testFolder, "test.txt"), 0777)
	if err != nil {
		t.Errorf("Error changing permission: %v", err)
	}

}
func TestWriteToFile(t *testing.T) {
	err := WriteToFile(filepath.Join(testFolder, "test.txt"), "Hello, World!")
	if err != nil {
		t.Errorf("Error writing to file: %v", err)
	}

}
func TestSlug(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"Hello, World!", "hello-world"},
		{"Hello, World! 123", "hello-world-123"},
	}

	for _, test := range tests {
		output := Slug(test.input)
		if output != test.output {
			t.Errorf("Expected %s, got %s", test.output, output)
		}
	}
}
