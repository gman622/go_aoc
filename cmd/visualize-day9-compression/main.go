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
	// Parse input
	tiles, err := day9.PointsFromFile("inputs/day9_input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Creating compression comparison visualization...\n")

	// Find bounds for original space
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

	// Extract unique coordinates for compression
	xSet := make(map[int]bool)
	ySet := make(map[int]bool)
	for _, p := range tiles {
		xSet[p.X] = true
		ySet[p.Y] = true
	}

	xCoords := make([]int, 0, len(xSet))
	for x := range xSet {
		xCoords = append(xCoords, x)
	}
	sort.Ints(xCoords)

	yCoords := make([]int, 0, len(ySet))
	for y := range ySet {
		yCoords = append(yCoords, y)
	}
	sort.Ints(yCoords)

	// Create coordinate mappings
	xToCompressed := make(map[int]int)
	for i, x := range xCoords {
		xToCompressed[x] = i
	}

	yToCompressed := make(map[int]int)
	for i, y := range yCoords {
		yToCompressed[y] = i
	}

	fmt.Printf("Original space: %d × %d\n", maxX-minX+1, maxY-minY+1)
	fmt.Printf("Compressed space: %d × %d\n", len(xCoords), len(yCoords))
	fmt.Printf("Compression ratio: %.1f:1\n", float64((maxX-minX+1)*(maxY-minY+1))/float64(len(xCoords)*len(yCoords)))

	// Image dimensions - make left panel much larger to show scale difference
	leftPanelWidth := 1000
	leftPanelHeight := 1000
	rightPanelWidth := 250  // Much smaller to show compression!
	rightPanelHeight := 250
	padding := 50
	gap := 80
	totalWidth := leftPanelWidth + gap + rightPanelWidth + 2*padding
	totalHeight := max(leftPanelHeight, rightPanelHeight) + 2*padding

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{250, 250, 250, 255}}, image.Point{}, draw.Src)

	// Draw left panel (original space)
	drawOriginalSpace(img, tiles, minX, maxX, minY, maxY, padding, padding, leftPanelWidth, leftPanelHeight)

	// Draw right panel (compressed space) - positioned to align with top of left panel
	rightPanelX := padding + leftPanelWidth + gap
	rightPanelY := padding + (leftPanelHeight-rightPanelHeight)/2 // Center vertically
	drawCompressedSpace(img, tiles, xToCompressed, yToCompressed, len(xCoords), len(yCoords),
		rightPanelX, rightPanelY, rightPanelWidth, rightPanelHeight)

	// Add titles
	drawTitle(img, "Original Space", padding+leftPanelWidth/2, padding-30)
	drawTitle(img, fmt.Sprintf("%d × %d = %.1fB cells", maxX-minX+1, maxY-minY+1,
		float64((maxX-minX+1)*(maxY-minY+1))/1e9), padding+leftPanelWidth/2, padding-10)

	drawTitle(img, "Compressed Space", rightPanelX+rightPanelWidth/2, rightPanelY-30)
	drawTitle(img, fmt.Sprintf("%d × %d = %dK cells", len(xCoords), len(yCoords),
		len(xCoords)*len(yCoords)/1000), rightPanelX+rightPanelWidth/2, rightPanelY-10)

	// Save image
	outFile, err := os.Create("aoc/day9/day9_compression.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Saved compression visualization to aoc/day9/day9_compression.png")
}

func drawOriginalSpace(img *image.RGBA, tiles []day9.Point, minX, maxX, minY, maxY, offsetX, offsetY, width, height int) {
	// Draw background box
	drawBox(img, offsetX, offsetY, width, height, color.RGBA{255, 255, 255, 255})

	// Calculate scale
	scaleX := float64(width-20) / float64(maxX-minX)
	scaleY := float64(height-20) / float64(maxY-minY)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Helper to transform coordinates
	toImageX := func(x int) int {
		return offsetX + 10 + int(float64(x-minX)*scale)
	}
	toImageY := func(y int) int {
		return offsetY + height - 10 - int(float64(y-minY)*scale)
	}

	// Draw sparse grid indicators (showing it's mostly empty)
	gridColor := color.RGBA{240, 240, 240, 255}
	for i := 0; i < 10; i++ {
		x := minX + (maxX-minX)*i/10
		imgX := toImageX(x)
		drawLine(img, imgX, offsetY, imgX, offsetY+height, gridColor)
	}
	for i := 0; i < 10; i++ {
		y := minY + (maxY-minY)*i/10
		imgY := toImageY(y)
		drawLine(img, offsetX, imgY, offsetX+width, imgY, gridColor)
	}

	// Draw polygon edges
	for i := range len(tiles) {
		from := tiles[i]
		to := tiles[(i+1)%len(tiles)]
		drawLine(img, toImageX(from.X), toImageY(from.Y), toImageX(to.X), toImageY(to.Y),
			color.RGBA{0, 200, 0, 255})
	}

	// Draw red tiles
	for _, p := range tiles {
		imgX := toImageX(p.X)
		imgY := toImageY(p.Y)
		drawCircle(img, imgX, imgY, 2, color.RGBA{255, 0, 0, 255})
	}
}

func drawCompressedSpace(img *image.RGBA, tiles []day9.Point, xToComp, yToComp map[int]int,
	compWidth, compHeight, offsetX, offsetY, width, height int) {
	// Draw background box
	drawBox(img, offsetX, offsetY, width, height, color.RGBA{255, 255, 255, 255})

	// Calculate scale for compressed space
	cellWidth := float64(width-20) / float64(compWidth)
	cellHeight := float64(height-20) / float64(compHeight)
	cellSize := cellWidth
	if cellHeight < cellWidth {
		cellSize = cellHeight
	}

	// Helper to transform compressed coordinates
	toImageX := func(x int) int {
		return offsetX + 10 + int(float64(x)*cellSize)
	}
	toImageY := func(y int) int {
		return offsetY + height - 10 - int(float64(y)*cellSize)
	}

	// Draw dense grid (showing all coordinates are used)
	gridColor := color.RGBA{220, 220, 220, 255}
	step := max(1, compWidth/20)
	for i := 0; i < compWidth; i += step {
		imgX := toImageX(i)
		drawLine(img, imgX, offsetY, imgX, offsetY+height, gridColor)
	}
	step = max(1, compHeight/20)
	for i := 0; i < compHeight; i += step {
		imgY := toImageY(i)
		drawLine(img, offsetX, imgY, offsetX+width, imgY, gridColor)
	}

	// Draw polygon edges in compressed space
	for i := range len(tiles) {
		from := tiles[i]
		to := tiles[(i+1)%len(tiles)]
		compFromX := xToComp[from.X]
		compFromY := yToComp[from.Y]
		compToX := xToComp[to.X]
		compToY := yToComp[to.Y]
		drawLine(img, toImageX(compFromX), toImageY(compFromY), toImageX(compToX), toImageY(compToY),
			color.RGBA{0, 200, 0, 255})
	}

	// Draw red tiles in compressed space
	for _, p := range tiles {
		compX := xToComp[p.X]
		compY := yToComp[p.Y]
		imgX := toImageX(compX)
		imgY := toImageY(compY)
		drawCircle(img, imgX, imgY, 2, color.RGBA{255, 0, 0, 255})
	}
}

func drawBox(img *image.RGBA, x, y, width, height int, c color.Color) {
	// Fill background
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			img.Set(i, j, c)
		}
	}
	// Draw border
	borderColor := color.RGBA{200, 200, 200, 255}
	drawLine(img, x, y, x+width, y, borderColor)
	drawLine(img, x+width, y, x+width, y+height, borderColor)
	drawLine(img, x+width, y+height, x, y+height, borderColor)
	drawLine(img, x, y+height, x, y, borderColor)
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

func drawTitle(img *image.RGBA, text string, centerX, y int) {
	// Simple text rendering by drawing text-like patterns
	// For a real implementation, you'd use a font library
	// For now, just draw a marker to indicate where text would go
	textColor := color.RGBA{50, 50, 50, 255}
	// Draw a simple underline to represent text
	textWidth := len(text) * 6
	startX := centerX - textWidth/2
	drawLine(img, startX, y+2, startX+textWidth, y+2, textColor)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
