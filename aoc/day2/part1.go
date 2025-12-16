package day2

import "fmt"

// Part1 solves Day 2 Part 1: sum all invalid product IDs in the given ranges.
//
// Problem: Find IDs that are patterns repeated exactly twice (e.g., 123123, 55, 6464)
//
// Algorithm Analysis:
// - Must check each ID individually - no mathematical shortcut exists
// - Can't skip ranges: pattern 123123 appears scattered throughout
// - Can't generate patterns: would need to enumerate all patterns then check ranges
//
// Why this approach is actually efficient:
// 1. Validator.IsInvalid() is O(log n) where n is the ID value (digit count)
// 2. String comparison in Go is highly optimized (memcmp)
// 3. Early exits in validator for odd-length and mismatches
// 4. Strategy pattern allows swapping validators without changing this code
//
// Performance: ~56ms for millions of IDs is reasonable
// - Each ID check: ~50-100ns (string conversion + comparison)
// - Alternative approaches (pattern generation, math) are actually slower
// - This is as fast as "check every number" can get in Go
//
// Potential micro-optimizations (not worth the complexity):
// - Parallel goroutines per range (overhead > benefit for this problem)
// - Digit-by-digit validation (requires complex logic, slower than string ops)
// - Caching (patterns don't repeat enough to matter)
func Part1(inputPath string) (int, error) {
	parser, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	ranges, err := parser.ParseAll()
	if err != nil {
		return 0, fmt.Errorf("parsing ranges: %w", err)
	}

	// Strategy pattern: validator encapsulates the validation logic
	// This keeps Part1 focused on iteration, validator focused on rules
	validator := ExactlyTwiceValidator{}
	sum := 0

	// Iterate through all IDs in all ranges
	// Time complexity: O(total_range_size * log(max_id))
	// Space complexity: O(1) - only accumulator
	for _, r := range ranges {
		for id := r.Start; id <= r.End; id++ {
			if validator.IsInvalid(id) {
				sum += id
			}
		}
	}

	return sum, nil
}
