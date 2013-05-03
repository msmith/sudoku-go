package sudoku

import "strconv"

type Cell struct {
	possibles []bool
	solved    bool
}

func (c *Cell) Assign(val int) {
	c.solved = true
	for v, _ := range c.possibles {
		c.possibles[v] = (v == val-1)
	}
}

func (c *Cell) Value() int {
	for v, p := range c.possibles {
		if p {
			return (v + 1)
		}
	}
	return 0
}

func (c *Cell) Possible(val int) bool {
	return c.possibles[val-1]
}

func (c *Cell) NumPossible() int {
	n := 0
	for _, p := range c.possibles {
		if p {
			n++
		}
	}
	return n
}

func (c *Cell) Eliminate(val int) {
	c.possibles[val-1] = false
}

func (c *Cell) Invalid() bool {
	for _, p := range c.possibles {
		if p {
			return false
		}
	}
	return true
}

func (c *Cell) Copy() Cell {
	d := NewCell()
	d.solved = c.solved
	copy(d.possibles, c.possibles)
	return d
}

func (c *Cell) Solved() bool {
	return c.solved
}

func NewCell() Cell {
	cell := &Cell{make([]bool, DIM2), false}
	for v := 0; v < DIM2; v++ {
		cell.possibles[v] = true
	}
	return *cell
}

func (c *Cell) String(whenUnsolved, whenInvalid string) string {
	if c.Invalid() {
		return "X"
	} else if c.Solved() {
		return strconv.Itoa(c.Value())
	}

	return "."
}
