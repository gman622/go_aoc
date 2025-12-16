package day11

import "fmt"

// Checkpoint bitmasks for required nodes
const (
	checkpointDac  = 1 << 0 // bit 0: dac visited
	checkpointFft  = 1 << 1 // bit 1: fft visited
	checkpointBoth = checkpointDac | checkpointFft
)

// checkpointMask represents which checkpoints have been visited as a bitmask
type checkpointMask int

// Part2 solves Day 11 Part 2: Find paths from "svr" to "out" that visit both "dac" and "fft"
//
// Key insight: The massive path explosion (390 trillion paths) happens BEFORE reaching
// dac/fft. The graph has high connectivity in the early portion, causing exponential
// path divergence. Memoization caches these early states and reuses them billions of
// times, reducing 390T path enumerations to ~1300 unique state calculations.
func Part2(inputPath string) (int, error) {
	graph, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	memo := make(map[state]int)
	pathCount := countPathsWithCheckpoints(graph, "svr", "out", 0, make(map[string]bool), memo)

	return pathCount, nil
}

// state represents a node and which required checkpoints have been visited
type state struct {
	node        string
	checkpoints checkpointMask
}

// countPathsWithCheckpoints counts paths visiting both required checkpoints.
// The explosion happens early (before checkpoints), so memoization has massive impact.
func countPathsWithCheckpoints(graph Graph, current, target string, checkpoints checkpointMask, visited map[string]bool, memo map[state]int) int {
	// Update checkpoints if we're at a required node
	switch current {
	case "dac":
		checkpoints |= checkpointDac
	case "fft":
		checkpoints |= checkpointFft
	}

	// Base case: reached target
	if current == target {
		if checkpoints == checkpointBoth {
			return 1
		}
		return 0
	}

	// Check memoization cache, but only if current isn't in the active path.
	// This prevents reusing cached results during cycles (visited nodes in the
	// current path can't use memo because the path context differs).
	s := state{current, checkpoints}
	if count, ok := memo[s]; ok && !visited[current] {
		return count
	}

	// Explore neighbors with cycle detection
	visited[current] = true
	defer delete(visited, current)

	totalPaths := 0
	for _, neighbor := range graph[current] {
		if !visited[neighbor] {
			totalPaths += countPathsWithCheckpoints(graph, neighbor, target, checkpoints, visited, memo)
		}
	}

	memo[s] = totalPaths
	return totalPaths
}
