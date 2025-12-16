package day10

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parser reads and parses input for Day 10.
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

// ParseAll reads all lines from the input and returns raw line strings.
func (p *Parser) ParseAll() ([]string, error) {
	var lines []string

	for p.scanner.Scan() {
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return lines, nil
}

// FromFile creates a parser from a file path and parses all lines immediately.
func FromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	parser := NewParser(file)
	return parser.ParseAll()
}
