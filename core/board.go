package core

import (
	"bytes"
	"fmt"
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

func coordMask(row, col byte) uint64 {
	return 1 << (uint64(row)*8 + uint64(col))
}

func (b *Board) Clear(row, col byte) {
	b.occupied &^= coordMask(row, col)
}

func (b *Board) Set(row, col byte, c Color, k Kind) {
	x := coordMask(row, col)

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
	x := coordMask(row, col)
	b.king |= x
}

func (b *Board) Uncrown(row, col byte) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king &^= x
}

func (b *Board) IsOccupied(row, col byte) bool {
	x := coordMask(row, col)
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

// This is used for testing
func DecodeBoard(s string) *Board {
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

// This is for passing the board over a network
func (b *Board) Serialize() string {
	if b == nil {
		return ""
	}
	var sb strings.Builder
	for row := byte(0); row < 8; row++ {
		for c := byte(0); c < 8; c++ {
			if !b.IsOccupied(row, c) {
				continue
			}
			color, kind := b.Get(row, c)
			sb.WriteByte(row + '0')
			sb.WriteByte(c + '0')

			if color == WhiteColor {
				sb.WriteByte('w')
			} else {
				sb.WriteByte('b')
			}

			if kind == KingKind {
				sb.WriteByte('k')
			} else {
				sb.WriteByte('p')
			}
		}
	}
	return sb.String()
}

func UnserializeBoard(s string) (*Board, error) {
	var b Board
	if len(s)%4 != 0 {
		return nil, fmt.Errorf("unserialize board: invalid board string (length %d not multiple of 4)", len(s))
	}
	for len(s) > 0 {
		rowRune, colRune, colorRune, kindRune := s[0], s[1], s[2], s[3]
		row := byte(rowRune) - '0'
		col := byte(colRune) - '0'

		if row > 7 || col > 7 {
			return nil, fmt.Errorf("unserialize board: invalid position (row %d, col %d)", row, col)
		}

		var color Color
		colorByte := byte(colorRune)
		if colorByte == 'w' {
			color = WhiteColor
		} else if colorByte == 'b' {
			color = BlackColor
		} else {
			return nil, fmt.Errorf("unserialize board: invalid color %c", colorByte)
		}

		var kind Kind
		kindByte := byte(kindRune)
		if kindByte == 'k' {
			kind = KingKind
		} else if kindByte == 'p' {
			kind = PawnKind
		} else {
			return nil, fmt.Errorf("unserialize board: invalid kind %c", kindByte)
		}

		b.Set(row, col, color, kind)

		s = s[4:]
	}
	return &b, nil
}
