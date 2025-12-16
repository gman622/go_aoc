package day9

import "fmt"

// Part1 solves Day 9 Part 1 - finds the largest rectangle with red tiles at opposite corners.
//
// Algorithm:
// 1. Parse all red tile positions
// 2. Check all pairs of red tiles as potential opposite corners
// 3. Calculate the area of each rectangle
// 4. Return the maximum area found
//
// Time complexity: O(n²) where n is the number of red tiles
func Part1(inputPath string) (int, error) {
	points, err := PointsFromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	if len(points) < 2 {
		return 0, fmt.Errorf("need at least 2 red tiles to form a rectangle")
	}

	maxArea := 0

	// Check all pairs of red tiles
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			area := points[i].RectangleArea(points[j])
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea, nil
}
