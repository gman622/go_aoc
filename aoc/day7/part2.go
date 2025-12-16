package day7

import "fmt"

// Part2 solves Day 7 Part 2 - counts quantum timelines.
//
// In quantum mode, a single particle takes all paths through splitters.
// Each splitter creates a branching point where reality splits into two timelines.
// We need to count the total number of unique timelines (complete paths through the manifold).
//
// Algorithm:
// - Use recursive DFS to explore all possible paths from S to exit
// - When particle hits a splitter: branch into left and right timelines
// - When particle hits empty space: continue straight down
// - When particle exits the grid: count as 1 timeline
// - Use memoization to cache results for positions we've seen before
func Part2(inputPath string) (int, error) {
	lines, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	if len(lines) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	// Find starting position (S)
	startCol := -1
	for i, ch := range lines[0] {
		if ch == 'S' {
			startCol = i
			break
		}
	}
	if startCol == -1 {
		return 0, fmt.Errorf("no starting position 'S' found")
	}

	// Count all possible timelines using memoized recursion
	memo := make(map[string]int)
	timelines := countTimelines(lines, 0, startCol, memo)

	return timelines, nil
}

// countTimelines recursively counts the number of unique paths (timelines)
// from position (row, col) to exiting the manifold.
func countTimelines(grid []string, row, col int, memo map[string]int) int {
	// Base case: exited the bottom of the manifold
	if row >= len(grid) {
		return 1 // This path is one complete timeline
	}

	// Out of bounds horizontally - particle exits
	if col < 0 || col >= len(grid[row]) {
		return 1 // Exited sideways, counts as one timeline
	}

	// Check memoization cache
	key := fmt.Sprintf("%d,%d", row, col)
	if val, ok := memo[key]; ok {
		return val
	}

	var count int
	ch := grid[row][col]

	if ch == '^' {
		// Hit a splitter - reality splits into two timelines
		leftTimelines := countTimelines(grid, row+1, col-1, memo)
		rightTimelines := countTimelines(grid, row+1, col+1, memo)
		count = leftTimelines + rightTimelines
	} else {
		// Empty space or S - particle continues straight down
		count = countTimelines(grid, row+1, col, memo)
	}

	memo[key] = count
	return count
}
