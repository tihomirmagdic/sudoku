package solver

import (
	"fmt"
	"testing"
	"time"
)

func TestSudokuSolver1(t *testing.T) {
	s := SudokuMatrix{
		Sudoku: [][]int{
			{8, 0, 0, +0, 0, 7, +0, 9, 0},
			{0, 2, 9, +0, 0, 4, +0, 0, 6},
			{3, 0, 0, +2, 0, 0, +0, 0, 0},

			{0, 0, 0, +0, 0, 6, +5, 0, 0},
			{0, 1, 7, +4, 0, 0, +0, 3, 0},
			{2, 0, 0, +0, 0, 0, +0, 0, 0},

			{0, 9, 4, +1, 0, 0, +0, 7, 0},
			{0, 0, 8, +0, 0, 0, +0, 0, 0},
			{0, 0, 0, +0, 7, 0, +0, 0, 3},
		}}

	v, err := CheckSudoku(&s)
	if err != nil {
		panic(fmt.Sprintf("error: %v\n", err))
	} else {
		fmt.Println("initial values ")
	}
	fmt.Println("sudoku to solve:")
	Print(v)

	start := time.Now()
	solved := Solve(v)
	duration := time.Since(start)

	fmt.Printf("total %s (%d)\n", duration, duration.Nanoseconds())

	fmt.Printf("sudoku solved: %v\n", solved)
	Print(v)
	_, err = CheckSudoku(&v.Problem)
	if err == nil {
		fmt.Println("sudoku valid")
	} else {
		fmt.Println("sudoku not valid")
	}
	if !solved {
		t.Errorf("Sudoku not solved")
	}
}

func TestSudokuSolver2(t *testing.T) {
	s := SudokuMatrix{
		Sudoku: [][]int{
			{4, 7, 0, +3, 0, 0, +2, 1, 8},
			{0, 8, 2, +4, 0, 1, +7, 0, 3},
			{1, 3, 0, +0, 8, 0, +0, 4, 5},

			{0, 1, 0, +0, 0, 0, +3, 0, 0},
			{6, 0, 3, +0, 1, 5, +4, 0, 0},
			{7, 4, 0, +0, 3, 0, +0, 0, 0},

			{8, 0, 1, +0, 0, 0, +5, 3, 9},
			{0, 0, 7, +5, 0, 0, +1, 0, 4},
			{0, 5, 4, +1, 0, 0, +0, 7, 0},
		}}

	v, err := CheckSudoku(&s)
	if err != nil {
		panic(fmt.Sprintf("error: %v\n", err))
	} else {
		fmt.Println("initial values ")
	}
	fmt.Println("sudoku to solve:")
	Print(v)

	start := time.Now()
	solved := Solve(v)
	duration := time.Since(start)

	fmt.Printf("total %s (%d)\n", duration, duration.Nanoseconds())

	fmt.Printf("sudoku solved: %v\n", solved)
	Print(v)
	_, err = CheckSudoku(&v.Problem)
	if err == nil {
		fmt.Println("sudoku valid")
	} else {
		fmt.Println("sudoku not valid")
	}
	if !solved {
		t.Errorf("Sudoku not solved")
	}
}

func TestSudokuSolver3(t *testing.T) {
	s := SudokuMatrix{
		Sudoku: [][]int{
			{0, 1, 0, +0, 0, 3, +0, 6, 0},
			{0, 0, 0, +0, 9, 0, +2, 0, 0},
			{5, 4, 8, +0, 6, 0, +0, 0, 0},

			{0, 0, 0, +9, 0, 0, +0, 2, 0},
			{0, 0, 7, +0, 0, 0, +5, 0, 0},
			{1, 0, +0, 0, 5, +0, 0, 0, 0},

			{0, 0, 0, +0, 8, 0, +4, 5, 0},
			{0, 0, 9, +0, 4, 0, +8, 0, 0},
			{0, 2, 0, +6, 0, 0, +0, 9, 3},
		}}

	v, err := CheckSudoku(&s)
	if err != nil {
		panic(fmt.Sprintf("error: %v\n", err))
	} else {
		fmt.Println("initial values ")
	}
	fmt.Println("sudoku to solve:")
	Print(v)

	start := time.Now()
	solved := Solve(v)
	duration := time.Since(start)

	fmt.Printf("total %s (%d)\n", duration, duration.Nanoseconds())

	fmt.Printf("sudoku solved: %v\n", solved)
	Print(v)
	_, err = CheckSudoku(&v.Problem)
	if err == nil {
		fmt.Println("sudoku valid")
	} else {
		fmt.Println("sudoku not valid")
	}
	if solved {
		t.Errorf("Sudoku not solved")
	}
}

func TestSudokuSolver4(t *testing.T) {
	s := SudokuMatrix{
		Sudoku: [][]int{
			{1, 2, +3, 0},
			{0, 0, +0, 4},

			{0, 0, +0, 0},
			{0, 0, +0, 0},
		}}

	v, err := CheckSudoku(&s)
	if err != nil {
		panic(fmt.Sprintf("error: %v\n", err))
	} else {
		fmt.Println("initial values ")
	}
	fmt.Println("sudoku to solve:")
	Print(v)

	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\nOK: panic occurred: %v\n", err)
			}
		}()

		Solve(v)

		errorStr := "should raise panic"
		fmt.Printf("\n" + errorStr + "\n")
		t.Errorf(errorStr)
	}()
}

func TestSudokuSolver5(t *testing.T) {
	s := SudokuMatrix{
		Sudoku: make([][]int, 9),
	}
	for i := range s.Sudoku {
		s.Sudoku[i] = make([]int, 9)
	}

	v, err := CheckSudoku(&s)
	if err != nil {
		panic(fmt.Sprintf("error: %v\n", err))
	} else {
		fmt.Println("initial values ")
	}
	fmt.Println("sudoku to solve:")
	Print(v)

	start := time.Now()
	solved := Solve(v)
	duration := time.Since(start)

	fmt.Printf("total %s (%d)\n", duration, duration.Nanoseconds())

	fmt.Printf("sudoku solved: %v\n", solved)
	Print(v)
	_, err = CheckSudoku(&v.Problem)
	if err == nil {
		fmt.Println("sudoku valid")
	} else {
		fmt.Println("sudoku not valid")
	}
	if !solved {
		t.Errorf("Sudoku not solved")
	}
}
