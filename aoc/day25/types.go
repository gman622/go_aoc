package day25

import (
	"fmt"
	"strings"
)

// Node represents a chamber in the underground facility
type Node string

// IsReactor checks if this node is a reactor core
func (n Node) IsReactor() bool {
	return strings.HasPrefix(string(n), "REACTOR_")
}

// Edge represents a tunnel between two chambers
type Edge struct {
	From Node
	To   Node
	Cost int // tunnel length in meters (= seconds for signal travel)
}

func (e Edge) String() string {
	return fmt.Sprintf("%s-%s:%d", e.From, e.To, e.Cost)
}

// Graph represents the facility network as an adjacency list with weights
type Graph map[Node]map[Node]int

// AddEdge adds a bidirectional edge to the graph
func (g Graph) AddEdge(from, to Node, cost int) {
	if g[from] == nil {
		g[from] = make(map[Node]int)
	}
	if g[to] == nil {
		g[to] = make(map[Node]int)
	}
	g[from][to] = cost
	g[to][from] = cost
}

// AddDirectedEdge adds a one-way edge to the graph (for DAG)
func (g Graph) AddDirectedEdge(from, to Node, cost int) {
	if g[from] == nil {
		g[from] = make(map[Node]int)
	}
	g[from][to] = cost
}

// Neighbors returns all adjacent nodes
func (g Graph) Neighbors(node Node) []Node {
	neighbors := make([]Node, 0, len(g[node]))
	for neighbor := range g[node] {
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

// GetReactors returns all reactor nodes in the graph
func (g Graph) GetReactors() []Node {
	reactors := make([]Node, 0)
	seen := make(map[Node]bool)

	for from := range g {
		if from.IsReactor() && !seen[from] {
			reactors = append(reactors, from)
			seen[from] = true
		}
		for to := range g[from] {
			if to.IsReactor() && !seen[to] {
				reactors = append(reactors, to)
				seen[to] = true
			}
		}
	}
	return reactors
}

// Path represents a route through the network
type Path struct {
	Nodes []Node
	Cost  int
}

// String returns a readable representation of the path
func (p Path) String() string {
	nodeStrs := make([]string, len(p.Nodes))
	for i, n := range p.Nodes {
		nodeStrs[i] = string(n)
	}
	return fmt.Sprintf("%s (cost: %d)", strings.Join(nodeStrs, " → "), p.Cost)
}
