package main

import (
	"encoding/json"
	"fmt"
)

type GameResult byte

const (
	PlayingResult = GameResult(iota)
	WhiteWonResult
	BlackWonResult
	DrawResult
)

func (s GameResult) IsOver() bool {
	return s != PlayingResult
}

func (s GameResult) HasWinner() bool {
	return s == WhiteWonResult || s == BlackWonResult
}

func (s GameResult) Winner() Color {
	if s == WhiteWonResult {
		return WhiteColor
	} else {
		return BlackColor
	}
}

func (s GameResult) String() string {
	switch s {
	case PlayingResult:
		return "playing"
	case WhiteWonResult:
		return "white won"
	case BlackWonResult:
		return "black won"
	case DrawResult:
		return "draw"
	default:
		return "INVALID GameResult"
	}
}

type gameState struct {
	toPlay               Color
	turnsSinceCapture    int16
	turnsSincePawnMove   int16
	turnsInSpecialEnding int16
	plies                []Ply
}

type UndoInfo struct {
	plyDone   Ply
	prevState gameState
}

type Game struct {
	stagnantTurnsToDraw int16 // stagnant here means no captures and no pawn moves
	captureRule         CaptureRule
	bestRule            BestRule
	board               *Board
	state               gameState
}

// TODO with the `json:""` stuff it's not necessary
// for the actual golang struct field names to be
// abbreviated

type GameMessageOut struct {
	B   string `json:"b"`
	P   Color  `json:"p"`
	Tsc int16  `json:"tsc"`
	Tsp int16  `json:"tsp"`
	Tis int16  `json:"tis"`
}

type GameMessageIn struct {
	GameMessageOut
	Ttd int16       `json:"ttd"`
	Cr  CaptureRule `json:"cr"`
	Br  BestRule    `json:"br"`
}

func (g *Game) String() string {
	return fmt.Sprintf(
		"{ToPlay: %v, turnsSinceCapture: %v, turnsSincePawnMove: %v, turnsInSpecialEnding: %v, Board:\n%v\n}",
		g.state.toPlay,
		g.state.turnsSinceCapture,
		g.state.turnsSincePawnMove,
		g.state.turnsInSpecialEnding,
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

	g.state.toPlay = initalPlayer

	g.state.turnsSinceCapture = 0
	g.state.turnsSincePawnMove = 0
	g.state.turnsInSpecialEnding = 0
	// once we get in a special ending turnsInSpecialEnding becomes 1 and increases each turn

	g.BoardChanged(nil)

	return &g
}

func NewStandardGame(captureRule CaptureRule, bestRule BestRule) *Game {
	return NewCustomGame(captureRule, bestRule, 20, nil, WhiteColor)
}

func (g *Game) SerializeOut() (string, error) {
	msg := GameMessageOut{
		B:   g.board.Serialize(),
		P:   g.state.toPlay,
		Tsc: g.state.turnsSinceCapture,
		Tsp: g.state.turnsSincePawnMove,
		Tis: g.state.turnsInSpecialEnding,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("serialize game: %v", err)
	}
	return string(bytes), err
}

func UnserializeGameIn(s string) (*Game, error) {
	in := &GameMessageIn{}
	err := json.Unmarshal([]byte(s), in)
	if err != nil {
		err = fmt.Errorf("unserialize game: %v", err)
		return nil, err
	}

	g := new(Game)

	g.board, err = UnserializeBoard(in.B)
	if err != nil {
		err = fmt.Errorf("unserialize game: %v", err)
		return nil, err
	}

	if in.Tsc < 0 || in.Tsp < 0 || in.Tis < 0 || in.Ttd < 0 {
		return nil, fmt.Errorf("unserialize game: negative number of turns")
	}

	g.state.toPlay = in.P
	g.state.turnsSinceCapture = in.Tsc
	g.state.turnsSincePawnMove = in.Tsp
	g.state.turnsInSpecialEnding = in.Tis
	g.stagnantTurnsToDraw = in.Ttd
	g.captureRule = in.Cr
	g.bestRule = in.Br
	return g, nil
}

func (g *Game) Board() *Board {
	return g.board
}

func (g *Game) ToPlay() Color {
	return g.state.toPlay
}

func (g *Game) DoPly(p Ply) UndoInfo {
	PerformInstructions(g.board, p)

	prevState := g.state
	g.state.toPlay = g.state.toPlay.Opposite()
	g.BoardChanged(p)

	return UndoInfo{plyDone: p, prevState: prevState}
}

func (g *Game) Result() GameResult {
	count := g.board.PieceCount()
	whiteCount := count.WhiteKings + count.WhitePawns
	blackCount := count.BlackKings + count.BlackPawns

	if whiteCount == 0 {
		return BlackWonResult
	} else if blackCount == 0 {
		return WhiteWonResult
	}

	if g.state.turnsInSpecialEnding == 5 {
		return DrawResult
	}

	if g.state.turnsSincePawnMove >= g.stagnantTurnsToDraw && g.state.turnsSinceCapture >= g.stagnantTurnsToDraw {
		return DrawResult
	}

	if len(g.Plies()) == 0 {
		if g.state.toPlay == WhiteColor {
			return BlackWonResult
		} else {
			return WhiteWonResult
		}
	}

	return PlayingResult
}

func (g *Game) UndoPly(undo UndoInfo) {
	UndoInstructions(g.board, undo.plyDone)
	g.state = undo.prevState
}

func (g *Game) Copy() *Game {
	// plies shallow-copied
	// board deep-copied
	return &Game{
		state: gameState{
			toPlay:               g.state.toPlay,
			turnsSinceCapture:    g.state.turnsSinceCapture,
			turnsSincePawnMove:   g.state.turnsSincePawnMove,
			turnsInSpecialEnding: g.state.turnsInSpecialEnding,
			plies:                g.state.plies,
		},
		stagnantTurnsToDraw: g.stagnantTurnsToDraw,
		captureRule:         g.captureRule,
		bestRule:            g.bestRule,
		board:               g.board.Copy(),
	}
}

func (g *Game) Equals(o *Game) bool {
	if g == nil && o == nil {
		return true
	}
	if g == nil || o == nil {
		return false
	}

	return g.captureRule == o.captureRule &&
		g.bestRule == o.bestRule &&
		g.state.toPlay == o.state.toPlay &&
		g.state.turnsInSpecialEnding == o.state.turnsInSpecialEnding &&
		g.state.turnsSinceCapture == o.state.turnsSinceCapture &&
		g.state.turnsSincePawnMove == o.state.turnsSincePawnMove &&
		g.board.Equals(o.board)
}

func (g *Game) BoardChanged(ply Ply) {
	count := g.board.PieceCount()

	if inSpecialEnding(count) {
		g.state.turnsInSpecialEnding++
	} else {
		g.state.turnsInSpecialEnding = 0
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
			g.state.turnsSinceCapture = 0
		} else {
			g.state.turnsSinceCapture++
		}

		if isPawnMove {
			g.state.turnsSincePawnMove = 0
		} else {
			g.state.turnsSincePawnMove++
		}
	}

	g.state.plies = nil
}

func (g *Game) generatePlies() []Ply {
	return GeneratePlies(make([]Ply, 0, 10), g.board, g.state.toPlay, g.captureRule, g.bestRule)
}

func (g *Game) Plies() []Ply {
	// Generated on demand, then cached
	if g.state.plies == nil {
		g.state.plies = g.generatePlies()
	}
	return g.state.plies
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
