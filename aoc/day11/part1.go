package day11

import "fmt"

// Part1 solves Day 11 Part 1: Find all paths from "you" to "out"
func Part1(inputPath string) (int, error) {
	graph, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	// Count all paths from "you" to "out" using DFS
	pathCount := countPaths(graph, "you", "out", make(map[string]bool))
	return pathCount, nil
}

// countPaths uses DFS to count all paths from start to end.
// visited tracks nodes in the current path to prevent cycles.
func countPaths(graph Graph, current, target string, visited map[string]bool) int {
	// Base case: reached the target
	if current == target {
		return 1
	}

	// Mark current node as visited in this path
	visited[current] = true
	defer delete(visited, current) // Backtrack when returning

	totalPaths := 0

	// Explore all neighbors
	for _, neighbor := range graph[current] {
		// Skip if already in current path (cycle detection)
		if visited[neighbor] {
			continue
		}

		// Recursively count paths from this neighbor
		totalPaths += countPaths(graph, neighbor, target, visited)
	}

	return totalPaths
}
