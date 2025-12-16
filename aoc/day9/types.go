package day9

import "fmt"

// Point represents a red tile position in 2D space.
type Point struct {
	X, Y int
}

// String implements fmt.Stringer for debugging.
func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

// RectangleArea calculates the area of a rectangle formed by this point
// and another point as opposite corners.
// The rectangle is inclusive of both corners, so we add 1 to each dimension.
func (p Point) RectangleArea(other Point) int {
	width := abs(p.X-other.X) + 1
	height := abs(p.Y-other.Y) + 1
	return width * height
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
