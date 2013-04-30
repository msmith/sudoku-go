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

	var count int64
	start := time.Now()

	for {
		b, err := sudoku.ReadBoardLine(reader)
		if err != nil {
			break;
		}
		count++
		b2, _ := b.Solve()
		elapsed := time.Since(start)
		rate := (float64(count) / elapsed.Seconds())
		fmt.Printf("Solved %v (%0.2f per second)\n", count, rate)
		fmt.Println(b2.String())
	}

	file.Close()
}