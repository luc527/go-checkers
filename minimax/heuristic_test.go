package minimax

import (
	"testing"

	c "github.com/luc527/go_checkers/core"
)

func assertHeuristicValue(t *testing.T, h Heuristic, g *c.Game, player c.Color, want float64) {
	if got := h(g.Board(), player); got != want {
		t.Errorf("heuristic %v fails: want %g got %g", h, want, got)
	}
}

func TestUnweightedCountHeuristic(t *testing.T) {
	b := c.DecodeBoard(`
	  ..@...#.
		.o.o.@.x
		.
		.
		.@.x
	`)
	t.Log("\n" + b.String())
	g := c.NewCustomGame(c.CapturesMandatory, c.BestMandatory, 5, b, c.WhiteColor)

	// 5 whites - 3 blacks = 2
	assertHeuristicValue(t, UnweightedCountHeuristic, g, c.WhiteColor, 2)

	assertHeuristicValue(t, UnweightedCountHeuristic, g, c.BlackColor, -2)
}

func TestWeightedCountHeuristic(t *testing.T) {
	b := c.DecodeBoard(`
	  ..#...@.
		.x.x.#.o
		.
		.
		.#.o
	`)
	t.Log("\n" + b.String())
	g := c.NewCustomGame(c.CapturesMandatory, c.BestMandatory, 5, b, c.WhiteColor)

	// 2 white pawns + 1 white king - 3 black kings - 2 black pawns
	// 2 + 2 - 6 - 2
	// -4
	assertHeuristicValue(t, WeightedCountHeuristic, g, c.WhiteColor, -4)

	assertHeuristicValue(t, WeightedCountHeuristic, g, c.BlackColor, 4)
}
