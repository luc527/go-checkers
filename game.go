package main

import "fmt"

type GameState byte

const (
	PlayingState = GameState(iota)
	WhiteWonState
	BlackWonState
	DrawState
)

func (s GameState) IsOver() bool {
	return s != PlayingState
}

func (s GameState) HasWinner() bool {
	return s == WhiteWonState || s == BlackWonState
}

func (s GameState) Winner() Color {
	if s == WhiteWonState {
		return WhiteColor
	} else {
		return BlackColor
	}
}

func (s GameState) String() string {
	switch s {
	case PlayingState:
		return "playing"
	case WhiteWonState:
		return "white won"
	case BlackWonState:
		return "black won"
	case DrawState:
		return "draw"
	default:
		return "INVALID GameState"
	}
}

// the part of the state that you need to remember in order to undo
type rememberedState struct {
	plies                []Ply
	lastPly              Ply
	turnsSinceCapture    int16
	turnsSincePawnMove   int16
	turnsInSpecialEnding int16
}

type Game struct {
	rememberedState
	stagnantTurnsToDraw int16 // stagnant here means no captures and no pawn moves
	captureRule         CaptureRule
	bestRule            BestRule
	board               *Board
	toPlay              Color
	history             []rememberedState
}

func (g *Game) String() string {
	return fmt.Sprintf(
		"{ToPlay: %v, LastPly: %v, Board:\n%v\n}",
		g.toPlay,
		g.lastPly,
		g.board,
	)
}

func NewCustomGame(captureRule CaptureRule, bestRule BestRule, stagnantTurnsToDraw int16, initialBoard *Board, initalPlayer Color) *Game {
	var g Game

	if initialBoard == nil {
		g.board = new(Board)
		PlaceInitialPieces(g.board)
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

	g.BoardChanged(nil)

	return &g

}

func NewStandardGame(captureRule CaptureRule, bestRule BestRule) *Game {
	return NewCustomGame(captureRule, bestRule, 20, nil, WhiteColor)
}

func (g *Game) Board() *Board {
	return g.board
}

func (g *Game) Plies() []Ply {
	return g.plies
}

func (g *Game) ToPlay() Color {
	return g.toPlay
}

func (g *Game) DoPly(p Ply) {
	PerformInstructions(g.board, p)
	g.toPlay = g.toPlay.Opposite()

	// save current state in the history
	g.history = append(g.history, g.rememberedState)
	g.lastPly = p

	g.BoardChanged(p)
}

func (g *Game) ComputeState() GameState {
	count := g.board.PieceCount()
	whiteCount := count.WhiteKings + count.WhitePawns
	blackCount := count.BlackKings + count.BlackPawns

	if whiteCount == 0 {
		return BlackWonState
	} else if blackCount == 0 {
		return WhiteWonState
	}

	if len(g.plies) == 0 {
		if g.toPlay == WhiteColor {
			return BlackWonState
		} else {
			return WhiteWonState
		}
	}

	if g.turnsInSpecialEnding == 5 {
		return DrawState
	}

	if g.turnsSincePawnMove >= g.stagnantTurnsToDraw && g.turnsSinceCapture >= g.stagnantTurnsToDraw {
		return DrawState
	}

	return PlayingState
}

func (g *Game) HasLastPly() bool {
	return len(g.history) > 0
}

// TODO maybe should return err if there's no lastply
func (g *Game) UndoLastPly() {
	if g.lastPly == nil {
		return
	}

	UndoInstructions(g.board, g.lastPly)
	g.toPlay = g.toPlay.Opposite()
	g.rememberedState = g.history[len(g.history)-1]

	g.history = g.history[:len(g.history)-1]
}

// These two copy methods have kind of a crappy/leaky interface

func (g *Game) Copy() *Game {
	// plies, lastPly, history all shallow-copied
	// board deep-copied
	return &Game{
		rememberedState: rememberedState{
			plies:   g.plies,
			lastPly: g.lastPly,
		},
		stagnantTurnsToDraw: g.stagnantTurnsToDraw,
		captureRule:         g.captureRule,
		bestRule:            g.bestRule,
		board:               g.board.Copy(),
		toPlay:              g.toPlay,
		history:             g.history,
	}
}

func (g *Game) Equals(o *Game) bool {
	if g == nil && o == nil {
		return true
	}
	if g == nil || o == nil {
		return false
	}

	pliesEq := func(a []Ply, b []Ply) bool {
		if len(a) != len(b) {
			return false
		}
		n := len(a)
		for i := 0; i < n; i++ {
			if !a[i].Equals(b[i]) {
				return false
			}
		}
		return true
	}

	return g.captureRule == o.captureRule &&
		g.bestRule == o.bestRule &&
		g.toPlay == o.toPlay &&
		g.board.Equals(o.board) &&
		g.lastPly.Equals(o.lastPly) &&
		pliesEq(g.plies, o.plies)
}

func (g *Game) BoardChanged(ply Ply) {
	count := g.board.PieceCount()

	if inSpecialEnding(count) {
		g.turnsInSpecialEnding++
	} else {
		g.turnsInSpecialEnding = 0
	}

	if ply != nil {
		isCapture := false
		isPawnMove := false

		for _, ins := range ply {
			if ins.t == captureInstruction {
				isCapture = true
			}
			if ins.t == moveInstruction {
				_, kind := g.board.Get(ins.row, ins.col)
				if kind == PawnKind {
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
	}

	g.plies = g.generatePlies()
}

func (g *Game) generatePlies() []Ply {
	return GeneratePlies(make([]Ply, 0, 10), g.board, g.toPlay, g.captureRule, g.bestRule)
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

func inSpecialEnding(c PieceCount) bool {
	wk, wp := c.WhiteKings, c.WhitePawns
	bk, bp := c.BlackKings, c.BlackPawns
	return oneColorSpecialEnding(wk, wp, bk, bp) || oneColorSpecialEnding(bk, bp, wk, wp)
}
