package solver

import (
	"fmt"
	"math"
)

func CheckSudoku(m *SudokuMatrix) (*Solver, error) {
	solver := Solver{}
	length := len((*m).Sudoku)
	fDim := math.Sqrt(float64(length))
	if math.Floor(fDim) != fDim {
		return &solver, fmt.Errorf("ERROR: Invalid sudoku matrix")
	}
	dim := int(fDim)
	solver.Length = length
	solver.Dim = int(dim)
	solver.Problem = *m

	for rowIndex, row := range (*m).Sudoku {

		if len(row) != length { // check whether the matrix is square
			return &solver, fmt.Errorf("ERROR: Sudoku matrix is not square")
		}

		for colIndex, search := range row {
			if search == 0 {
				continue
			}

			for c, colValue := range row { // search for duplicates in row
				if (c != colIndex) && (search == colValue) {
					return &solver, fmt.Errorf("ERROR: Same value in row %v", rowIndex)
				}
			}

			for r, rowValue := range m.Sudoku { // search for duplicates in cols
				if (r != rowIndex) && (search == rowValue[colIndex]) {
					return &solver, fmt.Errorf("ERROR: Same value in col %v", colIndex)
				}
			}

			startRow := (rowIndex / dim) * dim
			endRow := startRow + dim
			startCol := (colIndex / dim) * dim
			endCol := startCol + dim

			for r := startRow; r < endRow; r++ { // search for duplicates in block
				if r == rowIndex {
					continue
				}
				for c := startCol; c < endCol; c++ {
					if (c != colIndex) && (search == (*m).Sudoku[r][c]) {
						return &solver, fmt.Errorf("ERROR: Same value in block in pos [%v, %v]", r, c)
					}
				}
			}
		}
	}

	return &solver, nil
}
