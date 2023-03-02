package main

type board struct {
	occupied uint64
	white    uint64
	king     uint64
}

func newEmptyBoard() *board {
	return &board{}
}

func tileColor(row, col uint8) color {
	if (row+col)%2 == 0 {
		return kWhite
	} else {
		return kBlack
	}
}

func placeInitialPieces(b *board) {
	for row := uint8(0); row <= 2; row++ {
		for col := uint8(0); col < 8; col++ {
			if tileColor(row, col) == kBlack {
				b.set(row, col, kWhite, kPawn)
			}
		}
	}
	for row := uint8(5); row <= 7; row++ {
		for col := uint8(0); col < 8; col++ {
			if tileColor(row, col) == kBlack {
				b.set(row, col, kBlack, kPawn)
			}
		}
	}
}

func (b *board) clear(row, col uint8) {
	b.occupied &^= uint64(1 << (uint64(row)*8 + uint64(col)))
}

func (b *board) set(row, col uint8, c color, k kind) {
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

func (b *board) crown(row, col uint8) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king |= x
}

func (b *board) uncrown(row, col uint8) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	b.king &^= x
}

func (b *board) isOccupied(row, col uint8) bool {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	return b.occupied&x != 0
}

func (b *board) get(row, col uint8) (c color, k kind) {
	x := uint64(1 << (uint64(row)*8 + uint64(col)))
	k = kind(b.king & x)
	c = color(b.white & x)
	return
}

func (b *board) copy() *board {
	var c board
	c.occupied = b.occupied
	c.white = b.white
	c.king = b.king
	return &c
}
