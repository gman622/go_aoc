package day12

import (
	"sort"
	"strings"
	"time"
)

// Solver attempts to fit presents into a region using backtracking with optimizations.
type Solver struct {
	shapes         []Shape
	region         Region
	grid           *Grid
	transformCache map[int][]Shape // Cache of all transformations per shape ID
	memo           map[string]bool  // Memoization cache
	timeout        time.Time
	timedOut       bool
}

// NewSolver creates a solver for a given region and available shapes.
func NewSolver(shapes []Shape, region Region) *Solver {
	grid := NewGrid(region.Width, region.Height)

	// Pre-compute all transformations for each shape
	transformCache := make(map[int][]Shape)
	for _, shape := range shapes {
		transformCache[shape.ID] = shape.AllTransformations()
	}

	return &Solver{
		shapes:         shapes,
		region:         region,
		grid:           grid,
		transformCache: transformCache,
		memo:           make(map[string]bool),
	}
}

// CanFit determines if all required presents can be placed in the region.
func (s *Solver) CanFit() bool {
	// Build list of presents to place
	var presentsToPlace []int
	totalArea := 0
	for shapeID, count := range s.region.Presents {
		for range count {
			presentsToPlace = append(presentsToPlace, shapeID)
			totalArea += len(s.shapes[shapeID].Points)
		}
	}

	// Quick check: if total area needed exceeds region area, impossible
	regionArea := s.region.Width * s.region.Height
	if totalArea > regionArea {
		return false
	}

	// Optimization: Sort by size (largest first) to prune search tree earlier
	sort.Slice(presentsToPlace, func(i, j int) bool {
		return len(s.shapes[presentsToPlace[i]].Points) > len(s.shapes[presentsToPlace[j]].Points)
	})

	// For small problems, use exact backtracking
	if len(presentsToPlace) <= 15 {
		s.timeout = time.Now().Add(50 * time.Millisecond)
		s.timedOut = false
		result := s.backtrack(presentsToPlace, 0)
		if s.timedOut {
			return false
		}
		return result
	}

	// For large problems, use greedy placement
	return s.greedyFit(presentsToPlace)
}

// backtrack attempts to place presents recursively with pruning and memoization.
func (s *Solver) backtrack(presents []int, index int) bool {
	// Check timeout periodically
	if index%10 == 0 && time.Now().After(s.timeout) {
		s.timedOut = true
		return false
	}

	// Base case: all presents placed successfully
	if index >= len(presents) {
		return true
	}

	// Memoization: Check if we've seen this state before
	// State = grid hash + remaining presents
	stateKey := s.computeStateKey(presents[index:])
	if result, exists := s.memo[stateKey]; exists {
		return result
	}

	// Pruning: Check if remaining space can fit remaining presents
	remainingArea := 0
	for i := index; i < len(presents); i++ {
		remainingArea += len(s.shapes[presents[i]].Points)
	}
	emptyCells := s.countEmptyCells()
	if emptyCells < remainingArea {
		s.memo[stateKey] = false
		return false // Not enough space left
	}

	shapeID := presents[index]

	// Get all transformations of this shape
	transformations := s.transformCache[shapeID]

	// Optimization: Find first empty cell and only try placements there
	// This avoids trying equivalent positions
	firstEmpty := s.findFirstEmpty()
	if firstEmpty.X == -1 {
		return false // No empty space but still have presents to place
	}

	// Try placing each transformation at positions near the first empty cell
	for _, transform := range transformations {
		// Try positions that could cover the first empty cell
		for y := max(0, firstEmpty.Y-2); y <= min(s.region.Height-1, firstEmpty.Y+2); y++ {
			for x := max(0, firstEmpty.X-2); x <= min(s.region.Width-1, firstEmpty.X+2); x++ {
				if s.grid.CanPlace(transform, x, y) && s.coversPoint(transform, x, y, firstEmpty) {
					// Place the shape
					label := rune('A' + index)
					s.grid.Place(transform, x, y, label)

					// Recursively try to place remaining presents
					if s.backtrack(presents, index+1) {
						return true
					}

					// Backtrack: remove the shape
					s.grid.Remove(transform, x, y)
				}
			}
		}
	}

	// Could not place this present
	s.memo[stateKey] = false
	return false
}

// computeStateKey creates a hash of the current grid state + remaining presents.
func (s *Solver) computeStateKey(remainingPresents []int) string {
	// Simple hash: grid cells + remaining present IDs
	var key strings.Builder
	for _, row := range s.grid.Cells {
		key.WriteString(string(row))
	}
	key.WriteByte('|')
	for _, id := range remainingPresents {
		key.WriteByte(byte('0' + id))
	}
	return key.String()
}

// greedyFit attempts to place all presents using a greedy first-fit strategy.
func (s *Solver) greedyFit(presents []int) bool {
	for _, shapeID := range presents {
		placed := false
		transformations := s.transformCache[shapeID]

		// Try to place this shape somewhere
		for _, transform := range transformations {
			if placed {
				break
			}
			for y := 0; y < s.region.Height && !placed; y++ {
				for x := 0; x < s.region.Width && !placed; x++ {
					if s.grid.CanPlace(transform, x, y) {
						s.grid.Place(transform, x, y, '#')
						placed = true
					}
				}
			}
		}

		if !placed {
			return false // Couldn't place this present
		}
	}
	return true // All presents placed
}

// countEmptyCells counts the number of empty cells in the grid.
func (s *Solver) countEmptyCells() int {
	count := 0
	for y := 0; y < s.region.Height; y++ {
		for x := 0; x < s.region.Width; x++ {
			if s.grid.Cells[y][x] == '.' {
				count++
			}
		}
	}
	return count
}

// findFirstEmpty finds the first empty cell (top-left to bottom-right).
func (s *Solver) findFirstEmpty() Point {
	for y := 0; y < s.region.Height; y++ {
		for x := 0; x < s.region.Width; x++ {
			if s.grid.Cells[y][x] == '.' {
				return Point{X: x, Y: y}
			}
		}
	}
	return Point{X: -1, Y: -1}
}

// coversPoint checks if placing a shape at (sx, sy) would cover point p.
func (s *Solver) coversPoint(shape Shape, sx, sy int, p Point) bool {
	for _, pt := range shape.Points {
		if sx+pt.X == p.X && sy+pt.Y == p.Y {
			return true
		}
	}
	return false
}
