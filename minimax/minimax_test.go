package minimax

import (
	"testing"

	c "github.com/luc527/go_checkers/core"
)

func TestDoUndoMinimax(t *testing.T) {
	g := c.NewStandardGame(c.CapturesMandatory, c.BestMandatory)

	whiteMm := Minimax{
		ToMaximize: c.WhiteColor,
		Cutoff:     5,
		Heuristic:  UnweightedCountHeuristic,
	}

	blackMm := Minimax{
		ToMaximize: c.BlackColor,
		Cutoff:     6,
		Heuristic:  WeightedCountHeuristic,
	}

	var states []*c.Game
	var undoInfos []*c.UndoInfo

	for !g.Result().IsOver() {
		states = append(states, g.Copy())

		var ply c.Ply
		if g.ToPlay() == c.WhiteColor {
			_, ply = whiteMm.Search(g)
		} else {
			_, ply = blackMm.Search(g)
		}
		undo, err := g.DoPly(ply)
		if err != nil {
			t.Fail()
		}
		undoInfos = append(undoInfos, undo)
	}

	for i := len(states) - 1; i >= 0; i-- {
		g.UndoPly(undoInfos[i])
		if !g.Equals(states[i]) {
			t.Fail()
			break
		}
	}
}
