package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sudoku"
)

func main() {
	var fName string

	if len(os.Args) > 1 {
		fName = os.Args[1]
	} else {
		fName = "puzzles/al_escargot"
	}

	file, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	b := sudoku.ReadBoard(reader)

	file.Close()

	fmt.Println(b.String())
	b, _ = b.Solve()
	fmt.Println(b.String())
}