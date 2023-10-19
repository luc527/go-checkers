package conc

import (
	"errors"
	"sync"

	"github.com/luc527/go_checkers/core"
)

type GameState struct {
	Board   core.Board
	ToPlay  core.Color
	Result  core.GameResult
	Plies   []core.Ply
	Version int
}

type Game struct {
	mu    sync.Mutex
	u     *core.Game              // underlying *core.Game
	cs    map[chan GameState]bool // observers
	state GameState               // don't use directly, always call .gameState() (otherwise might get a stale version)
	v     int                     // current iteration of the game (v for "version")
}

func NewConcurrentGame(cr core.CaptureRule, br core.BestRule) *Game {
	u := core.NewGame(cr, br)
	return newConcurrentGame(u)
}

func newConcurrentGame(u *core.Game) *Game {
	g := &Game{
		u:  u,
		cs: make(map[chan GameState]bool),
		v:  1,
	}
	return g
}

func (g *Game) CurrentState() GameState {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.gameState()
}

func (g *Game) NextStates() chan GameState {
	g.mu.Lock()
	defer g.mu.Unlock()

	c := make(chan GameState)
	s := g.gameState()
	if s.Result.Over() {
		close(c)
	} else {
		g.cs[c] = true
	}
	return c
}

func (g *Game) detach(c chan GameState) {
	// check to avoid closing twice (closing a closed channel panics)
	if _, ok := g.cs[c]; ok {
		delete(g.cs, c)
		close(c)
	}
}

func (g *Game) Detach(c chan GameState) {
	if c == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.detach(c)
}

func (g *Game) DetachAll() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for c := range g.cs {
		g.detach(c)
	}
}

func (g *Game) gameState() GameState {
	if g.v != g.state.Version {
		g.state = GameState{
			Board:   *g.u.Board(),
			ToPlay:  g.u.ToPlay(),
			Result:  g.u.Result(),
			Plies:   core.CopyPlies(g.u.Plies()),
			Version: g.v,
		}
	}
	return g.state
}

func (g *Game) DoPly(v int, i int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.gameState().Result.Over() {
		return errors.New("game already over")
	}
	if v != g.v {
		return errors.New("stale game state version")
	}

	plies := g.gameState().Plies
	if i < 0 || i >= len(plies) {
		return errors.New("ply index out of bounds")
	}

	ply := plies[i]
	if _, err := g.u.DoPly(ply); err != nil {
		return err
	}

	g.v++
	state := g.gameState()

	for c := range g.cs {
		go func(c chan GameState, s GameState) {
			c <- s
			if s.Result.Over() {
				g.Detach(c)
			}
		}(c, state)
	}

	return nil
}
