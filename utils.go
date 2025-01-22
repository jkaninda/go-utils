package goutils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// FileExists checks if the file does exist
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if the folder does exist
func FolderExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()

}

// GetStringEnvWithDefault gets the value of an environment variable or returns a default value
func GetStringEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// GetIntEnv gets the value of an environment variable or returns an empty string
func GetIntEnv(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue

	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue

	}
	return i

}

// GetBoolEnv gets the value of an environment variable or returns an empty string
func GetBoolEnv(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue

	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return b

}

// SetEnv Set env
func SetEnv(name, value string) {
	if len(value) != 0 {
		err := os.Setenv(name, value)
		if err != nil {
			return
		}
	}

}

// MergeSlices merges two slices of strings
func MergeSlices(slice1, slice2 []string) []string {
	return append(slice1, slice2...)
}

// ParseURLPath removes duplicated [//]
//
// Ensures the path starts with a single leading slash
func ParseURLPath(urlPath string) string {
	// Replace any double slashes with a single slash
	urlPath = strings.ReplaceAll(urlPath, "//", "/")

	// Ensure the path starts with a single leading slash
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	return urlPath
}

// ParseRoutePath parses the route path and returns the blocked path
func ParseRoutePath(path, blockedPath string) string {
	basePath := ParseURLPath(path)
	switch {
	case blockedPath == "":
		return basePath
	case strings.HasSuffix(blockedPath, "/*"):
		return basePath + blockedPath[:len(blockedPath)-2]
	case strings.HasSuffix(blockedPath, "*"):
		return basePath + blockedPath[:len(blockedPath)-1]
	default:
		return basePath + blockedPath
	}
}

// IsJson checks if the given string is valid JSON
func IsJson(s string) bool {
	var js interface{}
	err := json.Unmarshal([]byte(s), &js)
	return err == nil
}

// UrlParsePath parses the URL path and returns the path
func UrlParsePath(uri string) string {
	parse, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return parse.Path
}

// HasWhitespace checks if the string contains whitespace
func HasWhitespace(s string) bool {
	return regexp.MustCompile(`\s`).MatchString(s)
}

// Slug converts a string to a URL-friendly slug
func Slug(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Replace spaces and special characters with hyphens
	re := regexp.MustCompile(`\W+`)
	text = re.ReplaceAllString(text, "-")

	// Remove leading and trailing hyphens
	text = strings.Trim(text, "-")

	return text
}

// TruncateText truncates the text to the specified limit and appends "..." if the text exceeds the limit
func TruncateText(text string, limit int) string {
	if len(text) > limit {
		return text[:limit] + "..."
	}
	return text
}

// ParseStringRanges converts a list of range strings to a slice of integers
func ParseStringRanges(rangeStrings []string) ([]int, error) {
	var result []int
	for _, rs := range rangeStrings {
		// Parse the range string
		r, err := ParseStringRange(rs)
		if err != nil {
			return nil, err
		}
		// Append the parsed range to the result slice
		result = append(result, r...)
	}
	return result, nil
}

// ParseStringRange converts a range string to a slice of integers
func ParseStringRange(rs string) ([]int, error) {
	var result []int

	// Check if the string contains a range (indicated by a hyphen)
	if strings.Contains(rs, "-") {
		// Split the range string by the delimiter (hyphen)
		parts := strings.Split(rs, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", rs)
		}

		// Convert the start and end of the range to integers
		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start value in range: %s", rs)
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end value in range: %s", rs)
		}

		// Ensure the start is less than or equal to the end
		if start > end {
			return nil, fmt.Errorf("start value is greater than end value in range: %s", rs)
		}

		// Append all integers in the range to the result slice
		for i := start; i <= end; i++ {
			result = append(result, i)
		}
	} else {
		// If it's a single integer, convert it directly
		num, err := strconv.Atoi(strings.TrimSpace(rs))
		if err != nil {
			return nil, fmt.Errorf("invalid integer value: %s", rs)
		}
		result = append(result, num)
	}

	return result, nil
}

// CopyFile copies a file from the source to the destination
func CopyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			return
		}
	}(sourceFile)

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			return

		}
	}(destinationFile)

	// Copy the content from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// Flush the buffer to ensure all data is written
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %v", err)
	}

	return nil
}

// ChangePermission changes the permission of the file
func ChangePermission(filePath string, mod int) error {
	if err := os.Chmod(filePath, fs.FileMode(mod)); err != nil {
		return err
	}
	return nil

}

// IsDirEmpty checks if the directory is empty
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	_, err = f.Readdirnames(1)
	if err == nil {
		return false, nil
	}
	return true, nil
}

// WriteToFile writes the given content to a file at the specified filePath.
// It returns an error if the operation fails.
func WriteToFile(filePath, content string) error {
	// Use os.WriteFile to handle file creation, writing, and closing in one call.
	// The file is created with read-write permissions for the user (0600).
	err := os.WriteFile(filePath, []byte(content), 0600)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}
