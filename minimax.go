package main

import (
	"math"
)

type Minimax struct {
	ToMaximize Color
	Cutoff     int
	Heuristic
}

const (
	drawValue = 0
	winValue  = +1_000_000
	lossValue = -1_000_000
)

// TODO should the receiver be (m *Minimax) or (m Minimax)?
// Minimax is really just a config, it's read-only
// compare performance between the two

func (m Minimax) Search(g *Game) (float64, Ply) {
	return m.searchAt(g, 0)
}

func (m Minimax) searchAt(g *Game, depth int) (float64, Ply) {
	return m.searchImpl(g, depth, math.Inf(-1), math.Inf(1))
}

func (m Minimax) searchImpl(g *Game, depth int, alpha, beta float64) (float64, Ply) {
	if g.IsOver() {
		if !g.HasWinner() {
			return drawValue, nil
		} else if m.ToMaximize == g.Winner() {
			return winValue, nil
		} else {
			return lossValue, nil
		}
	}

	if depth >= m.Cutoff {
		return m.Heuristic(g, m.ToMaximize), nil
	}

	plies := g.Plies()

	if depth == 0 && len(plies) == 1 {
		// No choice, so don't even explore the tree
		return 0, plies[0]
	}

	maximizeTurn := g.ToPlay() == m.ToMaximize

	value := math.Inf(1)
	if maximizeTurn {
		value = math.Inf(-1)
	}

	// no randomness introduced yet

	var ply Ply

	for _, subPly := range plies {
		g.DoPly(subPly)

		subValue, _ := m.searchImpl(g, depth+1, alpha, beta)

		if maximizeTurn {
			if subValue >= value {
				ply = subPly
				value = subValue
			}
			alpha = math.Max(alpha, subValue)
			if subValue >= beta {
				g.UndoLastPly()
				return value, ply
			}
		} else {
			if subValue <= value {
				ply = subPly
				value = subValue
			}
			beta = math.Min(beta, subValue)
			if subValue <= alpha {
				g.UndoLastPly()
				return value, ply
			}
		}

		g.UndoLastPly()
	}

	return value, ply
}
