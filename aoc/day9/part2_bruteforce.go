package day9

import "fmt"

// Part2BruteForce solves Day 9 Part 2 using brute force approach (5+ minutes).
// Kept for educational comparison with the optimized coordinate compression version.
//
// Algorithm:
// 1. Parse red tiles in order (they form a closed polygon when connected)
// 2. Build set of green tiles:
//    - Tiles on edges between consecutive red tiles
//    - Tiles inside the polygon
// 3. For each pair of red tiles as corners, check if rectangle only contains red/green
// 4. Return maximum valid area
func Part2BruteForce(inputPath string) (int, error) {
	redTiles, err := PointsFromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	if len(redTiles) < 2 {
		return 0, fmt.Errorf("need at least 2 red tiles to form a rectangle")
	}

	// Build red tile set once
	redSet := make(map[Point]bool)
	for _, r := range redTiles {
		redSet[r] = true
	}

	// Build ONLY edge tiles (fast - just perimeter, not interior)
	// Can't pre-compute interior - coordinate space is 97K × 97K = 9.4B points!
	fmt.Println("Building edge tile set...")
	edgeTiles := make(map[Point]bool)
	n := len(redTiles)
	for i := 0; i < n; i++ {
		from := redTiles[i]
		to := redTiles[(i+1)%n]
		addLineTiles(edgeTiles, from, to)
	}
	fmt.Printf("Found %d edge tiles\n", len(edgeTiles))

	// Cache for polygon checks - will be shared across all rectangles
	polygonCache := make(map[Point]bool)

	maxArea := 0
	totalPairs := len(redTiles) * (len(redTiles) - 1) / 2
	checked := 0
	validRects := 0

	fmt.Printf("Checking %d rectangle pairs...\n", totalPairs)

	// Check all pairs of red tiles as opposite corners
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			// Lazy validation: only check points in THIS rectangle
			// Uses cached polygon checks to avoid redundant computation
			if isRectangleValid(redTiles[i], redTiles[j], redSet, edgeTiles, redTiles, polygonCache) {
				validRects++
				area := redTiles[i].RectangleArea(redTiles[j])
				if area > maxArea {
					maxArea = area
					fmt.Printf("\n  New max: %d (from %v to %v)\n", maxArea, redTiles[i], redTiles[j])
				}
			}
			checked++

			// Progress output every 1000 pairs
			if checked%1000 == 0 {
				percent := float64(checked) * 100 / float64(totalPairs)
				fmt.Printf("\r  Progress: %d/%d (%.1f%%) - Valid: %d - Max: %d - Cache: %d",
					checked, totalPairs, percent, validRects, maxArea, len(polygonCache))
			}
		}
	}
	fmt.Printf("\r  Progress: %d/%d (100.0%%) - Valid: %d - Max: %d - Cache: %d\n",
		totalPairs, totalPairs, validRects, maxArea, len(polygonCache))
	fmt.Printf("Final answer: %d\n", maxArea)

	return maxArea, nil
}
