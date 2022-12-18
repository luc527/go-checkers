package board

import (
	"fmt"
	"math/rand"
	"strings"
)

type CompactBoard struct {
	HasPiece uint64
	Cells0   uint64
	Cells1   uint64
}

type Color bool

const (
	Black = Color(false)
	White = Color(true)
)

type Kind bool

const (
	Pawn = Kind(false)
	King = Kind(true)
)

func (b *CompactBoard) Occupied(row, col int) bool {
	i := row*8 + col
	return b.HasPiece&(1<<i) != 0
}

func (b *CompactBoard) Clear(row, col int) {
	i := row*8 + col
	b.HasPiece &^= (1 << i)
}

func (b *CompactBoard) Get(row, col int) (color Color, kind Kind) {
	// hopefully the compiler turns this into 2 CMOVs
	cells := b.Cells0
	r := row
	if row >= 4 {
		cells = b.Cells1
		r = row - 4
	}

	pow := uint64(1 << (r*16 + col*2))

	color = Color(cells&pow != 0)
	kind = Kind(cells&(pow<<1) != 0)
	return
}

func (b *CompactBoard) Set(row, col int, color Color, kind Kind) {

	b.HasPiece |= 1 << (row*8 + col)

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
		b.Cells0 &^= 3 << base
		b.Cells0 |= pat
	} else {
		b.Cells1 &^= 3 << base
		b.Cells1 |= pat
	}
}

func CellColor(row, col int) Color {
	if (row+col)%2 == 0 {
		return White
	} else {
		return Black
	}
}

func (board *CompactBoard) String() string {
	var sb strings.Builder
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if !board.Occupied(row, col) {
				if CellColor(row, col) == Black {
					sb.WriteRune('.')
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
		sb.WriteRune('\n')
	}
	return sb.String()
}

func InitialBoard() *CompactBoard {
	var b CompactBoard

	for row := 0; row < 3; row++ {
		for col := 0; col < 8; col++ {
			if CellColor(row, col) == Black {
				b.Set(row, col, Black, Pawn)
			}
		}
	}

	for row := 5; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if CellColor(row, col) == Black {
				b.Set(row, col, White, Pawn)
			}
		}
	}

	return &b
}

func RandomBoard() *CompactBoard {
	var b CompactBoard

	b.HasPiece = rand.Uint64()
	b.Cells0 = rand.Uint64()
	b.Cells1 = rand.Uint64()

	return &b
}

func (b *CompactBoard) Debug() {
	fmt.Println("HasPiece")
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			has := (b.HasPiece & (1 << (r*8 + c))) != 0
			if has {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}

	fmt.Println("Cells0")
	for r := 0; r < 4; r++ {
		for c := 0; c < 8; c++ {
			a := (b.Cells0 & (1 << (r*16 + c*2))) != 0
			b := (b.Cells0 & (1 << (1 + (r*16 + c*2)))) != 0
			if a {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
			if b {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}

	fmt.Println("Cells1")
	for r := 0; r < 4; r++ {
		for c := 0; c < 8; c++ {
			a := (b.Cells1 & (1 << (r*16 + c*2))) != 0
			b := (b.Cells1 & (1 << (1 + (r*16 + c*2)))) != 0
			if a {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
			if b {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}
