package core

import (
	"strings"
	"testing"
)

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

func TestKindOpposite(t *testing.T) {
	if PawnKind.Opposite() != KingKind {
		t.Fail()
	}
	if KingKind.Opposite() != PawnKind {
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

	b := DecodeBoard(s)

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

func TestBoardNotEquals(t *testing.T) {
	nilBoard := (*Board)(nil)
	b := new(Board)
	if b.Equals(nilBoard) || nilBoard.Equals(b) {
		t.Fail()
	}

	c := new(Board)

	b.Set(0, 0, WhiteColor, PawnKind)
	if c.Equals(b) {
		t.Fail()
	}

	c.Set(0, 0, BlackColor, PawnKind)
	if c.Equals(b) {
		t.Fail()
	}

	c.Set(0, 0, WhiteColor, KingKind)
	if c.Equals(b) {
		t.Fail()
	}

	emptyBoard := (*Board)(nil)
	if emptyBoard.Equals(b) {
		t.Fail()
	}
}

func assertPieceCount(t *testing.T, c PieceCount, wp, wk, bp, bk int8) {
	if wp != c.WhitePawns {
		t.Fail()
	}
	if wk != c.WhiteKings {
		t.Fail()
	}
	if bp != c.BlackPawns {
		t.Fail()
	}
	if bk != c.BlackKings {
		t.Fail()
	}
}

func TestPieceCount(t *testing.T) {
	b := new(Board)

	var wp, wk, bp, bk int8

	assertPieceCount(t, b.PieceCount(), wp, wk, bp, bk)

	b.Set(0, 0, BlackColor, PawnKind)
	bp++
	b.Set(0, 1, BlackColor, PawnKind)
	bp++
	b.Set(0, 2, BlackColor, PawnKind)
	bp++

	assertPieceCount(t, b.PieceCount(), wp, wk, bp, bk)

	b.Set(2, 3, WhiteColor, PawnKind)
	wp++

	assertPieceCount(t, b.PieceCount(), wp, wk, bp, bk)

	b.Set(3, 3, WhiteColor, KingKind)
	wk++

	assertPieceCount(t, b.PieceCount(), wp, wk, bp, bk)

	b.Set(4, 4, BlackColor, KingKind)
	bk++
	assertPieceCount(t, b.PieceCount(), wp, wk, bp, bk)

	b.Set(7, 1, WhiteColor, PawnKind)
	wp++
	assertPieceCount(t, b.PieceCount(), wp, wk, bp, bk)
}

func TestSerializeBoard(t *testing.T) {
	b := new(Board)

	if (*Board)(nil).Serialize() != "" {
		t.Log("serializing the nil board should return the empty string")
		t.Fail()
	}

	if b.Serialize() != "" {
		t.Log("serializing empty board should return the empty string")
		t.Fail()
	}

	b.Set(4, 6, WhiteColor, KingKind)
	b.Set(2, 7, WhiteColor, PawnKind)
	b.Set(5, 1, BlackColor, PawnKind)
	b.Set(6, 0, BlackColor, KingKind)

	if b.Serialize() != "27wp46wk51bp60bk" {
		t.Log("serializing failed")
		t.Fail()
	}
}

func TestUnserializeBoard(t *testing.T) {
	assertUnserializeErr := func(err error) {
		if err == nil {
			t.Log("unserializing succeeded but should've returned an error")
			t.Fail()
		}
		if !strings.HasPrefix(err.Error(), "unserialize board: ") {
			t.Log("unseralize error message should start with 'unserialize board: '")
			t.Fail()
		}
	}

	{
		b, err := UnserializeBoard("")
		if err != nil {
			t.Logf("unserializing returned err when it should've succeeded: %v", err)
			t.Fail()
		}
		count := b.PieceCount()
		if count.WhiteKings > 0 || count.WhitePawns > 0 || count.BlackKings > 0 || count.BlackPawns > 0 {
			t.Log("unserializing empty string failed, should've returned an empty board")
			t.Fail()
		}
	}

	{
		b0, err := UnserializeBoard("11wp22wk33wk14bp27bk07bk")
		if err != nil {
			t.Logf("unserializing returned err when it should've succeeded: %v", err)
			t.Fail()
		}
		b1 := new(Board)
		b1.Set(1, 1, WhiteColor, PawnKind)
		b1.Set(2, 2, WhiteColor, KingKind)
		b1.Set(3, 3, WhiteColor, KingKind)
		b1.Set(1, 4, BlackColor, PawnKind)
		b1.Set(2, 7, BlackColor, KingKind)
		b1.Set(0, 7, BlackColor, KingKind)
		if !b0.Equals(b1) {
			t.Log("unserializing failed")
			t.Fail()
		}
	}

	{
		_, err := UnserializeBoard("11wp22w")
		assertUnserializeErr(err)
	}

	{
		_, err := UnserializeBoard("11wp80wk")
		assertUnserializeErr(err)
	}

	{
		_, err := UnserializeBoard("12wp37wl")
		assertUnserializeErr(err)
	}

	{
		_, err := UnserializeBoard("12wp37mp")
		assertUnserializeErr(err)
	}
}
