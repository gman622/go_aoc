package day8

import (
	"fmt"
	"sort"
)

// Part1 solves Day 8 Part 1 - connects junction boxes and finds largest circuits.
//
// Algorithm:
// 1. Parse 3D coordinates of junction boxes
// 2. Calculate distances between all pairs of boxes
// 3. Sort pairs by distance (ascending)
// 4. Use Union-Find to connect the 1000 closest pairs
// 5. Find the three largest circuits and multiply their sizes
func Part1(inputPath string) (int, error) {
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
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dist := points[i].DistanceTo(points[j])
			pairs = append(pairs, pair{i, j, dist})
		}
	}

	// Sort pairs by distance
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].distance < pairs[j].distance
	})

	// Connect the 1000 closest pairs
	connectionsNeeded := 1000
	if len(pairs) < connectionsNeeded {
		connectionsNeeded = len(pairs)
	}

	for _, p := range pairs[:connectionsNeeded] {
		uf.Union(p.i, p.j)
	}

	// Get all component sizes
	sizes := uf.ComponentSizes()

	// Sort sizes in descending order to find the three largest
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	// Multiply the three largest circuits
	if len(sizes) < 3 {
		return 0, fmt.Errorf("fewer than 3 circuits found")
	}

	result := sizes[0] * sizes[1] * sizes[2]
	return result, nil
}
