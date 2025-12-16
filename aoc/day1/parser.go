package day1

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Rotation represents a dial rotation instruction (L10, R25, etc.)
type Rotation struct {
	Direction rune // 'L' or 'R'
	Distance  int
}

// RotationParser reads and parses dial rotation instructions from input
type RotationParser struct {
	scanner *bufio.Scanner
}

// NewRotationParser creates a parser from an io.Reader
func NewRotationParser(r io.Reader) *RotationParser {
	return &RotationParser{
		scanner: bufio.NewScanner(r),
	}
}

// Parse reads all rotations and applies a function to each one
func (p *RotationParser) Parse(fn func(Rotation) error) error {
	lineNum := 0
	for p.scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}

		rotation, err := parseRotation(line)
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}

		if err := fn(rotation); err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
	}

	if err := p.scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	return nil
}

// FromFile creates a parser from a file path
func FromFile(path string) (*RotationParser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	return NewRotationParser(f), nil
}

// ProcessFile is a convenience function that opens a file, parses it, and processes each rotation
func ProcessFile(path string, fn func(Rotation) error) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	parser := NewRotationParser(f)
	return parser.Parse(fn)
}

// parseRotation parses a rotation string like "L68" or "R48"
func parseRotation(s string) (Rotation, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return Rotation{}, fmt.Errorf("invalid rotation: too short")
	}

	dir := rune(s[0])
	if dir != 'L' && dir != 'R' {
		return Rotation{}, fmt.Errorf("invalid direction: %c", s[0])
	}

	distance, err := strconv.Atoi(s[1:])
	if err != nil {
		return Rotation{}, fmt.Errorf("invalid distance in %q: %w", s, err)
	}

	return Rotation{Direction: dir, Distance: distance}, nil
}
