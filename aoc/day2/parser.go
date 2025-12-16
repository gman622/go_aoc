package day2

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Range represents a product ID range with start and end values
type Range struct {
	Start, End int
}

// RangeParser reads and parses product ID ranges from input
type RangeParser struct {
	reader io.Reader
}

// NewRangeParser creates a parser from an io.Reader
func NewRangeParser(r io.Reader) *RangeParser {
	return &RangeParser{reader: r}
}

// FromFile creates a parser from a file path
func FromFile(path string) (*RangeParser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// Read the entire file content since we need to close the file
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty input file")
	}

	line := scanner.Text()
	return NewRangeParser(strings.NewReader(line)), scanner.Err()
}

// ParseAll reads and parses all ranges from the input
func (p *RangeParser) ParseAll() ([]Range, error) {
	scanner := bufio.NewScanner(p.reader)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty input")
	}

	line := scanner.Text()
	ranges, err := parseRanges(line)
	if err != nil {
		return nil, fmt.Errorf("parsing ranges: %w", err)
	}

	return ranges, scanner.Err()
}

// parseRanges parses comma-separated ranges like "11-22,95-115"
func parseRanges(line string) ([]Range, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	// Remove trailing comma if present
	line = strings.TrimSuffix(line, ",")

	parts := strings.Split(line, ",")
	ranges := make([]Range, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		nums := strings.Split(part, "-")
		if len(nums) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", part)
		}

		start, err := strconv.Atoi(strings.TrimSpace(nums[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start number: %w", err)
		}

		end, err := strconv.Atoi(strings.TrimSpace(nums[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end number: %w", err)
		}

		ranges = append(ranges, Range{Start: start, End: end})
	}

	return ranges, nil
}
