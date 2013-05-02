package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sudoku"
	"compress/gzip"
	"time"
)

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

	var count int64
	start := time.Now()

	for {
		b, err := sudoku.ReadBoardLine(reader)
		if err != nil {
			break;
		}
		count++

		boardStart := time.Now()
		b2, _ := b.Solve()
		t := time.Since(boardStart)

		solution := sudoku.Solution{&b, &b2, t}

		fmt.Println(solution.String())
	}

	elapsed := time.Since(start)
	rate := (float64(count) / elapsed.Seconds())
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
}