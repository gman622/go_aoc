package day3

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// BankParser reads and parses battery banks from input
type BankParser struct {
	scanner *bufio.Scanner
}

// NewBankParser creates a parser from an io.Reader
func NewBankParser(r io.Reader) *BankParser {
	return &BankParser{
		scanner: bufio.NewScanner(r),
	}
}

// ParseAll reads all battery banks from the input
func (p *BankParser) ParseAll() ([]string, error) {
	var banks []string
	lineNum := 0

	for p.scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}

		// Validate that the line contains only digits
		for _, ch := range line {
			if ch < '0' || ch > '9' {
				return nil, fmt.Errorf("line %d: invalid character %q, expected digits only", lineNum, ch)
			}
		}

		banks = append(banks, line)
	}

	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return banks, nil
}

// FromFile creates a parser from a file path and parses all banks immediately
func FromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	parser := NewBankParser(file)
	return parser.ParseAll()
}
