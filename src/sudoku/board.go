package sudoku

import (
	"bytes"
	"time"
)

const (
	DIM  = 3
	DIM2 = DIM * DIM
	SZ   = DIM2 * DIM2
)

type Board struct {
	cells [SZ]cell
}

// Peers is a lookup table that provides the peer cells of each cell
var Peers [][]int = createPeers()

// Groups contains the cell indexes of each group (row, column, and region)
var Groups [][]int = createGroups()

func createPeers() [][]int {
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

func createGroups() [][]int {
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

func (b *Board) Valid() bool {
	for _, c := range b.cells {
		if c.Invalid() {
			return false
		}
	}
	return true
}

func (b Board) Set(idx int, val int) Board {
	b.cells[idx].Assign(val)
	for _, i := range Peers[idx] {
		b.eliminate(i, val)
	}
	return b
}

func (b *Board) eliminate(idx int, val int) {
	b.cells[idx].Eliminate(val)
}

func NewBoard() Board {
	// initialize Board
	board := new(Board)
	for i := 0; i < SZ; i++ {
		board.cells[i] = NewCell()
	}
	return *board
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
	for i, c := range b.cells {
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
	for _, c := range b.cells {
		if !c.Solved() {
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
	// check groups of peers to look for a value that only appears once
	for _, group := range Groups {
		for v := 1; v <= DIM2; v++ {
			idx := -1  // index of the cell which can be solved
			count := 0 // number of times that v is seen within this group
			for _, c_idx := range group {
				c := b.cells[c_idx]
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
	c := b.cells[idx]
	for v := 1; v <= DIM2; v++ {
		if c.Possible(v) {
			b2 := b.Set(idx, v)
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
					cell := b.cells[idx]
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

func (b *Board) ShortString() string {
	var buffer bytes.Buffer

	for _, cell := range b.cells {
		ch := cell.String(".", "x")
		buffer.WriteString(ch)
	}
	return buffer.String()
}
