package main

import (
	"fmt"
	"os"
	"runtime"
	"sudoku"
	"time"
)

func solver(boards <-chan *sudoku.Board, solutions chan<- *sudoku.Solution, done chan bool) {
	for board := range boards {
		s := board.Solve()
		solutions <- &s
	}
	// signal that this solver is done
	done <- true
}

func waitForSolvers(workerCount int, done <-chan bool, solutions chan *sudoku.Solution) {
	// wait for workers to be done
	for i := 0; i < workerCount; i++ {
		<-done
	}
	// close the solutions channel, which will allow collectResults to return
	close(solutions)
}

func collectResults(solutions <-chan *sudoku.Solution) {
	count := 0
	start := time.Now()
	var hardest sudoku.Solution

	for s := range solutions {
		fmt.Println(s.String())

		count++
		if s.Elapsed > hardest.Elapsed {
			hardest = *s
		}
	}

	elapsed := time.Since(start)
	rate := float64(count) / elapsed.Seconds()
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
	fmt.Printf("Hardest puzzle took %v\n", hardest.Elapsed)
}

func loadBoards(fName string, unsolved chan<- *sudoku.Board) {
	sudoku.ReadBoardSet(fName, func(b sudoku.Board) {
		unsolved <- &b
	})
	close(unsolved)
}

func main() {
	if len(os.Args) < 2 {
		println("Usage: sudoku_parallel [file]")
		os.Exit(1)
	}

	workerCount := runtime.NumCPU()
	runtime.GOMAXPROCS(workerCount)

	fName := os.Args[1]

	boards := make(chan *sudoku.Board)
	solutions := make(chan *sudoku.Solution)
	done := make(chan bool)

	go loadBoards(fName, boards)

	for i := 0; i < workerCount; i++ {
		go solver(boards, solutions, done)
	}

	go waitForSolvers(workerCount, done, solutions)

	collectResults(solutions)
}
