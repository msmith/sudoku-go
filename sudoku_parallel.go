package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sudoku"
	"compress/gzip"
	"time"
	"runtime"
)

func solver(unsolved <-chan *sudoku.Board, solved chan<- *sudoku.Solution, done chan bool) {
	for board := range unsolved {
		solution := board.Solution()
		solved <- &solution
	}
	done <- true
}

func collectResults(solutions chan *sudoku.Solution) {
	var count int64
	start := time.Now()
	for s := range solutions {
		count++
		fmt.Println(s)
	}
	elapsed := time.Since(start)
	rate := (float64(count) / elapsed.Seconds())
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
}

func waitForSolvers(workers int, done chan bool, toClose chan *sudoku.Solution) {
	for i := 0; i < workers; i++ {
		<- done
	}
	close(toClose)
}

func loadBoards(fName string, unsolved chan *sudoku.Board) {
	file, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gz_reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(gz_reader)

	for {
		b, err := sudoku.ReadBoardLine(reader)
		if err != nil {
			break
		}
		unsolved <- &b
	}
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

	for i := 0; i < workers; i++ {
		go solver(unsolved, solved, done)
	}

	go loadBoards(fName, unsolved)

	go waitForSolvers(workers, done, solved)

	collectResults(solved)
}