package go_checkers

import (
	"strings"
	"testing"
)

func compareGeneratedPlies(
	got []Ply,
	want []Ply,
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

func assertEqualPlies(t *testing.T, got []Ply, want []Ply) {
	extra, missing := compareGeneratedPlies(got, want)
	if len(extra) > 0 {
		t.Errorf("generated extra instruction lists:\n%s", strings.Join(extra, "\n"))
	}
	if len(missing) > 0 {
		t.Errorf("missing instruction lists:\n%s", strings.Join(missing, "\n"))
	}
}

func TestSimplePawnMove(t *testing.T) {
	b := new(Board)

	b.Set(1, 1, BlackColor, PawnKind)
	// no possibilities
	b.Set(0, 0, BlackColor, PawnKind)
	// one possibility occupied
	b.Set(0, 2, BlackColor, PawnKind)
	// against the (left and right) walls
	b.Set(1, 7, BlackColor, PawnKind)
	b.Set(3, 0, BlackColor, PawnKind)
	// crowning
	b.Set(6, 6, BlackColor, PawnKind)

	//
	// black
	//

	b.Set(6, 1, WhiteColor, PawnKind)
	// no possibilities
	b.Set(7, 0, WhiteColor, PawnKind)
	// one possibility occupied
	b.Set(7, 2, WhiteColor, PawnKind)
	// against the (left and right) walls
	b.Set(4, 0, WhiteColor, PawnKind)
	b.Set(4, 7, WhiteColor, PawnKind)
	// crowning
	b.Set(1, 5, WhiteColor, PawnKind)

	t.Log("\n" + b.String())

	blackPliesGot := generateSimplePlies(nil, b, BlackColor)

	blackPliesWant := []Ply{
		{MoveInstruction(1, 1, 2, 2)},
		{MoveInstruction(1, 1, 2, 0)},
		{MoveInstruction(0, 2, 1, 3)},
		{MoveInstruction(1, 7, 2, 6)},
		{MoveInstruction(3, 0, 4, 1)},
		{
			MoveInstruction(6, 6, 7, 7),
			CrownInstruction(7, 7),
		},
		{
			MoveInstruction(6, 6, 7, 5),
			CrownInstruction(7, 5),
		},
	}
	assertEqualPlies(t, blackPliesGot, blackPliesWant)

	whitePliesWant := []Ply{
		{MoveInstruction(6, 1, 5, 0)},
		{MoveInstruction(6, 1, 5, 2)},
		{MoveInstruction(7, 2, 6, 3)},
		{MoveInstruction(4, 0, 3, 1)},
		{MoveInstruction(4, 7, 3, 6)},
		{
			MoveInstruction(1, 5, 0, 6),
			CrownInstruction(0, 6),
		},
		{
			MoveInstruction(1, 5, 0, 4),
			CrownInstruction(0, 4),
		},
	}
	whitePliesGot := generateSimplePlies(nil, b, WhiteColor)

	assertEqualPlies(t, whitePliesGot, whitePliesWant)
}

func TestSimpleKingMove(t *testing.T) {

	b := new(Board)

	//white
	b.Set(5, 5, WhiteColor, KingKind)
	b.Set(0, 7, WhiteColor, KingKind)

	//black
	b.Set(2, 2, BlackColor, KingKind)

	t.Log("\n" + b.String())

	whitePliesGot := generateSimplePlies(nil, b, WhiteColor)

	whitePliesWant := []Ply{
		//
		// moving the white king at (5, 5)
		//
		// down, right
		{MoveInstruction(5, 5, 6, 6)},
		{MoveInstruction(5, 5, 7, 7)},
		// down, left
		{MoveInstruction(5, 5, 6, 4)},
		{MoveInstruction(5, 5, 7, 3)},
		// up, left
		{MoveInstruction(5, 5, 4, 4)},
		{MoveInstruction(5, 5, 3, 3)}, // gets stopped by the black king (no (2,2), (1,1), (0,0))
		// up, right
		{MoveInstruction(5, 5, 4, 6)},
		{MoveInstruction(5, 5, 3, 7)},
		//
		// moving the white king at (0, 7)
		//
		{MoveInstruction(0, 7, 1, 6)},
		{MoveInstruction(0, 7, 2, 5)},
		{MoveInstruction(0, 7, 3, 4)},
		{MoveInstruction(0, 7, 4, 3)},
		{MoveInstruction(0, 7, 5, 2)},
		{MoveInstruction(0, 7, 6, 1)},
		{MoveInstruction(0, 7, 7, 0)},
	}

	assertEqualPlies(t, whitePliesGot, whitePliesWant)

	blackPliesWant := []Ply{
		//
		// moving the black king at (2, 2)
		//
		// down, right
		{MoveInstruction(2, 2, 3, 3)},
		{MoveInstruction(2, 2, 4, 4)}, // gets stopped by the white king (no (5,5) etc.)
		// down, left
		{MoveInstruction(2, 2, 3, 1)},
		{MoveInstruction(2, 2, 4, 0)},
		// up, right
		{MoveInstruction(2, 2, 1, 3)},
		{MoveInstruction(2, 2, 0, 4)},
		// up, left
		{MoveInstruction(2, 2, 1, 1)},
		{MoveInstruction(2, 2, 0, 0)},
	}
	blackPliesGot := generateSimplePlies(nil, b, BlackColor)

	assertEqualPlies(t, blackPliesGot, blackPliesWant)
}

func TestCapturePawnMove(t *testing.T) {
	b := new(Board)

	b.Set(4, 6, WhiteColor, PawnKind)
	b.Set(3, 5, BlackColor, PawnKind)
	b.Set(3, 3, BlackColor, PawnKind)
	b.Set(1, 3, BlackColor, PawnKind)
	b.Set(1, 1, WhiteColor, PawnKind)

	b.Set(5, 3, BlackColor, PawnKind)
	b.Set(6, 4, WhiteColor, PawnKind)

	t.Log("\n" + b.String())

	blackPliesGot := generateCapturePlies(nil, b, BlackColor)

	blackPliesWant := []Ply{
		{
			MoveInstruction(3, 5, 5, 7),
			CaptureInstruction(4, 6, WhiteColor, PawnKind),
		},
		{
			MoveInstruction(5, 3, 7, 5),
			CaptureInstruction(6, 4, WhiteColor, PawnKind),
			CrownInstruction(7, 5),
		},
	}

	assertEqualPlies(t, blackPliesGot, blackPliesWant)

	whitePliesWant := []Ply{
		{
			MoveInstruction(4, 6, 2, 4),
			CaptureInstruction(3, 5, BlackColor, PawnKind),
			MoveInstruction(2, 4, 0, 2),
			CaptureInstruction(1, 3, BlackColor, PawnKind),
			CrownInstruction(0, 2),
		},
		{
			MoveInstruction(4, 6, 2, 4),
			CaptureInstruction(3, 5, BlackColor, PawnKind),
			MoveInstruction(2, 4, 4, 2),
			CaptureInstruction(3, 3, BlackColor, PawnKind),
		},
		{
			MoveInstruction(6, 4, 4, 2),
			CaptureInstruction(5, 3, BlackColor, PawnKind),
			MoveInstruction(4, 2, 2, 4),
			CaptureInstruction(3, 3, BlackColor, PawnKind),
			MoveInstruction(2, 4, 0, 2),
			CaptureInstruction(1, 3, BlackColor, PawnKind),
			CrownInstruction(0, 2),
		},
	}
	whitePliesGot := generateCapturePlies(nil, b, WhiteColor)

	assertEqualPlies(t, whitePliesGot, whitePliesWant)
}

func TestCaptureThroughCrowningRowDoesntCrown(t *testing.T) {
	b := new(Board)

	// passes through (0, 3), with 0 being the crowning row for white pieces,
	// but it shouldn't crown because it doesn't *end* at that position,
	// just goes through it
	b.Set(4, 7, WhiteColor, PawnKind)
	b.Set(3, 6, BlackColor, PawnKind)
	b.Set(1, 4, BlackColor, PawnKind)
	b.Set(1, 2, BlackColor, PawnKind)

	t.Log("\n" + b.String())

	pliesWant := []Ply{
		{
			MoveInstruction(4, 7, 2, 5),
			CaptureInstruction(3, 6, BlackColor, PawnKind),
			MoveInstruction(2, 5, 0, 3),
			CaptureInstruction(1, 4, BlackColor, PawnKind),
			MoveInstruction(0, 3, 2, 1),
			CaptureInstruction(1, 2, BlackColor, PawnKind),
		},
	}
	pliesGot := generateCapturePlies(nil, b, WhiteColor)

	assertEqualPlies(t, pliesGot, pliesWant)
}

func TestCaptureKingMoveOneDiagonal(t *testing.T) {
	b := new(Board)

	b.Set(3, 3, WhiteColor, KingKind)
	b.Set(5, 5, BlackColor, PawnKind)

	pliesGot := generateCapturePlies(nil, b, WhiteColor)
	pliesWant := []Ply{
		{
			MoveInstruction(3, 3, 6, 6),
			CaptureInstruction(5, 5, BlackColor, PawnKind),
		},
		{
			MoveInstruction(3, 3, 7, 7),
			CaptureInstruction(5, 5, BlackColor, PawnKind),
		},
	}

	assertEqualPlies(t, pliesGot, pliesWant)
}

func TestAllowOverPreviousTile(t *testing.T) {
	b := new(Board)

	b.Set(4, 7, WhiteColor, PawnKind)
	b.Set(3, 6, BlackColor, PawnKind)
	b.Set(1, 6, BlackColor, PawnKind)
	b.Set(3, 4, BlackColor, PawnKind)
	b.Set(1, 4, BlackColor, PawnKind)
	b.Set(3, 2, BlackColor, PawnKind)
	b.Set(1, 2, BlackColor, PawnKind)

	t.Log("\n" + b.String())

	pliesWant := []Ply{
		{
			MoveInstruction(4, 7, 2, 5),
			CaptureInstruction(3, 6, BlackColor, PawnKind),
			MoveInstruction(2, 5, 0, 7),
			CaptureInstruction(1, 6, BlackColor, PawnKind),
			CrownInstruction(0, 7),
		},
		{
			MoveInstruction(4, 7, 2, 5),
			CaptureInstruction(3, 6, BlackColor, PawnKind),
			MoveInstruction(2, 5, 0, 3),
			CaptureInstruction(1, 4, BlackColor, PawnKind),
			MoveInstruction(0, 3, 2, 1),
			CaptureInstruction(1, 2, BlackColor, PawnKind),
			MoveInstruction(2, 1, 4, 3),
			CaptureInstruction(3, 2, BlackColor, PawnKind),
			MoveInstruction(4, 3, 2, 5),
			CaptureInstruction(3, 4, BlackColor, PawnKind),
			MoveInstruction(2, 5, 0, 7),
			CaptureInstruction(1, 6, BlackColor, PawnKind),
			CrownInstruction(0, 7),
		},
		{
			MoveInstruction(4, 7, 2, 5),
			CaptureInstruction(3, 6, BlackColor, PawnKind),
			MoveInstruction(2, 5, 4, 3),
			CaptureInstruction(3, 4, BlackColor, PawnKind),
			MoveInstruction(4, 3, 2, 1),
			CaptureInstruction(3, 2, BlackColor, PawnKind),
			MoveInstruction(2, 1, 0, 3),
			CaptureInstruction(1, 2, BlackColor, PawnKind),
			MoveInstruction(0, 3, 2, 5),
			CaptureInstruction(1, 4, BlackColor, PawnKind),
			MoveInstruction(2, 5, 0, 7),
			CaptureInstruction(1, 6, BlackColor, PawnKind),
			CrownInstruction(0, 7),
		},
	}
	pliesGot := generateCapturePlies(nil, b, WhiteColor)

	assertEqualPlies(t, pliesGot, pliesWant)
}
