package minimax

import (
	"testing"

	c "github.com/luc527/go_checkers/core"
)

func TestDoUndoMinimax(t *testing.T) {
	g := c.NewGame(c.CapturesMandatory, c.BestMandatory)

	whiteMm := DepthLimitedSearcher{
		ToMax:      c.WhiteColor,
		DepthLimit: 5,
		Heuristic:  UnweightedCountHeuristic,
	}

	blackMm := DepthLimitedSearcher{
		ToMax:      c.BlackColor,
		DepthLimit: 6,
		Heuristic:  WeightedCountHeuristic,
	}

	var states []*c.Game
	var undoInfos []*c.UndoInfo

	for !g.Result().Over() {
		states = append(states, g.Copy())

		var ply c.Ply
		if g.ToPlay() == c.WhiteColor {
			ply = whiteMm.Search(g)
		} else {
			ply = blackMm.Search(g)
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
