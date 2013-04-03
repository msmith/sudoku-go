# Sudoku Go

A Sudoku solver in [Go](http://golang.org/).

# Puzzle files

Two puzzle file formats are supported:

## Single puzzle file

A single puzzle. Each row is on a separate line of 9 chars. Some samples are included in __puzzles/*__

    1..6..2..
    .6.2.4...
    ..2...8.3
    75...8.1.
    ....1....
    .4.3...56
    2.4...1..
    ...5.3.4.
    ..3..9..7

## Puzzle set

A set of multiple puzzles. Each puzzle is defined on a separate line of 81 chars. The resulting file is gzipped. Some samples are included in __sets/*__

    000000010400000000020000000000050407008000300001090000300400200050100000000806000
    000000010400000000020000000000050604008000300001090000300400200050100000000807000

# Building

	  $ export GOPATH=$PWD

# Single puzzle solver

Solves a single puzzle

	  $ go run sudoku.go puzzles/p1

# Multiple puzzle solver

Solves multiple puzzles sequentially

	  $ export GOPATH=$PWD
	  $ go run sudoku_batch.go set/hard.gz

# Parallel puzzle solver

Uses multiple goroutines to solve puzzles in parallel

	  $ export GOPATH=$PWD
	  $ go run sudoku_batch.go set/hard.gz