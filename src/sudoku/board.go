package sudoku

import (
	"bytes"
	"strconv"
	"time"
)

const (
	DIM  = 3
	DIM2 = DIM * DIM
	SZ   = DIM2 * DIM2
)

type Board struct {
	Cells []Cell
}

// Peers is a lookup table that provides the peer cells for each cell
var Peers [][]int = make([][]int, SZ)

// Groups contains cell indexes for each row, column, and region
var Groups [][]int = make([][]int, DIM2*3)

// initialize Peers & Groups
func init() {
	byRow := Groups[0:DIM2]
	byCol := Groups[DIM2 : 2*DIM2]
	byReg := Groups[2*DIM2 : 3*DIM2]

	for i := 0; i < SZ; i++ {
		row, col, reg := posOf(i)

		byRow[row] = append(byRow[row], i)
		byCol[col] = append(byCol[col], i)
		byReg[reg] = append(byReg[reg], i)

		for j := 0; j < SZ; j++ {
			if j != i {
				row_j, col_j, reg_j := posOf(j)
				if row_j == row || col_j == col || reg_j == reg {
					Peers[i] = append(Peers[i], j)
				}
			}
		}
	}
}

func (b *Board) Set(idx int, val int) Board {
	newBoard := b.Copy()
	newBoard.Cells[idx].Assign(val)
	for _, i := range Peers[idx] {
		newBoard.eliminate(i, val)
	}
	return newBoard
}

func (b *Board) eliminate(idx int, val int) {
	b.Cells[idx].Eliminate(val)
}

func NewBoard() Board {
	// initialize Board
	board := &Board{make([]Cell, SZ)}
	for i := 0; i < SZ; i++ {
		board.Cells[i] = NewCell()
	}
	return *board
}

func (b *Board) Copy() Board {
	newBoard := &Board{make([]Cell, SZ)}
	for i := 0; i < SZ; i++ {
		newBoard.Cells[i] = b.Cells[i].Copy()
	}
	return *newBoard
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
func (b *Board) PickUnsolvedCell() int {
	idx := -1
	num_possible := DIM2 + 1
	for i, c := range b.Cells {
		if !c.Solved {
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
	for i := 0; i < SZ; i++ {
		c := b.Cells[i]
		if !c.Solved {
			return false
		}
	}
	return true
}

func (b *Board) Solve() Solution {
	start := time.Now()
	b2, _ := b.solve()
	t := time.Since(start)

	return Solution{b, &b2, t}
}

// Solve attempts to find a solution to the given Board.
func (b *Board) solve() (Board, bool) {
	if b.Solved() {
		return *b, true
	}

	// check groups of peers to look for a value that only appears once
	for _, group := range Groups {
		idx := -1
		for v := 1; v <= DIM2; v++ {
			count := 0
			for _, c_idx := range group {
				c := b.Cells[c_idx]
				if !c.Solved && c.Possible(v) {
					count++
					if count == 1 {
						idx = c_idx
					}
				}
			}
			if count == 1 {
				b2 := b.Set(idx, v)
				s, valid := b2.solve()
				if valid {
					return s, true
				} else {
					return s, false
				}
			}
		}
	}

	// guess
	c_idx := b.PickUnsolvedCell()
	c := b.Cells[c_idx]
	for v := 1; v <= DIM2; v++ {
		if c.Possible(v) {
			b2 := b.Set(c_idx, v)
			s, valid := b2.solve()
			if valid {
				return s, true
			}
		}
	}
	return *b, false
}

func (b *Board) String() string {
	var buffer bytes.Buffer

	for l := 0; l < DIM; l++ {
		for k := 0; k < DIM; k++ {
			row := l*DIM + k
			for i := 0; i < DIM; i++ {
				for j := 0; j < DIM; j++ {
					col := i*DIM + j
					idx := indexOf(row, col)
					cell := b.Cells[idx]
					var ch string
					if cell.Invalid() {
						ch = "X"
					} else if cell.Solved {
						ch = strconv.Itoa(cell.Value())
					} else {
						ch = "."
					}
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

func (b *Board) ShortString() string {
	var buffer bytes.Buffer

	for _, cell := range b.Cells {
		var ch string
		if cell.Invalid() {
			ch = "X"
		} else if cell.Solved {
			ch = strconv.Itoa(cell.Value())
		} else {
			ch = "."
		}
		buffer.WriteString(ch)
	}
	return buffer.String()
}

func (b *Board) DebugString() string {
	var buffer bytes.Buffer
	for row := 0; row < DIM2; row++ {
		for col := 0; col < DIM2; col++ {
			idx := indexOf(row, col)
			cell := b.Cells[idx]
			var ch string
			for v := 1; v <= DIM2; v++ {
				if cell.Possible(v) {
					ch = strconv.Itoa(v)
				} else {
					if cell.Solved {
						ch = "-"
					} else {
						ch = "."
					}
				}
				buffer.WriteString(ch)
			}
			buffer.WriteString(" ")
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}
