package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/msmith/sudoku-go/internal/sudoku"
)

const (
	DEFAULT_PORT  = "8080"
	CACHE_SECONDS = 24 * 60 * 60 // 24 hours
	SUDOKU_PATH   = "/sudoku/"
)

var port string = DEFAULT_PORT

func init() {
	flag.StringVar(&port, "port", DEFAULT_PORT, "HTTP server port")
	flag.Parse()
}

type CellValues [sudoku.DIM2][sudoku.DIM2]int

func httpError(w http.ResponseWriter, msg string, err error, code int) {
	http.Error(w, msg, code)
	if err != nil {
		msg = msg + ": " + err.Error()
	}
	log.Println(msg)
}

/*
Expects a POST of form data containing the unsolved `board` data as a JSON-encoded array of 81 ints.

Returns the solved board as a JSON-encoded array of 81 ints.
*/
func solve(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.Header().Add("Allow", "GET")
		httpError(w, req.Method+" is not allowed", nil, http.StatusMethodNotAllowed)
		return
	}

	path := req.URL.Path
	boardString := path[len(SUDOKU_PATH):]

	b, err := sudoku.ParseBoard(boardString)
	if err != nil {
		httpError(w, "Failed to parse board", err, http.StatusBadRequest)
		return
	}
	log.Println("Got board:", b.ShortString())

	if !b.Valid() {
		httpError(w, "Invalid board", nil, http.StatusBadRequest)
		return
	}

	solution := b.Solve()
	log.Println("Returning:", solution.Solved.ShortString())

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "max-age="+strconv.Itoa(CACHE_SECONDS))
	bytes, err := json.Marshal(solution.Solved.Values())
	if err != nil {
		httpError(w, "Failed to marshal JSON response", err, http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc(SUDOKU_PATH, solve)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
