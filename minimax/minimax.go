package minimax

import (
	"math"
	"math/rand"
	"time"

	c "github.com/luc527/go_checkers/core"
)

const (
	MinTimeLimit = 100 * time.Millisecond
	MaxTimeLimit = 10 * time.Second
)

// A Searcher searches for the best ply in the game tree of the given game.
type Searcher interface {
	Search(g *c.Game) c.Ply
}

// A DepthLimitedSearcher is a Searcher that stops searching the game tree
// when it reaches a certain depth. It evaluates the leaves according
// to the game result with regards to the player to maximize, and to
// a given heuristic function.
type DepthLimitedSearcher struct {
	ToMax c.Color
	Heuristic
	DepthLimit int
}

var _ Searcher = DepthLimitedSearcher{}

func (s DepthLimitedSearcher) Search(g *c.Game) c.Ply {
	ctx := searchContext{toMax: s.ToMax, h: s.Heuristic, stop: nil}
	_, ply := ctx.search(g, s.DepthLimit, math.Inf(-1), math.Inf(1))
	return ply
}

// A TimeLimitedSearcher is a Searcher that stops searching the game tree
// after a certain amount of time has elapsed. It uses iterative
// deepening search.
type TimeLimitedSearcher struct {
	ToMax c.Color
	Heuristic
	TimeLimit time.Duration
}

var _ Searcher = TimeLimitedSearcher{}

func (s TimeLimitedSearcher) Search(g *c.Game) c.Ply {
	tlim := s.TimeLimit
	if tlim < MinTimeLimit {
		tlim = MinTimeLimit
	}
	if tlim > MaxTimeLimit {
		tlim = MaxTimeLimit
	}

	stopTime := time.Now().Add(tlim)

	var ply c.Ply
	ctx := searchContext{toMax: s.ToMax, h: s.Heuristic, stop: closeAfter(tlim)}
	for dlim := 1; ; dlim++ {
		// We can only assign the result of a search (variable ply0) to the best known ply so far (variable ply)
		// if the ply0 search went all the way to the end. Otherwise, it's possible that the search has
		// stopped in a node at an early depth in the tree, *which might have a large heuristic value
		// without it actually being a good move*
		searchStart := time.Now()
		_, ply0 := ctx.search(g, dlim, math.Inf(-1), math.Inf(1))
		if ctx.cancelled() {
			break
		}
		ply = ply0

		searchDuration := time.Since(searchStart)
		timeLeft := time.Until(stopTime)
		delta := searchDuration - timeLeft
		if delta >= 0 || -delta < 100*time.Millisecond {
			break
		}
	}

	return ply
}

func closeAfter(d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(ch)
	}()
	return ch
}

type searchContext struct {
	toMax c.Color
	stop  <-chan struct{}
	h     Heuristic
}

func (ctx searchContext) cancelled() bool {
	select {
	case <-ctx.stop:
		return true
	default:
		return false
	}
}

const (
	drawValue = 0
	winValue  = +1_000_000
	lossValue = -1_000_000
)

func (ctx searchContext) search(g *c.Game, depthLeft int, alpha float64, beta float64) (float64, c.Ply) {
	res := g.Result()
	if res.Over() {
		if !res.HasWinner() {
			return drawValue, nil
		} else if ctx.toMax == res.Winner() {
			return winValue, nil
		} else {
			return lossValue, nil
		}
	}
	if ctx.cancelled() || depthLeft <= 0 {
		return ctx.h(g, ctx.toMax), nil
	}

	plies := g.Plies()
	plies = shuffle(plies)

	maximizeTurn := g.ToPlay() == ctx.toMax

	value := math.Inf(1)
	if maximizeTurn {
		value = math.Inf(-1)
	}

	var ply c.Ply

	// TODO: check whether the alpha-beta search is messing with the algorithm somehow
	// (like making one AI that should be better than the other always lose)

	for _, subPly := range plies {
		undoInfo, _ := g.DoPly(subPly)

		subValue, _ := ctx.search(g, depthLeft-1, alpha, beta)

		if maximizeTurn {
			if subValue > value {
				ply, value = subPly, subValue
				alpha = math.Max(alpha, subValue)
			}
			if subValue >= beta {
				g.UndoPly(undoInfo)
				return value, ply
			}
		} else {
			if subValue <= value {
				ply, value = subPly, subValue
				beta = math.Min(beta, subValue)
			}
			if subValue <= alpha {
				g.UndoPly(undoInfo)
				return value, ply
			}
		}

		g.UndoPly(undoInfo)
	}

	return value, ply
}

func shuffle(a []c.Ply) []c.Ply {
	n := len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
	return a
}
