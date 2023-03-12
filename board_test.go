package main

import "testing"

func TestColorString(t *testing.T) {
	if whiteColor.String() != "white" {
		t.Fail()
	}
	if blackColor.String() != "black" {
		t.Fail()
	}
}

func TestKindString(t *testing.T) {
	if kingKind.String() != "king" {
		t.Fail()
	}
	if pawnKind.String() != "pawn" {
		t.Fail()
	}
}

func TestNewEmptyBoard(t *testing.T) {
	board := new(board)
	for row := uint8(0); row < 8; row++ {
		for col := uint8(0); col < 8; col++ {
			if board.isOccupied(row, col) {
				t.Errorf("row %d col %d of empty board is occupied\n", row, col)
			}
		}
	}
}

func TestSet(t *testing.T) {
	board := new(board)

	type testCase struct {
		row, col byte
		c        color
		k        kind
	}

	cases := []testCase{
		{2, 2, whiteColor, kingKind},
		{3, 3, blackColor, kingKind},
		{1, 5, blackColor, pawnKind},
		{7, 4, whiteColor, pawnKind},
	}

	for _, c := range cases {
		row, col, color, kind := c.row, c.col, c.c, c.k
		board.set(row, col, color, kind)
	}

	if board.isOccupied(1, 1) {
		t.Error("(1, 1) should not be occupied but it is")
	}

	for _, c := range cases {
		row, col, color, kind := c.row, c.col, c.c, c.k
		if !board.isOccupied(row, col) {
			t.Errorf("(%d, %d) should be occupied but isn't", row, col)
		} else {
			actualColor, actualKind := board.get(row, col)
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
	board := new(board)

	board.set(5, 4, blackColor, kingKind)

	board.clear(5, 4)

	if board.isOccupied(5, 4) {
		c, k := board.get(5, 4)
		t.Errorf("expected (5, 4) to be empty after clear, but it's occupied with a %s %s\n", c, k)
	}

	board.clear(1, 1)

	if board.isOccupied(1, 1) {
		t.Error("expected (1, 1) to be empty after clear (was already empty before), but it's occupied")
	}
}

func TestInitialPieces(t *testing.T) {
	type piece struct {
		c color
		k kind
	}

	whi := &piece{c: whiteColor, k: pawnKind}
	bla := &piece{c: blackColor, k: pawnKind}

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

	b := new(board)
	placeInitialPieces(b)

	for row := uint8(0); row < 8; row++ {
		for col := uint8(0); col < 8; col++ {
			ptr := initial[row][col]

			wantOccupied := ptr != nil
			gotOccupied := b.isOccupied(row, col)

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
				gotColor, gotKind := b.get(row, col)

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

	type coord struct {
		row, col byte
	}

	type piece struct {
		color
		kind
	}

	expect := make(map[coord]piece)
	expect[coord{0, 2}] = piece{blackColor, kingKind}
	expect[coord{1, 3}] = piece{whiteColor, pawnKind}
	expect[coord{1, 6}] = piece{whiteColor, kingKind}
	expect[coord{3, 2}] = piece{blackColor, pawnKind}
	expect[coord{3, 2}] = piece{blackColor, pawnKind}
	expect[coord{5, 2}] = piece{blackColor, pawnKind}

	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			expectedPiece, expectOccupied := expect[coord{row, col}]

			if expectOccupied {
				color, kind := b.get(row, col)
				if !b.isOccupied(row, col) {
					t.Errorf("expected (%d %d) to containt something but it's empty", row, col)
				} else if color != expectedPiece.color || kind != expectedPiece.kind {
					t.Errorf(
						"expected a %s %s at (%d, %d) but it contains a %s %s",
						expectedPiece.color, expectedPiece.kind,
						row, col,
						color, kind,
					)
				}
			} else {
				if b.isOccupied(row, col) {
					color, kind := b.get(row, col)
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
