package sudoku

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"
	"regexp"
)

var whitespace *regexp.Regexp = regexp.MustCompile("\\s+")

// A file containing a single board
type SingleBoardFile string

func (f SingleBoardFile) Board() (*Board, error) {
	bytes, err := ioutil.ReadFile(string(f))
	if err != nil {
		return nil, err
	}
	return ParseBoard(string(bytes))
}

// A file containing multiple boards
type MultipleBoardFile string

type gotBoard func(*Board)

func (f MultipleBoardFile) EachBoard(callback gotBoard) error {
	file, err := os.Open(string(f))
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(gzReader)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		boardString := scanner.Text()
		b, err := ParseBoard(boardString)
		if err != nil {
			return err
		}
		callback(b)
	}

	return nil
}
