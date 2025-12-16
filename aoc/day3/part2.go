package day3

import "fmt"

// Part2 solves Day 3 Part 2: find the maximum 12-digit joltage from each battery bank
// and return the total output joltage
func Part2(inputPath string) (int, error) {
	banks, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	totalJoltage := 0
	for _, bank := range banks {
		maxJoltage := findMaxJoltage12(bank)
		totalJoltage += maxJoltage
	}

	return totalJoltage, nil
}

// findMaxJoltage12 finds the maximum 12-digit joltage from a battery bank
// by selecting exactly 12 batteries (maintaining their order)
func findMaxJoltage12(bank string) int {
	if len(bank) < 12 {
		return 0
	}

	const numDigits = 12
	selected := make([]int, 0, numDigits)

	// Greedily select the largest digit at each position
	startIdx := 0
	for len(selected) < numDigits {
		remaining := numDigits - len(selected)
		// We can search up to len(bank) - remaining (need to leave enough positions)
		searchEnd := len(bank) - remaining + 1

		// Find the maximum digit in the valid range
		maxDigit := -1
		maxIdx := -1
		for i := startIdx; i < searchEnd; i++ {
			digit := int(bank[i] - '0')
			if digit > maxDigit {
				maxDigit = digit
				maxIdx = i
			}
		}

		selected = append(selected, maxDigit)
		startIdx = maxIdx + 1
	}

	// Convert selected digits to a number
	result := 0
	for _, digit := range selected {
		result = result*10 + digit
	}

	return result
}
