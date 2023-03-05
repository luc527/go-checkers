package main

// TODO add log to tests, run with go test -v

import (
	"testing"
)

func TestMakeCrownIns(t *testing.T) {
	var row, col byte
	row, col = 5, 6
	i := makeCrownInstruction(row, col)
	if i.t != crownInstruction {
		t.Errorf("expected instruction to be of type crown, is %s", i.t)
		return
	}
	if i.row != row || i.col != col {
		t.Errorf("expected instruction to be crowning %d %d but is crowning %d %d", row, col, i.row, i.col)
	}
}

func TestMakeMoveIns(t *testing.T) {
	var srow, scol byte
	var drow, dcol byte

	// invalid, whatever; it's not what we're testing here
	srow, scol = 1, 3
	drow, dcol = 2, 5

	i := makeMoveInstruction(srow, scol, drow, dcol)

	if i.t != moveInstruction {
		t.Errorf("expected type move but is of type %s", i.t)
		return
	}

	if i.row != srow || i.col != scol {
		t.Errorf("expected source %d %d but is %d %d", srow, scol, i.row, i.col)
		return
	}

	if i.d[0] != drow || i.d[1] != dcol {
		t.Errorf("expected destination %d %d but is %d %d", drow, dcol, i.row, i.col)
		return
	}
}

func TestMakeCaptureIns(t *testing.T) {
	type testCase struct {
		row, col byte
		c        color
		k        kind
	}

	cases := []testCase{
		{1, 2, kWhite, kKing},
		{3, 1, kBlack, kPawn},
		{2, 2, kWhite, kPawn},
		{5, 7, kBlack, kKing},
	}

	for _, test := range cases {
		row, col, c, k := test.row, test.col, test.c, test.k
		i := makeCaptureInstruction(row, col, c, k)

		if i.t != captureInstruction {
			t.Errorf("expected type capture but got type %s", i.t)
			return
		}

		if row != i.row || col != i.col {
			t.Errorf("expected coord %d %d but got %d %d", row, col, i.row, i.col)
			return
		}

		actualC := color(i.d[0])
		if actualC != c {
			t.Errorf("expected color %s but got %s", c, actualC)
			return
		}

		actualK := kind(i.d[1])
		if actualK != k {
			t.Errorf("expected kind %s but got %s", k, actualK)
			return
		}
	}

}

func TestCrownIns(t *testing.T) {
	b := newEmptyBoard()

	var row, col byte
	row, col = 5, 4

	b.set(row, col, kWhite, kPawn)

	i := makeCrownInstruction(row, col)
	is := []instruction{i}

	performInstructions(b, is)

	_, newKind := b.get(row, col)
	if newKind != kKing {
		t.Errorf("crown instruction failed, %d %d still a pawn", row, col)
	}

	undoInstructions(b, is)

	_, oldKind := b.get(row, col)
	if oldKind != kPawn {
		t.Errorf("undo of crown instruction failed, %d %d still a king", row, col)
	}
}

func TestMoveIns(t *testing.T) {
	b := newEmptyBoard()

	var frow, fcol byte //from
	var trow, tcol byte //to

	frow, fcol = 3, 7
	trow, tcol = 4, 6
	c, k := kBlack, kKing

	b.set(frow, fcol, c, k)

	i := makeMoveInstruction(frow, fcol, trow, tcol)
	is := []instruction{i}

	performInstructions(b, is)

	if b.isOccupied(frow, fcol) {
		t.Errorf("after move, source should be empty")
	}

	if !b.isOccupied(trow, tcol) {
		t.Errorf("after move, destination should be occupied")
	} else {
		ac, ak := b.get(trow, tcol)
		if ac != c || ak != k {
			t.Errorf("piece changed after move, was %s %s now is %s %s", c, k, ac, ak)
		}
	}

	undoInstructions(b, is)

	if b.isOccupied(trow, tcol) {
		t.Errorf("after undo move, destination should be empty")
	}

	if !b.isOccupied(frow, fcol) {
		t.Errorf("after undo move, source should be occupied")
	} else {
		ac, ak := b.get(frow, fcol)
		if ac != c || ak != k {
			t.Errorf("piece changed after undo move, was %s %s now is %s %s", c, k, ac, ak)
		}
	}
}

func TestCaptureIns(t *testing.T) {
	b := newEmptyBoard()

	var row, col byte
	row, col = 3, 6
	color, kind := kWhite, kPawn

	b.set(row, col, color, kind)

	t.Log("Before capture:")
	t.Log(b)

	i := makeCaptureInstruction(row, col, color, kind)
	is := []instruction{i}

	performInstructions(b, is)

	t.Log("After capture:")
	t.Log(b)

	if b.isOccupied(row, col) {
		t.Errorf("(%d, %d) should be empty after capture, is occupied", row, col)
	}

	undoInstructions(b, is)

	t.Log("After undoing capture:")
	t.Log(b)

	if !b.isOccupied(row, col) {
		t.Errorf("(%d, %d) should be occupied after undoing the capture, is empty", row, col)
	} else {
		actualColor, actualKind := b.get(row, col)
		if actualColor != color || actualKind != kind {
			t.Errorf(
				"expected (%d, %d) to contain a %s %s after undoing the capture, but it contains a %s %s",
				row, col,
				color, kind,
				actualColor, actualKind,
			)
		}
	}
}

func TestInsSequence(t *testing.T) {

	b := newEmptyBoard()

	b.set(3, 5, kWhite, kPawn)
	b.set(1, 0, kBlack, kKing)
	b.set(2, 2, kBlack, kPawn)

	t.Log("Before:")
	t.Log("\n" + b.String())

	before := b.copy()

	is := []instruction{
		makeMoveInstruction(3, 5, 2, 4),
		makeCrownInstruction(2, 4),
		makeCaptureInstruction(2, 4, kWhite, kKing),
		makeMoveInstruction(1, 0, 4, 6),
		makeMoveInstruction(2, 2, 3, 5),
		makeCrownInstruction(3, 5),
	}

	performInstructions(b, is)

	t.Log("After:")
	t.Log("\n" + b.String())

	assertOccupied(b, t, 3, 5)
	assertContains(b, t, 3, 5, kBlack, kKing)
	assertOccupied(b, t, 4, 6)
	assertContains(b, t, 4, 6, kBlack, kKing)
	assertEmpty(b, t, 1, 0)
	assertEmpty(b, t, 2, 2)
	assertEmpty(b, t, 2, 4)

	undoInstructions(b, is)

	t.Log("After undo:")
	t.Log("\n" + b.String())

	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			wantOccupied := before.isOccupied(row, col)
			gotOccupied := b.isOccupied(row, col)

			if wantOccupied != gotOccupied {
				t.Errorf("row %d col %d should be occupied(%v) but is occupied(%v)", row, col, wantOccupied, gotOccupied)
			} else if gotOccupied {
				wantColor, wantKind := before.get(row, col)
				gotColor, gotKind := b.get(row, col)

				if wantColor != gotColor || wantKind != gotKind {
					t.Errorf("row %d col %d should contain %s %s but contains %s %s", row, col, wantColor, wantKind, gotColor, gotKind)
				}
			}
		}
	}
}

// TODO refactor other tests to use these assertions

func assertOccupied(b *board, t *testing.T, row, col byte) {
	if !b.isOccupied(row, col) {
		t.Errorf("row %d col %d should be occupied", row, col)
	}
}

func assertContains(b *board, t *testing.T, row, col byte, c color, k kind) {
	ac, ak := b.get(row, col)
	if ac != c || ak != k {
		t.Errorf("row %d col %d should contain %s %s but contains %s %s", row, col, c, k, ac, ak)
	}
}

func assertEmpty(b *board, t *testing.T, row, col byte) {
	if b.isOccupied(row, col) {
		t.Errorf("row %d col %d should be empty", row, col)
	}
}

// TODO test random instruction sequence
