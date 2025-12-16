package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	day9 "adv2025/aoc/day9"
)

func main() {
	// Parse input
	tiles, err := day9.PointsFromFile("inputs/day9_input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Visualizing %d red tiles...\n", len(tiles))

	// Find bounds
	minX, maxX := tiles[0].X, tiles[0].X
	minY, maxY := tiles[0].Y, tiles[0].Y
	for _, p := range tiles {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	fmt.Printf("X range: %d to %d (span: %d)\n", minX, maxX, maxX-minX)
	fmt.Printf("Y range: %d to %d (span: %d)\n", minY, maxY, maxY-minY)

	// Count unique coordinates
	xSet := make(map[int]bool)
	ySet := make(map[int]bool)
	for _, p := range tiles {
		xSet[p.X] = true
		ySet[p.Y] = true
	}
	fmt.Printf("Unique X coords: %d\n", len(xSet))
	fmt.Printf("Unique Y coords: %d\n", len(ySet))

	// Create image with padding
	padding := 50
	width := 1200
	height := 1000

	// Calculate scale
	scaleX := float64(width-2*padding) / float64(maxX-minX)
	scaleY := float64(height-2*padding) / float64(maxY-minY)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Helper to transform coordinates
	toImageX := func(x int) int {
		return padding + int(float64(x-minX)*scale)
	}
	toImageY := func(y int) int {
		// Flip Y axis (image coordinates go down, our coordinates go up)
		return height - padding - int(float64(y-minY)*scale)
	}

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw polygon edges in green
	for i := range len(tiles) {
		from := tiles[i]
		to := tiles[(i+1)%len(tiles)]
		drawLine(img, toImageX(from.X), toImageY(from.Y), toImageX(to.X), toImageY(to.Y), color.RGBA{0, 200, 0, 255})
	}

	// Draw red tiles
	for _, p := range tiles {
		imgX := toImageX(p.X)
		imgY := toImageY(p.Y)
		drawCircle(img, imgX, imgY, 3, color.RGBA{255, 0, 0, 255})
	}

	// Save image
	outFile, err := os.Create("day9_visualization.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Saved visualization to day9_visualization.png")
}

// drawLine draws a line between two points using Bresenham's algorithm
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

// drawCircle draws a filled circle
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
