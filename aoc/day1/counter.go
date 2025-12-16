package day1

// Counter defines a strategy for counting during dial rotations
type Counter interface {
	// Count processes a rotation and returns the count contribution
	Count(rotation Rotation, position int) int
}

// EndPositionCounter counts only when the dial ends at position 0
type EndPositionCounter struct{}

func (EndPositionCounter) Count(rotation Rotation, position int) int {
	if applyRotation(rotation, position) == 0 {
		return 1
	}
	return 0
}

// ZeroCrossingCounter counts every time the dial passes through 0
type ZeroCrossingCounter struct{}

func (ZeroCrossingCounter) Count(rotation Rotation, position int) int {
	return countZeroCrossings(rotation, position)
}

// Dial represents the safe's dial with a current position
type Dial struct {
	position int
	counter  Counter
	count    int
}

// NewDial creates a dial starting at position 50 with the given counter strategy
func NewDial(counter Counter) *Dial {
	return &Dial{
		position: 50,
		counter:  counter,
	}
}

// Rotate applies a rotation, updates the count, and returns the dial for chaining
func (d *Dial) Rotate(r Rotation) *Dial {
	d.count += d.counter.Count(r, d.position)
	d.position = applyRotation(r, d.position)
	return d
}

// Count returns the accumulated count
func (d *Dial) Count() int {
	return d.count
}

// applyRotation applies a rotation to a position and returns the new position (0-99)
func applyRotation(r Rotation, position int) int {
	var newPos int
	if r.Direction == 'L' {
		newPos = position - r.Distance
	} else { // 'R'
		newPos = position + r.Distance
	}

	// Normalize to 0-99 range
	newPos = newPos % 100
	if newPos < 0 {
		newPos += 100
	}
	return newPos
}

// countZeroCrossings counts how many times a rotation crosses position 0
func countZeroCrossings(r Rotation, from int) int {
	if r.Direction == 'L' {
		if from == 0 {
			// Starting at 0, count complete wraps
			return r.Distance / 100
		}
		// Going left from position p, we hit 0 after p steps
		if r.Distance >= from {
			return 1 + (r.Distance-from)/100
		}
		return 0
	} else { // 'R'
		// Going right, we cross 0 every 100 steps starting from (100 - position)
		return (from + r.Distance) / 100
	}
}
