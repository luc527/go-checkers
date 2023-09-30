package minimax

import (
	"math"

	c "github.com/luc527/go_checkers/core"
)

type Minimax struct {
	ToMaximize c.Color
	Cutoff     int
	Heuristic
}

const (
	drawValue = 0
	winValue  = +1_000_000
	lossValue = -1_000_000
)

func (m *Minimax) Search(g *c.Game) (float64, c.Ply) {
	return m.searchAt(g, 0)
}

func (m *Minimax) searchAt(g *c.Game, depth int) (float64, c.Ply) {
	return m.searchImpl(g, depth, math.Inf(-1), math.Inf(1))
}

func (m *Minimax) searchImpl(g *c.Game, depth int, alpha, beta float64) (float64, c.Ply) {
	state := g.Result()
	if state.IsOver() {
		if !state.HasWinner() {
			return drawValue, nil
		} else if m.ToMaximize == state.Winner() {
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

	var ply c.Ply

	for _, subPly := range plies {
		undoInfo, _ := g.DoPly(subPly)

		subValue, _ := m.searchImpl(g, depth+1, alpha, beta)

		if maximizeTurn {
			if subValue >= value {
				ply = subPly
				value = subValue
			}
			alpha = math.Max(alpha, subValue)
			if subValue >= beta {
				g.UndoPly(undoInfo)
				return value, ply
			}
		} else {
			if subValue <= value {
				ply = subPly
				value = subValue
			}
			beta = math.Min(beta, subValue)
			if subValue <= alpha {
				g.UndoPly(undoInfo)
				return value, ply
			}
		}

		g.UndoPly(undoInfo)
	}

	return value, ply
}
