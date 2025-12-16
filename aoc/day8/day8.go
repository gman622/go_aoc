package day8

// Parts contains all implemented parts for this day.
//
// This slice demonstrates Go's first-class function support - functions can be
// stored in slices, passed as parameters, and invoked dynamically. The main
// runner uses variadic arguments to register all parts at once:
//   register(8, day8.Parts...)
//
// Benefits of this pattern:
// - Each day package is self-describing (knows its own parts)
// - Adding Part3 requires only updating this slice, not main.go
// - Dynamic registration without reflection
// - Type-safe: compiler ensures all functions match the signature
var Parts = []func(string) (int, error){Part1, Part2}
