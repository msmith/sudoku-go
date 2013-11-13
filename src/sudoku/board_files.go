package sudoku

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

var whitespace *regexp.Regexp = regexp.MustCompile("\\s+")

// A file containing a single board
type BoardFile string

func (f BoardFile) Board() (*Board, error) {
	bytes, err := ioutil.ReadFile(string(f))
	if err != nil {
		return nil, err
	}
	return boardFromString(string(bytes))
}

// A file containing multiple boards
type BoardSet string

type gotBoard func(*Board)

func (f BoardSet) EachBoard(callback gotBoard) error {
	file, err := os.Open(string(f))
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
		b, err := ReadBoardLine(reader)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		callback(b)
	}

	return nil
}

func boardFromString(str string) (*Board, error) {
	// remove whitespace
	str = whitespace.ReplaceAllString(str, "")

	if len(str) != SZ {
		return nil, errors.New("line was incorrect length")
	}
	board := NewBoard()
	for i := 0; i < len(str); i++ {
		ch := str[i : i+1]
		val, _ := strconv.Atoi(ch)
		if val > 0 {
			board = board.Set(i, val)
		}
	}
	if !board.Valid() {
		return nil, errors.New("board is invalid")
	}
	return &board, nil
}

func ReadBoardLine(reader *bufio.Reader) (*Board, error) {
	str, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return boardFromString(str)
}
