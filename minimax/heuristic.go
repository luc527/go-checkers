package minimax

import (
	"fmt"
	"reflect"
	"runtime"

	c "github.com/luc527/go_checkers/core"
)

type Heuristic func(b *c.Board, player c.Color) float64

func (h Heuristic) String() string {
	return fmt.Sprintf("%q", runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
}

func HeuristicFromString(s string) Heuristic {
	switch s {
	case "UnweightedCount":
		return UnweightedCountHeuristic
	case "WeightedCount":
		return WeightedCountHeuristic
	default:
		return nil
	}
}

var _ Heuristic = UnweightedCountHeuristic
var _ Heuristic = WeightedCountHeuristic

func UnweightedCountHeuristic(b *c.Board, player c.Color) float64 {
	count := b.PieceCount()
	whites := int(count.WhitePawns + count.WhiteKings)
	blacks := int(count.BlackPawns + count.BlackKings)

	factor := 1
	if player == c.BlackColor {
		factor = -1
	}

	return float64(factor * (whites - blacks))
}

func WeightedCountHeuristic(b *c.Board, player c.Color) float64 {
	const (
		pawnWeight = 1
		kingWeight = 2
	)

	count := b.PieceCount()
	whites := int(count.WhitePawns*pawnWeight + count.WhiteKings*kingWeight)
	blacks := int(count.BlackPawns*pawnWeight + count.BlackKings*kingWeight)

	factor := 1
	if player == c.BlackColor {
		factor = -1
	}

	return float64(factor * (whites - blacks))
}

// TODO distance heuristic

// TODO "clusters" heuristic but simpler to compute, i.e. only look at neighbours

// TODO more general heuristic where you give a weight to each tile
// and return a heuristic that uses that weight map
// (distance heuristic is a specific version)
