package day10

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Machine represents a factory machine with lights and buttons.
type Machine struct {
	TargetLights []bool  // Target state for each light (true = on, false = off)
	Buttons      [][]int // Each button lists which lights it toggles
	Joltages     []int   // Target joltage levels for Part 2
}

// ParseMachine parses a single machine line.
// Format: [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
func ParseMachine(line string) (*Machine, error) {
	// Extract lights pattern [.##.]
	lightsRe := regexp.MustCompile(`\[([.#]+)\]`)
	lightsMatch := lightsRe.FindStringSubmatch(line)
	if len(lightsMatch) < 2 {
		return nil, fmt.Errorf("no lights pattern found in line: %s", line)
	}

	// Parse target lights
	lightsStr := lightsMatch[1]
	targetLights := make([]bool, len(lightsStr))
	for i, ch := range lightsStr {
		targetLights[i] = (ch == '#')
	}

	// Extract button patterns (1,3) (2) etc.
	buttonsRe := regexp.MustCompile(`\(([0-9,]+)\)`)
	buttonMatches := buttonsRe.FindAllStringSubmatch(line, -1)

	buttons := make([][]int, 0)
	for _, match := range buttonMatches {
		if len(match) < 2 {
			continue
		}
		numStrs := strings.Split(match[1], ",")
		button := make([]int, 0, len(numStrs))
		for _, numStr := range numStrs {
			num, err := strconv.Atoi(strings.TrimSpace(numStr))
			if err != nil {
				return nil, fmt.Errorf("parsing button number: %w", err)
			}
			button = append(button, num)
		}
		buttons = append(buttons, button)
	}

	// Extract joltage requirements {3,5,4,7}
	joltageRe := regexp.MustCompile(`\{([0-9,]+)\}`)
	joltageMatch := joltageRe.FindStringSubmatch(line)

	var joltages []int
	if len(joltageMatch) >= 2 {
		numStrs := strings.Split(joltageMatch[1], ",")
		joltages = make([]int, 0, len(numStrs))
		for _, numStr := range numStrs {
			num, err := strconv.Atoi(strings.TrimSpace(numStr))
			if err != nil {
				return nil, fmt.Errorf("parsing joltage number: %w", err)
			}
			joltages = append(joltages, num)
		}
	}

	return &Machine{
		TargetLights: targetLights,
		Buttons:      buttons,
		Joltages:     joltages,
	}, nil
}

// String implements fmt.Stringer for debugging.
func (m *Machine) String() string {
	lights := make([]rune, len(m.TargetLights))
	for i, on := range m.TargetLights {
		if on {
			lights[i] = '#'
		} else {
			lights[i] = '.'
		}
	}
	return fmt.Sprintf("[%s] %d buttons", string(lights), len(m.Buttons))
}
