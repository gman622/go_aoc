package day6

import (
	"fmt"
	"strconv"
	"strings"
)

// Problem represents a single vertical math problem
type Problem struct {
	Numbers   []int
	Operation rune // '+' or '*'
}

// Calculate computes the result of this problem
func (p Problem) Calculate() int {
	if len(p.Numbers) == 0 {
		return 0
	}

	result := p.Numbers[0]
	for i := 1; i < len(p.Numbers); i++ {
		switch p.Operation {
		case '+':
			result += p.Numbers[i]
		case '*':
			result *= p.Numbers[i]
		}
	}
	return result
}

// ReadingMode determines how to interpret the worksheet columns
type ReadingMode int

const (
	LeftToRight ReadingMode = iota // Part 1: each field is a complete number
	RightToLeft                    // Part 2: each column contains digits of one number
)

// ParseProblems extracts all vertical problems using the specified reading mode
func ParseProblems(lines []string, mode ReadingMode) ([]Problem, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("input too short: need at least 2 lines")
	}

	operationLine := lines[len(lines)-1]
	numberLines := lines[:len(lines)-1]

	if mode == LeftToRight {
		return parseLeftToRight(numberLines, operationLine)
	}
	return parseRightToLeft(numberLines, operationLine)
}

// parseLeftToRight interprets each space-separated field as a complete number (Part 1)
func parseLeftToRight(numberLines []string, operationLine string) ([]Problem, error) {
	// Parse each number line into fields
	var allFields [][]string
	for _, line := range numberLines {
		allFields = append(allFields, strings.Fields(line))
	}

	opFields := strings.Fields(operationLine)

	var problems []Problem
	for i := range len(opFields) {
		if len(opFields[i]) != 1 {
			continue
		}

		op := rune(opFields[i][0])
		if op != '+' && op != '*' {
			continue
		}

		// Extract the i-th field from each row
		var numbers []int
		for _, fields := range allFields {
			if i < len(fields) {
				if num, err := strconv.Atoi(fields[i]); err == nil {
					numbers = append(numbers, num)
				}
			}
		}

		if len(numbers) > 0 {
			problems = append(problems, Problem{Numbers: numbers, Operation: op})
		}
	}

	return problems, nil
}

// parseRightToLeft interprets each column as digits forming a number (Part 2)
func parseRightToLeft(numberLines []string, operationLine string) ([]Problem, error) {
	maxLen := len(operationLine)
	for _, line := range numberLines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Find problem column groups (separated by all-space columns)
	type group struct {
		cols []int
		op   rune
	}

	var groups []group
	var currentCols []int
	var currentOp rune

	for col := maxLen - 1; col >= 0; col-- {
		// Check if column is all spaces
		allSpaces := true
		for _, line := range numberLines {
			if col < len(line) && line[col] != ' ' {
				allSpaces = false
				break
			}
		}

		if allSpaces {
			if len(currentCols) > 0 {
				groups = append(groups, group{cols: currentCols, op: currentOp})
				currentCols = nil
				currentOp = 0
			}
		} else {
			currentCols = append(currentCols, col)
			if col < len(operationLine) {
				ch := rune(operationLine[col])
				if ch == '+' || ch == '*' {
					currentOp = ch
				}
			}
		}
	}

	if len(currentCols) > 0 {
		groups = append(groups, group{cols: currentCols, op: currentOp})
	}

	// Extract numbers from each group
	var problems []Problem
	for _, g := range groups {
		if g.op == 0 {
			continue
		}

		var numbers []int
		for _, col := range g.cols {
			// Read digits top-to-bottom in this column
			var digits strings.Builder
			for _, line := range numberLines {
				if col < len(line) && line[col] >= '0' && line[col] <= '9' {
					digits.WriteByte(line[col])
				}
			}

			if digits.Len() > 0 {
				if num, err := strconv.Atoi(digits.String()); err == nil {
					numbers = append(numbers, num)
				}
			}
		}

		if len(numbers) > 0 {
			problems = append(problems, Problem{Numbers: numbers, Operation: g.op})
		}
	}

	return problems, nil
}

// SolveWorksheet calculates the grand total using the specified reading mode
func SolveWorksheet(lines []string, mode ReadingMode) (int, error) {
	problems, err := ParseProblems(lines, mode)
	if err != nil {
		return 0, fmt.Errorf("parsing problems: %w", err)
	}

	grandTotal := 0
	for _, problem := range problems {
		grandTotal += problem.Calculate()
	}

	return grandTotal, nil
}

// For debugging: format a problem as a string
func (p Problem) String() string {
	nums := make([]string, len(p.Numbers))
	for i, n := range p.Numbers {
		nums[i] = strconv.Itoa(n)
	}
	return fmt.Sprintf("%s %c %d", strings.Join(nums, fmt.Sprintf(" %c ", p.Operation)), '=', p.Calculate())
}
