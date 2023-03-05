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

	b.set(1, 1, kWhite, kPawn)
	// no possibilities
	b.set(0, 0, kWhite, kPawn)
	// one possibility occupied
	b.set(0, 2, kWhite, kPawn)
	// against the (left and right) walls
	b.set(1, 7, kWhite, kPawn)
	b.set(3, 0, kWhite, kPawn)
	// crowning
	b.set(6, 6, kWhite, kPawn)

	//
	// black
	//

	b.set(6, 1, kBlack, kPawn)
	// no possibilities
	b.set(7, 0, kBlack, kPawn)
	// one possibility occupied
	b.set(7, 2, kBlack, kPawn)
	// against the (left and right) walls
	b.set(4, 0, kBlack, kPawn)
	b.set(4, 7, kBlack, kPawn)
	// crowning
	b.set(1, 5, kBlack, kPawn)

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

	b.set(5, 5, kWhite, kKing)
	b.set(2, 2, kBlack, kKing)
	b.set(0, 7, kWhite, kKing)

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

	b.set(5, 6, kWhite, kPawn)
	b.set(4, 5, kBlack, kPawn) // black pawn that can be captured
	b.set(4, 3, kBlack, kKing) // black king that can be captured in sequence
	b.set(2, 3, kBlack, kPawn) // alternative black pawn that can be captured in sequence
	b.set(2, 1, kWhite, kPawn) // can't capture this one

	b.set(6, 3, kBlack, kKing) // could be another capture
	b.set(7, 4, kWhite, kPawn) // were it not for this piece
	// this piece at (7, 4) can also capture

	t.Log("\n" + b.String())

	movesGot := generateCaptureMoves(nil, b)

	movesWant := [][]instruction{
		{
			makeMoveInstruction(4, 5, 6, 7),
			makeCaptureInstruction(5, 6, kWhite, kPawn),
		},
		{
			makeMoveInstruction(5, 6, 3, 4),
			makeCaptureInstruction(4, 5, kBlack, kPawn),
			makeMoveInstruction(3, 4, 1, 2),
			makeCaptureInstruction(2, 3, kBlack, kPawn),
		},
		{
			makeMoveInstruction(5, 6, 3, 4),
			makeCaptureInstruction(4, 5, kBlack, kPawn),
			makeMoveInstruction(3, 4, 5, 2),
			makeCaptureInstruction(4, 3, kBlack, kKing),
		},
		{
			makeMoveInstruction(7, 4, 5, 2),
			makeCaptureInstruction(6, 3, kBlack, kKing),
			makeMoveInstruction(5, 2, 3, 4),
			makeCaptureInstruction(4, 3, kBlack, kKing),
			makeMoveInstruction(3, 4, 1, 2),
			makeCaptureInstruction(2, 3, kBlack, kPawn),
		},
	}

	assertEqualInstructionLists(t, movesGot, movesWant)
}
