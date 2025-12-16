package day8

import (
	"fmt"
	"math"
)

// Point3D represents a junction box position in 3D space.
type Point3D struct {
	X, Y, Z int
}

// DistanceTo calculates the Euclidean distance to another point.
func (p Point3D) DistanceTo(other Point3D) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	dz := float64(p.Z - other.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// String implements fmt.Stringer for debugging.
func (p Point3D) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

// UnionFind is a disjoint-set data structure for tracking connected components (circuits).
// It uses path compression and union by rank for efficient operations.
type UnionFind struct {
	parent []int // parent[i] is the parent of node i
	rank   []int // rank[i] is the approximate depth of the tree rooted at i
	size   []int // size[i] is the size of the component containing i (only valid for root nodes)
}

// NewUnionFind creates a new UnionFind structure with n elements.
// Initially, each element is in its own set.
func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
		size:   make([]int, n),
	}
	for i := range n {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

// Find returns the root of the set containing x.
// Uses path compression to flatten the tree structure.
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

// Union merges the sets containing x and y.
// Uses union by rank to keep trees balanced.
// Returns true if x and y were in different sets, false otherwise.
func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // Already in the same set
	}

	// Union by rank: attach smaller tree under larger tree
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
		uf.rank[rootX]++
	}

	return true
}

// ComponentSizes returns the sizes of all connected components.
func (uf *UnionFind) ComponentSizes() []int {
	sizes := make(map[int]int)
	for i := 0; i < len(uf.parent); i++ {
		root := uf.Find(i)
		sizes[root] = uf.size[root]
	}

	result := make([]int, 0, len(sizes))
	for _, size := range sizes {
		result = append(result, size)
	}
	return result
}
