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

func waitForSolvers(n int, done <-chan bool, solutions chan *sudoku.Solution) {
	// wait for n workers to be done
	for i := 0; i < n; i++ {
		<-done
	}
	// close the solutions channel, which will allow collectResults to return
	close(solutions)
}

func collectResults(solutions <-chan *sudoku.Solution) {
	count := 0
	start := time.Now()

	for s := range solutions {
		count++
		fmt.Println(s)
	}

	elapsed := time.Since(start)
	rate := float64(count) / elapsed.Seconds()
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
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

	workers := runtime.NumCPU()
	runtime.GOMAXPROCS(workers)

	fName := os.Args[1]

	unsolved := make(chan *sudoku.Board)
	solved := make(chan *sudoku.Solution)
	done := make(chan bool)

	go loadBoards(fName, unsolved)

	for i := 0; i < workers; i++ {
		go solver(unsolved, solved, done)
	}

	go waitForSolvers(workers, done, solved)

	collectResults(solved)
}
