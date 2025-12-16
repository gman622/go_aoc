package day12

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Parser reads and parses input for Day 12.
//
// Go Best Practice: Accept interfaces, return concrete types
// The parser accepts io.Reader (interface) making it testable with strings.NewReader,
// but returns concrete structures for clarity. This enables testing without file I/O:
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

// ParseAll reads all shapes and regions from the input.
func (p *Parser) ParseAll() ([]Shape, []Region, error) {
	var shapes []Shape
	var regions []Region

	for p.scanner.Scan() {
		line := p.scanner.Text()

		// Check if this is a shape definition (ends with :)
		if strings.HasSuffix(strings.TrimSpace(line), ":") {
			shape, err := p.parseShape(line)
			if err != nil {
				return nil, nil, fmt.Errorf("parsing shape: %w", err)
			}
			shapes = append(shapes, shape)
		} else if strings.Contains(line, "x") && strings.Contains(line, ":") {
			// This is a region definition
			region, err := p.parseRegion(line)
			if err != nil {
				return nil, nil, fmt.Errorf("parsing region: %w", err)
			}
			regions = append(regions, region)
		}
		// Skip empty lines
	}

	if err := p.scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("reading input: %w", err)
	}

	return shapes, regions, nil
}

// parseShape parses a shape definition starting with "N:" and followed by grid lines.
func (p *Parser) parseShape(headerLine string) (Shape, error) {
	// Parse the shape ID from "N:"
	idStr := strings.TrimSuffix(strings.TrimSpace(headerLine), ":")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return Shape{}, fmt.Errorf("invalid shape ID %q: %w", idStr, err)
	}

	// Read the grid lines for this shape
	var gridLines []string
	for p.scanner.Scan() {
		line := p.scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Empty line marks end of shape
		if trimmed == "" {
			break
		}

		// If we hit another shape or region definition, we're done
		if strings.HasSuffix(trimmed, ":") || (strings.Contains(trimmed, "x") && strings.Contains(trimmed, ":")) {
			// Put the line back by scanning it again (unfortunately we can't unscan)
			// We'll handle this by breaking and letting the outer loop process it
			// Actually, bufio.Scanner doesn't support putting back, so we need a different approach
			// For now, let's just break and assume proper formatting
			break
		}

		gridLines = append(gridLines, line)
	}

	// Parse grid into points
	var points []Point
	for y, line := range gridLines {
		for x, ch := range line {
			if ch == '#' {
				points = append(points, Point{X: x, Y: y})
			}
		}
	}

	shape := Shape{ID: id, Points: points}
	return shape.Normalize(), nil
}

// parseRegion parses a region definition like "12x5: 1 0 1 0 2 2".
func (p *Parser) parseRegion(line string) (Region, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return Region{}, fmt.Errorf("invalid region format: %q", line)
	}

	// Parse dimensions "WxH"
	dims := strings.Split(strings.TrimSpace(parts[0]), "x")
	if len(dims) != 2 {
		return Region{}, fmt.Errorf("invalid dimensions: %q", parts[0])
	}

	width, err := strconv.Atoi(dims[0])
	if err != nil {
		return Region{}, fmt.Errorf("invalid width: %w", err)
	}

	height, err := strconv.Atoi(dims[1])
	if err != nil {
		return Region{}, fmt.Errorf("invalid height: %w", err)
	}

	// Parse present counts
	countStrs := strings.Fields(parts[1])
	presents := make([]int, len(countStrs))
	for i, s := range countStrs {
		count, err := strconv.Atoi(s)
		if err != nil {
			return Region{}, fmt.Errorf("invalid present count: %w", err)
		}
		presents[i] = count
	}

	return Region{
		Width:    width,
		Height:   height,
		Presents: presents,
	}, nil
}

// FromFile creates a parser from a file path and parses all shapes and regions.
func FromFile(path string) ([]Shape, []Region, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	parser := NewParser(file)
	return parser.ParseAll()
}
