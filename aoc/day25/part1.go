package day25

import "fmt"

// Part1 solves Day 25 Part 1: Find the time when the last reactor receives its activation signal.
//
// Algorithm:
// 1. Parse the facility network graph
// 2. Run Dijkstra's shortest path from START to all nodes
// 3. Find the maximum shortest path distance among all reactor cores
//
// Time complexity: O(E log V) where E = edges, V = nodes
func Part1(inputPath string) (int, error) {
	graph, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	start := Node("START")

	// Find maximum shortest path distance to any reactor
	maxDist := FindMaxReactorDistance(graph, start)

	return maxDist, nil
}
