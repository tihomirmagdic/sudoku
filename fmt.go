package solver

import (
	"fmt"
	"strconv"
)

const NullValue = "."
const Space = "  "

func Print(s *Solver) {
	l := len(strconv.Itoa(s.Length))

	for rInd, row := range s.Problem.Sudoku {
		if (rInd != 0) && (rInd%s.Dim) == 0 {
			fmt.Println()
		}
		for cInd, value := range row {
			if (cInd != 0) && (cInd%s.Dim) == 0 {
				fmt.Print(Space)
			}
			if value == 0 {
				fmt.Printf("%*v ", l, NullValue)
			} else {
				fmt.Printf("%*v ", l, value)
			}
		}
		fmt.Println()
	}
}
