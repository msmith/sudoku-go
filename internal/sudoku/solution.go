package sudoku

import (
	"bytes"
	"time"
)

type Solution struct {
	Original *Board
	Solved   *Board
	Elapsed  time.Duration
}

func (s *Solution) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(s.Original.ShortString())
	buffer.WriteString("\n")
	buffer.WriteString(s.Solved.ShortString())
	buffer.WriteString(" ")
	buffer.WriteString(s.Elapsed.String())
	return buffer.String()
}
