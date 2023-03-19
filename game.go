package main

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
		return "INVALID gameState"
	}
}

// the part of the state that you need to remember in order to undo
type rememberedState struct {
	state                gameState
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
	pieceCount          PieceCount
	// maybe don't cache the piececount
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

	g.BoardChanged()

	return &g

}

func NewStandardGame(captureRule CaptureRule, bestRule BestRule) *Game {
	return NewCustomGame(captureRule, bestRule, 20, nil, WhiteColor)
}

func (g *Game) IsOver() bool {
	return g.state != playingState
}

func (g *Game) HasWinner() bool {
	return g.state == whiteWonState || g.state == blackWonState
}

func (g *Game) Winner() Color {
	if g.state == whiteWonState {
		return WhiteColor
	}
	return BlackColor
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

func (g *Game) PieceCount() PieceCount {
	return g.pieceCount
}

func (g *Game) DoPly(p Ply) {
	// maybe it'd be safer for DoPly to receive an index into the g.plies slice
	// because there's no guarantee the given Ply is one in the slice,
	// so it can be invalid (e.g. moving from an empty position etc.)
	// but idk if it would work with minimax

	PerformInstructions(g.board, p)
	g.toPlay = g.toPlay.Opposite()

	// save current state in the history
	g.history = append(g.history, g.rememberedState)
	g.lastPly = p

	g.BoardChanged()

	if g.IsOver() {
		return
	}

	isCapture := false
	isPawnMove := false

	for _, ins := range p {
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

	if g.turnsSincePawnMove >= g.stagnantTurnsToDraw && g.turnsSinceCapture >= g.stagnantTurnsToDraw {
		g.state = drawState
	}
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

	// just count again
	// more cpu, less memory usage
	// TODO benchmark?
	g.pieceCount = g.board.PieceCount()

	g.history = g.history[:len(g.history)-1]
}

func (g *Game) Copy() *Game {
	// plies, lastPly, history all shallow-copied
	// board deep-copied
	return &Game{
		captureRule: g.captureRule,
		bestRule:    g.bestRule,
		rememberedState: rememberedState{
			state:   g.state,
			plies:   g.plies,
			lastPly: g.lastPly,
		},
		board:   g.board.Copy(),
		toPlay:  g.toPlay,
		history: g.history,
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
		g.board.Equals(o.board) &&
		sliceEq(g.lastPly, o.lastPly) &&
		pliesEq(g.plies, o.plies)
}

// If you modify the g.Board() manully, you need to call this function afterwards
// in order for the game state to update its internal state according to the updated board
func (g *Game) BoardChanged() {
	count := g.board.PieceCount()
	whiteCount := count.WhiteKings + count.WhitePawns
	blackCount := count.BlackKings + count.BlackPawns

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

	g.pieceCount = count

	// so if the game is over we don't say with the previous' state plies because of the early return
	g.plies = nil

	if g.IsOver() {
		return
	}

	g.plies = GeneratePlies(g.board, g.toPlay, g.captureRule, g.bestRule)
	if len(g.plies) == 0 {
		if g.toPlay == WhiteColor {
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

func inSpecialEnding(c PieceCount) bool {
	wk, wp := c.WhiteKings, c.WhitePawns
	bk, bp := c.BlackKings, c.BlackPawns
	return oneColorSpecialEnding(wk, wp, bk, bp) || oneColorSpecialEnding(bk, bp, wk, wp)
}
