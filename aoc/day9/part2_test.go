package day9

import (
	"strings"
	"testing"
)

func TestPart2Example(t *testing.T) {
	input := `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

	redTiles, err := ParsePoints(strings.Split(input, "\n"))
	if err != nil {
		t.Fatalf("Failed to parse points: %v", err)
	}

	// Build red set and edge tiles
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

	// Test individual rectangles from the problem
	tests := []struct {
		name   string
		p1, p2 Point
		valid  bool
		area   int
	}{
		{"7,3 to 11,1", Point{7, 3}, Point{11, 1}, true, 15},
		{"9,7 to 9,5", Point{9, 7}, Point{9, 5}, true, 3},
		{"9,5 to 2,3", Point{9, 5}, Point{2, 3}, true, 24}, // Maximum
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isRectangleValid(tt.p1, tt.p2, redSet, edgeTiles, redTiles, cache)
			if valid != tt.valid {
				t.Errorf("Rectangle %v to %v validity = %v, want %v", tt.p1, tt.p2, valid, tt.valid)
			}
			if valid {
				area := tt.p1.RectangleArea(tt.p2)
				if area != tt.area {
					t.Errorf("Rectangle %v to %v area = %d, want %d", tt.p1, tt.p2, area, tt.area)
				}
			}
		})
	}

	// Test that we find the maximum area
	maxArea := 0
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			if isRectangleValid(redTiles[i], redTiles[j], redSet, edgeTiles, redTiles, cache) {
				area := redTiles[i].RectangleArea(redTiles[j])
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	want := 24
	if maxArea != want {
		t.Errorf("Maximum valid area = %d, want %d", maxArea, want)
	}
}

func TestBuildGreenTiles(t *testing.T) {
	// Simple square
	redTiles := []Point{
		{0, 0},
		{2, 0},
		{2, 2},
		{0, 2},
	}

	greenTiles := BuildGreenTiles(redTiles)

	// Check edge tiles
	expectedGreen := []Point{
		{1, 0}, // top edge
		{2, 1}, // right edge
		{1, 2}, // bottom edge
		{0, 1}, // left edge
		{1, 1}, // interior
	}

	for _, p := range expectedGreen {
		if !greenTiles[p] {
			t.Errorf("Expected %v to be green", p)
		}
	}
}
