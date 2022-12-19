package utils

import (
	"fmt"
	"strings"
)

// ParseNumber parses a hex or decimal number
func ParseNumber(number string) (int, error) {
	number = strings.Trim(number, " \t")
	number = strings.ToLower(number)
	if len(number) == 0 {
		return 0, fmt.Errorf("Empty string is not a number")
	}

	var value int
	if strings.HasPrefix(number, "\\x") || strings.HasPrefix(number, "0x") {
		for i := 2; i < len(number); i++ {
			if (!(number[i] >= 48 && number[i] <= 57)) && !(number[i] >= 97 && number[i] <= 102) {
				return 0, fmt.Errorf("%s is not valid hexadecimal", number)
			}
		}
		_, err := fmt.Sscanf(number[2:], "%x", &value)
		if err != nil {
			return 0, err
		}
	} else {
		for i := 0; i < len(number); i++ {
			if !(number[i] >= 48 && number[i] <= 57) {
				return 0, fmt.Errorf("%s is not valid decimal", number)
			}
		}
		_, err := fmt.Sscanf(number, "%d", &value)
		if err != nil {
			return 0, err
		}
	}

	return value, nil
}

// IsHexString returns true if the input is a hex string
func IsHexString(number string) bool {
	lower := strings.ToLower(number)
	if strings.HasPrefix(lower, "0x") || strings.HasPrefix(lower, "\\x") && len(lower) >= 3 {
		for i := 2; i < len(lower); i++ {
			if (!(lower[i] >= 48 && lower[i] <= 57)) && !(lower[i] >= 97 && lower[i] <= 102) {
				return false
			}
		}
		return true
	}
	return false
}
