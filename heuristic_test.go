package checkers

import "testing"

func assertHeuristicValue(t *testing.T, h Heuristic, g *Game, player Color, want float64) {
	if got := h(g, player); got != want {
		t.Errorf("heuristic %v fails: want %g got %g", h, want, got)
	}
}

func TestUnweightedCountHeuristic(t *testing.T) {
	b := decodeBoard(`
	  ..@...#.
		.o.o.@.x
		.
		.
		.@.x
	`)
	t.Log("\n" + b.String())
	g := NewCustomGame(CapturesMandatory, BestMandatory, 5, b, WhiteColor)

	// 5 whites - 3 blacks = 2
	assertHeuristicValue(t, UnweightedCountHeuristic, g, WhiteColor, 2)

	assertHeuristicValue(t, UnweightedCountHeuristic, g, BlackColor, -2)
}

func TestWeightedCountHeuristic(t *testing.T) {
	b := decodeBoard(`
	  ..#...@.
		.x.x.#.o
		.
		.
		.#.o
	`)
	t.Log("\n" + b.String())
	g := NewCustomGame(CapturesMandatory, BestMandatory, 5, b, WhiteColor)

	// 2 white pawns + 1 white king - 3 black kings - 2 black pawns
	// 2 + 2 - 6 - 2
	// -4
	assertHeuristicValue(t, WeightedCountHeuristic, g, WhiteColor, -4)

	assertHeuristicValue(t, WeightedCountHeuristic, g, BlackColor, 4)
}
