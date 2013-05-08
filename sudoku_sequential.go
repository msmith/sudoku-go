package main

import (
	"fmt"
	"os"
	"sudoku"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: sudoku_batch [file]")
		os.Exit(1)
	}

	fName := os.Args[1]
	count := 0
	start := time.Now()

	sudoku.ReadBoardSet(fName, func(b sudoku.Board) {
		count++

		solution := b.Solve()
		fmt.Println(solution.String())
	})

	elapsed := time.Since(start)
	rate := float64(count) / elapsed.Seconds()
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
}
