package main

import (
	"strings"
	"testing"
)

func compareGeneratedInstructions(
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

func assertEqualInstructionLists(t *testing.T, got [][]instruction, want [][]instruction) {
	extra, missing := compareGeneratedInstructions(got, want)
	if len(extra) > 0 {
		t.Errorf("generated extra instruction lists:\n%s", strings.Join(extra, "\n"))
	}
	if len(missing) > 0 {
		t.Errorf("missing instruction lists:\n%s", strings.Join(missing, "\n"))
	}
}

func TestSimplePawnMove(t *testing.T) {
	b := new(board)

	b.set(1, 1, whiteColor, pawnKind)
	// no possibilities
	b.set(0, 0, whiteColor, pawnKind)
	// one possibility occupied
	b.set(0, 2, whiteColor, pawnKind)
	// against the (left and right) walls
	b.set(1, 7, whiteColor, pawnKind)
	b.set(3, 0, whiteColor, pawnKind)
	// crowning
	b.set(6, 6, whiteColor, pawnKind)

	//
	// black
	//

	b.set(6, 1, blackColor, pawnKind)
	// no possibilities
	b.set(7, 0, blackColor, pawnKind)
	// one possibility occupied
	b.set(7, 2, blackColor, pawnKind)
	// against the (left and right) walls
	b.set(4, 0, blackColor, pawnKind)
	b.set(4, 7, blackColor, pawnKind)
	// crowning
	b.set(1, 5, blackColor, pawnKind)

	t.Log("\n" + b.String())

	var movesGot [][]instruction
	movesGot = generateSimpleMoves(movesGot, b)

	movesWant := [][]instruction{
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
	assertEqualInstructionLists(t, movesGot, movesWant)
}

func TestSimpleKingMove(t *testing.T) {

	b := new(board)

	b.set(5, 5, whiteColor, kingKind)
	b.set(2, 2, blackColor, kingKind)
	b.set(0, 7, whiteColor, kingKind)

	t.Log("\n" + b.String())

	var movesGot [][]instruction
	movesGot = generateSimpleMoves(movesGot, b)

	movesWant := [][]instruction{
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

	assertEqualInstructionLists(t, movesGot, movesWant)
}

func TestCapturePawnMove(t *testing.T) {
	b := new(board)

	b.set(4, 6, blackColor, pawnKind)
	b.set(3, 5, whiteColor, pawnKind)
	b.set(3, 3, whiteColor, kingKind)
	b.set(1, 3, whiteColor, pawnKind)
	b.set(1, 1, blackColor, pawnKind)

	b.set(5, 3, whiteColor, kingKind)
	b.set(6, 4, blackColor, pawnKind)

	t.Log("\n" + b.String())

	movesGot := generateCaptureMoves(nil, b)

	movesWant := [][]instruction{
		{
			makeMoveInstruction(3, 5, 5, 7),
			makeCaptureInstruction(4, 6, blackColor, pawnKind),
		},
		{
			makeMoveInstruction(4, 6, 2, 4),
			makeCaptureInstruction(3, 5, whiteColor, pawnKind),
			makeMoveInstruction(2, 4, 0, 2),
			makeCaptureInstruction(1, 3, whiteColor, pawnKind),
			makeCrownInstruction(0, 2),
		},
		{
			makeMoveInstruction(4, 6, 2, 4),
			makeCaptureInstruction(3, 5, whiteColor, pawnKind),
			makeMoveInstruction(2, 4, 4, 2),
			makeCaptureInstruction(3, 3, whiteColor, kingKind),
		},
		{
			makeMoveInstruction(6, 4, 4, 2),
			makeCaptureInstruction(5, 3, whiteColor, kingKind),
			makeMoveInstruction(4, 2, 2, 4),
			makeCaptureInstruction(3, 3, whiteColor, kingKind),
			makeMoveInstruction(2, 4, 0, 2),
			makeCaptureInstruction(1, 3, whiteColor, pawnKind),
			makeCrownInstruction(0, 2),
		},
	}

	assertEqualInstructionLists(t, movesGot, movesWant)
}
