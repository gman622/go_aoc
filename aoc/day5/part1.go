package day5

import "fmt"

// Part1 solves Day 5 Part 1: Count how many available ingredient IDs are fresh.
// An ingredient ID is fresh if it falls within any of the fresh ranges (inclusive).
func Part1(inputPath string) (int, error) {
	db, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	freshCount := 0
	for _, id := range db.AvailableIDs {
		if isFresh(id, db.FreshRanges) {
			freshCount++
		}
	}

	return freshCount, nil
}

// isFresh checks if an ingredient ID is fresh (falls within any range).
func isFresh(id int, ranges []Range) bool {
	for _, r := range ranges {
		if r.Contains(id) {
			return true
		}
	}
	return false
}
