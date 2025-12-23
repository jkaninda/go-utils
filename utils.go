package goutils

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	binaryUnits = map[string]int64{
		"Ki": 1024, "KiB": 1024,
		"Mi": 1024 * 1024, "MiB": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024, "GiB": 1024 * 1024 * 1024,
		"Ti": 1024 * 1024 * 1024 * 1024, "TiB": 1024 * 1024 * 1024 * 1024,
		"Pi": 1024 * 1024 * 1024 * 1024 * 1024, "PiB": 1024 * 1024 * 1024 * 1024 * 1024,
		"Ei": 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "EiB": 1024 * 1024 * 1024 * 1024 * 1024 * 1024,
	}
	decimalUnits = map[string]int64{
		"K": 1000, "KB": 1000,
		"M": 1000 * 1000, "MB": 1000 * 1000,
		"G": 1000 * 1000 * 1000, "GB": 1000 * 1000 * 1000,
		"T": 1000 * 1000 * 1000 * 1000, "TB": 1000 * 1000 * 1000 * 1000,
		"P": 1000 * 1000 * 1000 * 1000 * 1000, "PB": 1000 * 1000 * 1000 * 1000 * 1000,
		"E": 1000 * 1000 * 1000 * 1000 * 1000 * 1000, "EB": 1000 * 1000 * 1000 * 1000 * 1000 * 1000,
	}
	defaultErrorWriter io.Writer = os.Stderr
	// Matches ${VAR_NAME} or {VAR_NAME}
	envPattern = regexp.MustCompile(`\$?\{([A-Za-z_][A-Za-z0-9_]*)\}`)

	// Matches {{function()}} or {{function(args)}}
	funcPattern = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\(([^)]*)\)\}\}`)
	// Extract the numeric part and the unit
	numberPart string
	unitPart   string
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

// ConvertToBytes converts a string with a size suffix (e.g., "1M", "1Mi", "1MiB", "1MB") to bytes.
func ConvertToBytes(input string) (int64, error) {
	if input == "" {
		return 0, errors.New("input cannot be empty")
	}

	numPart, unitPart := splitNumberAndUnit(input)
	if numPart == "" || unitPart == "" {
		return 0, errors.New("invalid format: missing number or unit")
	}

	value, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %w", err)
	}

	multiplier, err := findMultiplier(unitPart)
	if err != nil {
		return 0, err
	}

	return value * multiplier, nil
}

// splitNumberAndUnit extracts the numeric part and the unit from the input string.
func splitNumberAndUnit(input string) (string, string) {
	for i, r := range input {
		if r < '0' || r > '9' {
			return input[:i], input[i:]
		}
	}
	return input, ""
}

// findMultiplier determines the correct multiplier based on the unit.
func findMultiplier(unit string) (int64, error) {
	if m, exists := binaryUnits[unit]; exists {
		return m, nil
	}
	if m, exists := decimalUnits[unit]; exists {
		return m, nil
	}
	return 0, fmt.Errorf("invalid unit: %s", unit)
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

// RealIP extracts the real IP address of the client from the HTTP Request.
func RealIP(r *http.Request) string {
	// Check the X-Forwarded-For header for the client IP.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the comma-separated list.
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check the X-Real-IP header as a fallback.
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	// Use the remote address if headers are not set.
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	// Return the raw remote address as a last resort.
	return r.RemoteAddr
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
// Deprecated: Use Env instead
func GetStringEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// GetIntEnv gets the value of an environment variable or returns an empty string
// Deprecated: Use EnvInt instead
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
// Deprecated: Use EnvBool instead
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

// SetEnv sets the value of an environment variable if the value is not empty
func SetEnv(name, value string) error {
	if len(value) != 0 {
		return os.Setenv(name, value)
	}
	return nil
}

// Env retrieves the value of the environment variable named by the key.
func Env(envName string, defaultValue string) string {
	if value, exists := os.LookupEnv(envName); exists {
		return value
	}
	return defaultValue
}

// EnvInt retrieves the integer value of the environment variable named by the key.
func EnvInt(envName string, defaultValue int) int {
	if value, exists := os.LookupEnv(envName); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// EnvBool retrieves the boolean value of the environment variable named by the key.
func EnvBool(envName string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(envName); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
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

// Base64Encode encodes the input string to base64
func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64Decode decodes the base64 encoded string
func Base64Decode(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// RemoveDuplicates removes duplicate elements from a slice of any comparable type.
// It preserves the original order of elements, keeping the first occurrence of each duplicate.
//
// Type Parameters:
//   - T: any type that supports equality comparison (comparable constraint)
//
// Parameters:
//   - elements: slice of elements to process
//
// Returns:
//   - []T: new slice with duplicates removed
func RemoveDuplicates[T comparable](elements []T) []T {
	encountered := make(map[T]bool)
	result := make([]T, 0, len(elements))

	for _, elem := range elements {
		if !encountered[elem] {
			encountered[elem] = true
			result = append(result, elem)
		}
	}

	return result
}

// LoadTLSConfig creates a TLS configuration from certificate and key files
// Parameters:
//   - certFile: Path to the certificate file (PEM format)
//   - keyFile: Path to the private key file (PEM format)
//   - caFile: Optional path to CA certificate file for client verification (set to "" to disable)
//   - clientAuth: Whether to require client certificate verification
//
// Returns:
//   - *tls.Config configured with the certificate and settings
//   - error if any occurred during loading
func LoadTLSConfig(certFile, keyFile, caFile string, clientAuth bool) (*tls.Config, error) {
	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Enforce minimum TLS version 1.2
	}

	// If caFile is provided, set up client certificate verification
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			_, _ = fmt.Fprintf(defaultErrorWriter, "Warning: failed to append CA certs from PEM")
		}

		config.ClientCAs = caCertPool
		if clientAuth {
			config.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			config.ClientAuth = tls.VerifyClientCertIfGiven
		}
	}

	return config, nil
}

// IsAValidAddr checks if the entrypoint address is valid.
// A valid entrypoint address should be in the format ":<port>" or "<IP>:<port>",
// where <IP> is a valid IP address and <port> is a valid port number (1-65535).
func IsAValidAddr(addr string) bool {
	// Split the addr into IP and port parts
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	// If the host is empty, it means the addr is in the format ":<port>"
	// Otherwise, validate the IP address
	if host != "" {
		ip := net.ParseIP(host)
		if ip == nil {
			return false
		}
	}

	// Convert the port string to an integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}

	// Check if the port is within the valid range
	if port < 1 || port > 65535 {
		return false
	}
	return true
}

// IsValidHTTPMethods checks if all strings in the input slice are valid HTTP methods.
// The check is case-insensitive (converts to uppercase for comparison), but returns
// the original invalid methods in the response.
//
// Parameters:
//   - methods: slice of strings to validate
//
// Returns:
//   - bool: true if all methods are valid, false otherwise
//   - []string: slice of original invalid methods found (empty if all are valid)
func IsValidHTTPMethods(methods ...string) (bool, []string) {
	// Standard HTTP methods
	validMethods := map[string]struct{}{
		http.MethodGet:     {},
		http.MethodHead:    {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodPatch:   {},
		http.MethodDelete:  {},
		http.MethodConnect: {},
		http.MethodOptions: {},
		http.MethodTrace:   {},
	}

	var invalidMethods []string

	for _, method := range methods {
		upperMethod := strings.ToUpper(method)
		if _, ok := validMethods[upperMethod]; !ok {
			invalidMethods = append(invalidMethods, method)
		}
	}

	if len(invalidMethods) > 0 {
		return false, invalidMethods
	}
	return true, nil
}

// NormalizeHTTPMethods converts all methods to uppercase and validates them.
// Returns the normalized methods if all are valid, or an error if any are invalid.
//
// Parameters:
//   - methods: slice of strings to normalize and validate
//
// Returns:
//   - []string: slice of uppercase methods if all are valid
//   - error: description of invalid methods if any are found
func NormalizeHTTPMethods(methods ...string) ([]string, error) {
	normalized := make([]string, len(methods))
	valid, invalid := IsValidHTTPMethods(methods...)

	if !valid {
		return nil, fmt.Errorf("invalid HTTP methods: %v", invalid)
	}

	for i, method := range methods {
		normalized[i] = strings.ToUpper(method)
	}

	return normalized, nil
}

// ReplaceEnvVars replaces environment variables and built-in functions
// Supports:
// - ${VAR_NAME} or {VAR_NAME} - environment variables
// - {{randomString(length)}} - random alphanumeric string
// - {{randomHex(length)}} - random hex string
// - {{uuid()}} - UUID v4
// - {{timestamp()}} - Unix timestamp in seconds
// - {{timestampMs()}} - Unix timestamp in milliseconds
// - {{date(format)}} - formatted date (format optional, default: RFC3339)
// - {{now()}} - current time in RFC3339 format
func ReplaceEnvVars(s string) string {
	if s == "" {
		return s
	}

	// First, replace function calls
	s = replaceFunctions(s)

	// Then, replace environment variables
	s = replaceEnvVariables(s)

	return s
}

func replaceEnvVariables(s string) string {
	if !envPattern.MatchString(s) {
		return s
	}

	return envPattern.ReplaceAllStringFunc(s, func(match string) string {
		submatch := envPattern.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		name := submatch[1]
		if val, ok := os.LookupEnv(name); ok {
			return val
		}
		return match
	})
}

func replaceFunctions(s string) string {
	if !funcPattern.MatchString(s) {
		return s
	}

	return funcPattern.ReplaceAllStringFunc(s, func(match string) string {
		submatch := funcPattern.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		funcName := submatch[1]
		args := ""
		if len(submatch) >= 3 {
			args = strings.TrimSpace(submatch[2])
		}

		result, err := executeFunction(funcName, args)
		if err != nil {
			return match
		}

		return result
	})
}

func executeFunction(funcName, args string) (string, error) {
	switch strings.ToLower(funcName) {
	case "randomstring":
		length := 32 // default
		if args != "" {
			if l, err := strconv.Atoi(args); err == nil && l > 0 && l <= 1024 {
				length = l
			} else {
				return "", fmt.Errorf("invalid length for randomString: %s", args)
			}
		}
		return randomString(length), nil

	case "randomhex":
		length := 32 // default
		if args != "" {
			if l, err := strconv.Atoi(args); err == nil && l > 0 && l <= 1024 {
				length = l
			} else {
				return "", fmt.Errorf("invalid length for randomHex: %s", args)
			}
		}
		return randomHex(length), nil

	case "uuid":
		return generateUUID(), nil

	case "timestamp":
		return strconv.FormatInt(time.Now().Unix(), 10), nil

	case "timestampms":
		return strconv.FormatInt(time.Now().UnixMilli(), 10), nil

	case "now":
		return time.Now().Format(time.RFC3339), nil

	case "date":
		format := time.RFC3339 // default
		if args != "" {
			format = parseTimeFormat(args)
		}
		return time.Now().Format(format), nil

	default:
		return "", fmt.Errorf("unknown function: %s", funcName)
	}
}

// randomString generates a random alphanumeric string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return ""
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b)
}

// randomHex generates a random hexadecimal string of specified length
func randomHex(length int) string {
	bytes := make([]byte, (length+1)/2)

	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	hexStr := hex.EncodeToString(bytes)
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}

	return hexStr
}

// generateUUID generates a UUID v4
func generateUUID() string {
	uuid := make([]byte, 16)

	if _, err := rand.Read(uuid); err != nil {
		return ""
	}

	// Set version (4) and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:])
}

// parseTimeFormat converts common format names to Go time format strings
func parseTimeFormat(format string) string {
	// Support common format aliases
	formats := map[string]string{
		"rfc3339":     time.RFC3339,
		"rfc822":      time.RFC822,
		"iso8601":     time.RFC3339,
		"unix":        "2006-01-02 15:04:05",
		"date":        "2006-01-02",
		"time":        "15:04:05",
		"datetime":    "2006-01-02 15:04:05",
		"kitchen":     time.Kitchen,
		"ansic":       time.ANSIC,
		"unixdate":    time.UnixDate,
		"rubydate":    time.RubyDate,
		"rfc850":      time.RFC850,
		"rfc1123":     time.RFC1123,
		"rfc1123z":    time.RFC1123Z,
		"rfc3339nano": time.RFC3339Nano,
	}

	if f, ok := formats[strings.ToLower(format)]; ok {
		return f
	}
	return format
}

// IsBase64 checks if the input is valid Base64-encoded content.
func IsBase64(input string) bool {
	if input == "" {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(input)
	return err == nil
}
