package checkers

import "testing"

func TestDoUndoMinimax(t *testing.T) {
	g := NewStandardGame(CapturesMandatory, BestMandatory)

	whiteMm := Minimax{
		ToMaximize: WhiteColor,
		Cutoff:     5,
		Heuristic:  UnweightedCountHeuristic,
	}

	blackMm := Minimax{
		ToMaximize: BlackColor,
		Cutoff:     6,
		Heuristic:  WeightedCountHeuristic,
	}

	var states []*Game
	var undoInfos []UndoInfo

	for !g.Result().IsOver() {
		states = append(states, g.Copy())

		var ply Ply
		if g.ToPlay() == WhiteColor {
			_, ply = whiteMm.Search(g)
		} else {
			_, ply = blackMm.Search(g)
		}
		undoInfos = append(undoInfos, g.DoPly(ply))
	}

	for i := len(states) - 1; i >= 0; i-- {
		g.UndoPly(undoInfos[i])
		if !g.Equals(states[i]) {
			t.Fail()
			break
		}
	}
}
