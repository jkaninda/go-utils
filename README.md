# Go Utils
A collection of Go utilities for various tasks.

[![Tests](https://github.com/jkaninda/go-utils/actions/workflows/test.yml/badge.svg)](https://github.com/jkaninda/go-utils/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jkaninda/go-utils)](https://goreportcard.com/report/github.com/jkaninda/go-utils)
[![Go Reference](https://pkg.go.dev/badge/github.com/jkaninda/go-utils.svg)](https://pkg.go.dev/github.com/jkaninda/go-utils)
## Installation

To install the package, run the following command:
```bash
go get github.com/jkaninda/go-utils
```

## Utility Functions Documentation

### 1. **File and Folder Utilities**

#### `FileExists(filename string) bool`
- **Purpose**: Checks if a file exists at the specified path.
- **Parameters**:
    - `filename`: The path to the file.
- **Returns**: `true` if the file exists and is not a directory, otherwise `false`.
- **Example**:
  ```go
  exists := FileExists("example.txt")
  fmt.Println(exists) // Output: true or false
  ```

#### `FolderExists(name string) bool`
- **Purpose**: Checks if a folder exists at the specified path.
- **Parameters**:
    - `name`: The path to the folder.
- **Returns**: `true` if the folder exists and is a directory, otherwise `false`.
- **Example**:
  ```go
  exists := FolderExists("example_folder")
  fmt.Println(exists) // Output: true or false
  ```

#### `IsDirEmpty(name string) (bool, error)`
- **Purpose**: Checks if a directory is empty.
- **Parameters**:
    - `name`: The path to the directory.
- **Returns**: `true` if the directory is empty, otherwise `false`. Returns an error if the directory cannot be accessed.
- **Example**:
  ```go
  isEmpty, err := IsDirEmpty("example_folder")
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println("Is directory empty?", isEmpty)
  }
  ```

#### `CopyFile(src, dst string) error`
- **Purpose**: Copies a file from the source path to the destination path.
- **Parameters**:
    - `src`: The source file path.
    - `dst`: The destination file path.
- **Returns**: An error if the operation fails.
- **Example**:
  ```go
  err := CopyFile("source.txt", "destination.txt")
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println("File copied successfully!")
  }
  ```

#### `ChangePermission(filePath string, mod int) error`
- **Purpose**: Changes the file permissions of a file.
- **Parameters**:
    - `filePath`: The path to the file.
    - `mod`: The permission mode (e.g., `0644`).
- **Returns**: An error if the operation fails.
- **Example**:
  ```go
  err := ChangePermission("example.txt", 0644)
  if err != nil {
      fmt.Println("Error:", err)
  }
  ```

#### `WriteToFile(filePath, content string) error`
- **Purpose**: Writes content to a file at the specified path.
- **Parameters**:
    - `filePath`: The path to the file.
    - `content`: The content to write.
- **Returns**: An error if the operation fails.
- **Example**:
  ```go
  err := WriteToFile("example.txt", "Hello, world!")
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println("File written successfully!")
  }
  ```

---

### 2. **Environment Variable Utilities**

#### `GetStringEnvWithDefault(key, defaultValue string) string`
- **Purpose**: Retrieves the value of an environment variable or returns a default value if the variable is not set.
- **Parameters**:
    - `key`: The environment variable key.
    - `defaultValue`: The default value to return if the key is not set.
- **Returns**: The value of the environment variable or the default value.
- **Example**:
  ```go
  value := GetStringEnvWithDefault("MY_ENV_VAR", "default_value")
  fmt.Println(value) // Output: The value of MY_ENV_VAR or "default_value"
  ```

#### `GetIntEnv(key string, defaultValue int) int`
- **Purpose**: Retrieves the value of an environment variable as an integer or returns a default value if the variable is not set or invalid.
- **Parameters**:
    - `key`: The environment variable key.
    - `defaultValue`: The default value to return if the key is not set or invalid.
- **Returns**: The integer value of the environment variable or the default value.
- **Example**:
  ```go
  value := GetIntEnv("MY_INT_ENV_VAR", 42)
  fmt.Println(value) // Output: The value of MY_INT_ENV_VAR or 42
  ```

#### `GetBoolEnv(key string, defaultValue bool) bool`
- **Purpose**: Retrieves the value of an environment variable as a boolean or returns a default value if the variable is not set or invalid.
- **Parameters**:
    - `key`: The environment variable key.
    - `defaultValue`: The default value to return if the key is not set or invalid.
- **Returns**: The boolean value of the environment variable or the default value.
- **Example**:
  ```go
  value := GetBoolEnv("MY_BOOL_ENV_VAR", true)
  fmt.Println(value) // Output: The value of MY_BOOL_ENV_VAR or true
  ```

#### `SetEnv(name, value string)`
- **Purpose**: Sets an environment variable.
- **Parameters**:
    - `name`: The environment variable name.
    - `value`: The value to set.
- **Example**:
  ```go
  SetEnv("MY_ENV_VAR", "my_value")
  ```

---

### 3. **String and Data Utilities**

#### `MergeSlices(slice1, slice2 []string) []string`
- **Purpose**: Merges two slices of strings into one.
- **Parameters**:
    - `slice1`: The first slice.
    - `slice2`: The second slice.
- **Returns**: A new slice containing all elements from `slice1` and `slice2`.
- **Example**:
  ```go
  result := MergeSlices([]string{"a", "b"}, []string{"c", "d"})
  fmt.Println(result) // Output: [a b c d]
  ```

#### `ParseURLPath(urlPath string) string`
- **Purpose**: Normalizes a URL path by removing duplicate slashes and ensuring it starts with a single slash.
- **Parameters**:
    - `urlPath`: The URL path to normalize.
- **Returns**: The normalized URL path.
- **Example**:
  ```go
  path := ParseURLPath("//example//path//")
  fmt.Println(path) // Output: /example/path/
  ```

#### `ParseRoutePath(path, blockedPath string) string`


#### `IsJson(s string) bool`
- **Purpose**: Checks if a string is valid JSON.
- **Parameters**:
    - `s`: The string to check.
- **Returns**: `true` if the string is valid JSON, otherwise `false`.
- **Example**:
  ```go
  valid := IsJson(`{"key": "value"}`)
  fmt.Println(valid) // Output: true
  ```

#### `UrlParsePath(uri string) string`
- **Purpose**: Extracts the path from a URL.
- **Parameters**:
    - `uri`: The URL to parse.
- **Returns**: The path component of the URL.
- **Example**:
  ```go
  path := UrlParsePath("https://example.com/path")
  fmt.Println(path) // Output: /path
  ```

#### `HasWhitespace(s string) bool`
- **Purpose**: Checks if a string contains any whitespace.
- **Parameters**:
    - `s`: The string to check.
- **Returns**: `true` if the string contains whitespace, otherwise `false`.
- **Example**:
  ```go
  hasSpace := HasWhitespace("hello world")
  fmt.Println(hasSpace) // Output: true
  ```

#### `Slug(text string) string`
- **Purpose**: Converts a string into a URL-friendly slug.
- **Parameters**:
    - `text`: The string to convert.
- **Returns**: The slugified string.
- **Example**:
  ```go
  slug := Slug("Hello, World!")
  fmt.Println(slug) // Output: hello-world
  ```

#### `TruncateText(text string, limit int) string`
- **Purpose**: Truncates a string to a specified length and appends "..." if truncated.
- **Parameters**:
    - `text`: The string to truncate.
    - `limit`: The maximum length of the string.
- **Returns**: The truncated string.
- **Example**:
  ```go
  truncated := TruncateText("This is a long text", 10)
  fmt.Println(truncated) // Output: This is a...
  ```

---

### 4. **Data Conversion Utilities**

#### `ConvertBytes(bytes uint64) string`
- **Purpose**: Converts a byte size into a human-readable string (e.g., "1.23 MiB").
- **Parameters**:
    - `bytes`: The byte size to convert.
- **Returns**: A formatted string with the appropriate unit.
- **Example**:
  ```go
  result := ConvertBytes(1024 * 1024)
  fmt.Println(result) // Output: 1.00 MiB
  ```

#### `ConvertToBytes(input string) (int64, error)`
- **Purpose**: Converts a string with a size suffix (e.g., "1M", "1Mi") to bytes.
- **Parameters**:
    - `input`: The string to convert.
- **Returns**: The byte size or an error if the input is invalid.
- **Example**:
  ```go
  bytes, err := ConvertToBytes("1Mi")
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println(bytes) // Output: 1048576
  }
  ```

#### `ParseStringRanges(rangeStrings []string) ([]int, error)`
- **Purpose**: Parses a list of range strings (e.g., `["1-3", "5"]`) into a slice of integers.
- **Parameters**:
    - `rangeStrings`: The list of range strings.
- **Returns**: A slice of integers or an error if parsing fails.
- **Example**:
  ```go
  result, err := ParseStringRanges([]string{"1-3", "5"})
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println(result) // Output: [1 2 3 5]
  }
  ```

#### `ParseStringRange(rs string) ([]int, error)`
- **Purpose**: Parses a single range string (e.g., `"1-3"`) into a slice of integers.
- **Parameters**:
    - `rs`: The range string.
- **Returns**: A slice of integers or an error if parsing fails.
- **Example**:
  ```go
  result, err := ParseStringRange("1-3")
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println(result) // Output: [1 2 3]
  }
  ```

---

### 5. **Time and Duration Utilities**

#### `ParseDuration(durationStr string) (time.Duration, error)`
- **Purpose**: Parses a duration string (e.g., `"1h30m"`) into a `time.Duration`.
- **Parameters**:
    - `durationStr`: The duration string.
- **Returns**: A `time.Duration` or an error if parsing fails.
- **Example**:
  ```go
  duration, err := ParseDuration("1h30m")
  if err != nil {
      fmt.Println("Error:", err)
  } else {
      fmt.Println(duration) // Output: 1h30m0s
  }
  ```

#### `FormatDuration(d time.Duration, decimalCount int) string`
- **Purpose**: Formats a duration into a human-readable string (e.g., `"1.5s"`).
- **Parameters**:
    - `d`: The duration to format.
    - `decimalCount`: The number of decimal places to include.
- **Returns**: A formatted string.
- **Example**:
  ```go
  result := FormatDuration(90*time.Second, 1)
  fmt.Println(result) // Output: 1.5m
  ```
### 6. **Network Utilities**

### IP and CIDR Validation Utilities

### 1. **IsCIDR(cidr string) bool**
- **Purpose**: Checks if the input string is a valid CIDR (Classless Inter-Domain Routing) notation.
- **Parameters**:
  - `cidr`: The string to validate as a CIDR.
- **Returns**: `true` if the input is a valid CIDR, otherwise `false`.
- **Example**:
  ```go
  isValid := IsIPAddress("192.168.1.0/24")
  fmt.Println(isValid) // Output: true
  ```

---

### 2. **IsIPOrCIDR(input string) (isIP bool, isCIDR bool)**
- **Purpose**: Determines whether the input string is a valid IP address or a valid CIDR notation.
- **Parameters**:
  - `input`: The string to check.
- **Returns**:
  - `isIP`: `true` if the input is a valid IP address.
  - `isCIDR`: `true` if the input is a valid CIDR notation.
- **Example**:
  ```go
  isIP, isCIDR := IsIPOrCIDR("192.168.1.1")
  fmt.Println(isIP, isCIDR) // Output: true, false

  isIP, isCIDR = IsIPOrCIDR("192.168.1.0/24")
  fmt.Println(isIP, isCIDR) // Output: false, true

  isIP, isCIDR = IsIPOrCIDR("invalid")
  fmt.Println(isIP, isCIDR) // Output: false, false
  ```

---

### 3. **IsIPAddress(ip string) bool**
- **Purpose**: Checks if the input string is a valid IP address (IPv4 or IPv6).
- **Parameters**:
  - `ip`: The string to validate as an IP address.
- **Returns**: `true` if the input is a valid IP address, otherwise `false`.
- **Example**:
  ```go
  isValid := IsIPAddress("192.168.1.1")
  fmt.Println(isValid) // Output: true

  isValid = IsIPAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
  fmt.Println(isValid) // Output: true

  isValid = IsIPAddress("invalid")
  fmt.Println(isValid) // Output: false
  ```

---

## Usage Examples

### Example 1: Validate a CIDR
```go
cidr := "192.168.1.0/24"
if IsCIDR(cidr) {
    fmt.Println("Valid CIDR")
} else {
    fmt.Println("Invalid CIDR")
}
```

### Example 2: Check if Input is an IP or CIDR
```go
input := "192.168.1.1"
isIP, isCIDR := IsIPOrCIDR(input)
if isIP {
    fmt.Println("Input is a valid IP address")
} else if isCIDR {
    fmt.Println("Input is a valid CIDR")
} else {
    fmt.Println("Input is neither an IP nor a CIDR")
}
```

### Example 3: Validate an IP Address
```go
ip := "192.168.1.1"
if IsIPAddress(ip) {
    fmt.Println("Valid IP address")
} else {
    fmt.Println("Invalid IP address")
}
```

---