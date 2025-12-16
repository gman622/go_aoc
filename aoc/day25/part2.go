package day25

import "fmt"

// Part2 solves Day 25 Part 2: Count all quantum superposition paths to all reactors.
//
// Algorithm:
// 1. Parse the facility network graph
// 2. For each reactor, count all valid paths from START using DFS with memoization
// 3. Sum the path counts across all reactors
//
// Key insight: The dense connectivity creates exponential path explosion.
// Memoization is critical to avoid recalculating paths from the same state.
//
// This is similar to Day 11 Part 2 (390 trillion paths) but counts paths
// to ALL reactors combined, not just one target with checkpoints.
//
// Uses math/big to handle astronomical numbers (sextillions, septillions, and beyond!)
//
// Time complexity: O(V * 2^V) worst case, but memoization reduces it dramatically
// Expected runtime: 1-10 minutes depending on graph structure
func Part2(inputPath string) (int, error) {
	// Load as DAG (directional edges only, no cycles!)
	graph, err := FromFileDAG(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	// Count all paths from START to all reactors
	totalPaths := CountAllPathsToReactors(graph)

	// For compatibility with the int return type, we return 0 if the number
	// is too large. The actual result is printed during execution.
	if totalPaths.IsInt64() {
		return int(totalPaths.Int64()), nil
	}

	// Number too large for int - return special value
	// (The real answer was already printed during CountAllPathsToReactors)
	return -1, nil
}
