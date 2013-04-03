package sudoku

import (
	"strconv"
	"bufio"
	"errors"
)

// read a board where each row is on a separate line
func ReadBoard(reader *bufio.Reader) Board {
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

// read a board where the entire board is on one line
func ReadBoardLine(reader *bufio.Reader) (Board, error) {
	str, _ := reader.ReadString('\n')
	board := NewBoard()
	if str == "" {
		return board, errors.New("no more boards")
	}
	for i := 0; i < SZ; i++ {
		ch := str[i : i+1]
		val, _ := strconv.Atoi(ch)
		if (val > 0) {
			board = board.Set(i, val)
		}
	}

	return board, nil
}