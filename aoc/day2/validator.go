package day2

import "strconv"

// Validator defines the interface for product ID validation strategies
type Validator interface {
	IsInvalid(id int) bool
}

// ExactlyTwiceValidator checks if an ID is made of a pattern repeated exactly twice
// Examples: 55 (5 twice), 6464 (64 twice), 123123 (123 twice)
type ExactlyTwiceValidator struct{}

// IsInvalid returns true if the ID is made of a sequence repeated exactly twice
func (v ExactlyTwiceValidator) IsInvalid(id int) bool {
	s := strconv.Itoa(id)

	// Must have even length to be repeated exactly twice
	if len(s)%2 != 0 {
		return false
	}

	// Check for leading zeros (numbers like 0101 are not valid IDs)
	if s[0] == '0' {
		return false
	}

	// Split in half and check if both halves are equal
	mid := len(s) / 2
	firstHalf := s[:mid]
	secondHalf := s[mid:]

	return firstHalf == secondHalf
}

// AtLeastTwiceValidator checks if an ID is made of a pattern repeated at least twice
// Examples: 111 (1 three times), 123123123 (123 three times), 1212 (12 twice)
type AtLeastTwiceValidator struct{}

// IsInvalid returns true if the ID is made of a pattern repeated at least twice
func (v AtLeastTwiceValidator) IsInvalid(id int) bool {
	s := strconv.Itoa(id)

	// Check for leading zeros (numbers like 0101 are not valid IDs)
	if s[0] == '0' {
		return false
	}

	n := len(s)

	// Try all possible pattern lengths from 1 to n/2
	// The pattern must repeat at least twice
	for patternLen := 1; patternLen <= n/2; patternLen++ {
		// The string length must be divisible by pattern length
		if n%patternLen != 0 {
			continue
		}

		// Check if the string is made by repeating the first patternLen characters
		pattern := s[:patternLen]
		isRepeating := true

		for i := patternLen; i < n; i += patternLen {
			if s[i:i+patternLen] != pattern {
				isRepeating = false
				break
			}
		}

		if isRepeating {
			return true
		}
	}

	return false
}
