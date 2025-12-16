package day25

import (
	"container/heap"
	"math"
)

// DijkstraState represents a node in Dijkstra's priority queue
type DijkstraState struct {
	node Node
	dist int
	index int // for heap interface
}

// PriorityQueue implements heap.Interface for Dijkstra's algorithm
type PriorityQueue []*DijkstraState

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*DijkstraState)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra finds shortest paths from start to all reachable nodes
func Dijkstra(graph Graph, start Node) map[Node]int {
	dist := make(map[Node]int)
	dist[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Push(&pq, &DijkstraState{node: start, dist: 0})

	visited := make(map[Node]bool)

	for pq.Len() > 0 {
		state := heap.Pop(&pq).(*DijkstraState)
		current := state.node

		if visited[current] {
			continue
		}
		visited[current] = true

		// Explore neighbors
		for neighbor, edgeCost := range graph[current] {
			if visited[neighbor] {
				continue
			}

			newDist := dist[current] + edgeCost
			if oldDist, ok := dist[neighbor]; !ok || newDist < oldDist {
				dist[neighbor] = newDist
				heap.Push(&pq, &DijkstraState{node: neighbor, dist: newDist})
			}
		}
	}

	return dist
}

// FindMaxReactorDistance finds the maximum shortest path distance to any reactor
func FindMaxReactorDistance(graph Graph, start Node) int {
	dist := Dijkstra(graph, start)
	reactors := graph.GetReactors()

	maxDist := 0
	for _, reactor := range reactors {
		if d, ok := dist[reactor]; ok {
			if d > maxDist {
				maxDist = d
			}
		} else {
			// Reactor is unreachable
			return math.MaxInt32
		}
	}

	return maxDist
}
