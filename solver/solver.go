package solver

import (
	"github.com/tihomirmagdic/sudoku/types"
	"github.com/tihomirmagdic/sudoku/validator"
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

func IsInRow(s *types.Solver, row int, search int) bool {
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

func IsInCol(s *types.Solver, col int, search int) bool {
	for _, row := range s.Problem.Sudoku {
		if search == row[col] {
			return true
		}
	}
	return false
}

func IsInBlockC(m *[][][]int, dim int, row int, col int, search int) bool {
	for r := (row / dim) * dim; r < ((row/dim)+1)*dim; r++ {
		for c := (col / dim) * dim; c < ((col/dim)+1)*dim; c++ {
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

func IsInBlock(s *types.Solver, row int, col int, search int, excludeRowsCols bool) bool {
	dim := s.Dim
	for r := (row / dim) * dim; r < ((row/dim)+1)*dim; r++ {
		if excludeRowsCols && (r == row) {
			continue
		}
		for c := (col / dim) * dim; c < ((col/dim)+1)*dim; c++ {
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

func GetCandidates(s *types.Solver) types.Solver {
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
	return *s
}

func UpdateCandidates(s *types.Solver, row int, col int, solvedCandidate int) (types.Solver, bool, bool) {
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
	for r := (row / dim) * dim; r < ((row/dim)+1)*dim; r++ {
		if r == row {
			continue
		}
		for c := (col / dim) * dim; c < ((col/dim)+1)*dim; c++ {
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
func SolveNakedSingle(s *types.Solver) (types.Solver, bool) {
	updated := false
	var cUpdatedPrev bool

	for r := 0; r < s.Length; r++ {
		for c := 0; c < s.Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}

			if len(s.Candidates[r][c]) == 1 {
				s.Problem.Sudoku[r][c] = s.Candidates[r][c][0]
				s.Candidates[r][c] = nil
				*s, _, cUpdatedPrev = UpdateCandidates(s, r, c, s.Problem.Sudoku[r][c])
				updated = updated || cUpdatedPrev
			}
		}
	}

	if updated {
		*s, updated = SolveNakedSingle(s)
	}
	return *s, updated
}

// Hidden Single
func SolveHiddenSingle(s *types.Solver) (types.Solver, bool) {
	updated := false
	var cUpdatedPrev bool

	for r := 0; r < s.Length; r++ {
		for c := 0; c < s.Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}

			for i, candidate := range s.Candidates[r][c] {
				if !IsInRowC(&s.Candidates, r, c, candidate) || !IsInColC(&s.Candidates, r, c, candidate) || !IsInBlockC(&s.Candidates, s.Dim, r, c, candidate) {
					s.Problem.Sudoku[r][c] = candidate
					s.Candidates[r][c] = append(s.Candidates[r][c][:i], s.Candidates[r][c][i+1:]...)
					*s, _, cUpdatedPrev = UpdateCandidates(s, r, c, s.Problem.Sudoku[r][c])
					updated = updated || cUpdatedPrev

				}
			}
		}
	}

	if updated {
		*s, _ = SolveNakedSingle(s)
		if updated {
			*s, updated = SolveNakedSingle(s)
		}
	}
	return *s, updated
}

func findForwardPairInRow(m *[][][]int, row int, col int) int {
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
				return c
			}
		}
	}
	return -1
}

func findForwardPairInCol(m *[][][]int, row int, col int) int {
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
				return r
			}
		}
	}
	return -1
}

func findForwardPairInBlock(m *[][][]int, dim int, row int, col int) (int, int) {
	pair := (*m)[row][col]
	for r := (row / dim) * dim; r < ((row/dim)+1)*dim; r++ {
		if r < row {
			continue
		}
		for c := (col / dim) * dim; c < ((col/dim)+1)*dim; c++ {
			if (r == row) && (c <= col) {
				continue
			}
			if len((*m)[r][c]) == 2 {
				pairFound := true
				for i := 0; i < 2; i++ {
					if (*m)[r][c][i] != pair[i] {
						pairFound = false
						break
					}
				}
				if pairFound {
					return r, c
				}
			}
		}
	}
	return -1, -1
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
				updated = updated || (rd < exceptRow1)
				i--
			}
		}
	}
	return updated
}

func removeCandidatesInBlock(m *[][][]int, dim int, row1 int, col1 int, row2 int, col2 int, pair *[]int) bool {
	updated := false
	for r := (row1 / dim) * dim; r < ((row1/dim)+1)*dim; r++ {
		for c := (col1 / dim) * dim; c < ((col1/dim)+1)*dim; c++ {
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
					updated = updated || (r < row1) || ((r == row1) && (c < col1))
					i--
				}
			}
		}
	}
	return updated
}

// Naked Pair
func SolveNakedPair(s *types.Solver) (types.Solver, bool) {
	updated := false

	for r := 0; r < (*s).Length; r++ {
		for c := 0; c < (*s).Length; c++ {
			if s.Problem.Sudoku[r][c] != 0 {
				continue
			}

			if len(s.Candidates[r][c]) == 2 {
				block := false
				var row int
				col := findForwardPairInRow(&s.Candidates, r, c)
				if col != -1 {
					row = r
				} else {
					row = findForwardPairInCol(&s.Candidates, r, c)
					if row != -1 {
						col = c
					} else {
						row, col = findForwardPairInBlock(&s.Candidates, s.Dim, r, c)
						block = true
					}
				}
				if (col != -1) && (row != -1) {

					// remove from rows
					updatedPrev := removeCandidatesInRow(&s.Candidates, r, c, col, &s.Candidates[r][c])
					updatedPrev = updatedPrev || removeCandidatesInRow(&s.Candidates, row, c, col, &s.Candidates[r][c])

					updatedPrev = updatedPrev || removeCandidatesInCol(&s.Candidates, c, r, row, &s.Candidates[r][c])
					updatedPrev = updatedPrev || removeCandidatesInCol(&s.Candidates, col, r, row, &s.Candidates[r][c])

					if block {
						updatedPrev = updatedPrev || removeCandidatesInBlock(&s.Candidates, s.Dim, r, row, c, col, &s.Candidates[r][c])
					}
					updated = updated || updatedPrev
				}
			}
		}
	}
	return *s, updated
}

func findCandidateOnlyInBlockRow(m *[][][]int, dim int, row int, col int, search int) bool {
	foundValid := false
	foundInvalid := false

	for r := (row / dim) * dim; r < ((row/dim)+1)*dim; r++ {
		for c := (col / dim) * dim; c < ((col/dim)+1)*dim; c++ {
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
				}
			}
		}
	}
	return foundValid
}

func findCandidateOnlyInBlockCol(m *[][][]int, dim int, row int, col int, search int) bool {
	foundValid := false
	foundInvalid := false

	for r := (row / dim) * dim; r < ((row/dim)+1)*dim; r++ {
		for c := (col / dim) * dim; c < ((col/dim)+1)*dim; c++ {
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
	blockRowEnd := ((row / dim) + 1) * dim
	blockColStart := (col / dim) * dim
	blockColEnd := ((col / dim) + 1) * dim
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
	blockRowEnd := ((row / dim) + 1) * dim
	blockColStart := (col / dim) * dim
	blockColEnd := ((col / dim) + 1) * dim
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
func SolvePointingPair(s *types.Solver) (types.Solver, bool) {
	updated := false

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
				}
			}
		}
	}
	return *s, updated
}

func SolveDepthFirstSearch(s *types.Solver, r int) (types.Solver, bool) {
	//fmt.Printf("r:%v\n", r)
	solved := false
	row := -1
	col := -1

	for r := 0; r < (*s).Length; r++ {
		for c := 0; c < (*s).Length; c++ {
			if s.Problem.Sudoku[r][c] == 0 {
				row = r
				col = c
				break
			}
		}
		if row != -1 {
			break
		}
	}
	if row == -1 { // sudoku is solved
		solved = true
	} else {

		candidates := s.Candidates[row][col]
		oldCandidates := make([]int, len(candidates))
		copy(oldCandidates, candidates)

		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]
			s.Problem.Sudoku[row][col] = candidate
			//fmt.Printf("i: %v\n", i)
			//fmt.Printf("canidates: %v\n", candidates)
			//_, err := validator.CheckSudoku(&s.Problem)
			valid := validator.CheckValue(s, row, col)
			//copy(candidates, oldCandidates)
			//if err == nil {
			if valid {
				candidates = append(candidates[:i], candidates[i+1:]...)
				s.Candidates[row][col] = candidates
				i--
				*s, solved = SolveDepthFirstSearch(s, r+1)
				if solved {
					break
				}
			}
		}

		if !solved {
			s.Problem.Sudoku[row][col] = 0
			s.Candidates[row][col] = oldCandidates
		}
	}
	return *s, solved
}

func Solve(s *types.Solver) (types.Solver, bool) {
	return SolveDepthFirstSearch(s, 1)
}
