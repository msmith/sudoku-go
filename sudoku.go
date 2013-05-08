package main

import (
	"fmt"
	"os"
	"sudoku"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: sudoku [puzzle_file]")
		os.Exit(1)
	}

	fName := os.Args[1]
	b := sudoku.ReadBoardFile(fName)

	solution := b.Solve()
	fmt.Println(solution.Original)
	fmt.Println(solution.Solved)
	fmt.Println(solution.Elapsed)
}
