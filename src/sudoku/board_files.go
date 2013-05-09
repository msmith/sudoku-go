package sudoku

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// read a single puzzle file
func ReadBoardFile(fName string) (Board, error) {
	bytes, err := ioutil.ReadFile(fName)
	if err != nil {
		return NewBoard(), err
	}
	return boardFromString(string(bytes))
}

type gotboard func(Board)

// read a puzzle set file
func ReadBoardSet(fName string, f gotboard) error {
	file, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer file.Close()

	gz_reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(gz_reader)

	for {
		b, err := readBoardLine(reader)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		f(b)
	}

	return nil
}

func boardFromString(str string) (Board, error) {
	// remove spaces & newlines
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)

	board := NewBoard()
	if len(str) != SZ {
		return board, errors.New("line was incorrect length")
	}
	for i := 0; i < len(str); i++ {
		ch := str[i : i+1]
		val, _ := strconv.Atoi(ch)
		if val > 0 {
			board = board.Set(i, val)
		}
	}
	if !board.Valid() {
		return board, errors.New("board is invalid")
	}
	return board, nil
}

func readBoardLine(reader *bufio.Reader) (Board, error) {
	str, err := reader.ReadString('\n')
	if err != nil {
		return NewBoard(), err
	}

	return boardFromString(str)
}
