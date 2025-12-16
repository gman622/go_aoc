# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is an Advent of Code 2025 solution repository in Go. Solutions are organized by day, with a centralized runner that executes them.

## Running Solutions

```bash
# Run all implemented solutions
go run cmd/main.go

# Run all parts for a specific day
go run cmd/main.go -day 10

# Run a specific day and part
go run cmd/main.go -day 10 -part 1
```

## Testing

```bash
# Run all tests
go test ./...

# Test specific day
go test ./aoc/day10/...

# Run benchmarks (some days have performance tests)
go test -bench=. ./aoc/day9/...
```

## Architecture

### Day Package Structure

Each day is in `aoc/dayN/` with this standard pattern:

- `dayN.go` - Exports a `Parts` slice: `var Parts = []func(string) (int, error){Part1, Part2}`
- `part1.go`, `part2.go` - Solution implementations
- `parser.go` - Input parsing logic
- Additional files as needed (types.go, solution.go, etc.)

### Registration System

The main runner uses function slices for dynamic registration:

1. Each day package exports `Parts []func(string) (int, error)`
2. `cmd/main.go` registers days in `init()`: `register(10, day10.Parts...)`
3. To add a new day: create the package, implement Parts, and add one line to main.go's init

### Input Files

Solutions expect input files at `inputs/dayN_input.txt`. The runner:
- Passes the file path to each solution function
- Handles timing and error reporting
- Gracefully skips missing input files

### Solution Function Signature

All part functions must match: `func(string) (int, error)`
- Parameter: path to input file
- Returns: integer result and error
