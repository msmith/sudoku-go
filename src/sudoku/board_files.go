package sudoku

import (
	"bufio"
	"compress/gzip"
	"errors"
	"log"
	"os"
	"strconv"
)

// read a board where each row is on a separate line
func readBoard(reader *bufio.Reader) Board {
	board := NewBoard()

	var row int

	for {
		str, _ := reader.ReadString('\n')
		if str == "" {
			break
		}
		for col := 0; col < DIM2; col++ {
			ch := str[col : col+1]
			val, _ := strconv.Atoi(ch)
			idx := indexOf(row, col)
			if val > 0 {
				board = board.Set(idx, val)
			}
		}
		row++
	}

	return board
}

func ReadBoardFile(fName string) Board {
	file, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return readBoard(reader)
}

// read a board where the entire board is on one line
func readBoardLine(reader *bufio.Reader) (Board, error) {
	str, _ := reader.ReadString('\n')
	board := NewBoard()
	if str == "" {
		return board, errors.New("no more boards")
	}
	for i := 0; i < SZ; i++ {
		ch := str[i : i+1]
		val, _ := strconv.Atoi(ch)
		if val > 0 {
			board = board.Set(i, val)
		}
	}

	return board, nil
}

type gotboard func(Board)

func ReadBoardSet(fName string, f gotboard) {
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
		b, err := readBoardLine(reader)
		if err != nil {
			break
		}
		f(b)
	}
}