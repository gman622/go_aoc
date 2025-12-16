package day1

import "fmt"

// Part1 solves part 1: count how many times the dial ends at position 0
func Part1(inputPath string) (int, error) {
	dial := NewDial(EndPositionCounter{})

	err := ProcessFile(inputPath, func(r Rotation) error {
		dial.Rotate(r)
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("processing rotations: %w", err)
	}

	return dial.Count(), nil
}
