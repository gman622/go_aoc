package day10

import (
	"strings"
	"testing"
)

func TestPart1Example(t *testing.T) {
	input := `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}`

	lines := strings.Split(input, "\n")
	totalPresses := 0

	expected := []int{2, 3, 2} // Expected presses for each machine

	for i, line := range lines {
		machine, err := ParseMachine(line)
		if err != nil {
			t.Fatalf("Failed to parse machine %d: %v", i+1, err)
		}

		presses := SolveMinPresses(machine)
		if presses < 0 {
			t.Fatalf("Machine %d has no solution", i+1)
		}

		t.Logf("Machine %d: %s -> %d presses (expected %d)", i+1, machine, presses, expected[i])

		if presses != expected[i] {
			t.Errorf("Machine %d: got %d presses, want %d", i+1, presses, expected[i])
		}

		totalPresses += presses
	}

	expectedTotal := 7
	if totalPresses != expectedTotal {
		t.Errorf("Total presses: got %d, want %d", totalPresses, expectedTotal)
	}
}

func TestParseMachine(t *testing.T) {
	tests := []struct {
		line          string
		numLights     int
		numButtons    int
		targetPattern string
	}{
		{
			line:          "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			numLights:     4,
			numButtons:    6,
			targetPattern: ".##.",
		},
		{
			line:          "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			numLights:     5,
			numButtons:    5,
			targetPattern: "...#.",
		},
		{
			line:          "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			numLights:     6,
			numButtons:    4,
			targetPattern: ".###.#",
		},
	}

	for i, tt := range tests {
		machine, err := ParseMachine(tt.line)
		if err != nil {
			t.Fatalf("Test %d: failed to parse: %v", i, err)
		}

		if len(machine.TargetLights) != tt.numLights {
			t.Errorf("Test %d: got %d lights, want %d", i, len(machine.TargetLights), tt.numLights)
		}

		if len(machine.Buttons) != tt.numButtons {
			t.Errorf("Test %d: got %d buttons, want %d", i, len(machine.Buttons), tt.numButtons)
		}

		// Check target pattern
		pattern := make([]rune, len(machine.TargetLights))
		for j, on := range machine.TargetLights {
			if on {
				pattern[j] = '#'
			} else {
				pattern[j] = '.'
			}
		}
		gotPattern := string(pattern)
		if gotPattern != tt.targetPattern {
			t.Errorf("Test %d: got pattern %s, want %s", i, gotPattern, tt.targetPattern)
		}
	}
}

func TestSolveMinPresses(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		want   int
		verify func(*testing.T, *Machine, []int)
	}{
		{
			name: "first example machine",
			line: "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			want: 2,
		},
		{
			name: "second example machine",
			line: "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			want: 3,
		},
		{
			name: "third example machine",
			line: "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine, err := ParseMachine(tt.line)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			got := SolveMinPresses(machine)
			if got != tt.want {
				t.Errorf("SolveMinPresses() = %d, want %d", got, tt.want)
			}
		})
	}
}
