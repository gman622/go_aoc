package day10

import "fmt"

// Part2 solves Day 10 Part 2 - finds minimum button presses to achieve joltage levels.
//
// Algorithm:
// 1. Parse each machine line (ignore lights, focus on joltages and button wirings)
// 2. For each machine, solve integer linear programming problem:
//    - Each button increments certain counters by 1
//    - Find minimum button presses to reach target joltage levels
// 3. Sum the minimum button presses across all machines
//
// This is an integer linear programming problem where we need non-negative
// integer solutions that minimize the objective function (total presses).
func Part2(inputPath string) (int, error) {
	lines, err := FromFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("loading input: %w", err)
	}

	totalPresses := 0

	for i, line := range lines {
		machine, err := ParseMachine(line)
		if err != nil {
			return 0, fmt.Errorf("parsing machine %d: %w", i+1, err)
		}

		presses := SolveMinJoltage(machine)
		if presses < 0 {
			return 0, fmt.Errorf("machine %d has no solution", i+1)
		}

		totalPresses += presses
	}

	return totalPresses, nil
}
