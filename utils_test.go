package goutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testFolder = "tests"

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
		fmt.Printf("%d Bytes -> %s\n", size, ConvertBytes(size))
	}
}

func TestConvertToBytes(t *testing.T) {
	byteSizes := []string{
		"1KB", "1MB", "1GB", "1Ki", "1Mi", "1Gi", "1TB", "1PB", "1EB",
	}

	for _, size := range byteSizes {
		bytes, err := ConvertToBytes(size)
		if err != nil {
			t.Errorf("Error converting %s to bytes: %v", size, err)
		}
		fmt.Printf("%s -> %d bytes\n", size, bytes)
	}
}
func TestValidateIPAddress(t *testing.T) {
	tests := []string{
		"192.168.1.100",
		"192.168.1.120",
	}
	for _, test := range tests {
		if IsIPAddress(test) {
			fmt.Println("Ip is valid")
		} else {
			fmt.Println("Ip is invalid")
		}
	}

}
func TestValidateIPOrCIDR(t *testing.T) {
	tests := []string{
		"192.168.1.100",
		"192.168.1.100",
		"192.168.1.100/32",
		"invalid-input",
		"192.168.1.100/33",
	}
	for _, test := range tests {
		isIP, isCIDR := IsIPOrCIDR(test)
		if isIP {
			fmt.Printf("%s is an IP address\n", test)
		} else if isCIDR {
			fmt.Printf("%s is a CIDR\n", test)
		} else {
			fmt.Printf("%s is neither an IP address nor a CIDR\n", test)
		}
	}

}
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

func TestDeepCopy(t *testing.T) {
	type Source struct {
		Name  string
		Age   int
		Email string
	}

	type Destination struct {
		Name  string
		Age   int
		Email string
	}
	src := Source{Name: "John", Age: 30, Email: "john@example.com"}
	dest := Destination{}

	err := DeepCopy(&dest, src)
	if err != nil {
		t.Errorf("Error copying struct: %v", err)
	} else {
		fmt.Printf("Destination: %+v\n", dest)
	}

}

func TestDeepCopyBetween(t *testing.T) {
	type Source struct {
		Name string
		Age  int
	}

	type Destination struct {
		Name  string
		Email string
	}

	src := Source{Name: "Bob", Age: 30}
	dest := Destination{Email: "bob@example.com"}

	err := DeepCopy(&dest, src)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Destination: %+v\n", dest)
	}
}
func TestCopyNested(t *testing.T) {
	type Address struct {
		City  string
		State string
	}

	type Source struct {
		Name    string
		Age     int
		Address Address
	}

	type Destination struct {
		Name    string
		Age     int
		Address Address
	}

	src := Source{
		Name: "Bennett",
		Age:  30,
		Address: Address{
			City:  "New York",
			State: "NY",
		},
	}

	dest := Destination{}

	err := DeepCopy(&dest, src)
	if err != nil {
		t.Errorf("Error copying struct: %v", err)
	} else {
		fmt.Printf("Destination: %+v\n", dest)
	}
	// Output: Destination: {Name:Dave Age:40 Address:{City:New York State:NY}}
}
