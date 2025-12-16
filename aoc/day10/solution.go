package day10

// SolveMinPresses finds the minimum number of button presses to configure a machine.
//
// This is a system of linear equations over GF(2) (binary field):
// - Each light is a variable (on=1, off=0)
// - Each button toggles certain lights (XOR operation)
// - We need to find which buttons to press (and how many times)
//
// Key insight: In GF(2), pressing a button twice is the same as not pressing it
// (XOR is self-inverse). So we only need to decide: press or don't press.
//
// Algorithm: Gaussian elimination over GF(2) + enumeration of free variables
// to find the solution with minimum hamming weight (fewest 1s).
func SolveMinPresses(m *Machine) int {
	numLights := len(m.TargetLights)
	numButtons := len(m.Buttons)

	if numButtons == 0 {
		// Check if all lights are already in target state (all off initially)
		for _, target := range m.TargetLights {
			if target {
				return -1 // Impossible: need lights on but no buttons
			}
		}
		return 0
	}

	// Build augmented matrix [A|b] where:
	// - A[i][j] = 1 if button j toggles light i
	// - b[i] = target state of light i
	matrix := make([][]int, numLights)
	for i := range matrix {
		matrix[i] = make([]int, numButtons+1)
		// Set target state in augmented column
		if m.TargetLights[i] {
			matrix[i][numButtons] = 1
		}
	}

	// Fill in button toggle information
	for j, button := range m.Buttons {
		for _, lightIdx := range button {
			if lightIdx < numLights {
				matrix[lightIdx][j] = 1
			}
		}
	}

	// Gaussian elimination to find all solutions
	minPresses := solveMinimalGF2(matrix, numButtons)

	return minPresses
}

// solveMinimalGF2 finds the solution with minimum number of 1s (minimum button presses).
// It performs Gaussian elimination, then enumerates all combinations of free variables
// to find the one that minimizes the hamming weight.
func solveMinimalGF2(matrix [][]int, numVars int) int {
	rows := len(matrix)
	if rows == 0 {
		return 0
	}

	// Make a copy of matrix since we'll modify it
	mat := make([][]int, rows)
	for i := range matrix {
		mat[i] = make([]int, len(matrix[i]))
		copy(mat[i], matrix[i])
	}

	cols := len(mat[0]) - 1 // Exclude augmented column

	// Forward elimination to reduced row echelon form (RREF)
	pivotRow := 0
	pivotCols := make([]int, 0)     // Columns that have pivots (basic variables)
	freeVars := make(map[int]bool)  // Columns without pivots (free variables)

	for col := 0; col < cols && pivotRow < rows; col++ {
		// Find pivot
		foundPivot := false
		for row := pivotRow; row < rows; row++ {
			if mat[row][col] == 1 {
				mat[pivotRow], mat[row] = mat[row], mat[pivotRow]
				foundPivot = true
				break
			}
		}

		if !foundPivot {
			freeVars[col] = true
			continue
		}

		pivotCols = append(pivotCols, col)

		// Eliminate all other 1s in this column
		for row := 0; row < rows; row++ {
			if row != pivotRow && mat[row][col] == 1 {
				for c := 0; c < len(mat[row]); c++ {
					mat[row][c] ^= mat[pivotRow][c]
				}
			}
		}

		pivotRow++
	}

	// Mark any remaining columns as free variables
	for col := 0; col < cols; col++ {
		isPivot := false
		for _, pc := range pivotCols {
			if pc == col {
				isPivot = true
				break
			}
		}
		if !isPivot {
			freeVars[col] = true
		}
	}

	// Check for inconsistency
	for row := 0; row < rows; row++ {
		allZero := true
		for col := 0; col < cols; col++ {
			if mat[row][col] == 1 {
				allZero = false
				break
			}
		}
		if allZero && mat[row][cols] == 1 {
			return -1 // No solution
		}
	}

	// Convert free variables to slice for enumeration
	freeVarList := make([]int, 0, len(freeVars))
	for v := range freeVars {
		freeVarList = append(freeVarList, v)
	}

	// If no free variables, there's only one solution
	if len(freeVarList) == 0 {
		solution := make([]int, numVars)
		for row := 0; row < len(pivotCols) && row < rows; row++ {
			solution[pivotCols[row]] = mat[row][cols]
		}
		count := 0
		for _, v := range solution {
			count += v
		}
		return count
	}

	// Enumerate all 2^k combinations of free variables
	minPresses := numVars + 1
	numCombinations := 1 << len(freeVarList)

	for combo := 0; combo < numCombinations; combo++ {
		solution := make([]int, numVars)

		// Set free variables according to this combination
		for i, varIdx := range freeVarList {
			if (combo>>i)&1 == 1 {
				solution[varIdx] = 1
			}
		}

		// Compute pivot variables based on free variables
		for row := 0; row < len(pivotCols) && row < rows; row++ {
			pivotCol := pivotCols[row]
			value := mat[row][cols]

			// XOR with contributions from free variables
			for col := 0; col < cols; col++ {
				if mat[row][col] == 1 && col != pivotCol {
					value ^= solution[col]
				}
			}

			solution[pivotCol] = value
		}

		// Count button presses for this solution
		count := 0
		for _, v := range solution {
			count += v
		}

		if count < minPresses {
			minPresses = count
		}
	}

	return minPresses
}

// SolveMinJoltage finds the minimum number of button presses to achieve target joltages.
//
// This is an integer linear programming problem:
// - Each button increments certain counters
// - We need each counter to reach its target value
// - Minimize total button presses
//
// Algorithm: This is a system of linear equations over integers (not GF(2)):
// For each counter i: sum of (button_j presses * (1 if button j affects counter i)) = target[i]
//
// We can solve this using Gaussian elimination over rationals, then check if the
// solution is non-negative integers. If there are free variables, we need to find
// the minimal solution.
func SolveMinJoltage(m *Machine) int {
	numCounters := len(m.Joltages)
	numButtons := len(m.Buttons)

	if numButtons == 0 {
		// Check if all joltages are already at target (all 0)
		for _, target := range m.Joltages {
			if target != 0 {
				return -1 // Impossible
			}
		}
		return 0
	}

	// Build matrix [A|b] where:
	// - A[i][j] = 1 if button j affects counter i, 0 otherwise
	// - b[i] = target joltage for counter i
	//
	// We want to solve: A * x = b where x >= 0 and minimize sum(x)

	// Since buttons only add 1 to counters they affect, we can use a greedy approach
	// or solve as an integer linear program. For simplicity, we'll use Gaussian elimination
	// and check if the solution is valid (non-negative integers).

	return solveMinJoltageInteger(m.Buttons, m.Joltages, numButtons, numCounters)
}

// solveMinJoltageInteger solves the integer linear programming problem.
// This is more complex than the GF(2) case because we need non-negative integer solutions.
func solveMinJoltageInteger(buttons [][]int, targets []int, numButtons, numCounters int) int {
	// Build coefficient matrix
	// matrix[counter][button] = 1 if button affects counter
	matrix := make([][]int, numCounters)
	for i := range matrix {
		matrix[i] = make([]int, numButtons)
	}

	for buttonIdx, button := range buttons {
		for _, counterIdx := range button {
			if counterIdx < numCounters {
				matrix[counterIdx][buttonIdx] = 1
			}
		}
	}

	// Day 9 lesson: Try simple approaches first before complex algorithms
	// For this problem, we can use a greedy backtracking search

	// First try: simple greedy approach (works for many cases)
	greedy := solveGreedySimple(matrix, targets, numButtons, numCounters)
	if greedy >= 0 {
		return greedy
	}

	// If greedy fails, try Gaussian elimination
	gaussian := solveJoltageGreedy(matrix, targets, numButtons, numCounters)
	if gaussian >= 0 {
		return gaussian
	}

	// Last resort: use floating point Gaussian elimination to find solution structure
	// This avoids integer coefficient explosion
	return solveWithFloatGaussian(matrix, targets, numButtons, numCounters)
}

// solveWithFloatGaussian uses floating point arithmetic to solve the system,
// avoiding the integer coefficient explosion problem.
func solveWithFloatGaussian(matrix [][]int, targets []int, numButtons, numCounters int) int {
	// Convert to float64 for Gaussian elimination
	aug := make([][]float64, numCounters)
	for i := range aug {
		aug[i] = make([]float64, numButtons+1)
		for j := 0; j < numButtons; j++ {
			aug[i][j] = float64(matrix[i][j])
		}
		aug[i][numButtons] = float64(targets[i])
	}

	// Gaussian elimination with partial pivoting to RREF
	pivotCols := make([]int, 0)
	currentRow := 0

	for col := 0; col < numButtons && currentRow < numCounters; col++ {
		// Find best pivot (largest absolute value) in this column
		maxRow := -1
		maxVal := 0.0
		for row := currentRow; row < numCounters; row++ {
			if abs64(aug[row][col]) > maxVal {
				maxVal = abs64(aug[row][col])
				maxRow = row
			}
		}

		// Skip if no good pivot found
		if maxVal < 1e-10 {
			continue
		}

		// Swap rows
		aug[currentRow], aug[maxRow] = aug[maxRow], aug[currentRow]
		pivotCols = append(pivotCols, col)

		// Scale pivot row to make pivot = 1
		pivot := aug[currentRow][col]
		for c := 0; c <= numButtons; c++ {
			aug[currentRow][c] /= pivot
		}

		// Eliminate all other rows (forward and backward)
		for row := 0; row < numCounters; row++ {
			if row == currentRow {
				continue
			}
			if abs64(aug[row][col]) < 1e-10 {
				continue
			}
			factor := aug[row][col]
			for c := 0; c <= numButtons; c++ {
				aug[row][c] -= factor * aug[currentRow][c]
			}
		}

		currentRow++
	}

	// Now we're in RREF. Identify free variables
	freeVarList := make([]int, 0)
	for i := 0; i < numButtons; i++ {
		isFree := true
		for _, pivotCol := range pivotCols {
			if pivotCol == i {
				isFree = false
				break
			}
		}
		if isFree {
			freeVarList = append(freeVarList, i)
		}
	}


	// Enumerate free variable assignments to find minimum
	// Find maximum target to bound search
	maxTarget := 0
	for _, t := range targets {
		if t > maxTarget {
			maxTarget = t
		}
	}

	// Limit free variable search based on count
	if len(freeVarList) > 8 {
		// Too many free variables - use greedy with free vars = 0
		return solveWithFreeVarsZero(aug, pivotCols, freeVarList, matrix, targets, numButtons, numCounters)
	}

	minTotal := -1
	// Need to allow larger free variable values to avoid negative pivot variables
	maxFreeVal := maxTarget // Try up to maxTarget

	var enumerate func(int, []float64)
	enumerate = func(idx int, currentSol []float64) {
		if idx == len(freeVarList) {
			// Compute pivot variables from RREF
			sol := make([]float64, numButtons)
			copy(sol, currentSol)

			// Read pivot values from RREF (accounting for free variables)
			for rowIdx, pivotCol := range pivotCols {
				val := aug[rowIdx][numButtons]
				// Subtract contributions from free variables
				for freeVar := range freeVarList {
					freeCol := freeVarList[freeVar]
					val -= aug[rowIdx][freeCol] * sol[freeCol]
				}
				sol[pivotCol] = val
			}

			// Round and verify
			intSol := make([]int, numButtons)
			for i := range sol {
				rounded := int(sol[i] + 0.5)
				if rounded < 0 {
					return // Invalid - negative button presses not allowed
				}
				intSol[i] = rounded
			}

			// Verify against original constraints
			valid := true
			for row := 0; row < numCounters; row++ {
				sum := 0
				for col := 0; col < numButtons; col++ {
					sum += matrix[row][col] * intSol[col]
				}
				if sum != targets[row] {
					valid = false
					break
				}
			}

			if valid {
				total := 0
				for _, presses := range intSol {
					total += presses
				}
				if minTotal < 0 || total < minTotal {
					minTotal = total
				}
			}
			return
		}

		// Try values for this free variable
		freeVar := freeVarList[idx]
		for val := 0; val <= maxFreeVal; val++ {
			newSol := make([]float64, numButtons)
			copy(newSol, currentSol)
			newSol[freeVar] = float64(val)
			enumerate(idx+1, newSol)
		}
	}

	enumerate(0, make([]float64, numButtons))
	return minTotal
}

func solveWithFreeVarsZero(aug [][]float64, pivotCols, freeVarList []int, matrix [][]int, targets []int, numButtons, numCounters int) int {
	// Simple case: set all free variables to 0
	solution := make([]float64, numButtons)
	for rowIdx, col := range pivotCols {
		if rowIdx < len(aug) {
			solution[col] = aug[rowIdx][numButtons]
		}
	}

	// Round to integers
	intSolution := make([]int, numButtons)
	for i := range solution {
		rounded := int(solution[i] + 0.5)
		if rounded < 0 {
			return -1
		}
		intSolution[i] = rounded
	}

	// Verify
	for row := 0; row < numCounters; row++ {
		sum := 0
		for col := 0; col < numButtons; col++ {
			sum += matrix[row][col] * intSolution[col]
		}
		if sum != targets[row] {
			return -1
		}
	}

	total := 0
	for _, presses := range intSolution {
		total += presses
	}
	return total
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// tryRoundingStrategies tries different rounding approaches when direct rounding fails.
func tryRoundingStrategies(solution []float64, matrix [][]int, targets []int, numButtons, numCounters int) int {
	// Try all combinations of floor/ceil for values close to 0.5
	closeToHalf := make([]int, 0)
	for i, v := range solution {
		frac := v - float64(int(v))
		if frac > 0.3 && frac < 0.7 {
			closeToHalf = append(closeToHalf, i)
		}
	}

	if len(closeToHalf) > 10 {
		return -1 // Too many ambiguous values
	}

	// Try all 2^n combinations
	numCombos := 1 << len(closeToHalf)
	for combo := 0; combo < numCombos; combo++ {
		intSolution := make([]int, numButtons)
		for i := range solution {
			intSolution[i] = int(solution[i] + 0.5)
		}

		// Apply this combination
		for i, idx := range closeToHalf {
			if (combo>>i)&1 == 1 {
				intSolution[idx] = int(solution[idx] + 1.0) // Ceiling
			} else {
				intSolution[idx] = int(solution[idx]) // Floor
			}
		}

		// Check if negative
		valid := true
		for _, v := range intSolution {
			if v < 0 {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		// Verify
		matches := true
		for row := 0; row < numCounters; row++ {
			sum := 0
			for col := 0; col < numButtons; col++ {
				sum += matrix[row][col] * intSolution[col]
			}
			if sum != targets[row] {
				matches = false
				break
			}
		}

		if matches {
			total := 0
			for _, presses := range intSolution {
				total += presses
			}
			return total
		}
	}

	return -1
}

// solveDFS uses depth-first search with pruning to find a solution.
// This is a last resort for cases where Gaussian elimination fails due to coefficient issues.
func solveDFS(matrix [][]int, targets []int, numButtons, numCounters int) int {
	// Use BFS/iterative deepening to find minimum solution
	for maxPresses := 0; maxPresses <= 1000; maxPresses++ {
		solution := make([]int, numButtons)
		if dfSearch(matrix, targets, numButtons, numCounters, 0, 0, maxPresses, solution) {
			total := 0
			for _, p := range solution {
				total += p
			}
			return total
		}
	}
	return -1
}

// dfSearch performs DFS with a budget constraint.
func dfSearch(matrix [][]int, targets []int, numButtons, numCounters, btnIdx, pressesUsed, maxPresses int, solution []int) bool {
	if btnIdx == numButtons {
		// Check if we've met all targets
		counters := make([]int, numCounters)
		for btn := 0; btn < numButtons; btn++ {
			for ctr := 0; ctr < numCounters; ctr++ {
				if matrix[ctr][btn] == 1 {
					counters[ctr] += solution[btn]
				}
			}
		}

		for i := 0; i < numCounters; i++ {
			if counters[i] != targets[i] {
				return false
			}
		}
		return true
	}

	// Try different numbers of presses for this button
	maxForButton := minInt(maxPresses-pressesUsed, 500)

	for presses := 0; presses <= maxForButton; presses++ {
		solution[btnIdx] = presses
		if dfSearch(matrix, targets, numButtons, numCounters, btnIdx+1, pressesUsed+presses, maxPresses, solution) {
			return true
		}
	}

	solution[btnIdx] = 0
	return false
}

// solveGreedySimple uses a simple greedy strategy: for each counter in order,
// press buttons to reach the target.
func solveGreedySimple(matrix [][]int, targets []int, numButtons, numCounters int) int {
	counters := make([]int, numCounters)
	presses := make([]int, numButtons)

	// For each counter, find buttons that ONLY affect this counter (or few others)
	// and press them to meet the target
	for ctr := 0; ctr < numCounters; ctr++ {
		if counters[ctr] >= targets[ctr] {
			continue // Already met
		}

		// Find the best button to press for this counter
		bestBtn := -1
		bestScore := -1

		for btn := 0; btn < numButtons; btn++ {
			if matrix[ctr][btn] == 0 {
				continue // Doesn't affect this counter
			}

			// Score = how many counters this button affects that still need incrementing
			score := 0
			willOvershoot := false

			for c := 0; c < numCounters; c++ {
				if matrix[c][btn] == 1 {
					if counters[c] < targets[c] {
						score++
					} else if counters[c] >= targets[c] {
						// Would overshoot a counter that's already done
						willOvershoot = true
						break
					}
				}
			}

			if !willOvershoot && score > bestScore {
				bestScore = score
				bestBtn = btn
			}
		}

		if bestBtn == -1 {
			return -1 // Can't make progress
		}

		// Press this button enough times to meet the current counter's target
		needed := targets[ctr] - counters[ctr]
		presses[bestBtn] += needed

		for c := 0; c < numCounters; c++ {
			if matrix[c][bestBtn] == 1 {
				counters[c] += needed
			}
		}
	}

	// Verify solution
	for i := 0; i < numCounters; i++ {
		if counters[i] != targets[i] {
			return -1
		}
	}

	total := 0
	for _, p := range presses {
		total += p
	}
	return total
}

// solveJoltageGreedy solves using Gaussian elimination over integers.
// Unlike GF(2), we work with actual integers and find non-negative solutions.
func solveJoltageGreedy(matrix [][]int, targets []int, numButtons, numCounters int) int {
	// Create augmented matrix for Gaussian elimination
	// We'll work with the transpose: each row represents a button, each column a counter
	// Actually, let's keep it as: rows = counters, cols = buttons

	// Make a copy
	aug := make([][]int, numCounters)
	for i := range aug {
		aug[i] = make([]int, numButtons+1)
		copy(aug[i][:numButtons], matrix[i])
		aug[i][numButtons] = targets[i]
	}

	// Gaussian elimination to find RREF
	pivotRow := 0
	pivotCols := []int{}

	for col := 0; col < numButtons && pivotRow < numCounters; col++ {
		// Find a non-zero pivot
		found := false
		for row := pivotRow; row < numCounters; row++ {
			if aug[row][col] != 0 {
				// Swap
				aug[pivotRow], aug[row] = aug[row], aug[pivotRow]
				found = true
				break
			}
		}

		if !found {
			continue // Free variable
		}

		pivotCols = append(pivotCols, col)

		// Make pivot = 1 and eliminate column (using integer arithmetic)
		pivotVal := aug[pivotRow][col]

		// Eliminate other rows
		for row := 0; row < numCounters; row++ {
			if row == pivotRow {
				continue
			}
			if aug[row][col] == 0 {
				continue
			}

			// Row operation: row -= (aug[row][col] / pivotVal) * pivotRow
			// To avoid fractions, use: row * pivotVal -= aug[row][col] * pivotRow
			multiplier := aug[row][col]
			for c := 0; c <= numButtons; c++ {
				aug[row][c] = aug[row][c]*pivotVal - multiplier*aug[pivotRow][c]
			}
		}

		pivotRow++
	}

	// Find free variables
	freeVars := make(map[int]bool)
	for col := 0; col < numButtons; col++ {
		isPivot := false
		for _, pc := range pivotCols {
			if pc == col {
				isPivot = true
				break
			}
		}
		if !isPivot {
			freeVars[col] = true
		}
	}

	// If there are free variables with reasonable range, enumerate them
	freeVarList := make([]int, 0, len(freeVars))
	for v := range freeVars {
		freeVarList = append(freeVarList, v)
	}

	// Determine the maximum reasonable value for free variables based on targets
	maxTarget := 0
	for _, t := range targets {
		if t > maxTarget {
			maxTarget = t
		}
	}

	// If no free variables, there's a unique solution (or no solution)
	if len(freeVarList) == 0 {
		return tryJoltageSolution(aug, pivotCols, numButtons, numCounters, make([]int, numButtons))
	}

	// If too many free variables, use a limited search with pruning
	// Day 9 lesson: Don't give up on large spaces - use smart search!
	if len(freeVarList) > 6 {
		// Use iterative deepening: try small values first, expand if needed
		return solveWithIterativeDeepening(aug, pivotCols, freeVarList, numButtons, numCounters, maxTarget)
	}

	// Enumerate combinations of free variable values (for small number of free vars)
	minTotal := -1

	// Use a reasonable upper bound for each free variable
	maxFreeVal := minInt(maxTarget, 30) // Don't try more than 30 for each

	var enumerate func(int, []int)
	enumerate = func(idx int, solution []int) {
		if idx == len(freeVarList) {
			// Try this assignment
			result := tryJoltageSolution(aug, pivotCols, numButtons, numCounters, solution)
			if result >= 0 && (minTotal < 0 || result < minTotal) {
				minTotal = result
			}
			return
		}

		// Try different values for this free variable
		freeVar := freeVarList[idx]
		for val := 0; val <= maxFreeVal; val++ {
			newSol := make([]int, numButtons)
			copy(newSol, solution)
			newSol[freeVar] = val
			enumerate(idx+1, newSol)
		}
	}

	enumerate(0, make([]int, numButtons))
	return minTotal
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// solveWithIterativeDeepening uses iterative deepening search to find minimum solution.
// Inspired by Day 9: don't enumerate billions of combinations - search smartly!
func solveWithIterativeDeepening(aug [][]int, pivotCols, freeVarList []int, numButtons, numCounters, maxTarget int) int {
	// Try increasing bounds on free variables until we find a solution
	// Start small and expand - most solutions have small free variable values
	maxBounds := []int{0, 2, 5, 10, 15, 20, 30, 50, 75, 100, 150, 200, maxTarget}

	for _, bound := range maxBounds {
		result := searchWithBound(aug, pivotCols, freeVarList, numButtons, numCounters, bound, 0, make([]int, numButtons))
		if result >= 0 {
			return result
		}
		if bound >= maxTarget {
			break // No point going higher
		}
	}

	return -1 // No solution found
}

// searchWithBound searches for solutions with free variables bounded by maxVal.
// Uses branch-and-bound pruning to avoid exploring bad paths.
func searchWithBound(aug [][]int, pivotCols, freeVarList []int, numButtons, numCounters, maxVal, idx int, solution []int) int {
	if idx == len(freeVarList) {
		// Try this assignment
		return tryJoltageSolution(aug, pivotCols, numButtons, numCounters, solution)
	}

	// Try different values for this free variable
	freeVar := freeVarList[idx]
	minResult := -1

	for val := 0; val <= maxVal; val++ {
		newSol := make([]int, numButtons)
		copy(newSol, solution)
		newSol[freeVar] = val

		result := searchWithBound(aug, pivotCols, freeVarList, numButtons, numCounters, maxVal, idx+1, newSol)
		if result >= 0 && (minResult < 0 || result < minResult) {
			minResult = result
		}

		// Pruning: if we found a solution with val=0, trying higher values won't help
		// (more free variable presses = more total presses)
		if idx == 0 && minResult >= 0 && val == 0 {
			break // Free variable 0 can be 0, that's optimal for this variable
		}
	}

	return minResult
}

// tryJoltageSolution tries a specific assignment of free variables and solves for pivot variables.
func tryJoltageSolution(aug [][]int, pivotCols []int, numButtons, numCounters int, solution []int) int {
	// Back-substitution with given free variable values
	for i := len(pivotCols) - 1; i >= 0; i-- {
		rowIdx := i
		if rowIdx >= numCounters {
			continue
		}
		colIdx := pivotCols[i]

		pivotVal := aug[rowIdx][colIdx]
		rhs := aug[rowIdx][numButtons]

		// Subtract contributions from all other variables
		for c := 0; c < numButtons; c++ {
			if c != colIdx {
				rhs -= aug[rowIdx][c] * solution[c]
			}
		}

		if rhs%pivotVal != 0 {
			return -1 // No integer solution
		}

		solution[colIdx] = rhs / pivotVal

		if solution[colIdx] < 0 {
			return -1 // Negative not allowed
		}
	}

	// Verify solution satisfies all constraints
	for row := 0; row < numCounters; row++ {
		sum := 0
		for col := 0; col < numButtons; col++ {
			sum += aug[row][col] * solution[col]
		}
		if sum != aug[row][numButtons] {
			return -1 // Solution doesn't satisfy constraints
		}
	}

	// Count total
	total := 0
	for _, presses := range solution {
		total += presses
	}
	return total
}
