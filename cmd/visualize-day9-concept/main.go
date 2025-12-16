package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sort"

	day9 "adv2025/aoc/day9"
)

func main() {
	tiles, err := day9.PointsFromFile("inputs/day9_input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Extract unique X coordinates
	xSet := make(map[int]bool)
	for _, p := range tiles {
		xSet[p.X] = true
	}
	xCoords := make([]int, 0, len(xSet))
	for x := range xSet {
		xCoords = append(xCoords, x)
	}
	sort.Ints(xCoords)

	fmt.Printf("Creating compression concept diagram...\n")
	fmt.Printf("Total unique X coordinates: %d\n", len(xCoords))
	fmt.Printf("Range: %d to %d\n", xCoords[0], xCoords[len(xCoords)-1])

	// Image dimensions
	width := 1400
	height := 600
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw two diagrams stacked vertically
	topY := 80
	bottomY := 350
	leftMargin := 50
	rightMargin := width - 50
	diagramWidth := rightMargin - leftMargin

	// Title
	drawText(img, "Coordinate Compression: Why It Works", width/2, 30, color.Black, true)

	// Top diagram: Original sparse space
	drawText(img, "Original Coordinate Space (Sparse)", width/2, topY-30, color.RGBA{100, 100, 100, 255}, true)

	// Draw axis
	axisY := topY + 100
	drawThickLine(img, leftMargin, axisY, rightMargin, axisY, color.Black, 2)

	// Show first ~20 coordinates to illustrate sparseness
	numToShow := min(20, len(xCoords))
	minX := xCoords[0]
	maxX := xCoords[numToShow-1]

	// Draw tick marks and values for original coordinates
	for i := 0; i < numToShow; i++ {
		x := xCoords[i]
		pixelX := leftMargin + int(float64(x-minX)/float64(maxX-minX)*float64(diagramWidth))

		// Tick mark
		drawLine(img, pixelX, axisY-10, pixelX, axisY+10, color.RGBA{200, 0, 0, 255})

		// Red dot for coordinate
		drawCircle(img, pixelX, axisY, 5, color.RGBA{255, 0, 0, 255})

		// Value label (every 3rd to avoid overlap)
		if i%3 == 0 {
			label := fmt.Sprintf("%d", x)
			drawText(img, label, pixelX, axisY+30, color.RGBA{100, 100, 100, 255}, true)
		}
	}

	// Highlight gaps with arrows
	drawText(img, "Huge gaps with no tiles!", leftMargin+diagramWidth/2, topY+40, color.RGBA{200, 100, 0, 255}, true)

	// Show gap between coordinates
	if numToShow >= 2 {
		x1 := xCoords[0]
		x2 := xCoords[1]
		gap := x2 - x1
		px1 := leftMargin + int(float64(x1-minX)/float64(maxX-minX)*float64(diagramWidth))
		px2 := leftMargin + int(float64(x2-minX)/float64(maxX-minX)*float64(diagramWidth))
		midX := (px1 + px2) / 2

		// Arrow pointing to gap
		drawLine(img, midX, topY+60, midX, axisY-30, color.RGBA{200, 100, 0, 255})
		drawText(img, fmt.Sprintf("Gap: %d unused", gap), midX, topY+55, color.RGBA{200, 100, 0, 255}, true)
	}

	// Bottom diagram: Compressed space
	drawText(img, "Compressed Space (Dense)", width/2, bottomY-30, color.RGBA{100, 100, 100, 255}, true)

	// Draw axis
	axisY2 := bottomY + 100
	drawThickLine(img, leftMargin, axisY2, rightMargin, axisY2, color.Black, 2)

	// Draw compressed coordinates (sequential)
	for i := 0; i < numToShow; i++ {
		// Even spacing in compressed space
		pixelX := leftMargin + int(float64(i)/float64(numToShow-1)*float64(diagramWidth))

		// Tick mark
		drawLine(img, pixelX, axisY2-10, pixelX, axisY2+10, color.RGBA{0, 150, 0, 255})

		// Green dot
		drawCircle(img, pixelX, axisY2, 5, color.RGBA{0, 200, 0, 255})

		// Index label (every 3rd)
		if i%3 == 0 {
			drawText(img, fmt.Sprintf("%d", i), pixelX, axisY2+30, color.RGBA{100, 100, 100, 255}, true)
		}
	}

	// Show compression benefit
	drawText(img, "No gaps! Every index is used!", leftMargin+diagramWidth/2, bottomY+40, color.RGBA{0, 150, 0, 255}, true)

	// Show mapping arrows between a few coordinates
	for i := 0; i < min(5, numToShow); i += 2 {
		x := xCoords[i]
		px1 := leftMargin + int(float64(x-minX)/float64(maxX-minX)*float64(diagramWidth))
		px2 := leftMargin + int(float64(i)/float64(numToShow-1)*float64(diagramWidth))

		// Dotted arrow
		drawDottedLine(img, px1, axisY+20, px2, axisY2-20, color.RGBA{100, 100, 200, 255})
	}

	// Summary text
	summaryY := height - 80
	drawText(img, fmt.Sprintf("Result: %d coordinates spanning 0-%d compressed to indices 0-%d",
		len(xCoords), xCoords[len(xCoords)-1], len(xCoords)-1),
		width/2, summaryY, color.RGBA{50, 50, 50, 255}, true)
	drawText(img, fmt.Sprintf("Same applies to Y axis: %d sparse coordinates → %d dense indices",
		len(xCoords), len(xCoords)),
		width/2, summaryY+25, color.RGBA{50, 50, 50, 255}, true)

	// Save image
	outFile, err := os.Create("aoc/day9/day9_compression_concept.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Saved concept diagram to aoc/day9/day9_compression_concept.png")
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func drawThickLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color, thickness int) {
	for i := -thickness / 2; i <= thickness/2; i++ {
		drawLine(img, x0, y0+i, x1, y1+i, c)
	}
}

func drawDottedLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	count := 0

	for {
		if count%10 < 5 {
			img.Set(x0, y0, c)
		}
		count++
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, radius int, c color.Color) {
	for x := cx - radius; x <= cx+radius; x++ {
		for y := cy - radius; y <= cy+radius; y++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, c)
			}
		}
	}
}

func drawText(img *image.RGBA, text string, x, y int, c color.Color, center bool) {
	// Simple text representation - draw a box where text would be
	textWidth := len(text) * 8
	startX := x
	if center {
		startX = x - textWidth/2
	}

	// Draw underline to represent text
	for i := 0; i < textWidth; i++ {
		img.Set(startX+i, y, c)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
