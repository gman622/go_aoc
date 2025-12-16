package day5

import (
	"fmt"
	"sort"
)

// Part2 solves Day 5 Part 2: Count total unique ingredient IDs covered by all fresh ranges.
// We need to merge overlapping ranges and sum their sizes.
func Part2(inputPath string) (int, error) {
	db, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	merged := mergeRanges(db.FreshRanges)

	totalCount := 0
	for _, r := range merged {
		totalCount += r.End - r.Start + 1
	}

	return totalCount, nil
}

// mergeRanges merges overlapping or adjacent ranges into non-overlapping ranges.
// Returns a sorted list of merged ranges.
func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return nil
	}

	// Sort ranges by start position
	sorted := make([]Range, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start < sorted[j].Start
	})

	merged := []Range{sorted[0]}

	for i := 1; i < len(sorted); i++ {
		current := sorted[i]
		last := &merged[len(merged)-1]

		// Check if current range overlaps or is adjacent to the last merged range
		if current.Start <= last.End+1 {
			// Merge by extending the end if necessary
			if current.End > last.End {
				last.End = current.End
			}
		} else {
			// No overlap, add as new range
			merged = append(merged, current)
		}
	}

	return merged
}
