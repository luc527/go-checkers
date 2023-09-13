package go_checkers

import (
	"math/rand"
	"testing"
)

func TestDoUndoState(t *testing.T) {
	g := NewStandardGame(CapturesMandatory, BestNotMandatory)

	var states []*Game
	var undos []UndoInfo

	for !g.Result().IsOver() {
		states = append(states, g.Copy())
		plies := g.Plies()
		r := rand.Int() % len(plies)
		t.Log(plies[r])
		undos = append(undos, g.DoPly(plies[r]))
	}

	t.Log("\n" + g.Board().String())

	for i := len(states) - 1; i >= 0; i-- {
		undo := undos[len(undos)-1]
		undos = undos[:len(undos)-1]
		g.UndoPly(undo)
		if !g.Equals(states[i]) {
			t.Log("\n Failed, expected")
			t.Log(g)
			t.Log("got")
			t.Log(states[i])
			t.Fail()
			break
		}
	}
}

func assertGameResult(t *testing.T, g *Game, want GameResult) {
	got := g.Result()
	if got != want {
		t.Errorf("expected game to be in the %v state, but it's in the %v state", want, got)
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
	assertGameResult(t, g, WhiteWonResult)
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
	assertGameResult(t, g, BlackWonResult)
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
	assertGameResult(t, g, WhiteWonResult)
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
	assertGameResult(t, g, BlackWonResult)
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
	assertGameResult(t, g, PlayingResult)

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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, PlayingResult)
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
	assertGameResult(t, g, DrawResult)
}

func assertSpecialEnding(t *testing.T, b *Board) {
	g := NewCustomGame(CapturesNotMandatory, BestNotMandatory, 20, b, WhiteColor)
	t.Log("\n" + g.Board().String())
	// 1 turn in special ending
	assertGameResult(t, g, PlayingResult)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 2 turns in special ending
	assertGameResult(t, g, PlayingResult)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 3 turns in special ending
	assertGameResult(t, g, PlayingResult)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 4 turns in special ending
	assertGameResult(t, g, PlayingResult)

	g.DoPly(randomInoffensiveMove(g.Board(), g.ToPlay()))
	t.Log("\n" + g.Board().String())
	// 5 turns in special ending
	assertGameResult(t, g, DrawResult)
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

func TestGameResult(t *testing.T) {
	if PlayingResult.IsOver() {
		t.Fail()
	}
	if PlayingResult.HasWinner() {
		t.Fail()
	}

	if !WhiteWonResult.IsOver() {
		t.Fail()
	}
	if !WhiteWonResult.HasWinner() {
		t.Fail()
	}
	if WhiteWonResult.Winner() != WhiteColor {
		t.Fail()
	}

	if !BlackWonResult.IsOver() {
		t.Fail()
	}
	if !BlackWonResult.HasWinner() {
		t.Fail()
	}
	if BlackWonResult.Winner() != BlackColor {
		t.Fail()
	}

	if !DrawResult.IsOver() {
		t.Fail()
	}
	if DrawResult.HasWinner() {
		t.Fail()
	}
}

func TestGameEquals(t *testing.T) {
	nilGame := (*Game)(nil)
	if !nilGame.Equals(nilGame) {
		t.Log("Nil game should be equal to nil game")
		t.Fail()
	}

	g := NewStandardGame(CapturesMandatory, BestNotMandatory)
	if nilGame.Equals(g) || g.Equals(nilGame) {
		t.Log("Nil game should not be equal to actual game")
		t.Fail()
	}

	if !g.Equals(g) {
		t.Log("Game should be equal to itself")
		t.Fail()
	}

	h := g.Copy()
	if !g.Equals(h) {
		t.Log("Game should be equal to a copy of iteself")
		t.Fail()
	}

	undoInfo := h.DoPly(h.Plies()[0])
	if g.Equals(h) {
		t.Log("Game should not be equal after a ply")
		t.Fail()
	}

	h.UndoPly(undoInfo)
	if !g.Equals(h) {
		t.Log("Game should be back to equal after undoing ply")
		t.Fail()
	}
}

func TestGameResultString(t *testing.T) {
	// Just to make the coverage tool happy
	if PlayingResult.String() != "playing" {
		t.Fail()
	}
	if WhiteWonResult.String() != "white won" {
		t.Fail()
	}
	if BlackWonResult.String() != "black won" {
		t.Fail()
	}
	if DrawResult.String() != "draw" {
		t.Fail()
	}
	if GameResult(10).String() != "INVALID GameResult" {
		t.Fail()
	}
}
