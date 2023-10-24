package minimax

import (
	"math/rand"
	"testing"
	"time"

	"github.com/luc527/go_checkers/core"
	c "github.com/luc527/go_checkers/core"
)

func TestDoUndoMinimax(t *testing.T) {
	g := c.NewGame()

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

		var ply core.Ply
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

func TestCloseAfter(t *testing.T) {
	c := closeAfter(200 * time.Millisecond)

	<-time.After(100 * time.Millisecond)
	if c.closed() {
		t.Fail()
	}

	<-time.After(200 * time.Millisecond)
	if !c.closed() {
		t.Fail()
	}
}

func TestTimeLimitedSearcher(t *testing.T) {
	ai := TimeLimitedSearcher{
		ToMax:     c.BlackColor,
		Heuristic: UnweightedCountHeuristic,
		TimeLimit: 100 * time.Millisecond,
	}
	g := c.NewGame()

	sig := make(chan struct{})

	ply := make(chan core.Ply)
	go func() {
		for range sig {
			ply <- ai.Search(g)
		}
	}()

	for {
		plies := g.Plies()
		r := rand.Intn(len(plies))
		g.DoPly(plies[r])

		if g.Result().Over() {
			close(sig)
			break
		}

		sig <- struct{}{}
		select {
		case <-time.After(150 * time.Millisecond):
			t.Logf("Time limited searcher took too long!")
			t.Fail()
		case p := <-ply:
			g.DoPly(p)
		}

		if g.Result().Over() {
			close(sig)
			break
		}
	}
}
