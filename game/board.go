package game

import (
	"fmt"
	"strings"
)

type Color bool

const (
	Black = Color(false)
	White = Color(true)
)

func (c Color) String() string {
	if c == Black {
		return "Black"
	} else {
		return "White"
	}
}

type Kind bool

const (
	Pawn = Kind(false)
	King = Kind(true)
)

func (k Kind) String() string {
	if k == Pawn {
		return "Pawn"
	} else {
		return "King"
	}
}

func CellColor(row, col int8) Color {
	if (row+col)%2 == 0 {
		return White
	} else {
		return Black
	}
}

//
// Board
//

type Board struct {
	hasPiece uint64
	cells0   uint64
	cells1   uint64
}

func (b *Board) Occupied(row, col int8) bool {
	if !Inbounds(row, col) {
		panic(fmt.Sprintf("(%v, %v) out of bounds", row, col))
	}
	i := row*8 + col
	return b.hasPiece&(1<<i) != 0
}

func (b *Board) Clear(row, col int8) {
	i := row*8 + col
	b.hasPiece &^= (1 << i)
}

func (b *Board) Get(row, col int8) (color Color, kind Kind) {

	// for catching programming errors, could be removed later
	if !b.Occupied(row, col) {
		panic(fmt.Sprintf("calling Get on empty board cell (%v, %v)", row, col))
	}

	// hopefully cmovs
	cells := b.cells0
	r := row
	if row >= 4 {
		cells = b.cells1
		r = row - 4
	}

	pow := uint64(1 << (r*16 + col*2))

	color = Color(cells&pow != 0)
	kind = Kind(cells&(pow<<1) != 0)
	return
}

// get and clear
func (b *Board) Take(row, col int8) (color Color, kind Kind) {
	// TODO could inline + optimize?
	color, kind = b.Get(row, col)
	b.Clear(row, col)
	return
}

func (b *Board) Set(row, col int8, color Color, kind Kind) {

	b.hasPiece |= 1 << (row*8 + col)

	on0 := row < 4
	row %= 4

	base := row*16 + col*2

	pat := uint64(0)
	if color {
		pat |= 1 << base
	}
	if kind {
		pat |= 1 << (base + 1)
	}

	if on0 {
		b.cells0 &^= 3 << base
		b.cells0 |= pat
	} else {
		b.cells1 &^= 3 << base
		b.cells1 |= pat
	}
}

func (board *Board) String() string {
	var sb strings.Builder
	for row := int8(0); row < 8; row++ {
		for col := int8(0); col < 8; col++ {
			if !board.Occupied(row, col) {
				if CellColor(row, col) == Black {
					sb.WriteRune('_')
				} else {
					sb.WriteRune(' ')
				}
			} else {
				var char rune
				color, kind := board.Get(row, col)
				if color == Black && kind == Pawn {
					char = '*'
				} else if color == White && kind == Pawn {
					char = 'o'
				} else if color == Black && kind == King {
					char = '#'
				} else if color == White && kind == King {
					char = '@'
				}
				sb.WriteRune(char)
			}
		}
		if row < 7 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func InitialBoard() *Board {
	var b Board

	for row := int8(0); row < 3; row++ {
		for col := int8(0); col < 8; col++ {
			if CellColor(row, col) == Black {
				b.Set(row, col, Black, Pawn)
			}
		}
	}

	for row := int8(5); row < 8; row++ {
		for col := int8(0); col < 8; col++ {
			if CellColor(row, col) == Black {
				b.Set(row, col, White, Pawn)
			}
		}
	}

	return &b
}

func IsCrowningRow(color Color, row int8) bool {
	return (row == 0 && color == White) || (row == 7 && color == Black)
}

func Inbounds(row, col int8) bool {
	// hopefully gets inlined
	return row >= 0 && row < 8 && col >= 0 && col < 8
}
