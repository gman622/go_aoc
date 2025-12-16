package day12

import "fmt"

// Part1 solves Day 12 Part 1 - counts how many regions can fit all their required presents.
func Part1(inputPath string) (int, error) {
	shapes, regions, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	count := 0
	for _, region := range regions {
		solver := NewSolver(shapes, region)
		if solver.CanFit() {
			count++
		}
	}

	return count, nil
}
