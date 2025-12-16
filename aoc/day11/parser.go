package day11

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parser reads and parses input for Day 11.
//
// Go Best Practice: Accept interfaces, return concrete types
// The parser accepts io.Reader (interface) making it testable with strings.NewReader,
// but returns concrete []string for clarity. This enables testing without file I/O:
//   parser := NewParser(strings.NewReader("test input"))
type Parser struct {
	scanner *bufio.Scanner
}

// NewParser creates a parser from an io.Reader.
//
// Constructor Pattern: New* functions are the idiomatic way to create instances
// in Go. This allows initialization logic and ensures fields are set correctly.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

// Graph represents the device connection graph as an adjacency list.
// Each device maps to a list of devices it outputs to.
type Graph map[string][]string

// ParseAll reads all lines from the input and builds a directed graph.
func (p *Parser) ParseAll() (Graph, error) {
	graph := make(Graph)
	lineNum := 0

	for p.scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}

		// Parse "device: output1 output2 ..." format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format, expected 'device: outputs'", lineNum)
		}

		device := strings.TrimSpace(parts[0])
		outputsStr := strings.TrimSpace(parts[1])

		// Split outputs by whitespace
		outputs := strings.Fields(outputsStr)

		graph[device] = outputs
	}

	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return graph, nil
}

// FromFile creates a parser from a file path and parses the graph immediately.
func FromFile(path string) (Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	parser := NewParser(file)
	return parser.ParseAll()
}
