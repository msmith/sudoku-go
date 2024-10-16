package sudoku

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

const (
	DIM  = 3
	DIM2 = DIM * DIM
	SZ   = DIM2 * DIM2
)

type Board [SZ]cell

// Peers is a lookup table that provides the peer cells of each cell. A
// peer is a cell that lives in the same row, column, or region as this one.
var Peers [][]int = initPeers()

// Groups contains the cell indexes of each group (row, column, and region)
var Groups [][]int = initGroups()

func initPeers() [][]int {
	peers := make([][]int, SZ)
	for i := 0; i < SZ; i++ {
		row, col, reg := posOf(i)

		for j := 0; j < SZ; j++ {
			if j != i {
				row_j, col_j, reg_j := posOf(j)
				if row_j == row || col_j == col || reg_j == reg {
					peers[i] = append(peers[i], j)
				}
			}
		}
	}
	return peers
}

func initGroups() [][]int {
	groups := make([][]int, DIM2*3)
	byRow := groups[0:DIM2]
	byCol := groups[DIM2 : 2*DIM2]
	byReg := groups[2*DIM2 : 3*DIM2]

	for i := 0; i < SZ; i++ {
		row, col, reg := posOf(i)

		byRow[row] = append(byRow[row], i)
		byCol[col] = append(byCol[col], i)
		byReg[reg] = append(byReg[reg], i)
	}

	return groups
}

// NewBoard returns a new empty Board.
func NewBoard() Board {
	// initialize Board
	board := new(Board)
	for i := 0; i < SZ; i++ {
		board[i] = NewCell()
	}
	return *board
}

func NewBoardFromValues(rows [DIM2][DIM2]int) Board {
	b := NewBoard()
	for rowIdx, row := range rows {
		for colIdx, v := range row {
			if v > 0 {
				i := indexOf(rowIdx, colIdx)
				b = b.Set(i, v)
			}
		}
	}
	return b
}

// Values returns the cell values as a 9x9 array.
func (b Board) Values() [DIM2][DIM2]int {
	var values [DIM2][DIM2]int
	for i, cell := range b {
		if cell.Solved() {
			row := i / DIM2
			col := i % DIM2
			values[row][col] = cell.Value()
		}
	}
	return values
}

// ParseBoard returns a new Board by parsing the given string.
func ParseBoard(str string) (*Board, error) {
	// remove whitespace
	str = whitespace.ReplaceAllString(str, "")

	if len(str) != SZ {
		return nil, errors.New("line was incorrect length")
	}
	board := NewBoard()
	for i, rune := range str {
		val, _ := strconv.Atoi(string(rune))
		if val > 0 {
			board = board.Set(i, val)
		}
	}
	if !board.Valid() {
		return nil, errors.New("board is invalid")
	}
	return &board, nil
}

// Valid checks that the value in each cell is legal, and that it does not
// conflict with the value in other cells.
func (b *Board) Valid() bool {
	for _, c := range b {
		if c.Invalid() {
			return false
		}
	}
	return true
}

// Set assigns a value (1-9) to the cell at the given index.
func (b Board) Set(idx int, val int) Board {
	b[idx].Assign(val)
	for _, peerIdx := range Peers[idx] {
		b[peerIdx].Eliminate(val)
	}
	return b
}

func indexOf(row, col int) int {
	return row*DIM2 + col
}

func posOf(idx int) (row, col, region int) {
	row = idx / DIM2
	col = idx % DIM2
	region = (row/DIM)*DIM + (col / DIM)
	return row, col, region
}

// PickUnsolvedCell finds a Cell with the fewest possible values and returns its index.
func (b *Board) pickUnsolvedCell() int {
	idx := -1
	num_possible := DIM2 + 1
	for i, c := range b {
		if !c.Solved() {
			n := c.NumPossible()
			if n < num_possible {
				idx = i
				num_possible = n
			}
		}
	}
	return idx
}

// Solved returns true if the board is solved.
func (b *Board) Solved() bool {
	for _, c := range b {
		if !c.Solved() {
			return false
		}
	}
	return true
}

// Solve will return a solution to this Board. Every valid Board has at least
// one solution. Only the first discovered solution will be returned.
func (b *Board) Solve() Solution {
	start := time.Now()
	b2, _ := b.solve()
	elapsed := time.Since(start)

	return Solution{b, &b2, elapsed}
}

// Solve attempts to find a solution to the given Board.
func (b *Board) solve() (Board, bool) {
	// check groups of peers to look for a value that only appears once
	for _, group := range Groups {
		for v := 1; v <= DIM2; v++ {
			idx := -1  // index of the cell which can be solved
			count := 0 // number of times that v is seen within this group
			for _, c_idx := range group {
				c := b[c_idx]
				if !c.Solved() && c.Possible(v) {
					count++
					if count > 1 {
						// no need to continue looking. v can't be solved in this group
						break
					}
					idx = c_idx
				}
			}
			if count == 1 {
				// v only appeared once, so we can solve it
				b2 := b.Set(idx, v)
				b = &b2
			}
		}
	}

	// we have to make a guess
	idx := b.pickUnsolvedCell()
	if idx == -1 {
		// we must be done!
		return *b, true
	}

	// try each possible value until we find a solution
	c := b[idx]
	for v := 1; v <= DIM2; v++ {
		if c.Possible(v) {
			b2 := b.Set(idx, v)
			b3, valid := b2.solve()
			if valid {
				return b3, true
			}
		}
	}

	return *b, false
}

/*
String returns a nicely formatted representation of this Board. For example:

	1 . .   6 . .   2 . .
	. 6 .   2 . 4   . . .
	. . 2   . . .   8 . 3

	7 5 .   . . 8   . 1 .
	. . .   . 1 .   . . .
	. 4 .   3 . .   . 5 6

	2 . 4   . . .   1 . .
	. . .   5 . 3   . 4 .
	. . 3   . . 9   . . 7
*/
func (b *Board) String() string {
	var buffer bytes.Buffer

	for l := 0; l < DIM; l++ {
		for k := 0; k < DIM; k++ {
			row := l*DIM + k
			for i := 0; i < DIM; i++ {
				for j := 0; j < DIM; j++ {
					col := i*DIM + j
					idx := indexOf(row, col)
					cell := b[idx]
					ch := cell.String(".", "x")
					buffer.WriteString(ch)
					buffer.WriteString(" ")
				}
				buffer.WriteString("  ")
			}
			buffer.WriteString("\n")
		}
		if l < DIM-1 {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

/*
ShortString returns a one-line representation of this Board. For example:

	1..6..2...6.2.4.....2...8.375...8.1.....1.....4.3...562.4...1.....5.3.4...3..9..7
*/
func (b *Board) ShortString() string {
	var buffer bytes.Buffer

	for _, cell := range b {
		ch := cell.String(".", "x")
		buffer.WriteString(ch)
	}
	return buffer.String()
}
