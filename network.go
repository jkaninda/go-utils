package goutils

import "net"

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
