package day7

import "fmt"

// Part1 solves Day 7 Part 1 - counts how many times tachyon beams are split.
//
// Algorithm:
// - Beam starts at 'S' and moves downward
// - When a beam hits a splitter ('^'), it stops and creates two new beams at left and right positions
// - Simulate the beam propagation row by row, tracking active beam columns
// - Count each split event
func Part1(inputPath string) (int, error) {
	lines, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	if len(lines) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	// Find starting position (S)
	startCol := -1
	for i, ch := range lines[0] {
		if ch == 'S' {
			startCol = i
			break
		}
	}
	if startCol == -1 {
		return 0, fmt.Errorf("no starting position 'S' found")
	}

	// Simulate beam propagation
	// activeBeams tracks the column positions of beams at the current row
	activeBeams := make(map[int]bool)
	activeBeams[startCol] = true
	splitCount := 0

	// Process each row after the starting row
	for rowIdx := 1; rowIdx < len(lines); rowIdx++ {
		row := lines[rowIdx]
		nextBeams := make(map[int]bool)

		// Process each active beam
		for col := range activeBeams {
			// Check bounds
			if col < 0 || col >= len(row) {
				continue
			}

			// Check what's at this position
			if row[col] == '^' {
				// Beam hits a splitter - count the split and create two new beams
				splitCount++
				// Add left and right beams (will be validated next iteration)
				nextBeams[col-1] = true
				nextBeams[col+1] = true
			} else {
				// Empty space or S - beam continues
				nextBeams[col] = true
			}
		}

		activeBeams = nextBeams

		// If no beams remain, we're done
		if len(activeBeams) == 0 {
			break
		}
	}

	return splitCount, nil
}
