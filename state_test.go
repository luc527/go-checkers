package main

import (
	"math/rand"
	"testing"
)

func TestDoUndoState(t *testing.T) {
	g := newGame(capturesMandatory, bestNotMandatory)

	var states []*game

	for !g.isOver() {
		states = append(states, g.deepCopy())
		r := rand.Int() % len(g.plies)
		t.Log(g.plies[r])
		g.doPly(g.plies[r])
	}

	t.Log("\n" + g.board.String())

	for i := len(states) - 1; i >= 0; i-- {
		g.undoLastPly()
		if !g.equals(states[i]) {
			t.Fail()
			break
		}
	}
}
