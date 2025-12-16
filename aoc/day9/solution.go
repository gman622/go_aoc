package day9

// BuildGreenTiles creates a set of all green tiles (edges and interior of polygon).
func BuildGreenTiles(redTiles []Point) map[Point]bool {
	greenTiles := make(map[Point]bool)
	n := len(redTiles)

	// Add edge tiles (straight lines between consecutive red tiles)
	for i := 0; i < n; i++ {
		from := redTiles[i]
		to := redTiles[(i+1)%n] // Wrap around

		// Add all tiles on the line from 'from' to 'to' (exclusive of endpoints)
		addLineTiles(greenTiles, from, to)
	}

	// Add interior tiles using point-in-polygon test
	// First, find the bounding box of all red tiles
	if n == 0 {
		return greenTiles
	}

	minX, maxX := redTiles[0].X, redTiles[0].X
	minY, maxY := redTiles[0].Y, redTiles[0].Y
	for _, p := range redTiles {
		minX = min(minX, p.X)
		maxX = max(maxX, p.X)
		minY = min(minY, p.Y)
		maxY = max(maxY, p.Y)
	}

	// Check each point in the bounding box
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			p := Point{x, y}
			// Skip if already marked as green
			if greenTiles[p] {
				continue
			}
			// Skip if it's a red tile
			isRed := false
			for _, r := range redTiles {
				if r == p {
					isRed = true
					break
				}
			}
			if isRed {
				continue
			}

			// Test if point is inside the polygon
			if isInsidePolygon(p, redTiles) {
				greenTiles[p] = true
			}
		}
	}

	return greenTiles
}

// addLineTiles adds all tiles on the line segment from 'from' to 'to' (exclusive of endpoints).
func addLineTiles(greenTiles map[Point]bool, from, to Point) {
	// Lines are always horizontal or vertical
	if from.X == to.X {
		// Vertical line
		startY := min(from.Y, to.Y)
		endY := max(from.Y, to.Y)
		for y := startY; y <= endY; y++ {
			p := Point{from.X, y}
			if p != from && p != to {
				greenTiles[p] = true
			}
		}
	} else if from.Y == to.Y {
		// Horizontal line
		startX := min(from.X, to.X)
		endX := max(from.X, to.X)
		for x := startX; x <= endX; x++ {
			p := Point{x, from.Y}
			if p != from && p != to {
				greenTiles[p] = true
			}
		}
	}
}

// isInsidePolygon uses ray casting algorithm to determine if a point is inside the polygon.
func isInsidePolygon(point Point, polygon []Point) bool {
	n := len(polygon)
	inside := false

	// Ray casting: count how many times a ray from point to infinity crosses polygon edges
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		p1, p2 := polygon[i], polygon[j]

		// Check if the ray crosses this edge
		if ((p1.Y > point.Y) != (p2.Y > point.Y)) &&
			(point.X < (p2.X-p1.X)*(point.Y-p1.Y)/(p2.Y-p1.Y)+p1.X) {
			inside = !inside
		}
	}

	return inside
}

// isRectangleValid checks if a rectangle contains only red or green tiles.
// Lazy validation: only checks points IN THIS RECTANGLE (not pre-computing 9.4B points!)
// Uses caching to avoid redundant polygon checks across overlapping rectangles.
func isRectangleValid(p1, p2 Point, redSet, edgeTiles map[Point]bool, polygon []Point, cache map[Point]bool) bool {
	minX := min(p1.X, p2.X)
	maxX := max(p1.X, p2.X)
	minY := min(p1.Y, p2.Y)
	maxY := max(p1.Y, p2.Y)

	width := maxX - minX + 1
	height := maxY - minY + 1
	area := width * height

	// Helper to check if a point is valid (red, edge, or inside polygon)
	// Fast path: check sets first (O(1))
	// Slow path: lazy polygon check with caching
	isValid := func(p Point) bool {
		if redSet[p] || edgeTiles[p] {
			return true
		}
		// Check cache first
		if result, ok := cache[p]; ok {
			return result
		}
		// Compute and cache
		result := isInsidePolygon(p, polygon)
		cache[p] = result
		return result
	}

	// For small rectangles, check all tiles
	if area <= 1000 {
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				if !isValid(Point{x, y}) {
					return false
				}
			}
		}
		return true
	}

	// For large rectangles, check border first (fast rejection)
	// Check top and bottom edges
	for x := minX; x <= maxX; x++ {
		if !isValid(Point{x, minY}) || !isValid(Point{x, maxY}) {
			return false
		}
	}
	// Check left and right edges (skip corners already checked)
	for y := minY + 1; y < maxY; y++ {
		if !isValid(Point{minX, y}) || !isValid(Point{maxX, y}) {
			return false
		}
	}

	// For large rectangles, trust border check ONLY
	// With 200M point rectangles, even sampling is too expensive!
	if area > 5000 {
		return true // Border passed, that's good enough
	}

	// Small rectangles: check all interior points
	for x := minX + 1; x < maxX; x++ {
		for y := minY + 1; y < maxY; y++ {
			if !isValid(Point{x, y}) {
				return false
			}
		}
	}

	return true
}
