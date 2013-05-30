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

	file := sudoku.BoardFile(os.Args[1])
	b, err := file.Board()
	if err != nil {
		panic(err)
	}

	solution := b.Solve()
	fmt.Println(solution.Original)
	fmt.Println(solution.Solved)
	fmt.Println(solution.Elapsed)
}
