package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sudoku"
)

const (
	DEFAULT_PORT = "8080"
)

var port string = DEFAULT_PORT

func init() {
	flag.StringVar(&port, "port", DEFAULT_PORT, "Number of concurrent clients")
	flag.Parse()
}

var Running map[string]bool = make(map[string]bool)
var Cache map[string]sudoku.Solution = make(map[string]sudoku.Solution)

func cacheGet(b sudoku.Board) (sudoku.Solution, bool) {
	s, ok := Cache[b.ShortString()]
	return s, ok
}

func cacheSet(s sudoku.Solution) {
	s2 := sudoku.Solution{s.Original, s.Solved, s.Elapsed}
	Cache[s.Original.ShortString()] = s2
}

func markRunning(b string) {
	Running[b] = true
}

func isRunning(b string) bool {
	_, present := Running[b]
	return present
}

func markDone(b string) {
	delete(Running, b)
}

func solve(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	boardString := req.Form.Get("board")
	r := strings.NewReader(boardString)
	br := bufio.NewReader(r)
	b, err := sudoku.ReadBoardLine(br)
	if err != nil {
		panic(err)
	}
	log.Println("Got board:", b.ShortString())

	if !b.Valid() {
		log.Println("Invalid!")
		fmt.Fprintln(res, "Invalid!")
		return
	}

	if isRunning(boardString) {
		log.Println("Already running!")
		fmt.Fprintln(res, "Already running!")
		return
	}

	s, ok := cacheGet(b)
	if ok {
		log.Println("Cache hit!")
	} else {
		log.Println("Cache miss! Solving")
		markRunning(boardString)
		s = b.Solve()
		markDone(boardString)
		cacheSet(s)
	}
	log.Println("Returning:", s.Solved.ShortString())
	fmt.Fprintln(res, s.Solved.ShortString(), s.Elapsed)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/solve", solve)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
