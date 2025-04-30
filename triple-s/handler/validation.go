package handler

import (
	"regexp"
	"strings"
)

// isValidBucketName validates if a bucket name meets the specified requirements
func isValidBucketName(bucketName string) bool {
	// Check length constraints
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return false
	}

	// Check if the bucket name matches an IP address format (rejected if it's a valid IP)
	ipAddressPattern := `^(?:\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(ipAddressPattern, bucketName)
	if match {
		return false // Reject if the name is an IP address
	}

	// Regex pattern for valid bucket names (lowercase letters, digits, dots, and dashes)
	validNamePattern := `^[a-z0-9.-]+$`
	match, _ = regexp.MatchString(validNamePattern, bucketName)
	if !match {
		return false // Reject if the name contains invalid characters
	}

	// Ensure the name doesn't start or end with a dot or dash
	if bucketName[0] == '.' || bucketName[0] == '-' || bucketName[len(bucketName)-1] == '.' || bucketName[len(bucketName)-1] == '-' {
		return false
	}

	// Ensure there are no consecutive dots or dashes
	if strings.Contains(bucketName, "..") || strings.Contains(bucketName, "-.") || strings.Contains(bucketName, ".-") || strings.Contains(bucketName, "--") {
		return false
	}

	// If all checks pass, the bucket name is valid
	return true
}
