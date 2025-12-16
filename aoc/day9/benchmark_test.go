package day9

import (
	"testing"
	"time"
)

func TestPerformanceEstimate(t *testing.T) {
	// Load actual input
	redTiles, err := PointsFromFile("../../inputs/day9_input.txt")
	if err != nil {
		t.Fatalf("Failed to load input: %v", err)
	}

	totalTiles := len(redTiles)
	t.Logf("Total red tiles: %d", totalTiles)
	t.Logf("Total pairs to check: %d", totalTiles*(totalTiles-1)/2)

	// Build shared data structures
	redSet := make(map[Point]bool)
	for _, r := range redTiles {
		redSet[r] = true
	}

	edgeTiles := make(map[Point]bool)
	n := len(redTiles)
	for i := 0; i < n; i++ {
		from := redTiles[i]
		to := redTiles[(i+1)%n]
		addLineTiles(edgeTiles, from, to)
	}
	t.Logf("Built edge tiles set: %d tiles", len(edgeTiles))

	cache := make(map[Point]bool)

	// Test with increasing subset sizes
	subsetSizes := []int{10, 20, 50, 100}

	for _, size := range subsetSizes {
		if size > totalTiles {
			break
		}

		subset := redTiles[:size]
		pairCount := size * (size - 1) / 2

		start := time.Now()
		maxArea := 0
		checked := 0

		for i := 0; i < len(subset); i++ {
			for j := i + 1; j < len(subset); j++ {
				if isRectangleValid(subset[i], subset[j], redSet, edgeTiles, redTiles, cache) {
					area := subset[i].RectangleArea(subset[j])
					if area > maxArea {
						maxArea = area
					}
				}
				checked++
			}
		}

		elapsed := time.Since(start)
		perPair := elapsed / time.Duration(pairCount)

		// Extrapolate to full dataset
		totalPairs := totalTiles * (totalTiles - 1) / 2
		estimatedTotal := time.Duration(totalPairs) * perPair

		t.Logf("\nSubset size: %d tiles (%d pairs)", size, pairCount)
		t.Logf("  Time: %v", elapsed)
		t.Logf("  Per pair: %v", perPair)
		t.Logf("  Max area found: %d", maxArea)
		t.Logf("  Extrapolated total time: %v (%.1f minutes)", estimatedTotal, estimatedTotal.Minutes())
	}
}

func BenchmarkRectangleValidation(b *testing.B) {
	// Load actual input
	redTiles, err := PointsFromFile("../../inputs/day9_input.txt")
	if err != nil {
		b.Fatalf("Failed to load input: %v", err)
	}

	// Build shared data structures
	redSet := make(map[Point]bool)
	for _, r := range redTiles {
		redSet[r] = true
	}

	edgeTiles := make(map[Point]bool)
	n := len(redTiles)
	for i := 0; i < n; i++ {
		from := redTiles[i]
		to := redTiles[(i+1)%n]
		addLineTiles(edgeTiles, from, to)
	}

	cache := make(map[Point]bool)

	// Benchmark first 100 pairs
	subset := redTiles[:20]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(subset); j++ {
			for k := j + 1; k < len(subset); k++ {
				isRectangleValid(subset[j], subset[k], redSet, edgeTiles, redTiles, cache)
			}
		}
	}
}
