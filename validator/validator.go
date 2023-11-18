package validator

import (
	"fmt"
	"math"

	"github.com/tihomirmagdic/sudoku/types"
)

func CheckValue(s *types.Solver, row int, col int) bool {
	m := &s.Problem
	dim := s.Dim
	search := (*m).Sudoku[row][col]

	for c, colValue := range (*m).Sudoku[row] { // search for duplicates in row
		if c == col {
			continue
		}
		if search == colValue {
			return false
		}
	}

	for r, rowValue := range m.Sudoku { // search for duplicates in cols
		if r == row {
			continue
		}
		if search == rowValue[col] {
			return false
		}
	}

	startRow := (row / dim) * dim
	endRow := ((row / dim) + 1) * dim
	startCol := (col / dim) * dim
	endCol := ((col / dim) + 1) * dim

	for r := startRow; r < endRow; r++ { // search for duplicates in block
		if r == row {
			continue
		}
		for c := startCol; c < endCol; c++ {
			if c == col {
				continue
			}
			if search == (*m).Sudoku[r][c] {
				return false
			}
		}
	}

	return true
}

func CheckSudoku(m *types.SudokuMatrix) (*types.Solver, error) {
	solver := types.Solver{}
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
				if c == colIndex {
					continue
				}
				if search == colValue {
					return &solver, fmt.Errorf("ERROR: Same value in row %v", rowIndex)
				}
			}

			for r, rowValue := range m.Sudoku { // search for duplicates in cols
				if r == rowIndex {
					continue
				}
				if search == rowValue[colIndex] {
					return &solver, fmt.Errorf("ERROR: Same value in col %v", colIndex)
				}
			}

			startRow := (rowIndex / dim) * dim
			endRow := ((rowIndex / dim) + 1) * dim
			startCol := (colIndex / dim) * dim
			endCol := ((colIndex / dim) + 1) * dim

			for r := startRow; r < endRow; r++ { // search for duplicates in block
				if r == rowIndex {
					continue
				}
				for c := startCol; c < endCol; c++ {
					if c == colIndex {
						continue
					}
					if search == (*m).Sudoku[r][c] {
						return &solver, fmt.Errorf("ERROR: Same value in block in pos [%v, %v]", r, c)
					}
				}
			}
		}
	}

	return &solver, nil
}
