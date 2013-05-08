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
	var hardest sudoku.Solution

	sudoku.ReadBoardSet(fName, func(b sudoku.Board) {
		s := b.Solve()
		fmt.Println(s.String())

		count++
		if s.Elapsed > hardest.Elapsed {
			hardest = s
		}
	})

	elapsed := time.Since(start)
	rate := float64(count) / elapsed.Seconds()
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
	fmt.Printf("Hardest puzzle took %v\n", hardest.Elapsed)
}
