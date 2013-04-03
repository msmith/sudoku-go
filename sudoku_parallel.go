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

func solver(in <-chan *sudoku.Board, out chan<- *sudoku.Board, done chan bool) {
	for b := range in {
		b2, _ := b.Solve()
		out <- &b2
	}
	done <- true
}

func collectResults(boards chan *sudoku.Board) {
	var count int64
	start := time.Now()
	for b := range boards {
		count++
		elapsed := time.Since(start)
		rate := (elapsed.Nanoseconds()/count) / 1000000
		fmt.Printf("Solved %v (avg. %d ms)\n", count, rate)
		fmt.Println(b.String())
	}
}

func waitForDone(workers int, done chan bool, toClose chan *sudoku.Board) {
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

	gz_reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(gz_reader)

	workers := runtime.NumCPU()
	runtime.GOMAXPROCS(workers)

	unsolved := make(chan *sudoku.Board)
	solved := make(chan *sudoku.Board)
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

	file.Close()
}