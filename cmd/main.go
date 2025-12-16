package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	day1 "adv2025/aoc/day1"
	day2 "adv2025/aoc/day2"
	day3 "adv2025/aoc/day3"
	day4 "adv2025/aoc/day4"
	day5 "adv2025/aoc/day5"
	day6 "adv2025/aoc/day6"
	day7 "adv2025/aoc/day7"
	day8 "adv2025/aoc/day8"
	day9 "adv2025/aoc/day9"
	day10 "adv2025/aoc/day10"
	day11 "adv2025/aoc/day11"
	day12 "adv2025/aoc/day12"
	day25 "adv2025/aoc/day25"
)

type solver struct {
	day   int
	part  int
	solve func(string) (int, error)
}

var solvers []solver

func register(day int, parts ...func(string) (int, error)) {
	for i, part := range parts {
		solvers = append(solvers, solver{day, i + 1, part})
	}
}

func init() {
	register(1, day1.Parts...)
	register(2, day2.Parts...)
	register(3, day3.Parts...)
	register(4, day4.Parts...)
	register(5, day5.Parts...)
	register(6, day6.Parts...)
	register(7, day7.Parts...)
	register(8, day8.Parts...)
	register(9, day9.Parts...)
	register(10, day10.Parts...)
	register(11, day11.Parts...)
	register(12, day12.Parts...)
	register(25, day25.Parts...)
}

func main() {
	day := flag.Int("day", 0, "Day to run (0 for all)")
	part := flag.Int("part", 0, "Part to run (0 for all parts of the day)")
	flag.Parse()

	toRun := filterSolvers(*day, *part)
	if len(toRun) == 0 {
		log.Fatalf("No solutions found for day %d part %d", *day, *part)
	}

	printHeader()
	totalStart := time.Now()

	for _, s := range toRun {
		runSolver(s)
	}

	fmt.Printf("\n⏱️  Total time: %v\n", time.Since(totalStart))
}

func filterSolvers(day, part int) []solver {
	if day == 0 {
		return solvers
	}

	var filtered []solver
	for _, s := range solvers {
		if s.day == day && (part == 0 || s.part == part) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func runSolver(s solver) {
	inputPath := filepath.Join("inputs", fmt.Sprintf("day%d_input.txt", s.day))

	if _, err := os.Stat(inputPath); err != nil {
		fmt.Printf("❌ Day %d Part %d: Input file not found\n", s.day, s.part)
		return
	}

	start := time.Now()
	result, err := s.solve(inputPath)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Day %d Part %d: %v\n", s.day, s.part, err)
	} else {
		fmt.Printf("✅ Day %d Part %d: %d (%v)\n", s.day, s.part, result, elapsed)
	}
}

func printHeader() {
	fmt.Println("🎄 Advent of Code 2025 Runner")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
}
