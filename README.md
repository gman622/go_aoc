# Advent of Code 2025 - Go Solutions

Clean Go implementations for Advent of Code 2025 problems.

## Structure

```
go_aoc/
├── cmd/
│   └── main.go          # Centralized runner
├── aoc/
│   ├── day1/            # Day 1 solution
│   ├── day2/            # Day 2 solution
│   └── ...
└── go.mod
```

## Running Solutions

```bash
# Run all solutions
go run cmd/main.go

# Run specific day
go run cmd/main.go -day 1

# Run specific part
go run cmd/main.go -day 1 -part 1
```

## Days Implemented

- Day 1-12: Various problems
- Day 23: Polynomial multiplication using FFT
- Day 24: Frequency analysis using FFT (Bell tuning)
- Day 25: Reactor core synchronization

Each day package exports a `Parts` slice containing solution functions.
