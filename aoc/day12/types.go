package day12

import (
	"fmt"
	"strings"
)

// Point represents a coordinate in 2D space.
// Using a custom type prevents mixing with unrelated integers.
type Point struct {
	X, Y int
}

// Shape represents a present shape as a set of relative coordinates.
// The coordinates are normalized so the minimum X and Y are 0.
type Shape struct {
	ID     int
	Points []Point
}

// Region represents a space under a tree where presents need to be placed.
type Region struct {
	Width    int
	Height   int
	Presents []int // presents[i] = count of shape i needed
}

// Grid represents a placement grid for checking if presents fit.
type Grid struct {
	Width   int
	Height  int
	Cells   [][]rune // '.' for empty, 'A'-'Z' for placed presents
	Shapes  []Shape
	Regions []Region
}

// NewGrid creates an empty grid of the given dimensions.
func NewGrid(width, height int) *Grid {
	cells := make([][]rune, height)
	for i := range cells {
		cells[i] = make([]rune, width)
		for j := range cells[i] {
			cells[i][j] = '.'
		}
	}
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// String implements fmt.Stringer for debugging.
func (g *Grid) String() string {
	var sb strings.Builder
	for _, row := range g.Cells {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

// CanPlace checks if a shape can be placed at position (x, y).
func (g *Grid) CanPlace(shape Shape, x, y int) bool {
	for _, p := range shape.Points {
		nx, ny := x+p.X, y+p.Y
		if nx < 0 || nx >= g.Width || ny < 0 || ny >= g.Height {
			return false
		}
		if g.Cells[ny][nx] != '.' {
			return false
		}
	}
	return true
}

// Place puts a shape on the grid at position (x, y) with the given label.
func (g *Grid) Place(shape Shape, x, y int, label rune) {
	for _, p := range shape.Points {
		g.Cells[y+p.Y][x+p.X] = label
	}
}

// Remove removes a shape from the grid at position (x, y).
func (g *Grid) Remove(shape Shape, x, y int) {
	for _, p := range shape.Points {
		g.Cells[y+p.Y][x+p.X] = '.'
	}
}

// Normalize shifts all points so the minimum X and Y are 0.
func (s *Shape) Normalize() Shape {
	if len(s.Points) == 0 {
		return *s
	}

	minX, minY := s.Points[0].X, s.Points[0].Y
	for _, p := range s.Points {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}

	normalized := Shape{ID: s.ID, Points: make([]Point, len(s.Points))}
	for i, p := range s.Points {
		normalized.Points[i] = Point{X: p.X - minX, Y: p.Y - minY}
	}
	return normalized
}

// Rotate90 rotates the shape 90 degrees clockwise.
func (s *Shape) Rotate90() Shape {
	rotated := Shape{ID: s.ID, Points: make([]Point, len(s.Points))}
	for i, p := range s.Points {
		// Rotation: (x, y) -> (y, -x)
		rotated.Points[i] = Point{X: p.Y, Y: -p.X}
	}
	return rotated.Normalize()
}

// FlipHorizontal flips the shape horizontally.
func (s *Shape) FlipHorizontal() Shape {
	flipped := Shape{ID: s.ID, Points: make([]Point, len(s.Points))}
	for i, p := range s.Points {
		flipped.Points[i] = Point{X: -p.X, Y: p.Y}
	}
	return flipped.Normalize()
}

// Hash returns a string representation for deduplication.
func (s *Shape) Hash() string {
	var sb strings.Builder
	for i, p := range s.Points {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("(%d,%d)", p.X, p.Y))
	}
	return sb.String()
}

// AllTransformations returns all unique rotations and flips of the shape.
func (s *Shape) AllTransformations() []Shape {
	seen := make(map[string]Shape)

	// Try all 4 rotations
	current := *s
	for i := 0; i < 4; i++ {
		hash := current.Hash()
		if _, exists := seen[hash]; !exists {
			seen[hash] = current
		}
		current = current.Rotate90()
	}

	// Try all 4 rotations of the flipped shape
	current = s.FlipHorizontal()
	for i := 0; i < 4; i++ {
		hash := current.Hash()
		if _, exists := seen[hash]; !exists {
			seen[hash] = current
		}
		current = current.Rotate90()
	}

	// Convert map to slice
	result := make([]Shape, 0, len(seen))
	for _, shape := range seen {
		result = append(result, shape)
	}
	return result
}
