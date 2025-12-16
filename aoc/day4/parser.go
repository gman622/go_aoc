package day4

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parser reads and parses input for Day 4.
//
// Go Best Practice: Accept interfaces, return concrete types
// The parser accepts io.Reader (interface) making it testable with strings.NewReader,
// but returns concrete []string for clarity. This enables testing without file I/O:
//   parser := NewParser(strings.NewReader("..@@..\n@@..@@"))
type Parser struct {
	// scanner is unexported (lowercase) - encapsulation principle
	// Callers interact through methods, not direct field access
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
//
// This demonstrates several Go patterns:
// - Input validation with helpful error messages
// - Error wrapping with %w (preserves error chain for errors.Is/As)
// - Line number tracking for debugging
// - Defensive programming: validate data meets expectations
func (p *Parser) ParseAll() ([]string, error) {
	var lines []string
	lineNum := 0

	// bufio.Scanner pattern: efficient line-by-line reading
	// Better than reading entire file into memory at once
	for p.scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue // Skip empty lines gracefully
		}

		// Input validation: Catch malformed data early with clear errors
		// This prevents cryptic failures later in the solution logic
		for _, ch := range line {
			if ch != '@' && ch != '.' {
				return nil, fmt.Errorf("line %d: invalid character %q, expected '@' or '.'", lineNum, ch)
			}
		}

		lines = append(lines, line)
	}

	// Always check scanner.Err() - Scan() returns false on error OR EOF
	if err := p.scanner.Err(); err != nil {
		// Error wrapping with %w: preserves original error for debugging
		// Caller can use errors.Is(err, io.EOF) to check specific errors
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return lines, nil
}

// FromFile creates a parser from a file path and parses all lines immediately.
//
// Convenience Function: Combines common operations (open + parse) into one call.
// Uses defer for guaranteed cleanup - file closes even if parsing fails.
//
// Error Handling: Each layer adds context to errors, making debugging easier:
//   - os.Open error: "opening file: no such file"
//   - parser.ParseAll error: "line 5: invalid character 'x'"
func FromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	// defer ensures file.Close() runs when function exits
	// Critical for resource management - prevents file descriptor leaks
	defer file.Close()

	parser := NewParser(file)
	return parser.ParseAll()
}
