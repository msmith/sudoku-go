package sudoku

import "strconv"

type cell struct {
	possibles [DIM2]bool
	solved    bool
}

func (c *cell) Assign(val int) {
	c.solved = true
	for v, _ := range c.possibles {
		c.possibles[v] = (v == val-1)
	}
}

func (c *cell) Value() int {
	for v, p := range c.possibles {
		if p {
			return (v + 1)
		}
	}
	return 0
}

func (c *cell) Possible(val int) bool {
	return c.possibles[val-1]
}

func (c *cell) NumPossible() int {
	n := 0
	for _, p := range c.possibles {
		if p {
			n++
		}
	}
	return n
}

func (c *cell) Eliminate(val int) {
	c.possibles[val-1] = false
}

func (c *cell) Invalid() bool {
	for _, p := range c.possibles {
		if p {
			return false
		}
	}
	return true
}

func (c *cell) Solved() bool {
	return c.solved
}

func NewCell() cell {
	cell := new(cell)
	for v := 0; v < DIM2; v++ {
		cell.possibles[v] = true
	}
	return *cell
}

func (c *cell) String(whenUnsolved, whenInvalid string) string {
	if c.Invalid() {
		return "X"
	} else if c.Solved() {
		return strconv.Itoa(c.Value())
	}

	return "."
}
