package day4

import "fmt"

// Part1 solves Day 4 Part 1: count rolls of paper accessible by forklifts.
//
// Problem: Count how many '@' symbols have fewer than 4 adjacent '@' symbols
// (checking all 8 surrounding positions: horizontal, vertical, and diagonal).
//
// This demonstrates a common pattern in AoC and real-world problems:
// - Grid/matrix traversal
// - Adjacency checking (neighbors)
// - Boundary validation
//
// Approach: Brute force - check every cell. For larger grids, consider:
// - Spatial indexing (quadtree, R-tree)
// - Only checking '@' positions (skip '.')
// But for this problem size (~140x150), simple iteration is fastest and clearest.
func Part1(inputPath string) (int, error) {
	// Delegate parsing to FromFile - separation of concerns
	// Part1 focuses on solving, not file I/O details
	grid, err := FromFile(inputPath)
	if err != nil {
		// Error wrapping adds context at each layer
		// Final error might be: "loading input: opening file: no such file"
		return 0, fmt.Errorf("loading input: %w", err)
	}

	count := 0
	// Nested loop pattern for 2D grid traversal
	// Time complexity: O(rows * cols * 8) = O(n) where n is total cells
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			// Short-circuit evaluation: check '@' first (cheaper than function call)
			if grid[row][col] == '@' && isAccessible(grid, row, col) {
				count++
			}
		}
	}

	return count, nil
}

// isAccessible returns true if a roll at (row, col) has fewer than 4 adjacent rolls.
//
// Helper Function Pattern: Extract complex logic into named functions for:
// - Readability: Function name documents intent
// - Testability: Can test isAccessible() independently
// - Reusability: Used by both Part1 (if needed elsewhere)
// - Single Responsibility: Each function does one thing well
//
// Adjacency Checking: Common pattern in grid problems (Conway's Game of Life, etc.)
func isAccessible(grid []string, row, col int) bool {
	adjacentCount := 0

	// Direction vectors: Mathematical approach to neighbor checking
	// Each [2]int is [rowOffset, colOffset] relative to current position
	// This is more maintainable than 8 separate if statements
	//
	// Layout visualization:
	//   [-1,-1] [-1,0] [-1,1]    NW  N  NE
	//   [ 0,-1]  [X,Y] [ 0,1]     W  @   E
	//   [ 1,-1] [ 1,0] [ 1,1]    SW  S  SE
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // top row
		{0, -1}, {0, 1},            // left and right (skip center)
		{1, -1}, {1, 0}, {1, 1},    // bottom row
	}

	// Check each of the 8 surrounding cells
	for _, dir := range directions {
		newRow := row + dir[0]
		newCol := col + dir[1]

		// Bounds checking: critical for grid problems to avoid panics
		// Go will panic on out-of-bounds slice access, unlike some languages
		// Order matters: check row bounds before accessing grid[newRow]
		if newRow >= 0 && newRow < len(grid) &&
			newCol >= 0 && newCol < len(grid[newRow]) &&
			grid[newRow][newCol] == '@' {
			adjacentCount++
		}
	}

	// Problem constraint: accessible if FEWER than 4 adjacent (0-3 is accessible)
	return adjacentCount < 4
}
