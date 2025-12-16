package day9

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Parser reads and parses input for Day 9.
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

// ParseAll reads all lines from the input.
func (p *Parser) ParseAll() ([]string, error) {
	var lines []string
	lineNum := 0

	for p.scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}

		// TODO: Add validation for expected input format
		
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

// ParsePoints parses lines into Point coordinates.
// Each line should be in the format: X,Y
func ParsePoints(lines []string) ([]Point, error) {
	points := make([]Point, 0, len(lines))

	for i, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 coordinates, got %d", i+1, len(parts))
		}

		x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid X coordinate: %w", i+1, err)
		}

		y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid Y coordinate: %w", i+1, err)
		}

		points = append(points, Point{X: x, Y: y})
	}

	return points, nil
}

// PointsFromFile is a convenience function that loads and parses points from a file.
func PointsFromFile(path string) ([]Point, error) {
	lines, err := FromFile(path)
	if err != nil {
		return nil, err
	}
	return ParsePoints(lines)
}
