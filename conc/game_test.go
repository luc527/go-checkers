package conc

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/luc527/go_checkers/core"
)

// TODO: redo tests using CurrState() and NextStates()

func assertMatches(t *testing.T, s GameState, g *core.Game) {
	if !s.Board.Equals(g.Board()) {
		t.Log("boards don't match")
		t.Fail()
	}
	if s.Result != g.Result() {
		t.Log("results don't match")
		t.Fail()
	}
	if s.ToPlay != g.ToPlay() {
		t.Log("current players don't match")
		t.Fail()
	}
	if !core.PliesEquals(s.Plies, g.Plies()) {
		t.Log("plies don't match")
		t.Fail()
	}
}

func receiveState(t *testing.T, o <-chan GameState) (GameState, bool) {
	select {
	case <-time.After(5 * time.Second):
		t.Log("failed to receive from channel in time")
		t.FailNow()
		return GameState{}, false
	case s, ok := <-o:
		return s, ok
	}
}

func assertHasPendingState(t *testing.T, o <-chan GameState) GameState {
	s, ok := receiveState(t, o)
	if !ok {
		t.Log("expected channel to be open")
		t.FailNow()
	}
	return s
}

func assertClosed(t *testing.T, o <-chan GameState) {
	_, ok := receiveState(t, o)
	if ok {
		t.Log("expected channel to be open")
		t.FailNow()
	}
}

func TestAttachDetach(t *testing.T) {
	g := newConcurrentGame(core.NewStandardGame())
	var s GameState

	s = g.CurrentState()
	assertMatches(t, s, g.u)

	c := g.NextStates()
	if err := g.DoPly(s.Version, 0); err != nil {
		t.Log(err)
		t.FailNow()
	}

	s = assertHasPendingState(t, c)
	assertMatches(t, s, g.u)

	g.Detach(c)
	assertClosed(t, c)
}

func TestAttachDetachAll(t *testing.T) {
	g := newConcurrentGame(core.NewStandardGame())

	const n = 10

	os := make([]chan GameState, n)
	for i := 0; i < n; i++ {
		o := g.NextStates()
		os[i] = o
	}

	if err := g.DoPly(1, 0); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, o := range os {
		s := assertHasPendingState(t, o)
		assertMatches(t, s, g.u)
	}

	// TODO: not receiving in time after DetachAll
	g.DetachAll()

	for _, o := range os {
		assertClosed(t, o)
	}
}

func TestPlayUntilOver(t *testing.T) {
	g := newConcurrentGame(core.NewStandardGame())

	o := g.NextStates()

	s := g.CurrentState()
	r := rand.Intn(len(s.Plies))
	g.DoPly(s.Version, r)

	i := 0
	for {
		i++
		if i == 1000 {
			t.Log("game took too long to finish")
			t.FailNow()
		}

		s = assertHasPendingState(t, o)
		assertMatches(t, s, g.u)

		if s.Result.Over() {
			break
		}

		r := rand.Intn(len(s.Plies))
		g.DoPly(s.Version, r)
	}

	assertClosed(t, o)
}

func TestPlyErrors(t *testing.T) {
	g := newConcurrentGame(core.NewStandardGame())
	o := g.NextStates()

	// Don't really care what the states are here,
	// but don't want a goroutine leak
	go func() {
		for range o {
		}
	}()

	var err error

	err = g.DoPly(5, 0)
	if err == nil || !strings.Contains(err.Error(), "stale") {
		t.Log("expected stale version error")
		t.Fail()
	}

	err = g.DoPly(1, -4)
	if err == nil || !strings.Contains(err.Error(), "bounds") {
		t.Log("expected out of bounds ply error")
		t.Fail()
	}

	err = g.DoPly(1, 200)
	if err == nil || !strings.Contains(err.Error(), "bounds") {
		t.Log("expected out of bounds ply error")
		t.Fail()
	}

	g.Detach(o)
}

func TestConcurrentObservers(t *testing.T) {
	g := newConcurrentGame(core.NewStandardGame())

	const n = 8

	seqC := make(chan []GameState, n)

	for i := 0; i < n; i++ {
		o := g.NextStates()
		go func(o chan GameState) {
			var seq []GameState

			s := g.CurrentState()
			r := rand.Intn(len(s.Plies))
			g.DoPly(s.Version, r)

			for s := range o {
				seq = append(seq, s)
				if s.Result.Over() {
					break
				}

				r := rand.Intn(len(s.Plies))
				g.DoPly(s.Version, r)

				ms := 0 + rand.Intn(40)
				<-time.After(time.Duration(ms * int(time.Millisecond)))
			}
			seqC <- seq
		}(o)
	}

	var seqs [][]GameState
	for i := 0; i < n; i++ {
		seqs = append(seqs, <-seqC)
	}

	// All observers must observe the game states IN ORDER
	// We will NOT guarantee that every observer receives the FIRST game state
	// That doesn't seem possible given how Attach works, and it seems very unlikely going to be a problem

	for _, seq := range seqs {
		if len(seq) == 0 {
			continue
		}
		for _, s := range seq {
			fmt.Print(s.Version, " ")
		}
		fmt.Println()
		prev := seq[0].Version
		for _, s := range seq[1:] {
			curr := s.Version
			if curr != prev+1 {
				t.FailNow()
			}
			prev = curr
		}
	}
}
