package day9

import (
	"strings"
	"testing"
)

func TestPart1Example(t *testing.T) {
	input := `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

	points, err := ParsePoints(strings.Split(input, "\n"))
	if err != nil {
		t.Fatalf("Failed to parse points: %v", err)
	}

	// Test individual rectangles from the problem
	tests := []struct {
		p1, p2 Point
		want   int
	}{
		{Point{2, 5}, Point{9, 7}, 24},   // width=8, height=3
		{Point{7, 1}, Point{11, 7}, 35},  // width=5, height=7
		{Point{7, 3}, Point{2, 3}, 6},    // width=6, height=1
		{Point{2, 5}, Point{11, 1}, 50},  // width=10, height=5 (largest)
	}

	for _, tt := range tests {
		got := tt.p1.RectangleArea(tt.p2)
		if got != tt.want {
			t.Errorf("RectangleArea(%v, %v) = %d, want %d", tt.p1, tt.p2, got, tt.want)
		}
	}

	// Test that we find the maximum area
	maxArea := 0
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			area := points[i].RectangleArea(points[j])
			if area > maxArea {
				maxArea = area
			}
		}
	}

	want := 50
	if maxArea != want {
		t.Errorf("Maximum area = %d, want %d", maxArea, want)
	}
}

func TestRectangleAreaCalculation(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 Point
		want   int
	}{
		{"same point", Point{5, 5}, Point{5, 5}, 1},          // 1×1 = 1
		{"horizontal line", Point{2, 5}, Point{7, 5}, 6},     // 6×1 = 6
		{"vertical line", Point{5, 2}, Point{5, 7}, 6},       // 1×6 = 6
		{"square 2x2", Point{0, 0}, Point{1, 1}, 4},          // 2×2 = 4
		{"rectangle 3x5", Point{0, 0}, Point{2, 4}, 15},      // 3×5 = 15
		{"negative coords", Point{-5, -5}, Point{5, 5}, 121}, // 11×11 = 121
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p1.RectangleArea(tt.p2)
			if got != tt.want {
				t.Errorf("RectangleArea(%v, %v) = %d, want %d", tt.p1, tt.p2, got, tt.want)
			}
		})
	}
}
