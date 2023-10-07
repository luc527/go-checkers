package core

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
		{MakeMoveInstruction(1, 1, 2, 2)},
		{MakeMoveInstruction(1, 1, 2, 0)},
		{MakeMoveInstruction(0, 2, 1, 3)},
		{MakeMoveInstruction(1, 7, 2, 6)},
		{MakeMoveInstruction(3, 0, 4, 1)},
		{
			MakeMoveInstruction(6, 6, 7, 7),
			MakeCrownInstruction(7, 7),
		},
		{
			MakeMoveInstruction(6, 6, 7, 5),
			MakeCrownInstruction(7, 5),
		},
	}
	assertEqualPlies(t, blackPliesGot, blackPliesWant)

	whitePliesWant := []Ply{
		{MakeMoveInstruction(6, 1, 5, 0)},
		{MakeMoveInstruction(6, 1, 5, 2)},
		{MakeMoveInstruction(7, 2, 6, 3)},
		{MakeMoveInstruction(4, 0, 3, 1)},
		{MakeMoveInstruction(4, 7, 3, 6)},
		{
			MakeMoveInstruction(1, 5, 0, 6),
			MakeCrownInstruction(0, 6),
		},
		{
			MakeMoveInstruction(1, 5, 0, 4),
			MakeCrownInstruction(0, 4),
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
		{MakeMoveInstruction(5, 5, 6, 6)},
		{MakeMoveInstruction(5, 5, 7, 7)},
		// down, left
		{MakeMoveInstruction(5, 5, 6, 4)},
		{MakeMoveInstruction(5, 5, 7, 3)},
		// up, left
		{MakeMoveInstruction(5, 5, 4, 4)},
		{MakeMoveInstruction(5, 5, 3, 3)}, // gets stopped by the black king (no (2,2), (1,1), (0,0))
		// up, right
		{MakeMoveInstruction(5, 5, 4, 6)},
		{MakeMoveInstruction(5, 5, 3, 7)},
		//
		// moving the white king at (0, 7)
		//
		{MakeMoveInstruction(0, 7, 1, 6)},
		{MakeMoveInstruction(0, 7, 2, 5)},
		{MakeMoveInstruction(0, 7, 3, 4)},
		{MakeMoveInstruction(0, 7, 4, 3)},
		{MakeMoveInstruction(0, 7, 5, 2)},
		{MakeMoveInstruction(0, 7, 6, 1)},
		{MakeMoveInstruction(0, 7, 7, 0)},
	}

	assertEqualPlies(t, whitePliesGot, whitePliesWant)

	blackPliesWant := []Ply{
		//
		// moving the black king at (2, 2)
		//
		// down, right
		{MakeMoveInstruction(2, 2, 3, 3)},
		{MakeMoveInstruction(2, 2, 4, 4)}, // gets stopped by the white king (no (5,5) etc.)
		// down, left
		{MakeMoveInstruction(2, 2, 3, 1)},
		{MakeMoveInstruction(2, 2, 4, 0)},
		// up, right
		{MakeMoveInstruction(2, 2, 1, 3)},
		{MakeMoveInstruction(2, 2, 0, 4)},
		// up, left
		{MakeMoveInstruction(2, 2, 1, 1)},
		{MakeMoveInstruction(2, 2, 0, 0)},
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
			MakeMoveInstruction(3, 5, 5, 7),
			MakeCaptureInstruction(4, 6, WhiteColor, PawnKind),
		},
		{
			MakeMoveInstruction(5, 3, 7, 5),
			MakeCaptureInstruction(6, 4, WhiteColor, PawnKind),
			MakeCrownInstruction(7, 5),
		},
	}

	assertEqualPlies(t, blackPliesGot, blackPliesWant)

	whitePliesWant := []Ply{
		{
			MakeMoveInstruction(4, 6, 2, 4),
			MakeCaptureInstruction(3, 5, BlackColor, PawnKind),
			MakeMoveInstruction(2, 4, 0, 2),
			MakeCaptureInstruction(1, 3, BlackColor, PawnKind),
			MakeCrownInstruction(0, 2),
		},
		{
			MakeMoveInstruction(4, 6, 2, 4),
			MakeCaptureInstruction(3, 5, BlackColor, PawnKind),
			MakeMoveInstruction(2, 4, 4, 2),
			MakeCaptureInstruction(3, 3, BlackColor, PawnKind),
		},
		{
			MakeMoveInstruction(6, 4, 4, 2),
			MakeCaptureInstruction(5, 3, BlackColor, PawnKind),
			MakeMoveInstruction(4, 2, 2, 4),
			MakeCaptureInstruction(3, 3, BlackColor, PawnKind),
			MakeMoveInstruction(2, 4, 0, 2),
			MakeCaptureInstruction(1, 3, BlackColor, PawnKind),
			MakeCrownInstruction(0, 2),
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
			MakeMoveInstruction(4, 7, 2, 5),
			MakeCaptureInstruction(3, 6, BlackColor, PawnKind),
			MakeMoveInstruction(2, 5, 0, 3),
			MakeCaptureInstruction(1, 4, BlackColor, PawnKind),
			MakeMoveInstruction(0, 3, 2, 1),
			MakeCaptureInstruction(1, 2, BlackColor, PawnKind),
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
			MakeMoveInstruction(3, 3, 6, 6),
			MakeCaptureInstruction(5, 5, BlackColor, PawnKind),
		},
		{
			MakeMoveInstruction(3, 3, 7, 7),
			MakeCaptureInstruction(5, 5, BlackColor, PawnKind),
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
			MakeMoveInstruction(4, 7, 2, 5),
			MakeCaptureInstruction(3, 6, BlackColor, PawnKind),
			MakeMoveInstruction(2, 5, 0, 7),
			MakeCaptureInstruction(1, 6, BlackColor, PawnKind),
			MakeCrownInstruction(0, 7),
		},
		{
			MakeMoveInstruction(4, 7, 2, 5),
			MakeCaptureInstruction(3, 6, BlackColor, PawnKind),
			MakeMoveInstruction(2, 5, 0, 3),
			MakeCaptureInstruction(1, 4, BlackColor, PawnKind),
			MakeMoveInstruction(0, 3, 2, 1),
			MakeCaptureInstruction(1, 2, BlackColor, PawnKind),
			MakeMoveInstruction(2, 1, 4, 3),
			MakeCaptureInstruction(3, 2, BlackColor, PawnKind),
			MakeMoveInstruction(4, 3, 2, 5),
			MakeCaptureInstruction(3, 4, BlackColor, PawnKind),
			MakeMoveInstruction(2, 5, 0, 7),
			MakeCaptureInstruction(1, 6, BlackColor, PawnKind),
			MakeCrownInstruction(0, 7),
		},
		{
			MakeMoveInstruction(4, 7, 2, 5),
			MakeCaptureInstruction(3, 6, BlackColor, PawnKind),
			MakeMoveInstruction(2, 5, 4, 3),
			MakeCaptureInstruction(3, 4, BlackColor, PawnKind),
			MakeMoveInstruction(4, 3, 2, 1),
			MakeCaptureInstruction(3, 2, BlackColor, PawnKind),
			MakeMoveInstruction(2, 1, 0, 3),
			MakeCaptureInstruction(1, 2, BlackColor, PawnKind),
			MakeMoveInstruction(0, 3, 2, 5),
			MakeCaptureInstruction(1, 4, BlackColor, PawnKind),
			MakeMoveInstruction(2, 5, 0, 7),
			MakeCaptureInstruction(1, 6, BlackColor, PawnKind),
			MakeCrownInstruction(0, 7),
		},
	}
	pliesGot := generateCapturePlies(nil, b, WhiteColor)

	assertEqualPlies(t, pliesGot, pliesWant)
}

func TestPlyMarshal(t *testing.T) {
	type test struct {
		p Ply
		s string
	}
	tests := []test{
		{Ply{}, ""},
		{Ply{MakeMoveInstruction(3, 3, 5, 5), MakeMoveInstruction(5, 5, 3, 3)}, "m3355,m5533"},
		{Ply{MakeMoveInstruction(1, 2, 3, 4), MakeCrownInstruction(3, 4), MakeCaptureInstruction(6, 6, WhiteColor, KingKind)}, "m1234,k34,c66wk"},
	}
	for _, test := range tests {
		if got, err := test.p.MarshalJSON(); err != nil {
			t.Logf("error: %v", err)
			t.Fail()
		} else if string(got) != test.s {
			t.Logf("wanted %v got %v", test.s, string(got))
			t.Fail()
		}
	}
}

func TestPlyUnmarshalCorrectly(t *testing.T) {
	type test struct {
		s string
		p Ply
	}
	tests := []test{
		{"", Ply{}},
		{"k55", Ply{MakeCrownInstruction(5, 5)}},
		{"c77wp", Ply{MakeCaptureInstruction(7, 7, WhiteColor, PawnKind)}},
		{"m1245,c34bk,k45", Ply{MakeMoveInstruction(1, 2, 4, 5), MakeCaptureInstruction(3, 4, BlackColor, KingKind), MakeCrownInstruction(4, 5)}},
	}
	for _, test := range tests {
		p := Ply(nil)
		if err := p.UnmarshalJSON([]byte(test.s)); err != nil {
			t.Logf("error: %v", err)
			t.Fail()
		}
		if !p.Equals(test.p) {
			t.Logf("wanted %v got %v", test.p, p)
		}
	}
}

func TestPlyUnmarshalIncorrectly(t *testing.T) {
	tests := []string{
		"m1234,,",
		"c12wk,  m4455, k12",
	}
	for _, test := range tests {
		p := Ply{}
		if err := p.UnmarshalJSON([]byte(test)); err == nil {
			t.Logf("expected error for %v", test)
			t.Fail()
		}
	}
}
