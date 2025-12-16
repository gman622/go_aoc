package day1

import "fmt"

// Part2 solves part 2: count how many times the dial passes through position 0
func Part2(inputPath string) (int, error) {
	dial := NewDial(ZeroCrossingCounter{})

	err := ProcessFile(inputPath, func(r Rotation) error {
		dial.Rotate(r)
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("processing rotations: %w", err)
	}

	return dial.Count(), nil
}
