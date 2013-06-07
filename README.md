# Sudoku Go

A Sudoku solver in [Go](http://golang.org/).

Uses a combination of constraint propagation and search, with an algorithm
that's pretty similar to [Peter Norvig's solver](http://norvig.com/sudoku.html).

This was mostly an excuse for me to learn Go and play with its concurrency
features. Both sequential and parallel solver implementations are included.
They both use the same algoritm, but the parallel version takes advantage of
multi-core CPUs.

# Running the Solvers

There are 3 variants of the solvers. Make sure to set your `GOPATH` before running

    $ export GOPATH=$PWD

## Single puzzle solver

    $ go run sudoku.go puzzles/p1

## Sequential puzzle solver

	$ go run sudoku_sequential.go sets/hard.gz

## Parallel puzzle solver

	$ go run sudoku_parallel.go sets/hard.gz