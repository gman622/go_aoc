package main

import (
	"fmt"
	"sort"

	day9 "adv2025/aoc/day9"
)

func main() {
	tiles, _ := day9.PointsFromFile("inputs/day9_input.txt")

	// Get first few tiles to see the pattern
	fmt.Println("First 10 original tiles:")
	for i := 0; i < 10 && i < len(tiles); i++ {
		fmt.Printf("  %d: %v\n", i, tiles[i])
	}

	// Compress
	xSet := make(map[int]bool)
	ySet := make(map[int]bool)
	for _, p := range tiles {
		xSet[p.X] = true
		ySet[p.Y] = true
	}

	xCoords := make([]int, 0, len(xSet))
	for x := range xSet {
		xCoords = append(xCoords, x)
	}
	sort.Ints(xCoords)

	yCoords := make([]int, 0, len(ySet))
	for y := range ySet {
		yCoords = append(yCoords, y)
	}
	sort.Ints(yCoords)

	xToComp := make(map[int]int)
	for i, x := range xCoords {
		xToComp[x] = i
	}
	yToComp := make(map[int]int)
	for i, y := range yCoords {
		yToComp[y] = i
	}

	fmt.Println("\nFirst 10 compressed tiles:")
	for i := 0; i < 10 && i < len(tiles); i++ {
		orig := tiles[i]
		comp := day9.Point{X: xToComp[orig.X], Y: yToComp[orig.Y]}
		fmt.Printf("  %d: %v -> %v\n", i, orig, comp)
	}

	// Check distances between consecutive points
	fmt.Println("\nDistances between consecutive points:")
	for i := 0; i < 5; i++ {
		p1 := tiles[i]
		p2 := tiles[i+1]
		origDist := abs(p1.X-p2.X) + abs(p1.Y-p2.Y)

		c1 := day9.Point{X: xToComp[p1.X], Y: yToComp[p1.Y]}
		c2 := day9.Point{X: xToComp[p2.X], Y: yToComp[p2.Y]}
		compDist := abs(c1.X-c2.X) + abs(c1.Y-c2.Y)

		fmt.Printf("  %d->%d: original=%d, compressed=%d (%.1fx)\n",
			i, i+1, origDist, compDist, float64(origDist)/float64(compDist))
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
