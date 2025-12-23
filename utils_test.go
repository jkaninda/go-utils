package goutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		"1KB", "1MB", "1GB", "1Ki", "125MiB", "1Mi", "1Gi", "1TB", "1PB", "1EB",
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

func TestBase64Encode(t *testing.T) {
	input := "Hello, World!"
	output := Base64Encode(input)
	fmt.Println(output)
}

func TestBase64Decode(t *testing.T) {
	input := "SGVsbG8sIFdvcmxkIQ=="
	output, err := Base64Decode(input)
	if err != nil {
		t.Errorf("Error decoding base64: %v", err)
	}
	fmt.Println(output)
}

func TestRemoveDuplicates(t *testing.T) {
	names := []string{"John", "Bob", "Jane", "Bob"}
	result := RemoveDuplicates(names)
	if len(result) == len(names) {
		t.Errorf("Failed to remove duplicates: %v", names)
	}
	fmt.Println(result)
}

func TestIsAValidAddr(t *testing.T) {
	addr := ":80"
	if !IsAValidAddr(addr) {
		t.Errorf("Address should be valid")
	}
}

func TestIsBase64(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"SGVsbG8sIFdvcmxkIQ==", true},
		{"InvalidBase64String", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsBase64(test.input)
		if result != test.output {
			t.Errorf("Expected %v, got %v for input %s", test.output, result, test.input)
		}
	}
}

func TestNormalizeHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "TRACE", "CONNECT", "HEAD"}
	for _, method := range methods {
		normalized, _ := NormalizeHTTPMethods(method)
		if normalized[0] != method {
			t.Errorf("Expected %s, got %s", method, normalized)
		}
	}

	// Test with lowercase methods
	for _, method := range methods {
		normalized, _ := NormalizeHTTPMethods(strings.ToLower(method))
		if normalized[0] != strings.ToUpper(method) {
			t.Errorf("Expected %s, got %s", strings.ToUpper(method), normalized)
		}
	}
}
func TestHasWhitespace(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"Hello World", true},
		{"NoWhitespace", false},
		{"   Leading and trailing spaces   ", true},
		{"\tTab character", true},
	}

	for _, test := range tests {
		result := HasWhitespace(test.input)
		if result != test.output {
			t.Errorf("Expected %v, got %v for input %s", test.output, result, test.input)
		}
	}
}
func TestIsJson(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{`{"name": "John", "age": 30}`, true},
		{`[1, 2, 3]`, true},
		{`"Hello, World!"`, true},
		{`{name: "John", age: 30}`, false}, // Invalid JSON
		{`Hello, World!`, false},           // Not JSON
	}

	for _, test := range tests {
		result := IsJson(test.input)
		if result != test.output {
			t.Errorf("Expected %v, got %v for input %s", test.output, result, test.input)
		}
	}
}
func TestReplaceEnvVars(t *testing.T) {
	// Set environment variables for testing
	err := os.Setenv("TEST_VAR1", "value1")
	if err != nil {
		return
	}
	err = os.Setenv("TEST_VAR2", "value2")
	if err != nil {
		return
	}

	tests := []struct {
		input  string
		output string
	}{
		{"This is a ${TEST_VAR1}", "This is a value1"},
		{"${TEST_VAR2} is here", "value2 is here"},
		{"No variables here", "No variables here"},
	}
	for _, test := range tests {
		output := ReplaceEnvVars(test.input)
		if output != test.output {
			t.Errorf("Expected %s, got %s", test.output, output)
		}
	}
}

func TestEnv(t *testing.T) {
	// Set environment variables for testing
	err := SetEnv("EXISTING_VAR", "exists")
	if err != nil {
		return
	}
	tests := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{"EXISTING_VAR", "default", "exists"},
		{"NON_EXISTING_VAR", "default", "default"},
	}

	for _, test := range tests {
		result := Env(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("For key %s, expected %s but got %s", test.key, test.expected, result)
		}
	}
}
func TestEnvInt(t *testing.T) {
	// Set environment variables for testing
	err := SetEnv("EXISTING_INT_VAR", "42")
	if err != nil {
		return
	}
	tests := []struct {
		key          string
		defaultValue int
		expected     int
	}{
		{"EXISTING_INT_VAR", 10, 42},
		{"NON_EXISTING_INT_VAR", 10, 10},
		{"INVALID_INT_VAR", 20, 20},
	}

	for _, test := range tests {
		result := EnvInt(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("For key %s, expected %d but got %d", test.key, test.expected, result)
		}
	}
}
func TestEnvBool(t *testing.T) {
	// Set environment variables for testing
	err := SetEnv("EXISTING_BOOL_VAR", "true")
	if err != nil {
		return
	}
	tests := []struct {
		key          string
		defaultValue bool
		expected     bool
	}{
		{"EXISTING_BOOL_VAR", false, true},
		{"NON_EXISTING_BOOL_VAR", true, true},
		{"INVALID_BOOL_VAR", false, false},
	}

	for _, test := range tests {
		result := EnvBool(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("For key %s, expected %v but got %v", test.key, test.expected, result)
		}
	}
}
