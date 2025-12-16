package day3

import "fmt"

// Part1 solves Day 3 Part 1: find the maximum joltage from each battery bank
// and return the total output joltage
func Part1(inputPath string) (int, error) {
	banks, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	totalJoltage := 0
	for _, bank := range banks {
		maxJoltage := findMaxJoltage(bank)
		totalJoltage += maxJoltage
	}

	return totalJoltage, nil
}

// findMaxJoltage finds the maximum two-digit joltage from a battery bank
// by selecting any two batteries (maintaining their order).
//
// Optimization: O(n) instead of O(n²) brute force
//
// Key insight: For maximum two-digit number XY:
// - We want the largest possible X (tens digit)
// - Then the largest possible Y (ones digit) that comes AFTER X
//
// Algorithm:
// 1. Find the maximum digit in the string
// 2. Find the maximum digit after the first occurrence of max
// 3. Result = max1*10 + max2
//
// Example: "73924"
// - Brute force: try all pairs (7,3), (7,9), (7,2), (7,4), (3,9), etc.
// - Optimized: max=9 at index 2, then max after index 2 is 4 → 94
//
// Time complexity: O(n) - two passes through string
// Space complexity: O(1) - only track indices and values
func findMaxJoltage(bank string) int {
	if len(bank) < 2 {
		return 0
	}

	// Find the maximum digit and its first occurrence
	maxDigit := -1
	maxIdx := -1
	for i := 0; i < len(bank); i++ {
		digit := int(bank[i] - '0')
		if digit > maxDigit {
			maxDigit = digit
			maxIdx = i
		}
	}

	// Edge case: max digit is at the end, no second digit available
	if maxIdx == len(bank)-1 {
		// Find second largest before maxIdx
		secondMax := -1
		for i := 0; i < maxIdx; i++ {
			digit := int(bank[i] - '0')
			if digit > secondMax {
				secondMax = digit
			}
		}
		// Best is secondMax * 10 + maxDigit
		return secondMax*10 + maxDigit
	}

	// Find the maximum digit after maxIdx
	secondMax := -1
	for i := maxIdx + 1; i < len(bank); i++ {
		digit := int(bank[i] - '0')
		if digit > secondMax {
			secondMax = digit
		}
	}

	return maxDigit*10 + secondMax
}
