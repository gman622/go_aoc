package day25

import (
	"fmt"
	"math/big"
)

// CountPathsDAG counts all paths from start to target in a DAG (no cycle detection needed!)
// This is MUCH faster than general graph path counting because:
// 1. No cycles = no need to track visited set
// 2. Memo key is just (node) - very simple
// 3. Each memo entry gets reused millions of times
//
// Uses big.Int to handle astronomical path counts (quintillions, sextillions, and beyond!)
func CountPathsDAG(graph Graph, start, target Node, memo map[Node]*big.Int) *big.Int {
	// Base case: reached the target
	if start == target {
		return big.NewInt(1)
	}

	// Check memoization cache
	if count, ok := memo[start]; ok {
		return new(big.Int).Set(count) // Return a copy
	}

	totalPaths := big.NewInt(0)

	// Explore all neighbors (no cycle check needed - it's a DAG!)
	for neighbor := range graph[start] {
		paths := CountPathsDAG(graph, neighbor, target, memo)
		totalPaths.Add(totalPaths, paths)
	}

	// Cache the result
	memo[start] = new(big.Int).Set(totalPaths)
	return totalPaths
}

// CountAllPathsToReactors counts total paths from START to all reactors combined
// Returns the count as a big.Int to handle numbers beyond int64 limits
func CountAllPathsToReactors(graph Graph) *big.Int {
	start := Node("START")
	reactors := graph.GetReactors()

	if len(reactors) == 0 {
		return big.NewInt(0)
	}

	totalPaths := big.NewInt(0)

	// Shared memo across all reactors for maximum efficiency
	memo := make(map[Node]*big.Int)

	fmt.Println("Counting quantum superposition paths...")
	for i, reactor := range reactors {
		paths := CountPathsDAG(graph, start, reactor, memo)
		fmt.Printf("  REACTOR_%d: %s paths (memo size: %d)\n", i+1, paths.String(), len(memo))
		totalPaths.Add(totalPaths, paths)
	}

	fmt.Printf("\nTotal paths to all reactors: %s\n", formatBigInt(totalPaths))
	fmt.Printf("Memoization cache entries: %d\n", len(memo))
	return totalPaths
}

// formatBigInt formats a big.Int with commas and name
func formatBigInt(n *big.Int) string {
	s := n.String()

	// Add commas
	var result []byte
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(digit))
	}

	formatted := string(result)

	// Determine name
	digits := len(s)
	var name string
	switch {
	case digits <= 3:
		name = ""
	case digits <= 6:
		name = "thousand"
	case digits <= 9:
		name = "million"
	case digits <= 12:
		name = "billion"
	case digits <= 15:
		name = "trillion"
	case digits <= 18:
		name = "quadrillion"
	case digits <= 21:
		name = "quintillion"
	case digits <= 24:
		name = "sextillion"
	case digits <= 27:
		name = "septillion"
	case digits <= 30:
		name = "octillion"
	case digits <= 33:
		name = "nonillion"
	case digits <= 36:
		name = "decillion"
	case digits <= 39:
		name = "undecillion"
	case digits <= 42:
		name = "duodecillion"
	default:
		name = fmt.Sprintf("%d digits", digits)
	}

	if name != "" {
		return fmt.Sprintf("%s (%s)", formatted, name)
	}
	return formatted
}
