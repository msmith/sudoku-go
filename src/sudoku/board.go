package sudoku

import (
	"strconv"
	"bytes"
)

const (
	DIM = 3
	DIM2 = DIM * DIM
	SZ = DIM2 * DIM2
)

type Board struct {
	Cells []Cell
}

var Peers [][]int

func init() {
	Peers = make([][]int, SZ)
	for i := 0; i < SZ; i++ {
		row_i := i / DIM2
		col_i := i % DIM2
		region_i := (row_i / DIM) * DIM + (col_i / DIM)

		ps := make([]int, 0, 20)

		for j := 0; j < SZ; j++ {
			row := j / DIM2
			col := j % DIM2
			region := (row / DIM) * DIM + (col / DIM)

			if (j != i) && (row == row_i || col == col_i || region == region_i) {
				ps = append(ps, j)
			}
		}

		Peers[i] = ps
	}
}

func (b *Board) Set(idx int, val int) Board {
	newBoard := b.Copy()
	newBoard.Cells[idx].Assign(val)
	for _, i := range Peers[idx] {
		newBoard.Eliminate(i, val)
	}
	return newBoard
}

func (b *Board) Eliminate(idx int, val int) {
	b.Cells[idx].Eliminate(val)
}

func NewBoard() Board {
	// initialize Board
	board := &Board{ make([]Cell, SZ) }
	for i := 0; i < SZ; i++ {
		board.Cells[i] = NewCell()
	}
	return *board
}

func (b *Board) Copy() Board {
	newBoard := &Board{ make([]Cell, SZ) }
	for i := 0; i < SZ; i++ {
		newBoard.Cells[i] = b.Cells[i].Copy()
	}
	return *newBoard
}

func indexOf(row, col int) int {
	return row * DIM2 + col
}

func posOf(idx int) (row, col int) {
	return (idx / DIM2), (idx % DIM2)
}

func (b *Board) CellAt(row, col int) *Cell {
	idx := indexOf(row, col)
	return &b.Cells[idx]
}

// PickUnsolvedCell finds a Cell with the fewest possible values and returns its index.
func (b *Board) PickUnsolvedCell() int {
	idx := -1
	num_possible := DIM2 + 1
	for i, c := range b.Cells {
		if !c.Solved {
			n := c.NumPossible()
			if (n < num_possible) {
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

// Invalid returns true if at least one cell has 0 possibilities left.
func (b *Board) Invalid() bool {
	for _, c := range b.Cells {
		if c.Invalid() {
			return true
		}
	}
	return false
}

// Solve attempts to find a solution to the given Board.
func (b *Board) Solve() (Board, bool) {
	if b.Solved() {
		return *b, true
	}

	c_idx := b.PickUnsolvedCell()
	c := b.Cells[c_idx]
	for v := 1; v <= DIM2; v++ {
		if c.Possible(v) {
			b2 := b.Set(c_idx, v)
			s, valid := b2.Solve()
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
			row := l * DIM + k
			for i := 0; i < DIM; i++ {
				for j := 0; j < DIM; j++ {
					col := i * DIM + j
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
				if (cell.Possible(v)) {
					ch = strconv.Itoa(v)
				} else {
					if (cell.Solved) {
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