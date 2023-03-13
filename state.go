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
	state                gameState
	plies                []ply
	lastPly              ply
	turnsSinceCapture    int
	turnsSincePawnMove   int
	turnsInSpecialEnding int
	// TODO turn into int8 maybe ^
}

type game struct {
	captureRule
	bestRule
	rememberedState
	stagnantTurnsToDraw int // stagnant here means no captures and no pawn moves
	board               *board
	toPlay              color
	history             []rememberedState
}

// TODO? make game proper ADT; no direct access to fields
// idk might not be very idiomatic in Go

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
	g.turnsSinceCapture = 0
	g.turnsSincePawnMove = 0
	g.turnsInSpecialEnding = 0
	// once we get in a special ending turnsInSpecialEnding becomes 1 and increases each turn

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

	if g.isOver() {
		return
	}

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
		g.turnsSinceCapture = 0
	} else {
		g.turnsSinceCapture++
	}

	if isPawnMove {
		g.turnsSincePawnMove = 0
	} else {
		g.turnsSincePawnMove++
	}

	if g.turnsSincePawnMove >= g.stagnantTurnsToDraw && g.turnsSinceCapture >= g.stagnantTurnsToDraw {
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

func (g *game) copy() *game {
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

	if inSpecialEnding(count) {
		g.turnsInSpecialEnding++
		if g.turnsInSpecialEnding == 5 {
			g.state = drawState
		}
	} else {
		g.turnsInSpecialEnding = 0
	}

	// so if the game is over we don't say with the previous' state plies because of the early return
	g.plies = nil

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

func oneColorSpecialEnding(ourKings, ourPawns, theirKings, theirPawns int8) bool {
	// a) 2 damas vs 2 damas
	// b) 2 damas vs 1 dama
	// c) 2 damas vs 1 dama e 1 pedra
	// d) 1 dama  vs 1 dama
	// e) 1 dama  vs 1 dama e 1 pedra
	//    ^ our   vs ^ their
	if ourPawns > 0 {
		return false
	}
	if ourKings == 2 {
		return (theirPawns == 0 && (theirKings == 2 || theirKings == 1)) || // a ou b
			(theirPawns == 1 && theirKings == 1) // c
	}
	if ourKings == 1 {
		return theirKings == 1 && (theirPawns == 0 || theirPawns == 1) // d or e
	}
	return false

	// let's check whether:
	// once we get in a special ending any further capture still leaves us in another special ending

	// a -> b by losing 1 king
	// b -> (win) by losing 1 king
	// b -> d by losing 1 king
	// c -> e by losing 1 king
	// c -> b by losing 1 pawn
	// c -> 2 damas vs 1 pedra, not an special ending!
	// d -> (win) by losing either king
	// e -> d by losing 1 pawn
	// e -> 1 dama vs 1 pedra, again not an special ending!

	// this means we need to check every time
}

func inSpecialEnding(c pieceCount) bool {
	wk, wp := c.whiteKings, c.whitePawns
	bk, bp := c.blackKings, c.blackPawns
	return oneColorSpecialEnding(wk, wp, bk, bp) || oneColorSpecialEnding(bk, bp, wk, wp)
}
