package main

import (
	"bytes"
	"math/bits"
	"strings"
)

type Color byte

type Kind byte

const (
	BlackColor = Color(0)
	WhiteColor = Color(1)
	PawnKind   = Kind(0)
	KingKind   = Kind(1)
)

// Used mostly for testing
type coord struct {
	row, col byte
}
type piece struct {
	Color
	Kind
}

var crowningRow = [2]byte{
	int(BlackColor): 7,
	int(WhiteColor): 0,
}

var forward = [2]int8{
	int(BlackColor): +1,
	int(WhiteColor): -1,
}

func (c Color) String() string {
	if c == WhiteColor {
		return "white"
	} else {
		return "black"
	}
}

func (c Color) Opposite() Color {
	if c == WhiteColor {
		return BlackColor
	}
	return WhiteColor
}

func (k Kind) String() string {
	if k == KingKind {
		return "king"
	} else {
		return "pawn"
	}
}

func (k Kind) Opposite() Kind {
	if k == KingKind {
		return PawnKind
	}
	return KingKind
}

type Board struct {
	occupied uint64
	white    uint64
	king     uint64
}

func pieceToRune(c Color, k Kind) rune {
	if c == WhiteColor {
		if k == KingKind {
			return '@'
		}
		return 'o'
	}
	//black
	if k == KingKind {
		return '#'
	}
	//pawn
	return 'x'
}

func (b *Board) String() string {
	buf := new(bytes.Buffer)

	buf.WriteRune(' ')
	for col := byte(0); col < 8; col++ {
		buf.WriteRune('0' + rune(col))
	}
	buf.WriteRune(' ')
	// for alignment when writing side by side

	for row := byte(0); row < 8; row++ {
		buf.WriteString("\n")
		buf.WriteRune('0' + rune(row))
		for col := byte(0); col < 8; col++ {
			if b.IsOccupied(row, col) {
				buf.WriteRune(pieceToRune(b.Get(row, col)))
			} else if TileColor(row, col) == BlackColor {
				buf.WriteRune('_')
			} else {
				buf.WriteRune(' ')
			}
		}
		buf.WriteRune('0' + rune(row))
	}

	buf.WriteRune('\n')
	buf.WriteRune(' ')
	for col := byte(0); col < 8; col++ {
		buf.WriteRune('0' + rune(col))
	}
	buf.WriteRune(' ')
	buf.WriteRune('\n')

	return buf.String()
}

func TileColor(row, col byte) Color {
	if (row+col)%2 == 0 {
		return WhiteColor
	} else {
		return BlackColor
	}
}

func PlaceInitialPieces(b *Board) {
	for row := byte(0); row <= 2; row++ {
		for col := byte(0); col < 8; col++ {
			if TileColor(row, col) == BlackColor {
				b.Set(row, col, BlackColor, PawnKind)
			}
		}
	}
	for row := byte(5); row <= 7; row++ {
		for col := byte(0); col < 8; col++ {
			if TileColor(row, col) == BlackColor {
				b.Set(row, col, WhiteColor, PawnKind)
			}
		}
	}
}

func (b *Board) Clear(row, col byte) {
	b.occupied &^= uint64(1 << (uint64(row)*8 + uint64(col)))
}

func (b *Board) Set(row, col byte, c Color, k Kind) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))

	b.occupied |= x

	if c == WhiteColor {
		b.white |= x
	} else {
		b.white &^= x
	}

	if k == KingKind {
		b.king |= x
	} else {
		b.king &^= x
	}
}

func (b *Board) Move(srow, scol, drow, dcol byte) {
	c, k := b.Get(srow, scol)
	b.Clear(srow, scol)
	b.Set(drow, dcol, c, k)
}

func (b *Board) Crown(row, col byte) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king |= x
}

func (b *Board) Uncrown(row, col byte) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king &^= x
}

func (b *Board) IsOccupied(row, col byte) bool {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	return b.occupied&x != 0
}

func (b *Board) Get(row, col byte) (c Color, k Kind) {
	n := uint64(row)*8 + uint64(col)
	x := uint64(1 << n)
	k = Kind((b.king & x) >> n)
	c = Color((b.white & x) >> n)
	return
}

func (b *Board) Copy() *Board {
	var c Board
	c.occupied = b.occupied
	c.white = b.white
	c.king = b.king
	return &c
}

type PieceCount struct {
	WhitePawns int8
	BlackPawns int8
	WhiteKings int8
	BlackKings int8
}

func (b *Board) PieceCount() PieceCount {
	var c PieceCount

	king := b.occupied & b.king
	pawn := b.occupied &^ b.king

	kings := bits.OnesCount64(king)
	pawns := bits.OnesCount64(pawn)

	whitePawns := bits.OnesCount64(pawn & b.white)
	c.WhitePawns = int8(whitePawns)
	c.BlackPawns = int8(pawns - whitePawns)

	whiteKings := bits.OnesCount64(king & b.white)
	c.WhiteKings = int8(whiteKings)
	c.BlackKings = int8(kings - whiteKings)

	return c
}

func (b *Board) Equals(o *Board) bool {
	if b == nil && o == nil {
		return true
	}
	if b == nil || o == nil {
		return false
	}
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if b.IsOccupied(row, col) != o.IsOccupied(row, col) {
				return false
			}
			if b.IsOccupied(row, col) {
				bc, bk := b.Get(row, col)
				oc, ok := o.Get(row, col)
				if bc != oc || bk != ok {
					return false
				}
			}
		}
	}
	return true
}

func decodeBoard(s string) *Board {
	rawLines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")

	// trim all liens and filter empty ones
	var lines []string
	for _, line := range rawLines {
		line = strings.Trim(line, " \t")
		if line != "" {
			lines = append(lines, line)
		}
	}

	b := new(Board)

	// parse lines rawLines
	maxRow := 8
	if len(lines) < 8 {
		maxRow = len(lines)
	}

	for row := 0; row < maxRow; row++ {
		line := lines[row]

		// can't count on len(line) because it counts bytes and not unicode runes
		col := 0
		for _, cell := range line {
			if col >= 8 {
				break
			}

			if cell == 'x' {
				b.Set(byte(row), byte(col), BlackColor, PawnKind)
			} else if cell == '#' {
				b.Set(byte(row), byte(col), BlackColor, KingKind)
			} else if cell == 'o' {
				b.Set(byte(row), byte(col), WhiteColor, PawnKind)
			} else if cell == '@' {
				b.Set(byte(row), byte(col), WhiteColor, KingKind)
			}

			col++
		}
	}

	return b
}
