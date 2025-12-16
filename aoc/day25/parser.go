package day25

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ParseGraph parses the facility network from an io.Reader (bidirectional edges)
func ParseGraph(r io.Reader) (Graph, error) {
	graph := make(Graph)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse line format: "A-B:10"
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		nodes := strings.Split(parts[0], "-")
		if len(nodes) != 2 {
			return nil, fmt.Errorf("invalid node format: %s", parts[0])
		}

		cost, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid cost: %s", parts[1])
		}

		from := Node(strings.TrimSpace(nodes[0]))
		to := Node(strings.TrimSpace(nodes[1]))

		graph.AddEdge(from, to, cost)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return graph, nil
}

// ParseDAG parses the facility network as a DAG (directional edges only)
func ParseDAG(r io.Reader) (Graph, error) {
	graph := make(Graph)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse line format: "A-B:10"
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		nodes := strings.Split(parts[0], "-")
		if len(nodes) != 2 {
			return nil, fmt.Errorf("invalid node format: %s", parts[0])
		}

		cost, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid cost: %s", parts[1])
		}

		from := Node(strings.TrimSpace(nodes[0]))
		to := Node(strings.TrimSpace(nodes[1]))

		// Only add edge in one direction (from → to)
		graph.AddDirectedEdge(from, to, cost)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return graph, nil
}

// FromFile reads and parses the graph from a file (bidirectional)
func FromFile(path string) (Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	return ParseGraph(file)
}

// FromFileDAG reads and parses the graph as a DAG (directional)
func FromFileDAG(path string) (Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	return ParseDAG(file)
}
