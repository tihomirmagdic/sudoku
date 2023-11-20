package solver

import (
	"fmt"
)

func IsInRowC(m *[][][]int, row int, col int, search int) bool {
	for c, cRow := range (*m)[row] {
		if c == col {
			continue
		}
		for _, candidate := range cRow {
			if candidate > search {
				break
			}
			if search == candidate {
				return true
			}
		}
	}
	return false
}

func IsInRow(s *Solver, row int, search int) bool {
	for _, value := range s.Problem.Sudoku[row] {
		if search == value {
			return true
		}
	}
	return false
}

func IsInColC(m *[][][]int, row int, col int, search int) bool {
	for r, cRow := range *m {
		for c, candidate := range cRow[col] {
			if candidate > search {
				break
			}
			if (r == row) && (c == col) {
				continue
			}
			if search == candidate {
				return true
			}
		}
	}
	return false
}

func IsInCol(s *Solver, col int, search int) bool {
	for _, row := range s.Problem.Sudoku {
		if search == row[col] {
			return true
		}
	}
	return false
}

func IsInBlockC(m *[][][]int, dim int, row int, col int, search int) bool {

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	for r := blockRowStart; r < blockRowEnd; r++ {
		for c := blockColStart; c < blockColEnd; c++ {
			if (r == row) && (c == col) {
				continue
			}
			for _, candidate := range (*m)[r][c] {
				if candidate > search {
					break
				}
				if search == candidate {
					return true
				}
			}
		}
	}
	return false
}

func IsInBlock(s *Solver, row int, col int, search int, excludeRowsCols bool) bool {
	dim := s.Dim

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	for r := blockRowStart; r < blockRowEnd; r++ {
		if excludeRowsCols && (r == row) {
			continue
		}
		for c := blockColStart; c < blockColEnd; c++ {
			if excludeRowsCols && (c == col) {
				continue
			}
			if search == s.Problem.Sudoku[r][c] {
				return true
			}
		}
	}
	return false
}

func getRowColCandidates(s *Solver, row int, col int) *[]int {
	candidates := make([]int, 0, s.Length)
	for i := 1; i <= s.Length; i++ { // add candidates in sorted order (ascending)
		if IsInRow(s, row, i) || IsInCol(s, col, i) || IsInBlock(s, row, col, i, true) {
			continue
		} else {
			candidates = append(candidates, i)
		}
	}
	return &candidates
}

func UpdateAllCandidates(s *Solver) {
	s.Candidates = make([][][]int, s.Length)
	for r := range s.Candidates {
		s.Candidates[r] = make([][]int, s.Length)
		for c := range s.Candidates[r] {
			s.Candidates[r][c] = make([]int, 0, s.Length)
		}
	}

	for r := 0; r < s.Length; r++ {
		for c := 0; c < s.Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}
			for i := 1; i <= s.Length; i++ { // add candidates in sorted order (ascending)
				if IsInRow(s, r, i) || IsInCol(s, c, i) || IsInBlock(s, r, c, i, true) {
					continue
				} else {
					s.Candidates[r][c] = append(s.Candidates[r][c], i)
				}
			}
			if len(s.Candidates[r][c]) == 0 {
				panic("Sudoku is unsolvable")
			}
		}
	}
}

func UpdateCandidates(s *Solver, row int, col int, solvedCandidate int) (Solver, bool, bool) {
	updated := false
	updatedPrev := false

	// update candidates in rows
	for c, cRow := range s.Candidates[row] {
		for i, candidate := range cRow {
			if candidate > solvedCandidate {
				break
			}
			if solvedCandidate == candidate {
				s.Candidates[row][c] = append(s.Candidates[row][c][:i], s.Candidates[row][c][i+1:]...)
				updated = true
				updatedPrev = updatedPrev || (c < col)
			}
		}
	}

	// update candidates in cols
	for r := 0; r < s.Length; r++ {
		for i, candidate := range s.Candidates[r][col] {
			if candidate > solvedCandidate {
				break
			}
			if solvedCandidate == candidate {
				s.Candidates[r][col] = append(s.Candidates[r][col][:i], s.Candidates[r][col][i+1:]...)
				updated = true
				updatedPrev = updatedPrev || (r < row)
			}
		}
	}

	// update candidates in block
	dim := s.Dim

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	for r := blockRowStart; r < blockRowEnd; r++ {
		if r == row {
			continue
		}
		for c := blockColStart; c < blockColEnd; c++ {
			if c == col {
				continue
			}
			for i, candidate := range s.Candidates[r][c] {
				if candidate > solvedCandidate {
					break
				}
				if solvedCandidate == candidate {
					s.Candidates[r][c] = append(s.Candidates[r][c][:i], s.Candidates[r][c][i+1:]...)
					updated = true
					updatedPrev = updatedPrev || (r < row) || (c < col)
				}
			}
		}
	}

	return *s, updated, updatedPrev
}

// Naked Single
func SolveNakedSingle(s *Solver) (bool, bool) {
	updated := false

	var cUpdatedPrev bool
	foundEmpty := false

	for r := 0; r < s.Length; r++ {
		for c := 0; c < s.Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}

			if len(s.Candidates[r][c]) == 1 {
				s.Problem.Sudoku[r][c] = s.Candidates[r][c][0]
				s.Candidates[r][c] = nil //[]int{}
				*s, _, cUpdatedPrev = UpdateCandidates(s, r, c, s.Problem.Sudoku[r][c])
				updated = updated || cUpdatedPrev
			} else {
				foundEmpty = true
			}
		}
	}

	return updated, !foundEmpty
}

// Hidden Single
func SolveHiddenSingle(s *Solver) (bool, bool) {
	updated := false

	var cUpdatedPrev bool
	foundEmpty := false

	for r := 0; r < s.Length; r++ {
		for c := 0; c < s.Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}

			for i, candidate := range s.Candidates[r][c] {
				if !IsInRowC(&s.Candidates, r, c, candidate) && !IsInColC(&s.Candidates, r, c, candidate) && !IsInBlockC(&s.Candidates, s.Dim, r, c, candidate) {
					s.Problem.Sudoku[r][c] = candidate
					s.Candidates[r][c] = append(s.Candidates[r][c][:i], s.Candidates[r][c][i+1:]...)
					*s, _, cUpdatedPrev = UpdateCandidates(s, r, c, s.Problem.Sudoku[r][c])
					updated = updated || cUpdatedPrev
				} else {
					foundEmpty = true
				}
			}
		}
	}

	return updated, !foundEmpty
}

func findForwardPairInRow(m *[][][]int, dim int, row int, col int) (int, bool) {
	pair := (*m)[row][col]
	for c, cCol := range (*m)[row] {
		if c <= col {
			continue
		}
		if len(cCol) == 2 {
			pairFound := true
			for i := 0; i < 2; i++ {
				if cCol[i] != pair[i] {
					pairFound = false
					break
				}
			}
			if pairFound {
				return c, isCellsInSameBlock(dim, row, col, row, c)
			}
		}
	}
	return -1, false
}

func findForwardPairInCol(m *[][][]int, dim int, row int, col int) (int, bool) {
	pair := (*m)[row][col]
	for r, cRow := range *m {
		if r <= row {
			continue
		}
		if len(cRow[col]) == 2 {
			pairFound := true
			for i := 0; i < 2; i++ {
				if cRow[col][i] != pair[i] {
					pairFound = false
					break
				}
			}
			if pairFound {
				return r, isCellsInSameBlock(dim, row, col, r, col)
			}
		}
	}
	return -1, false
}

func isCellsInSameBlock(dim int, row1 int, col1 int, row2 int, col2 int) bool {
	blockRowStart1 := (row1 / dim) * dim
	blockColStart1 := (col1 / dim) * dim

	blockRowStart2 := (row2 / dim) * dim
	blockColStart2 := (col2 / dim) * dim

	return (blockRowStart1 == blockRowStart2) && (blockColStart1 == blockColStart2)
}

func removeCandidatesInRow(m *[][][]int, row int, exceptCol1 int, exceptCol2 int, pair *[]int) bool {
	updated := false
	for cd, candidates := range (*m)[row] {
		if (cd == exceptCol1) || (cd == exceptCol2) {
			continue
		}
		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]
			if candidate > (*pair)[1] {
				break
			}
			if (candidate == (*pair)[0]) || (candidate == (*pair)[1]) {
				candidates = append(candidates[:i], candidates[i+1:]...)
				(*m)[row][cd] = candidates
				updated = updated || (cd < exceptCol1)
				i--
			}
		}
	}
	return updated
}

func removeCandidatesInCol(m *[][][]int, col int, exceptRow1 int, exceptRow2 int, pair *[]int) bool {
	updated := false
	for rd, cRow := range *m {
		if (rd == exceptRow1) || (rd == exceptRow2) {
			continue
		}
		candidates := cRow[col]
		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]
			if candidate > (*pair)[1] {
				break
			}
			if (candidate == (*pair)[0]) || (candidate == (*pair)[1]) {
				candidates = append(candidates[:i], candidates[i+1:]...)
				(*m)[rd][col] = candidates
				updated = updated || (rd < exceptRow1)
				i--
			}
		}
	}
	return updated
}

func removeCandidatesInBlock(m *[][][]int, dim int, row1 int, col1 int, row2 int, col2 int, pair *[]int) bool {
	updated := false

	blockRowStart := (row1 / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col1 / dim) * dim
	blockColEnd := blockColStart + dim

	for r := blockRowStart; r < blockRowEnd; r++ {
		for c := blockColStart; c < blockColEnd; c++ {
			if ((r == row1) && (c <= col1)) || ((r == row2) && (c <= col2)) {
				continue
			}
			candidates := (*m)[r][c]
			for i := 0; i < len(candidates); i++ {
				candidate := candidates[i]
				if candidate > (*pair)[1] {
					break
				}
				if (candidate == (*pair)[0]) || (candidate == (*pair)[1]) {
					candidates = append(candidates[:i], candidates[i+1:]...)
					(*m)[r][c] = candidates
					updated = updated || (r < row1) || ((r == row1) && (c < col1))
					i--
				}
			}
		}
	}
	return updated
}

// Naked Pair
func SolveNakedPair(s *Solver) (bool, bool) {
	updated := false
	foundEmpty := false

	for r := 0; r < (*s).Length; r++ {
		for c := 0; c < (*s).Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}

			if len(s.Candidates[r][c]) == 2 {

				row := r
				col := c
				block := false
				col, block = findForwardPairInRow(&s.Candidates, s.Dim, r, c)
				if col == -1 {
					row, block = findForwardPairInCol(&s.Candidates, s.Dim, r, c)
					col = c
				}
				if (col != -1) && (row != -1) {

					updatedPrevInBlock := false
					// remove from rows
					if block {
						updatedPrevInBlock = removeCandidatesInBlock(&s.Candidates, s.Dim, r, c, row, col, &s.Candidates[r][c])
					}
					updatedPrev := false
					if r == row {
						updatedPrev = removeCandidatesInRow(&s.Candidates, r, c, col, &s.Candidates[r][c])
					} else if c == col {
						updatedPrev = removeCandidatesInCol(&s.Candidates, c, r, row, &s.Candidates[r][c])
					}

					updated = updated || updatedPrevInBlock || updatedPrev
				}
			} else {
				foundEmpty = true
			}
		}
	}

	return updated, !foundEmpty
}

func findCandidateOnlyInBlockRow(m *[][][]int, dim int, row int, col int, search int) bool {
	foundValid := false
	foundInvalid := false

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	for r := blockRowStart; r < blockRowEnd; r++ {
		for c := blockColStart; c < blockColEnd; c++ {
			if (r == row) && (c == col) {
				continue
			}
			if foundValid && (r == row) {
				continue
			}
			candidates := (*m)[r][c]
			for i := 0; i < len(candidates); i++ {
				candidate := candidates[i]
				if candidate > search {
					break
				}
				if candidate == search {
					if !foundValid {
						foundValid = r == row
					}
					foundInvalid = r != row
					if foundInvalid {
						return false
					}
					break
				}
			}
		}
	}
	return foundValid
}

func findCandidateOnlyInBlockCol(m *[][][]int, dim int, row int, col int, search int) bool {
	foundValid := false
	foundInvalid := false

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	for r := blockRowStart; r < blockRowEnd; r++ {
		for c := blockColStart; c < blockColEnd; c++ {
			if (r == row) && (c == col) {
				continue
			}
			if foundValid && (c == col) {
				continue
			}
			candidates := (*m)[r][c]
			for i := 0; i < len(candidates); i++ {
				candidate := candidates[i]
				if candidate > search {
					break
				}
				if candidate == search {
					if !foundValid {
						foundValid = c == col
					}
					foundInvalid = c != col
					if foundInvalid {
						return false
					}
				}
			}
		}
	}
	return foundValid
}

func removeCandidateInRowOutsideBlock(m *[][][]int, length int, dim int, row int, col int, search int) bool {
	updated := false

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	r := row
	for c := 0; c < length; c++ {
		if (r >= blockRowStart) && (c >= blockColStart) && (r < blockRowEnd) && (c < blockColEnd) { // ignore the current block
			continue
		}
		candidates := (*m)[r][c]
		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]
			if candidate > search {
				break
			}
			if candidate == search {
				candidates = append(candidates[:i], candidates[i+1:]...)
				(*m)[r][c] = candidates
				updated = true
				break
			}
		}
	}
	return updated
}

func removeCandidateInColOutsideBlock(m *[][][]int, length int, dim int, row int, col int, search int) bool {
	updated := false

	blockRowStart := (row / dim) * dim
	blockRowEnd := blockRowStart + dim
	blockColStart := (col / dim) * dim
	blockColEnd := blockColStart + dim

	for r := 0; r < length; r++ {
		c := col
		if (r >= blockRowStart) && (c >= blockColStart) && (r < blockRowEnd) && (c < blockColEnd) { // ignore the current block
			continue
		}
		candidates := (*m)[r][c]
		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]
			if candidate > search {
				break
			}
			if candidate == search {
				candidates = append(candidates[:i], candidates[i+1:]...)
				(*m)[r][c] = candidates
				updated = true
				break
			}
		}
	}
	return updated
}

// Point Pair (Triple)
func SolvePointingPair(s *Solver) (bool, bool) {
	updated := false
	foundEmpty := false

	for r := 0; r < (*s).Length; r++ {
		for c := 0; c < (*s).Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}
			for _, candidate := range s.Candidates[r][c] {
				if findCandidateOnlyInBlockRow(&s.Candidates, s.Dim, r, c, candidate) {
					removeCandidateInRowOutsideBlock(&s.Candidates, s.Length, s.Dim, r, c, candidate)
				} else if findCandidateOnlyInBlockCol(&s.Candidates, s.Dim, r, c, candidate) {
					removeCandidateInColOutsideBlock(&s.Candidates, s.Length, s.Dim, r, c, candidate)
				} else {
					foundEmpty = true
				}
			}
		}
	}
	return updated, !foundEmpty
}

func SolveDepthFirstSearch(s *Solver, rInit int, cInit int, rec int) bool {
	//fmt.Printf("rec:%v\n", rec)
	solved := false
	emptyFound := false
	row := -1
	col := -1

	c := cInit
	for r := rInit; r < (*s).Length; r++ {
		for ; c < (*s).Length; c++ {
			if s.Problem.Sudoku[r][c] == 0 {
				row = r
				col = c
				emptyFound = true
				break
			}
		}
		if emptyFound {
			break
		}
		c = 0
	}

	if !emptyFound { // sudoku is solved
		solved = true
	} else {

		candidates := getRowColCandidates(s, row, col)

		for i := 0; i < len(*candidates); i++ {
			s.Problem.Sudoku[row][col] = (*candidates)[i]
			valid := true
			//valid := CheckValue(s, row, col)
			if valid {
				solved = SolveDepthFirstSearch(s, row, col+1, rec+1)
				if solved {
					break
				}
			}
		}

		if !solved {
			s.Problem.Sudoku[row][col] = 0
		}
	}
	return solved
}

func Solve(s *Solver) bool {

	UpdateAllCandidates(s)

	fmt.Println("solving with NakedSingle")
	updated, solved := SolveNakedSingle(s)

	var cUpdated bool
	exit := solved

	for !exit {
		cUpdated = updated
		updated = false
		for cUpdated && !solved {
			cUpdated, solved = SolveNakedSingle(s)
			fmt.Println("again solving with NakedSingle")
		}

		if !solved {
			fmt.Println("solving with HiddenSingle")
			cUpdated, solved = SolveHiddenSingle(s)
			updated = updated || cUpdated
		}

		if !solved {
			fmt.Println("solving with NakedPair")
			cUpdated, solved = SolveNakedPair(s)
			updated = updated || cUpdated
		}

		if !solved {
			fmt.Println("solving with PointingPair")
			cUpdated, solved = SolvePointingPair(s)
			updated = updated || cUpdated
		}

		exit = !updated || solved
	}

	fmt.Println("after strategies solved:")
	Print(s)

	if !solved {
		fmt.Println("not solved so far - using DepthFirstSearch")
		solved = SolveDepthFirstSearch(s, 0, 0, 1)
	}
	return solved
}
