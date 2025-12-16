package day10

import (
	"strings"
	"testing"
)

func TestPart2Example(t *testing.T) {
	input := `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}`

	lines := strings.Split(input, "\n")
	totalPresses := 0

	expected := []int{10, 12, 11} // Expected presses for each machine

	for i, line := range lines {
		machine, err := ParseMachine(line)
		if err != nil {
			t.Fatalf("Failed to parse machine %d: %v", i+1, err)
		}

		t.Logf("Machine %d: %d counters, %d buttons", i+1, len(machine.Joltages), len(machine.Buttons))
		t.Logf("  Joltages: %v", machine.Joltages)

		presses := SolveMinJoltage(machine)
		if presses < 0 {
			t.Fatalf("Machine %d has no solution for joltages", i+1)
		}

		t.Logf("  Result: %d presses (expected %d)", presses, expected[i])

		if presses != expected[i] {
			t.Errorf("Machine %d: got %d presses, want %d", i+1, presses, expected[i])
		}

		totalPresses += presses
	}

	expectedTotal := 33
	if totalPresses != expectedTotal {
		t.Errorf("Total presses: got %d, want %d", totalPresses, expectedTotal)
	}
}

func TestSolveMinJoltage(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		{
			name: "first example machine",
			line: "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			want: 10,
		},
		{
			name: "second example machine",
			line: "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			want: 12,
		},
		{
			name: "third example machine",
			line: "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			want: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine, err := ParseMachine(tt.line)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			got := SolveMinJoltage(machine)
			if got != tt.want {
				t.Errorf("SolveMinJoltage() = %d, want %d", got, tt.want)
			}
		})
	}
}
