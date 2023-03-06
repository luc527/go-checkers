package main

import (
	"strings"
	"testing"
)

func compareGeneratedMoves(
	got [][]instruction,
	want [][]instruction,
) (
	extra []string,
	missing []string,
) {

	gotMap := make(map[string]bool)
	for _, is := range got {
		gotMap[instructionsToString(is)] = true
	}

	wantMap := make(map[string]bool)
	for _, is := range want {
		wantMap[instructionsToString(is)] = true
	}

	for is := range gotMap {
		if _, ok := wantMap[is]; !ok {
			extra = append(extra, is)
		}
	}

	for is := range wantMap {
		if _, ok := gotMap[is]; !ok {
			missing = append(missing, is)
		}
	}

	return
}

func assertEqualMoves(t *testing.T, got [][]instruction, want [][]instruction) {
	extra, missing := compareGeneratedMoves(got, want)
	if len(extra) > 0 {
		t.Errorf("generated extra instruction lists:\n%s", strings.Join(extra, "\n"))
	}
	if len(missing) > 0 {
		t.Errorf("missing instruction lists:\n%s", strings.Join(missing, "\n"))
	}
}

func TestSimplePawnMove(t *testing.T) {
	b := new(board)

	b.set(1, 1, blackColor, pawnKind)
	// no possibilities
	b.set(0, 0, blackColor, pawnKind)
	// one possibility occupied
	b.set(0, 2, blackColor, pawnKind)
	// against the (left and right) walls
	b.set(1, 7, blackColor, pawnKind)
	b.set(3, 0, blackColor, pawnKind)
	// crowning
	b.set(6, 6, blackColor, pawnKind)

	//
	// black
	//

	b.set(6, 1, whiteColor, pawnKind)
	// no possibilities
	b.set(7, 0, whiteColor, pawnKind)
	// one possibility occupied
	b.set(7, 2, whiteColor, pawnKind)
	// against the (left and right) walls
	b.set(4, 0, whiteColor, pawnKind)
	b.set(4, 7, whiteColor, pawnKind)
	// crowning
	b.set(1, 5, whiteColor, pawnKind)

	t.Log("\n" + b.String())

	blackMovesGot := generateSimpleMoves(b, blackColor)

	blackMovesWant := [][]instruction{
		{makeMoveInstruction(1, 1, 2, 2)},
		{makeMoveInstruction(1, 1, 2, 0)},
		{makeMoveInstruction(0, 2, 1, 3)},
		{makeMoveInstruction(1, 7, 2, 6)},
		{makeMoveInstruction(3, 0, 4, 1)},
		{
			makeMoveInstruction(6, 6, 7, 7),
			makeCrownInstruction(7, 7),
		},
		{
			makeMoveInstruction(6, 6, 7, 5),
			makeCrownInstruction(7, 5),
		},
	}
	assertEqualMoves(t, blackMovesGot, blackMovesWant)

	whiteMovesWant := [][]instruction{
		{makeMoveInstruction(6, 1, 5, 0)},
		{makeMoveInstruction(6, 1, 5, 2)},
		{makeMoveInstruction(7, 2, 6, 3)},
		{makeMoveInstruction(4, 0, 3, 1)},
		{makeMoveInstruction(4, 7, 3, 6)},
		{
			makeMoveInstruction(1, 5, 0, 6),
			makeCrownInstruction(0, 6),
		},
		{
			makeMoveInstruction(1, 5, 0, 4),
			makeCrownInstruction(0, 4),
		},
	}
	whiteMovesGot := generateSimpleMoves(b, whiteColor)

	assertEqualMoves(t, whiteMovesGot, whiteMovesWant)
}

func TestSimpleKingMove(t *testing.T) {

	b := new(board)

	//white
	b.set(5, 5, whiteColor, kingKind)
	b.set(0, 7, whiteColor, kingKind)

	//black
	b.set(2, 2, blackColor, kingKind)

	t.Log("\n" + b.String())

	whiteMovesGot := generateSimpleMoves(b, whiteColor)

	whiteMovesWant := [][]instruction{
		//
		// moving the white king at (5, 5)
		//
		// down, right
		{makeMoveInstruction(5, 5, 6, 6)},
		{makeMoveInstruction(5, 5, 7, 7)},
		// down, left
		{makeMoveInstruction(5, 5, 6, 4)},
		{makeMoveInstruction(5, 5, 7, 3)},
		// up, left
		{makeMoveInstruction(5, 5, 4, 4)},
		{makeMoveInstruction(5, 5, 3, 3)}, // gets stopped by the black king (no (2,2), (1,1), (0,0))
		// up, right
		{makeMoveInstruction(5, 5, 4, 6)},
		{makeMoveInstruction(5, 5, 3, 7)},
		//
		// moving the white king at (0, 7)
		//
		{makeMoveInstruction(0, 7, 1, 6)},
		{makeMoveInstruction(0, 7, 2, 5)},
		{makeMoveInstruction(0, 7, 3, 4)},
		{makeMoveInstruction(0, 7, 4, 3)},
		{makeMoveInstruction(0, 7, 5, 2)},
		{makeMoveInstruction(0, 7, 6, 1)},
		{makeMoveInstruction(0, 7, 7, 0)},
	}

	assertEqualMoves(t, whiteMovesGot, whiteMovesWant)

	blackMovesWant := [][]instruction{
		//
		// moving the black king at (2, 2)
		//
		// down, right
		{makeMoveInstruction(2, 2, 3, 3)},
		{makeMoveInstruction(2, 2, 4, 4)}, // gets stopped by the white king (no (5,5) etc.)
		// down, left
		{makeMoveInstruction(2, 2, 3, 1)},
		{makeMoveInstruction(2, 2, 4, 0)},
		// up, right
		{makeMoveInstruction(2, 2, 1, 3)},
		{makeMoveInstruction(2, 2, 0, 4)},
		// up, left
		{makeMoveInstruction(2, 2, 1, 1)},
		{makeMoveInstruction(2, 2, 0, 0)},
	}
	blackMovesGot := generateSimpleMoves(b, blackColor)

	assertEqualMoves(t, blackMovesGot, blackMovesWant)
}

func TestCapturePawnMove(t *testing.T) {
	b := new(board)

	b.set(4, 6, whiteColor, pawnKind)
	b.set(3, 5, blackColor, pawnKind)
	b.set(3, 3, blackColor, pawnKind)
	b.set(1, 3, blackColor, pawnKind)
	b.set(1, 1, whiteColor, pawnKind)

	b.set(5, 3, blackColor, pawnKind)
	b.set(6, 4, whiteColor, pawnKind)

	t.Log("\n" + b.String())

	blackMovesGot := generateCaptureMoves(b, blackColor)

	blackMovesWant := [][]instruction{
		{
			makeMoveInstruction(3, 5, 5, 7),
			makeCaptureInstruction(4, 6, whiteColor, pawnKind),
		},
		{
			makeMoveInstruction(5, 3, 7, 5),
			makeCaptureInstruction(6, 4, whiteColor, pawnKind),
			makeCrownInstruction(7, 5),
		},
	}

	assertEqualMoves(t, blackMovesGot, blackMovesWant)

	whiteMovesWant := [][]instruction{
		{
			makeMoveInstruction(4, 6, 2, 4),
			makeCaptureInstruction(3, 5, blackColor, pawnKind),
			makeMoveInstruction(2, 4, 0, 2),
			makeCaptureInstruction(1, 3, blackColor, pawnKind),
			makeCrownInstruction(0, 2),
		},
		{
			makeMoveInstruction(4, 6, 2, 4),
			makeCaptureInstruction(3, 5, blackColor, pawnKind),
			makeMoveInstruction(2, 4, 4, 2),
			makeCaptureInstruction(3, 3, blackColor, pawnKind),
		},
		{
			makeMoveInstruction(6, 4, 4, 2),
			makeCaptureInstruction(5, 3, blackColor, pawnKind),
			makeMoveInstruction(4, 2, 2, 4),
			makeCaptureInstruction(3, 3, blackColor, pawnKind),
			makeMoveInstruction(2, 4, 0, 2),
			makeCaptureInstruction(1, 3, blackColor, pawnKind),
			makeCrownInstruction(0, 2),
		},
	}
	whiteMovesGot := generateCaptureMoves(b, whiteColor)

	assertEqualMoves(t, whiteMovesGot, whiteMovesWant)
}

func TestCaptureThroughCrowningRowDoesntCrown(t *testing.T) {
	b := new(board)

	// passes through (0, 3), with 0 being the crowning row for white pieces,
	// but it shouldn't crown because it doesn't *end* at that position,
	// just goes through it
	b.set(4, 7, whiteColor, pawnKind)
	b.set(3, 6, blackColor, pawnKind)
	b.set(1, 4, blackColor, pawnKind)
	b.set(1, 2, blackColor, pawnKind)

	t.Log("\n" + b.String())

	movesWant := [][]instruction{
		{
			makeMoveInstruction(4, 7, 2, 5),
			makeCaptureInstruction(3, 6, blackColor, pawnKind),
			makeMoveInstruction(2, 5, 0, 3),
			makeCaptureInstruction(1, 4, blackColor, pawnKind),
			makeMoveInstruction(0, 3, 2, 1),
			makeCaptureInstruction(1, 2, blackColor, pawnKind),
		},
	}
	movesGot := generateCaptureMoves(b, whiteColor)

	assertEqualMoves(t, movesGot, movesWant)
}

func TestCaptureKingMoveOneDiagonal(t *testing.T) {
	b := new(board)

	b.set(3, 3, whiteColor, kingKind)
	b.set(5, 5, blackColor, pawnKind)

	movesGot := generateCaptureMoves(b, whiteColor)
	movesWant := [][]instruction{
		{
			makeMoveInstruction(3, 3, 6, 6),
			makeCaptureInstruction(5, 5, blackColor, pawnKind),
		},
		{
			makeMoveInstruction(3, 3, 7, 7),
			makeCaptureInstruction(5, 5, blackColor, pawnKind),
		},
	}

	assertEqualMoves(t, movesGot, movesWant)
}

// TODO use some examples from http://www.damasciencias.com.br/regras-jogo-de-damas/ as tests
