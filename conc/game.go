package conc

import (
	"errors"
	"fmt"
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

func NewConcurrentGame() *Game {
	u := core.NewGame()
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

func (g *Game) doPlyInner(ply core.Ply) error {
	if _, err := g.u.DoPly(ply); err != nil {
		return fmt.Errorf("do ply: %v", err)
	}

	g.v++
	s1 := g.gameState()

	for c := range g.cs {
		go func(c chan GameState, s GameState) {
			c <- s
			if s.Result.Over() {
				g.Detach(c)
			}
		}(c, s1)
	}

	return nil
}

func (g *Game) validatePly(player core.Color, v int) error {
	s := g.gameState()
	if s.Result.Over() {
		return errors.New("do ply: game already over")
	}
	if v != g.v {
		return errors.New("do ply: stale game state version")
	}
	if s.ToPlay != player {
		return errors.New("do ply: not your turn")
	}
	return nil
}

func (g *Game) DoPlyGiven(player core.Color, v int, ply core.Ply) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.validatePly(player, v); err != nil {
		return err
	}
	if err := g.doPlyInner(ply); err != nil {
		return fmt.Errorf("do ply: %v", err)
	}
	return nil
}

func (g *Game) DoPlyIndex(player core.Color, v int, i int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err := g.validatePly(player, v); err != nil {
		return err
	}

	s := g.gameState()
	plies := s.Plies
	if i < 0 || i >= len(plies) {
		return errors.New("do ply: ply index out of bounds")
	}
	if err := g.doPlyInner(plies[i]); err != nil {
		return fmt.Errorf("do ply: %v", err)
	}

	return nil
}

func (g *Game) UnderlyingGame() *core.Game {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.u.Copy()
}
