package main

import (
	"bytes"
	"strings"
)

type color byte

type kind byte

const (
	blackColor = color(0)
	whiteColor = color(1)
	pawnKind   = kind(0)
	kingKind   = kind(1)
)

var crowningRow = [2]byte{
	int(blackColor): 7,
	int(whiteColor): 0,
}

var forward = [2]int8{
	int(blackColor): +1,
	int(whiteColor): -1,
}

func (c color) String() string {
	if c == whiteColor {
		return "white"
	} else {
		return "black"
	}
}

func (c color) opposite() color {
	if c == whiteColor {
		return blackColor
	}
	return whiteColor
}

func (k kind) String() string {
	if k == kingKind {
		return "king"
	} else {
		return "pawn"
	}
}

type board struct {
	occupied uint64
	white    uint64
	king     uint64
}

func pieceToRune(c color, k kind) rune {
	if c == whiteColor {
		if k == kingKind {
			return '@'
		}
		return 'o'
	}
	//black
	if k == kingKind {
		return '#'
	}
	//pawn
	return 'x'
}

func (b *board) String() string {
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
			if b.isOccupied(row, col) {
				buf.WriteRune(pieceToRune(b.get(row, col)))
			} else if tileColor(row, col) == blackColor {
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

func tileColor(row, col byte) color {
	if (row+col)%2 == 0 {
		return whiteColor
	} else {
		return blackColor
	}
}

func placeInitialPieces(b *board) {
	for row := byte(0); row <= 2; row++ {
		for col := byte(0); col < 8; col++ {
			if tileColor(row, col) == blackColor {
				b.set(row, col, blackColor, pawnKind)
			}
		}
	}
	for row := byte(5); row <= 7; row++ {
		for col := byte(0); col < 8; col++ {
			if tileColor(row, col) == blackColor {
				b.set(row, col, whiteColor, pawnKind)
			}
		}
	}
}

func (b *board) clear(row, col byte) {
	b.occupied &^= uint64(1 << (uint64(row)*8 + uint64(col)))
}

func (b *board) set(row, col byte, c color, k kind) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))

	b.occupied |= x

	if c == whiteColor {
		b.white |= x
	} else {
		b.white &^= x
	}

	if k == kingKind {
		b.king |= x
	} else {
		b.king &^= x
	}
}

func (b *board) move(srow, scol, drow, dcol byte) {
	c, k := b.get(srow, scol)
	b.clear(srow, scol)
	b.set(drow, dcol, c, k)
}

func (b *board) crown(row, col byte) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king |= x
}

func (b *board) uncrown(row, col byte) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king &^= x
}

func (b *board) isOccupied(row, col byte) bool {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	return b.occupied&x != 0
}

func (b *board) get(row, col byte) (c color, k kind) {
	n := uint64(row)*8 + uint64(col)
	x := uint64(1 << n)
	k = kind((b.king & x) >> n)
	c = color((b.white & x) >> n)
	return
}

func (b *board) copy() *board {
	var c board
	c.occupied = b.occupied
	c.white = b.white
	c.king = b.king
	return &c
}

type pieceCount struct {
	whitePawns int8
	blackPawns int8
	whiteKings int8
	blackKings int8
}

func (b *board) pieceCount() pieceCount {
	var c pieceCount

	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.isOccupied(row, col) {
				continue
			}

			color, kind := b.get(row, col)
			if color == whiteColor {
				if kind == pawnKind {
					c.whitePawns++
				} else {
					c.whiteKings++
				}
			} else {
				if kind == pawnKind {
					c.blackPawns++
				} else {
					c.blackKings++
				}
			}
		}
	}

	return c
}

func (b *board) equals(o *board) bool {
	if b == nil && o == nil {
		return true
	}
	if b == nil || o == nil {
		return false
	}
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if b.isOccupied(row, col) != o.isOccupied(row, col) {
				return false
			}
			if b.isOccupied(row, col) {
				bc, bk := b.get(row, col)
				oc, ok := o.get(row, col)
				if bc != oc || bk != ok {
					return false
				}
			}
		}
	}
	return true
}

func decodeBoard(s string) *board {
	rawLines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")

	// trim all liens and filter empty ones
	var lines []string
	for _, line := range rawLines {
		line = strings.Trim(line, " \t")
		if line != "" {
			lines = append(lines, line)
		}
	}

	b := new(board)

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
				b.set(byte(row), byte(col), blackColor, pawnKind)
			} else if cell == '#' {
				b.set(byte(row), byte(col), blackColor, kingKind)
			} else if cell == 'o' {
				b.set(byte(row), byte(col), whiteColor, pawnKind)
			} else if cell == '@' {
				b.set(byte(row), byte(col), whiteColor, kingKind)
			}

			col++
		}
	}

	return b
}
