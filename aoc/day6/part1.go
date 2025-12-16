package day6

import "fmt"

// Part1 solves Day 6 Part 1 (left-to-right field reading)
func Part1(inputPath string) (int, error) {
	lines, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	return SolveWorksheet(lines, LeftToRight)
}
