package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sudoku"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: sudoku [puzzle_file]")
		os.Exit(1)
	}

	fName := os.Args[1]
	file, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	b := sudoku.ReadBoard(reader)

	fmt.Println(b.String())
	b, _ = b.Solve()
	fmt.Println(b.String())
}