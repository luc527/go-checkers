package main

import (
	"math/rand"
	"testing"
)

func TestDoUndoState(t *testing.T) {

	// this should already fully test g.history and g.lastPly

	g := NewStandardGame(CapturesMandatory, BestNotMandatory)

	var states []*Game

	for !g.IsOver() {
		states = append(states, g.Copy())
		r := rand.Int() % len(g.Plies())
		t.Log(g.Plies()[r])
		g.DoPly(g.Plies()[r])
	}

	t.Log("\n" + g.Board().String())

	for i := len(states) - 1; i >= 0; i-- {
		g.UndoLastPly()
		if !g.Equals(states[i]) {
			t.Fail()
			break
		}
	}
}

func assertGameState(t *testing.T, g *Game, s gameState) {
	if g.state != s {
		t.Errorf("expected game to be in the %v state, but it's in the %v state", s, g.state)
	}
}

func TestWhiteWinsByNoBlackPieces(t *testing.T) {
	b := decodeBoard(`
		.
		...@
		.....o
		...o
	`)
	t.Log("\n" + b.String())
	g := NewCustomGame(CapturesNotMandatory, BestNotMandatory, 20, b, WhiteColor)
	assertGameState(t, g, whiteWonState)
}

func TestBlackWinsByNoWhitePieces(t *testing.T) {
	b := decodeBoard(`
		.x....#
		....x..
		..#....
		.
		.
		.
	`)
	t.Log("\n" + b.String())
	g := NewCustomGame(CapturesNotMandatory, BestNotMandatory, 20, b, WhiteColor)
	assertGameState(t, g, blackWonState)
}

func TestWhiteWinsByNoBlackPlies(t *testing.T) {
	b := decodeBoard(`
		....x
		...@.o
		..o...o
		.
		.
		.
		.x
		x.o
	`)
	t.Log("\n" + b.String())
	g := NewCustomGame(CapturesNotMandatory, BestNotMandatory, 20, b, BlackColor)
	assertGameState(t, g, whiteWonState)
}

func TestBlackWinsByNoWhitePlies(t *testing.T) {
	b := decodeBoard(`
		x.o
		.o
		.
		...x...x
		....x.#
		.....o
	`)
	t.Log("\n" + b.String())
	g := NewCustomGame(CapturesNotMandatory, BestNotMandatory, 20, b, WhiteColor)
	t.Log()
	assertGameState(t, g, blackWonState)
}

func assertEqualBoards(t *testing.T, a *Board, b *Board) {
	if !a.Equals(b) {
		t.Errorf("boards not equal: \n%vand\n%v", a, b)
	}
}

func TestDrawByNoCaptureNorKingMovesForNTurns(t *testing.T) {
	// too hard to make up 20 turns that will result in a draw
	// considering we also have to test edge cases to be sure
	// so let's just make N a parameter with default value 20

	// we put all those white pawns at the bottom just we don't accidentaly get into a special ending

	b := decodeBoard(`
	  ..x...#
		.
		.
		...o.o
		....@
		.
		.
		ooooooo
	`)

	g := NewCustomGame(CapturesMandatory, BestMandatory, 3, b, WhiteColor)
	assertGameState(t, g, playingState)

	g.DoPly(Ply{MoveInstruction(3, 3, 2, 2)})

	// just to make the code more legible by showing what each intermediary board looks like
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x...#
		.
		..o
		.....o
		....@
		.
		.
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// at this point turnsSincePawnMove=0, turnsSinceCapture=1

	g.DoPly(Ply{MoveInstruction(0, 6, 1, 5)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.....#
		..o
		.....o
		....@
		.
		.
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// at this point turnsSincePawnMove=1, turnsSinceCapture=2

	g.DoPly(Ply{MoveInstruction(4, 4, 6, 2)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.....#
		..o
		.....o
		.
		.
		..@
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// at this point turnsSincePawnMove=2, turnsSinceCapture=3

	// let's reset a counter, ply doesn't have to be legal
	g.DoPly(Ply{MoveInstruction(2, 2, 2, 4)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.....#
		....o
		.....o
		.
		.
		..@
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// at this point turnsSincePawnMove=0, turnsSinceCapture=4

	// let's reset another counter
	g.DoPly(Ply{MoveInstruction(1, 5, 3, 3), CaptureInstruction(2, 4, WhiteColor, PawnKind)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.
		.
		...#.o
		.
		.
		..@
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// at this point turnsSincePawnMove=1, turnsSinceCapture=0

	// now let's keep the state stagnant

	g.DoPly(Ply{MoveInstruction(6, 2, 5, 3)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.
		.
		...#.o
		.
		...@
		.
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// turnsSincePawnMove=2, turnsSinceCapture=1

	g.DoPly(Ply{MoveInstruction(3, 3, 2, 2)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.
		..#
		.....o
		.
		...@
		.
		ooooooo
	`))
	assertGameState(t, g, playingState)
	// turnsSincePawnMove=3, turnsSinceCapture=2

	g.DoPly(Ply{MoveInstruction(5, 3, 4, 4)})
	assertEqualBoards(t, g.Board(), decodeBoard(`
	  ..x
		.
		..#
		.....o
		....@
		.
		.
		ooooooo
	`))
	// turnsSincePawnMove=3, turnsSinceCapture=2, should draw now!
	assertGameState(t, g, drawState)
}

func assertSpecialEnding(t *testing.T, b *Board) {
	g := NewCustomGame(CapturesNotMandatory, BestNotMandatory, 20, b, WhiteColor)
	t.Log("\n" + g.Board().String())
	// 1 turn in special ending
	assertGameState(t, g, playingState)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 2 turns in special ending
	assertGameState(t, g, playingState)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 3 turns in special ending
	assertGameState(t, g, playingState)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 4 turns in special ending
	assertGameState(t, g, playingState)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 5 turns in special ending
	assertGameState(t, g, drawState)
}

func TestDrawBySpecialEnding(t *testing.T) {
	//a
	assertSpecialEnding(t, decodeBoard(`
	  ..@
		....@
		.
		.....#
		.#
	`))

	//b
	assertSpecialEnding(t, decodeBoard(`
	  ..@
		.
		.....#
		.#
	`))

	//c
	assertSpecialEnding(t, decodeBoard(`
	  ..@
		....@
		.
		.....x
		.#
	`))

	//d
	assertSpecialEnding(t, decodeBoard(`
	  ..@
		.
		.#
	`))

	//e
	assertSpecialEnding(t, decodeBoard(`
	  ..@.x
		.
		.#
	`))
}

func BenchmarkGame(b *testing.B) {
	trials := 100_000
	for t := 0; t < trials; t++ {
		g := NewStandardGame(CapturesNotMandatory, BestNotMandatory)
		for !g.IsOver() {
			plies := g.Plies()
			randomPly := plies[rand.Int()%len(plies)]
			g.DoPly(randomPly)
		}
		for g.HasLastPly() {
			g.UndoLastPly()
		}
	}

	// _M means _ mandatory
	// _X means _ NOT mandatory

	// BM probably slower because it allocates a new plies array
	// to which it adds the best ones

	// before IsRowEmpty
	// CM, BM: 13.42s
	// CM, BX: 10.09s
	// CX, BX: 29.67s
	// after IsRowEmpty
	// CM, BM: 12.38s
	// CM, BX: 10.46s
	// CX, BX: 28.71s

	// now with the correct implementation it's only marginally faster :(
}
