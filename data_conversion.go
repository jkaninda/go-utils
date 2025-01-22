package goutils

import (
	"fmt"
	"strconv"
	"strings"
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
		return fmt.Sprintf("%.2f EiB", float64(bytes)/float64(EiB))
	case bytes >= PiB:
		return fmt.Sprintf("%.2f PiB", float64(bytes)/float64(PiB))
	case bytes >= TiB:
		return fmt.Sprintf("%.2f TiB", float64(bytes)/float64(TiB))
	case bytes >= GiB:
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GiB))
	case bytes >= MiB:
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MiB))
	case bytes >= KiB:
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KiB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// ConvertToBytes converts a string with a size suffix (e.g., "1M", "1Mi") to bytes.
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
		"K": 1000,
		"M": 1000 * 1000,
		"G": 1000 * 1000 * 1000,
		"T": 1000 * 1000 * 1000 * 1000,
		"P": 1000 * 1000 * 1000 * 1000 * 1000,
		"E": 1000 * 1000 * 1000 * 1000 * 1000 * 1000,
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
