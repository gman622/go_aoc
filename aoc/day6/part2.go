package day6

import "fmt"

// Part2 solves Day 6 Part 2 (right-to-left column reading)
func Part2(inputPath string) (int, error) {
	lines, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	return SolveWorksheet(lines, RightToLeft)
}
