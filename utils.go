package goutils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ConvertBytes converts bytes to a human-readable string with the appropriate unit (bytes, MiB, GiB, TiB, PiB, or EiB).
func ConvertBytes(bytes uint64) string {
	const (
		KiB = 1024
		MiB = KiB * 1024
		GiB = MiB * 1024
		TiB = GiB * 1024
		PiB = TiB * 1024
		EiB = PiB * 1024
	)

	switch {
	case bytes >= EiB:
		return fmt.Sprintf("%.2f EB", float64(bytes)/float64(EiB))
	case bytes >= PiB:
		return fmt.Sprintf("%.2f PB", float64(bytes)/float64(PiB))
	case bytes >= TiB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TiB))
	case bytes >= GiB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GiB))
	case bytes >= MiB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MiB))
	case bytes >= KiB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KiB))
	default:
		return fmt.Sprintf("%d Bytes", bytes)
	}
}

// ConvertToBytes converts a string with a size suffix (e.g., "1M", "1Mi", "1MB") to bytes.
func ConvertToBytes(input string) (int64, error) {
	// Define the mapping for binary (Mi) and decimal (M) units
	binaryUnits := map[string]int64{
		"Ki": 1024,
		"Mi": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024,
		"Ti": 1024 * 1024 * 1024 * 1024,
		"Pi": 1024 * 1024 * 1024 * 1024 * 1024,
		"Ei": 1024 * 1024 * 1024 * 1024 * 1024 * 1024,
	}
	decimalUnits := map[string]int64{
		"K":  1000,
		"KB": 1000,
		"M":  1000 * 1000,
		"MB": 1000 * 1000,
		"G":  1000 * 1000 * 1000,
		"GB": 1000 * 1000 * 1000,
		"T":  1000 * 1000 * 1000 * 1000,
		"TB": 1000 * 1000 * 1000 * 1000,
		"P":  1000 * 1000 * 1000 * 1000 * 1000,
		"PB": 1000 * 1000 * 1000 * 1000 * 1000,
		"E":  1000 * 1000 * 1000 * 1000 * 1000 * 1000,
		"EB": 1000 * 1000 * 1000 * 1000 * 1000 * 1000,
	}

	// Extract the numeric part and the unit
	var numberPart string
	var unitPart string

	for i, r := range input {
		if r < '0' || r > '9' {
			numberPart = input[:i]
			unitPart = input[i:]
			break
		}
	}

	// Handle case where no valid unit is found
	if unitPart == "" {
		return 0, fmt.Errorf("invalid format: no unit provided")
	}

	// Convert the numeric part to an integer
	value, err := strconv.ParseInt(numberPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %w", err)
	}

	// Determine the multiplier
	var multiplier int64
	if strings.HasSuffix(unitPart, "i") {
		// Binary units
		multiplier, err = findMultiplier(unitPart, binaryUnits)
	} else {
		// Decimal units
		multiplier, err = findMultiplier(unitPart, decimalUnits)
	}

	if err != nil {
		return 0, err
	}

	// Calculate the bytes
	return value * multiplier, nil
}

// Helper function to find the multiplier for a given unit
func findMultiplier(unit string, units map[string]int64) (int64, error) {
	multiplier, ok := units[unit]
	if !ok {
		return 0, fmt.Errorf("invalid unit: %s", unit)
	}
	return multiplier, nil
}

// IsCIDR checks if the input is a valid CIDR notation
func IsCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// IsIPOrCIDR determines whether the input is an IP address or a CIDR
func IsIPOrCIDR(input string) (isIP bool, isCIDR bool) {
	// Check if it's a valid IP address
	if net.ParseIP(input) != nil {
		return true, false
	}

	// Check if it's a valid CIDR
	if _, _, err := net.ParseCIDR(input); err == nil {
		return false, true
	}

	// Neither IP nor CIDR
	return false, false
}

// IsIPAddress checks if the input is a valid IP address (IPv4 or IPv6)
func IsIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

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
		return fmt.Sprintf(format, float64(d.Microseconds())/1000, "ms")
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

// DeepCopy copies the fromValue struct to the toValue struct by marshalling and unmarshalling the data.
func DeepCopy(toValue interface{}, fromValue interface{}) error {
	// Ensure toValue is a pointer to a struct
	val := reflect.ValueOf(toValue)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("toValue must be a pointer to a struct")
	}

	// Marshal the fromValue to JSON, then unmarshal into the toValue
	data, err := json.Marshal(fromValue)
	if err != nil {
		return fmt.Errorf("failed to marshal fromValue: %v", err)
	}

	err = json.Unmarshal(data, toValue)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to toValue struct: %v", err)
	}

	return nil
}
