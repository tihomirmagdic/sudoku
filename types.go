package solver

type SudokuMatrix = struct {
	Sudoku [][]int
}

type Solver = struct {
	Problem    SudokuMatrix
	Candidates [][][]int
	Length     int
	Dim        int
}
