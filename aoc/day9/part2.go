package day9

import (
	"fmt"
	"sort"
)

// Part2 solves Day 9 Part 2 using coordinate compression (sub-second runtime).
//
// Key insight: The 496 red tiles occupy a 97K × 97K coordinate space (9.4B points).
// We can't pre-compute all interior points in that space, but we can compress
// coordinates down to a 496×496 grid and pre-fill it completely!
//
// Algorithm:
// 1. Extract unique X and Y coordinates from red tiles
// 2. Build compressed coordinate mappings (sparse → dense)
// 3. Pre-fill the entire compressed polygon interior (~246K cells)
// 4. Check all pairs of red tiles using compressed coordinates
//
// This transforms an infeasible O(97K²) space into a trivial O(496²) space.
func Part2(inputPath string) (int, error) {
	redTiles, err := PointsFromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	if len(redTiles) < 2 {
		return 0, fmt.Errorf("need at least 2 red tiles to form a rectangle")
	}

	// Step 1: Extract unique coordinates
	xSet := make(map[int]bool)
	ySet := make(map[int]bool)
	for _, p := range redTiles {
		xSet[p.X] = true
		ySet[p.Y] = true
	}

	// Convert to sorted slices
	xCoords := make([]int, 0, len(xSet))
	for x := range xSet {
		xCoords = append(xCoords, x)
	}
	sort.Ints(xCoords)

	yCoords := make([]int, 0, len(ySet))
	for y := range ySet {
		yCoords = append(yCoords, y)
	}
	sort.Ints(yCoords)

	// Step 2: Build coordinate mappings (sparse → dense)
	xToCompressed := make(map[int]int)
	for i, x := range xCoords {
		xToCompressed[x] = i
	}

	yToCompressed := make(map[int]int)
	for i, y := range yCoords {
		yToCompressed[y] = i
	}

	// Compress red tiles
	compressedRed := make([]Point, len(redTiles))
	for i, p := range redTiles {
		compressedRed[i] = Point{
			X: xToCompressed[p.X],
			Y: yToCompressed[p.Y],
		}
	}

	// Build red tile set in compressed space
	redSet := make(map[Point]bool)
	for _, p := range compressedRed {
		redSet[p] = true
	}

	// Step 3: Pre-fill the compressed polygon interior
	fmt.Printf("Compressed space: %d × %d = %d cells\n", len(xCoords), len(yCoords), len(xCoords)*len(yCoords))

	// Build edge tiles in compressed space
	edgeTiles := make(map[Point]bool)
	n := len(compressedRed)
	for i := 0; i < n; i++ {
		from := compressedRed[i]
		to := compressedRed[(i+1)%n]
		addLineTiles(edgeTiles, from, to)
	}

	// Pre-fill interior using point-in-polygon test on compressed space
	interior := make(map[Point]bool)
	for x := 0; x < len(xCoords); x++ {
		for y := 0; y < len(yCoords); y++ {
			p := Point{x, y}
			if redSet[p] || edgeTiles[p] {
				interior[p] = true
			} else if isInsidePolygon(p, compressedRed) {
				interior[p] = true
			}
		}
	}
	fmt.Printf("Pre-filled %d interior cells\n", len(interior))

	// Step 4: Check all pairs of red tiles
	maxArea := 0
	totalPairs := len(redTiles) * (len(redTiles) - 1) / 2
	checked := 0
	validRects := 0

	fmt.Printf("Checking %d rectangle pairs...\n", totalPairs)

	for i := 0; i < len(compressedRed); i++ {
		for j := i + 1; j < len(compressedRed); j++ {
			p1, p2 := compressedRed[i], compressedRed[j]

			// Check if rectangle is valid in compressed space
			if isRectangleValidCompressed(p1, p2, interior) {
				validRects++
				// Calculate area in ORIGINAL coordinate space
				origP1 := redTiles[i]
				origP2 := redTiles[j]
				area := origP1.RectangleArea(origP2)

				if area > maxArea {
					maxArea = area
					fmt.Printf("New max: %d (from %v to %v)\n", maxArea, origP1, origP2)
				}
			}
			checked++

			if checked%1000 == 0 {
				percent := float64(checked) * 100 / float64(totalPairs)
				fmt.Printf("\r  Progress: %d/%d (%.1f%%) - Valid: %d - Max: %d",
					checked, totalPairs, percent, validRects, maxArea)
			}
		}
	}
	fmt.Printf("\r  Progress: %d/%d (100.0%%) - Valid: %d - Max: %d\n",
		totalPairs, totalPairs, validRects, maxArea)

	return maxArea, nil
}

// isRectangleValidCompressed checks if all points in a compressed rectangle are interior points.
// This is trivial since we pre-computed all interior points!
func isRectangleValidCompressed(p1, p2 Point, interior map[Point]bool) bool {
	minX := min(p1.X, p2.X)
	maxX := max(p1.X, p2.X)
	minY := min(p1.Y, p2.Y)
	maxY := max(p1.Y, p2.Y)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if !interior[Point{x, y}] {
				return false
			}
		}
	}
	return true
}
