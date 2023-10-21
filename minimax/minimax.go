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
	Search(g *c.Game) int
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

func (s DepthLimitedSearcher) Search(g *c.Game) int {
	ctx := searchContext{toMax: s.ToMax, h: s.Heuristic, timedCloser: nil}
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

func (s TimeLimitedSearcher) Search(g *c.Game) int {
	tlim := s.TimeLimit
	if tlim < MinTimeLimit {
		tlim = MinTimeLimit
	}
	if tlim > MaxTimeLimit {
		tlim = MaxTimeLimit
	}

	stopTime := time.Now().Add(tlim)

	var ply int
	ctx := searchContext{toMax: s.ToMax, h: s.Heuristic, timedCloser: closeAfter(tlim)}
	for dlim := 1; ; dlim++ {
		// We can only assign the result of a search (variable ply0) to the best known ply so far (variable ply)
		// if the ply0 search went all the way to the end. Otherwise, it's possible that the search has
		// stopped in a node at an early depth in the tree, *which might have a large heuristic value
		// without it actually being a good move*
		searchStart := time.Now()
		_, ply0 := ctx.search(g, dlim, math.Inf(-1), math.Inf(1))
		if ctx.closed() {
			break
		}
		ply = ply0

		searchDuration := time.Since(searchStart)
		timeLeft := time.Until(stopTime)
		if searchDuration >= timeLeft {
			break
		}
	}

	return ply
}

type timedCloser <-chan struct{}

func closeAfter(d time.Duration) timedCloser {
	ch := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(ch)
	}()
	return ch
}

type searchContext struct {
	toMax c.Color
	timedCloser
	h Heuristic
}

func (c timedCloser) closed() bool {
	select {
	case <-c:
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

func (ctx searchContext) search(g *c.Game, depthLeft int, alpha float64, beta float64) (float64, int) {
	res := g.Result()
	if res.Over() {
		if !res.HasWinner() {
			return drawValue, 0
		} else if ctx.toMax == res.Winner() {
			return winValue, 0
		} else {
			return lossValue, 0
		}
	}
	if ctx.closed() || depthLeft <= 0 {
		return ctx.h(g.Board(), ctx.toMax), 0
	}

	plies := g.Plies()
	plies = shuffle(plies)

	maximizeTurn := g.ToPlay() == ctx.toMax

	value := math.Inf(1)
	if maximizeTurn {
		value = math.Inf(-1)
	}

	var plyIndex int

	// TODO: check whether the alpha-beta search is messing with the algorithm somehow
	// (like making one AI that should be better than the other always lose)

	for subIndex, subPly := range plies {
		undoInfo, _ := g.DoPly(subPly)

		subValue, _ := ctx.search(g, depthLeft-1, alpha, beta)

		if maximizeTurn {
			if subValue > value {
				plyIndex, value = subIndex, subValue
				alpha = math.Max(alpha, subValue)
			}
			if subValue >= beta {
				g.UndoPly(undoInfo)
				return value, plyIndex
			}
		} else {
			if subValue <= value {
				plyIndex, value = subIndex, subValue
				beta = math.Min(beta, subValue)
			}
			if subValue <= alpha {
				g.UndoPly(undoInfo)
				return value, plyIndex
			}
		}

		g.UndoPly(undoInfo)
	}

	return value, plyIndex
}

func shuffle(a []c.Ply) []c.Ply {
	n := len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
	return a
}
