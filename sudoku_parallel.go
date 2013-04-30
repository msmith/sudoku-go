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

func solver(in <-chan *sudoku.Board, out chan<- *sudoku.Solution, done chan bool) {
	for b := range in {
		start := time.Now()
		b2, _ := b.Solve()
		t := time.Since(start)

		solution := sudoku.Solution{b, &b2, t}

		out <- &solution
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

func waitForDone(workers int, done chan bool, toClose chan *sudoku.Solution) {
	for i := 0; i < workers; i++ {
		<- done
	}
	close(toClose)
}

func main() {
	var fName string

	if len(os.Args) > 1 {
		fName = os.Args[1]
	} else {
		fName = "sets/sudoku17.gz"
	}

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

	workers := runtime.NumCPU()
	runtime.GOMAXPROCS(workers)

	unsolved := make(chan *sudoku.Board)
	solved := make(chan *sudoku.Solution)
	done := make(chan bool)

	for i := 0; i < workers; i++ {
		go solver(unsolved, solved, done)
	}

	go func() {
		for {
			b, err := sudoku.ReadBoardLine(reader)
			if err != nil {
				break
			}
			unsolved <- &b
		}
		close(unsolved)
	}()

	go waitForDone(workers, done, solved)

	collectResults(solved)
}