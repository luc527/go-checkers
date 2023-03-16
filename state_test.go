package main

import (
	"math/rand"
	"testing"
)

func TestDoUndoState(t *testing.T) {

	// this should already fully test g.history and g.lastPly

	g := newGame(capturesMandatory, bestNotMandatory)

	var states []*game

	for !g.isOver() {
		states = append(states, g.copy())
		r := rand.Int() % len(g.plies)
		t.Log(g.plies[r])
		g.doPly(g.plies[r])
	}

	t.Log("\n" + g.board.String())

	for i := len(states) - 1; i >= 0; i-- {
		g.undoLastPly()
		if !g.equals(states[i]) {
			t.Fail()
			break
		}
	}
}

func assertGameState(t *testing.T, g *game, s gameState) {
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
	g := newCustomGame(capturesNotMandatory, bestNotMandatory, 20, b, whiteColor)
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
	g := newCustomGame(capturesNotMandatory, bestNotMandatory, 20, b, whiteColor)
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
	g := newCustomGame(capturesNotMandatory, bestNotMandatory, 20, b, blackColor)
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
	g := newCustomGame(capturesNotMandatory, bestNotMandatory, 20, b, whiteColor)
	t.Log()
	assertGameState(t, g, blackWonState)
}

func assertEqualBoards(t *testing.T, a *board, b *board) {
	if !a.equals(b) {
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

	g := newCustomGame(capturesMandatory, bestMandatory, 3, b, whiteColor)
	assertGameState(t, g, playingState)

	g.doPly(ply{makeMoveInstruction(3, 3, 2, 2)})

	// just to make the code more legible by showing what each intermediary board looks like
	assertEqualBoards(t, g.board, decodeBoard(`
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

	g.doPly(ply{makeMoveInstruction(0, 6, 1, 5)})
	assertEqualBoards(t, g.board, decodeBoard(`
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

	g.doPly(ply{makeMoveInstruction(4, 4, 6, 2)})
	assertEqualBoards(t, g.board, decodeBoard(`
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
	g.doPly(ply{makeMoveInstruction(2, 2, 2, 4)})
	assertEqualBoards(t, g.board, decodeBoard(`
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
	g.doPly(ply{makeMoveInstruction(1, 5, 3, 3), makeCaptureInstruction(2, 4, whiteColor, pawnKind)})
	assertEqualBoards(t, g.board, decodeBoard(`
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

	g.doPly(ply{makeMoveInstruction(6, 2, 5, 3)})
	assertEqualBoards(t, g.board, decodeBoard(`
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

	g.doPly(ply{makeMoveInstruction(3, 3, 2, 2)})
	assertEqualBoards(t, g.board, decodeBoard(`
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

	g.doPly(ply{makeMoveInstruction(5, 3, 4, 4)})
	assertEqualBoards(t, g.board, decodeBoard(`
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

func assertSpecialEnding(t *testing.T, b *board) {
	g := newCustomGame(capturesNotMandatory, bestNotMandatory, 20, b, whiteColor)
	t.Log("\n" + g.board.String())
	// 1 turn in special ending
	assertGameState(t, g, playingState)

	g.doPly(randomInoffensiveMove(g.board, g.toPlay))
	t.Log("\n" + g.board.String())
	// 2 turns in special ending
	assertGameState(t, g, playingState)

	g.doPly(randomInoffensiveMove(g.board, g.toPlay))
	t.Log("\n" + g.board.String())
	// 3 turns in special ending
	assertGameState(t, g, playingState)

	g.doPly(randomInoffensiveMove(g.board, g.toPlay))
	t.Log("\n" + g.board.String())
	// 4 turns in special ending
	assertGameState(t, g, playingState)

	g.doPly(randomInoffensiveMove(g.board, g.toPlay))
	t.Log("\n" + g.board.String())
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
