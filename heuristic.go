package main

import (
	"fmt"
	"reflect"
	"runtime"
)

// the heuristics take a game and not just a board
// because the game caches the piece count
// and some heuristics rely on the piece count
// -- maybe not a very nice abstraction

type Heuristic func(g *Game, player Color) int

func (h Heuristic) String() string {
	return fmt.Sprintf("%q", runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
}

var _ Heuristic = UnweightedCountHeuristic
var _ Heuristic = WeightedCountHeuristic

func UnweightedCountHeuristic(g *Game, player Color) int {
	count := g.PieceCount()
	whites := int(count.WhitePawns + count.WhiteKings)
	blacks := int(count.BlackPawns + count.BlackKings)

	fmt.Println("whites", whites)
	fmt.Println("blacks", blacks)

	factor := 1
	if player == BlackColor {
		factor = -1
	}

	return factor * (whites - blacks)
}

func WeightedCountHeuristic(g *Game, player Color) int {
	const (
		pawnWeight = 1
		kingWeight = 2
	)

	count := g.PieceCount()
	whites := int(count.WhitePawns*pawnWeight + count.WhiteKings*kingWeight)
	blacks := int(count.BlackPawns*pawnWeight + count.BlackKings*kingWeight)

	factor := 1
	if player == BlackColor {
		factor = -1
	}

	return factor * (whites - blacks)
}

// TODO distance heuristic

// TODO "clusters" heuristic but simpler to compute, i.e. only look at neighbours

// TODO more general heuristic where you give a weight to each tile
// and return a heuristic that uses that weight map
// (distance heuristic is a specific version)
