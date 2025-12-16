package day2

import "fmt"

// Part2 solves Day 2 Part 2: sum all invalid product IDs with relaxed rules.
//
// Problem: Find IDs that are patterns repeated at least twice (e.g., 111, 123123, 55)
//
// Same algorithmic approach as Part1, but with different validation rules.
// See Part1 comments for detailed analysis of why iteration is optimal here.
//
// Strategy Pattern Benefit: We only changed the validator, not the iteration logic.
// Part1 and Part2 share the same structure but different behavior - classic OOP pattern.
//
// Performance: Similar to Part1 (~56ms)
// - AtLeastTwiceValidator is slightly more complex (tries multiple pattern lengths)
// - Still O(log n) per ID check (where n is the ID value)
// - Early exits keep average case fast
func Part2(inputPath string) (int, error) {
	parser, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	ranges, err := parser.ParseAll()
	if err != nil {
		return 0, fmt.Errorf("parsing ranges: %w", err)
	}

	// Different validator, same iteration pattern
	// This demonstrates the power of the Strategy pattern
	validator := AtLeastTwiceValidator{}
	sum := 0

	for _, r := range ranges {
		for id := r.Start; id <= r.End; id++ {
			if validator.IsInvalid(id) {
				sum += id
			}
		}
	}

	return sum, nil
}
