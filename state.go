package main

import "fmt"

type gameState byte

const (
	playingState = gameState(iota)
	whiteWonState
	blackWonState
	drawState
)

func (s gameState) String() string {
	switch s {
	case playingState:
		return "playing"
	case whiteWonState:
		return "white won"
	case blackWonState:
		return "black won"
	case drawState:
		return "draw"
	default:
		panic(fmt.Sprintf("Invalid game state: %d", s))
	}
}

// the part of the state that you need to remember in order to undo
type rememberedState struct {
	state               gameState
	plies               []ply
	lastPly             ply
	roundsSinceCapture  int
	roundsSincePawnMove int
}

type game struct {
	captureRule
	bestRule
	rememberedState
	stagnantTurnsToDraw int // stagnant here means no captures and no king moves
	board               *board
	toPlay              color
	history             []rememberedState
}

func newCustomGame(captureRule captureRule, bestRule bestRule, stagnantTurnsToDraw int, initialBoard *board, initalPlayer color) *game {
	var g game

	if initialBoard == nil {
		g.board = new(board)
		placeInitialPieces(g.board)
	} else {
		g.board = initialBoard
	}

	g.captureRule = captureRule
	g.bestRule = bestRule
	g.stagnantTurnsToDraw = stagnantTurnsToDraw

	g.toPlay = initalPlayer

	g.lastPly = nil
	g.roundsSinceCapture = 0
	g.roundsSincePawnMove = 0

	g.boardChanged()

	return &g

}

func newGame(captureRule captureRule, bestRule bestRule) *game {
	return newCustomGame(captureRule, bestRule, 20, nil, whiteColor)
}

func (g *game) isOver() bool {
	return g.state != playingState
}

func (g *game) hasWinner() bool {
	return g.state == whiteWonState || g.state == blackWonState
}

func (g *game) winner() color {
	if g.state == whiteWonState {
		return whiteColor
	}
	return blackColor
}

func (g *game) doPly(p ply) {
	// maybe it'd be safer for doPly to receive an index into the g.plies slice
	// because there's no guarantee the given ply is one in the slice,
	// so it can be invalid (e.g. moving from an empty position etc.)
	// but idk if it would work with minimax

	performInstructions(g.board, p)
	g.toPlay = g.toPlay.opposite()

	// save current state in the history
	g.history = append(g.history, g.rememberedState)
	g.lastPly = p

	g.boardChanged()

	// Draw detection

	isCapture := false
	isPawnMove := false

	for _, ins := range p {
		if ins.t == captureInstruction {
			isCapture = true
		}
		if ins.t == moveInstruction {
			_, kind := g.board.get(ins.row, ins.col)
			if kind == pawnKind {
				isPawnMove = true
			}
		}
	}

	if isCapture {
		g.roundsSinceCapture = 0
	} else {
		g.roundsSinceCapture++
	}

	if isPawnMove {
		g.roundsSincePawnMove = 0
	} else {
		g.roundsSincePawnMove++
	}

	// TODO roundsInSpecialEnding

	if g.roundsSincePawnMove >= g.stagnantTurnsToDraw && g.roundsSinceCapture >= g.stagnantTurnsToDraw {
		g.state = drawState
	}
}

func (g *game) undoLastPly() {
	if g.lastPly == nil {
		return
	}

	undoInstructions(g.board, g.lastPly)
	g.toPlay = g.toPlay.opposite()
	g.rememberedState = g.history[len(g.history)-1]

	g.history = g.history[:len(g.history)-1]
}

func (g *game) Copy() *game {
	// plies, lastPly, history all shallow-copied
	// board deep-copied
	return &game{
		captureRule: g.captureRule,
		bestRule:    g.bestRule,
		rememberedState: rememberedState{
			state:   g.state,
			plies:   g.plies,
			lastPly: g.lastPly,
		},
		board:   g.board.copy(),
		toPlay:  g.toPlay,
		history: g.history,
	}
}

func (g *game) equals(o *game) bool {
	if g == nil && o == nil {
		return true
	}
	if g == nil || o == nil {
		return false
	}

	pliesEq := func(a []ply, b []ply) bool {
		if len(a) != len(b) {
			return false
		}
		for i, x := range a {
			y := b[i]
			if !sliceEq(x, y) {
				return false
			}
		}
		return true
	}

	return g.captureRule == o.captureRule &&
		g.bestRule == o.bestRule &&
		g.state == o.state &&
		g.toPlay == o.toPlay &&
		g.board.equals(o.board) &&
		sliceEq(g.lastPly, o.lastPly) &&
		pliesEq(g.plies, o.plies)
}

func (g *game) boardChanged() {
	count := g.board.pieceCount()
	whiteCount := count.whiteKings + count.whitePawns
	blackCount := count.blackKings + count.blackPawns

	if whiteCount == 0 {
		g.state = blackWonState
	} else if blackCount == 0 {
		g.state = whiteWonState
	}

	if g.isOver() {
		return
	}

	g.plies = generatePlies(g.board, g.toPlay, g.captureRule, g.bestRule)
	if len(g.plies) == 0 {
		if g.toPlay == whiteColor {
			g.state = blackWonState
		} else {
			g.state = whiteWonState
		}
	}
}
