package util

import "regexp"

func IsValidHalfWidthASCII(s string, min, max int) bool {
	regex := regexp.MustCompile(`^[!-~]+$`)
	return len(s) >= min && len(s) <= max && regex.MatchString(s)
}

func SpacesOrControlCodes(s string) bool {
	for _, char := range s {
		if char <= 32 || char >= 127 {
			return true
		}
	}
	return false
}
