package main

import "testing"

const compareRandomActions = 10_000_000

// I wonder how much is the cost of generating random numbers here...
// TODO make up a deterministic stress test that doesn't rely on generating random numbers

func BenchmarkNaiveRandom(b *testing.B) {
	board := newEmptyNaiveBoard()
	nRandomActions(board, compareRandomActions)
}

// naive   ~33s
// compact ~25s
// compact is about 1/3 times faster

func BenchmarkCompactRandom(b *testing.B) {
	board := new(board)
	nRandomActions(board, compareRandomActions)
}
