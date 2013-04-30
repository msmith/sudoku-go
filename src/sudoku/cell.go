package sudoku

type Cell struct {
	Possibles []bool
	Solved    bool
}

func (c *Cell) Assign(val int) {
	c.Solved = true
	for v, _ := range c.Possibles {
		c.Possibles[v] = (v == val-1)
	}
}

func (c *Cell) Value() int {
	for v, p := range c.Possibles {
		if p {
			return (v + 1)
		}
	}
	return 0
}

func (c *Cell) Possible(val int) bool {
	return c.Possibles[val-1]
}

func (c *Cell) NumPossible() int {
	n := 0
	for _, p := range c.Possibles {
		if p {
			n++
		}
	}
	return n
}

func (c *Cell) Eliminate(val int) {
	c.Possibles[val-1] = false
}

func (c *Cell) Invalid() bool {
	for _, p := range c.Possibles {
		if p {
			return false
		}
	}
	return true
}

func (c *Cell) Copy() Cell {
	d := NewCell()
	d.Solved = c.Solved
	copy(d.Possibles, c.Possibles)
	return d
}

func NewCell() Cell {
	cell := &Cell{make([]bool, DIM2), false}
	for v := 0; v < DIM2; v++ {
		cell.Possibles[v] = true
	}
	return *cell
}
