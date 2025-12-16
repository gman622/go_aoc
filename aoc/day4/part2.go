package day4

import "fmt"

// Part2 solves Day 4 Part 2: iteratively remove accessible rolls.
//
// Problem: Removing accessible rolls can make previously inaccessible rolls
// become accessible. Keep removing until no more can be removed.
//
// This demonstrates several important patterns:
// - Iterative algorithms (repeat until stable state)
// - Mutable data structures (modify in place for efficiency)
// - Simulation problems (Conway's Game of Life, cellular automata)
// - Convergence to fixed point (eventually nothing changes)
//
// Complexity Analysis:
// - Worst case: O(n * iterations) where n is grid size
// - Each iteration removes at least 1 roll (or terminates)
// - Maximum iterations = total number of '@' symbols
// - In practice: converges quickly (logarithmic-like)
func Part2(inputPath string) (int, error) {
	lines, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	// Convert to mutable grid: [][]byte instead of []string
	// Why [][]byte?
	// - Mutable: can modify individual cells (strings are immutable in Go)
	// - Efficient: avoid creating new strings on each modification
	// - Idiomatic: byte slices are standard for text manipulation
	grid := make([][]byte, len(lines))
	for i, line := range lines {
		// []byte(line) converts string to byte slice (makes a copy)
		grid[i] = []byte(line)
	}

	totalRemoved := 0

	// Infinite loop with explicit termination: common pattern for simulations
	// Alternative: while(condition) doesn't exist in Go, use for{} + break
	for {
		// Find all currently accessible rolls
		// Important: find ALL first, then remove ALL
		// If we removed one-by-one, we'd affect the counts mid-iteration
		accessible := findAccessibleRolls(grid)

		// Termination condition: no more accessible rolls (stable state reached)
		if len(accessible) == 0 {
			break
		}

		// Batch removal: remove all accessible rolls simultaneously
		// This simulates "one step" in the iterative process
		for _, pos := range accessible {
			grid[pos.row][pos.col] = '.' // Modify in place
		}

		// Accumulate total across all iterations
		totalRemoved += len(accessible)
	}

	return totalRemoved, nil
}

// position represents a 2D coordinate in the grid.
//
// Struct Pattern: Use structs to group related data
// - More readable than passing two separate ints
// - Type safety: can't accidentally swap row/col
// - Extensible: easy to add fields like "value" or "type" later
//
// Unexported (lowercase): This is an implementation detail
// External packages don't need to know about our position type
type position struct {
	row, col int
}

// findAccessibleRolls returns positions of all accessible rolls in the grid.
//
// This function demonstrates:
// - Separation of concerns: finding vs. removing are separate operations
// - Collecting results in a slice for batch processing
// - Working with mutable [][]byte grids
func findAccessibleRolls(grid [][]byte) []position {
	var accessible []position

	// Same traversal pattern as Part1, but collecting positions instead of counting
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			if grid[row][col] == '@' && isAccessibleMutable(grid, row, col) {
				// Struct literal: position{row, col} creates position with named fields
				accessible = append(accessible, position{row, col})
			}
		}
	}

	return accessible
}

// isAccessibleMutable checks if a roll is accessible in a mutable grid.
//
// Function Naming: "Mutable" suffix indicates this works with [][]byte
// Part1's isAccessible() works with []string (immutable)
// Both implement the same logic, but on different types
//
// Why duplicate instead of generics?
// - Different types ([]string vs [][]byte) for different use cases
// - Part1 doesn't need mutability (simpler, safer)
// - Part2 needs mutability (efficiency)
// - Small amount of duplication (< 20 lines) is acceptable in Go
// - "A little copying is better than a little dependency" - Go proverb
func isAccessibleMutable(grid [][]byte, row, col int) bool {
	adjacentCount := 0

	// Same direction vectors as Part1 - mathematical pattern for neighbors
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // top row
		{0, -1}, {0, 1},            // left and right
		{1, -1}, {1, 0}, {1, 1},    // bottom row
	}

	for _, dir := range directions {
		newRow := row + dir[0]
		newCol := col + dir[1]

		// Same bounds checking as Part1, but with [][]byte instead of []string
		if newRow >= 0 && newRow < len(grid) &&
			newCol >= 0 && newCol < len(grid[newRow]) &&
			grid[newRow][newCol] == '@' {
			adjacentCount++
		}
	}

	return adjacentCount < 4
}
