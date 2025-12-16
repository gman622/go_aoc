package day4

import (
	"strings"
	"testing"
)

func TestPart1Example(t *testing.T) {
	input := `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

	parser := NewParser(strings.NewReader(input))
	grid, err := parser.ParseAll()
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	count := 0
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			if grid[row][col] == '@' && isAccessible(grid, row, col) {
				count++
			}
		}
	}

	expected := 13
	if count != expected {
		t.Errorf("expected %d accessible rolls, got %d", expected, count)
	}
}

func TestPart2Example(t *testing.T) {
	input := `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

	parser := NewParser(strings.NewReader(input))
	lines, err := parser.ParseAll()
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Convert to mutable grid
	grid := make([][]byte, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}

	totalRemoved := 0

	// Keep removing accessible rolls until none remain
	for {
		accessible := findAccessibleRolls(grid)
		if len(accessible) == 0 {
			break
		}

		// Remove all accessible rolls
		for _, pos := range accessible {
			grid[pos.row][pos.col] = '.'
		}

		totalRemoved += len(accessible)
	}

	expected := 43
	if totalRemoved != expected {
		t.Errorf("expected %d total removed, got %d", expected, totalRemoved)
	}
}
