package day8

import (
	"fmt"
	"sort"
)

// Part2 solves Day 8 Part 2 - finds the last connection that unites all circuits.
//
// Algorithm:
// 1. Parse 3D coordinates of junction boxes
// 2. Calculate distances between all pairs of boxes
// 3. Sort pairs by distance (ascending)
// 4. Use Union-Find to connect pairs until all boxes are in one circuit
// 5. Return the product of X coordinates of the last two boxes connected
func Part2(inputPath string) (int, error) {
	points, err := PointsFromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	n := len(points)
	if n == 0 {
		return 0, fmt.Errorf("no junction boxes found")
	}

	// Create Union-Find structure
	uf := NewUnionFind(n)

	// Calculate all pairwise distances
	type pair struct {
		i, j     int
		distance float64
	}

	pairs := make([]pair, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			dist := points[i].DistanceTo(points[j])
			pairs = append(pairs, pair{i, j, dist})
		}
	}

	// Sort pairs by distance
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].distance < pairs[j].distance
	})

	// Connect pairs until all boxes are in one circuit
	componentsCount := n
	for _, p := range pairs {
		// Try to union - if they were in different sets, we merged them
		if uf.Union(p.i, p.j) {
			componentsCount--
			// Check if all boxes are now in one circuit
			if componentsCount == 1 {
				// This was the last connection!
				result := points[p.i].X * points[p.j].X
				return result, nil
			}
		}
	}

	return 0, fmt.Errorf("could not connect all junction boxes into one circuit")
}
