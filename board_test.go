package main

import "testing"

func TestColorString(t *testing.T) {
	if WhiteColor.String() != "white" {
		t.Fail()
	}
	if BlackColor.String() != "black" {
		t.Fail()
	}
}

func TestKindString(t *testing.T) {
	if KingKind.String() != "king" {
		t.Fail()
	}
	if PawnKind.String() != "pawn" {
		t.Fail()
	}
}

func TestNewEmptyBoard(t *testing.T) {
	board := new(Board)
	for row := uint8(0); row < 8; row++ {
		for col := uint8(0); col < 8; col++ {
			if board.IsOccupied(row, col) {
				t.Errorf("row %d col %d of empty board is occupied\n", row, col)
			}
		}
	}
}

func TestSet(t *testing.T) {
	board := new(Board)

	type testCase struct {
		row, col byte
		c        Color
		k        Kind
	}

	cases := []testCase{
		{2, 2, WhiteColor, KingKind},
		{3, 3, BlackColor, KingKind},
		{1, 5, BlackColor, PawnKind},
		{7, 4, WhiteColor, PawnKind},
	}

	for _, c := range cases {
		row, col, color, kind := c.row, c.col, c.c, c.k
		board.Set(row, col, color, kind)
	}

	if board.IsOccupied(1, 1) {
		t.Error("(1, 1) should not be occupied but it is")
	}

	for _, c := range cases {
		row, col, color, kind := c.row, c.col, c.c, c.k
		if !board.IsOccupied(row, col) {
			t.Errorf("(%d, %d) should be occupied but isn't", row, col)
		} else {
			actualColor, actualKind := board.Get(row, col)
			if actualColor != color || actualKind != kind {
				t.Errorf(
					"piece at (%d, %d) should be a %s %s but is a %s %s",
					row, col,
					color, kind,
					actualColor, actualKind,
				)
			}
		}
	}

}

func TestClear(t *testing.T) {
	board := new(Board)

	board.Set(5, 4, BlackColor, KingKind)

	board.Clear(5, 4)

	if board.IsOccupied(5, 4) {
		c, k := board.Get(5, 4)
		t.Errorf("expected (5, 4) to be empty after Clear, but it's occupied with a %s %s\n", c, k)
	}

	board.Clear(1, 1)

	if board.IsOccupied(1, 1) {
		t.Error("expected (1, 1) to be empty after Clear (was already empty before), but it's occupied")
	}
}

func TestInitialPieces(t *testing.T) {
	type piece struct {
		c Color
		k Kind
	}

	whi := &piece{c: WhiteColor, k: PawnKind}
	bla := &piece{c: BlackColor, k: PawnKind}

	initial := [8][8]*piece{
		{nil, bla, nil, bla, nil, bla, nil, bla},
		{bla, nil, bla, nil, bla, nil, bla, nil},
		{nil, bla, nil, bla, nil, bla, nil, bla},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{whi, nil, whi, nil, whi, nil, whi, nil},
		{nil, whi, nil, whi, nil, whi, nil, whi},
		{whi, nil, whi, nil, whi, nil, whi, nil},
	}

	b := new(Board)
	PlaceInitialPieces(b)

	for row := uint8(0); row < 8; row++ {
		for col := uint8(0); col < 8; col++ {
			ptr := initial[row][col]

			wantOccupied := ptr != nil
			gotOccupied := b.IsOccupied(row, col)

			if wantOccupied != gotOccupied {
				if wantOccupied {
					t.Errorf("expected row %d col %d to be occupied but it's empty\n", row, col)
					return
				} else {
					t.Errorf("expected row %d col %d to be empty but it's occupied\n", row, col)
					return
				}
			}

			if gotOccupied {
				wantColor, wantKind := ptr.c, ptr.k
				gotColor, gotKind := b.Get(row, col)

				if wantColor != gotColor || wantKind != gotKind {
					t.Errorf(
						"expected row %d col %d to contain %s %s but it contains %s %s\n",
						row, col,
						wantColor, wantKind,
						gotColor, gotKind,
					)
					return
				}
			}
		}
	}
}

func TestDecodeBoard(t *testing.T) {

	s := `
		..#
		...o..@
		.....
		..x...........o
		
		1234




		..x..
		...
		...
		...
		...@
	`

	b := decodeBoard(s)

	t.Log("\n" + b.String())

	expect := make(map[coord]piece)
	expect[coord{0, 2}] = piece{BlackColor, KingKind}
	expect[coord{1, 3}] = piece{WhiteColor, PawnKind}
	expect[coord{1, 6}] = piece{WhiteColor, KingKind}
	expect[coord{3, 2}] = piece{BlackColor, PawnKind}
	expect[coord{3, 2}] = piece{BlackColor, PawnKind}
	expect[coord{5, 2}] = piece{BlackColor, PawnKind}

	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			expectedPiece, expectOccupied := expect[coord{row, col}]

			if expectOccupied {
				color, kind := b.Get(row, col)
				if !b.IsOccupied(row, col) {
					t.Errorf("expected (%d %d) to containt something but it's empty", row, col)
				} else if color != expectedPiece.Color || kind != expectedPiece.Kind {
					t.Errorf(
						"expected a %s %s at (%d, %d) but it contains a %s %s",
						expectedPiece.Color, expectedPiece.Kind,
						row, col,
						color, kind,
					)
				}
			} else {
				if b.IsOccupied(row, col) {
					color, kind := b.Get(row, col)
					t.Errorf(
						"expected (%d, %d) to be empty but it contains a %s, %s",
						row, col,
						color, kind,
					)
				}
			}
		}
	}

}

func TestBoardEquals(t *testing.T) {
	if !((*Board)(nil)).Equals(nil) {
		t.Fail()
	}
	b := new(Board)
	if b.Equals(nil) || (*Board)(nil).Equals(b) {
		t.Fail()
	}

	c := new(Board)
	for i := 0; i < 10; i++ {
		row, col, color, kind := rn8(), rn8(), rnColor(), rnKind()
		b.Set(row, col, color, kind)
		c.Set(row, col, color, kind)
	}

	if !b.Equals(c) {
		t.Fail()
	}
}

func assertRowEmpty(t *testing.T, b *Board, row byte, want bool) {
	got := b.IsRowEmpty(row)
	if got != want {
		t.Errorf("expected rowEmpty(%d) to be %v", row, want)
	}
}

func TestRowEmpty(t *testing.T) {
	b := decodeBoard(`
	  .....x
		.
		..@...o
		.
		.
		..x
		.
		.o
	`)
	assertRowEmpty(t, b, 0, false)
	assertRowEmpty(t, b, 1, true)
	assertRowEmpty(t, b, 2, false)
	assertRowEmpty(t, b, 3, true)
	assertRowEmpty(t, b, 4, true)
	assertRowEmpty(t, b, 5, false)
	assertRowEmpty(t, b, 6, true)
	assertRowEmpty(t, b, 7, false)
}
