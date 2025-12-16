package day10

import "fmt"

// Part1 solves Day 10 Part 1 - finds minimum button presses for all machines.
//
// Algorithm:
// 1. Parse each machine line (target lights, button wirings)
// 2. For each machine, solve system of linear equations over GF(2)
//    using Gaussian elimination
// 3. Sum the minimum button presses across all machines
//
// This is a classic linear algebra problem in the binary field:
// - Each button press toggles certain lights (XOR)
// - Pressing a button twice = not pressing (XOR is self-inverse)
// - Find minimal subset of buttons to press
func Part1(inputPath string) (int, error) {
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

		presses := SolveMinPresses(machine)
		if presses < 0 {
			return 0, fmt.Errorf("machine %d has no solution", i+1)
		}

		totalPresses += presses
	}

	return totalPresses, nil
}
