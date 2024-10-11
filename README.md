# Sudoku Go

A Sudoku solver in [Go](http://golang.org/).

The puzzles are solved using constraint propagation and backtracking search, using
an algorithm that's pretty similar to
[Peter Norvig's solver](http://norvig.com/sudoku.html).

This was mostly an excuse to learn Go and play with its concurrency features. Both
sequential and parallel implementations are included. They both use the same
algorithm, but the parallel version takes advantage of multi-core CPUs.

# Running the Solvers

There are 3 variants of the solvers. Make sure to set your `GOPATH` before running.

    $ export GOPATH=$PWD

### Single puzzle solver

    $ go run ./bin/sudoku puzzles/p1

### Multiple puzzle solvers

	$ go run ./bin/sudoku_sequential sets/hard.gz
	$ go run ./bin/sudoku_parallel sets/hard.gz

# Board files

Some example puzzles are included in the `puzzles/` and `sets/` directories.

### Single puzzle files

These are in 9x9 format, where periods (`.`) represent blank cells. Whitespace is
ignored.

    1..6..2..
    .6.2.4...
    ..2...8.3
    75...8.1.
    ....1....
    .4.3...56
    2.4...1..
    ...5.3.4.
    ..3..9..7

These can be solved with the single puzzle solver: **sudoku**

### Puzzle sets

These are gzip'd text files where each 81-character line is a separate puzzle.

    1..6..2...6.2.4.....2...8.375...8.1.....1.....4.3...562.4...1.....5.3.4...3..9..7

These can be solved with either of the batch solvers: **sudoku_sequential** or
**sudoku_parallel**
