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
	if len(os.Args) < 2 {
		println("Usage: sudoku_batch [file]")
		os.Exit(1)
	}

	fName := os.Args[1]
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

		solution := b.Solution()
		fmt.Println(solution.String())
	}

	elapsed := time.Since(start)
	rate := (float64(count) / elapsed.Seconds())
	fmt.Printf("Solved %v puzzles in %v (%0.2f per second)\n", count, elapsed, rate)
}