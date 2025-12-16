package day10

import (
	"testing"
)

func TestMachine31(t *testing.T) {
	line := `[#...#....#] (1,2,4,7,8,9) (0,2,4,6) (0,1,2,4,5,9) (4,6) (0,2,3,5,6,7,8,9) (0,1,4,5,6,8,9) (6,7,8) (0,3,5,6,9) (0,1,3,4,5,6,8,9) (1,4,6,7,9) (9) {94,221,64,45,251,78,266,206,64,275}`

	machine, err := ParseMachine(line)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	t.Logf("Machine 31:")
	t.Logf("  Counters: %d", len(machine.Joltages))
	t.Logf("  Buttons: %d", len(machine.Buttons))
	t.Logf("  Joltages: %v", machine.Joltages)

	// Try each solver individually
	numCounters := len(machine.Joltages)
	numButtons := len(machine.Buttons)

	matrix := make([][]int, numCounters)
	for i := range matrix {
		matrix[i] = make([]int, numButtons)
	}
	for buttonIdx, button := range machine.Buttons {
		for _, counterIdx := range button {
			if counterIdx < numCounters {
				matrix[counterIdx][buttonIdx] = 1
			}
		}
	}

	t.Logf("Testing greedy solver...")
	greedy := solveGreedySimple(matrix, machine.Joltages, numButtons, numCounters)
	t.Logf("  Greedy: %d", greedy)

	t.Logf("Testing integer Gaussian...")
	gaussian := solveJoltageGreedy(matrix, machine.Joltages, numButtons, numCounters)
	t.Logf("  Integer Gaussian: %d", gaussian)

	t.Logf("Testing float Gaussian with partial pivoting...")
	floatResult := solveWithFloatGaussian(matrix, machine.Joltages, numButtons, numCounters)
	t.Logf("  Float Gaussian: %d", floatResult)

	result := SolveMinJoltage(machine)
	t.Logf("  Final Result: %d", result)

	if result < 0 {
		t.Errorf("No solution found, but there should be one")
	}
}
