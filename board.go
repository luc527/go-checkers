package main

import (
	"bytes"
)

type color byte

type kind byte

const (
	kBlack = color(0)
	kWhite = color(1)
	kPawn  = kind(0)
	kKing  = kind(1)
)

var crowningRow = [2]byte{
	int(kBlack): 0,
	int(kWhite): 7,
}

var forward = [2]int8{
	int(kBlack): -1,
	int(kWhite): +1,
}

func (c color) String() string {
	if c == kWhite {
		return "white"
	} else {
		return "black"
	}
}

func (k kind) String() string {
	if k == kKing {
		return "king"
	} else {
		return "pawn"
	}
}

// the board is 8x8, so we can just use a uint64 for each property that might be true or false for a tile
type board struct {
	occupied uint64
	white    uint64
	king     uint64
}

func pieceToRune(c color, k kind) rune {
	if c == kWhite {
		if k == kKing {
			return '@'
		}
		return 'o'
	}
	if k == kKing {
		return '#'
	}
	return 'x'
}

func (b *board) String() string {
	buf := new(bytes.Buffer)
	sep := ""
	for row := byte(0); row < 8; row++ {
		buf.WriteString(sep)
		for col := byte(0); col < 8; col++ {
			if b.isOccupied(row, col) {
				buf.WriteRune(pieceToRune(b.get(row, col)))
			} else {
				buf.WriteRune('.')
			}
		}
		sep = "\n"
	}
	return buf.String()
}

func tileColor(row, col byte) color {
	if (row+col)%2 == 0 {
		return kWhite
	} else {
		return kBlack
	}
}

func placeInitialPieces(b *board) {
	for row := byte(0); row <= 2; row++ {
		for col := byte(0); col < 8; col++ {
			if tileColor(row, col) == kBlack {
				b.set(row, col, kWhite, kPawn)
			}
		}
	}
	for row := byte(5); row <= 7; row++ {
		for col := byte(0); col < 8; col++ {
			if tileColor(row, col) == kBlack {
				b.set(row, col, kBlack, kPawn)
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

	if c == kWhite {
		b.white |= x
	} else {
		b.white &^= x
	}

	if k == kKing {
		b.king |= x
	} else {
		b.king &^= x
	}
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
