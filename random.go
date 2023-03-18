package main

import (
	"math/rand"
)

func rn8() uint8 {
	return uint8(rand.Uint32() % 8)
}

func rnColor() Color {
	return Color(rand.Uint32() % 2)
}

func rnKind() Kind {
	return Kind(rand.Uint32() % 2)
}

func randomAction(b boardInterface) {
	x := rand.Uint32() % 5
	if x == 0 { // happens 1/5 of the times
		b.Clear(rn8(), rn8())
	} else if x == 1 || x == 2 || x == 3 { // happens 3/5 of the times
		b.Set(rn8(), rn8(), rnColor(), rnKind())
	} else { // x == 4, x == 5, happens 2/5 of the time

		// flip

		r, c := rn8(), rn8()
		if b.IsOccupied(r, c) {
			color, kind := b.Get(r, c)

			if color == WhiteColor {
				color = BlackColor
			} else {
				color = WhiteColor
			}

			if kind == KingKind {
				kind = PawnKind
			} else {
				kind = KingKind
			}

			b.Set(r, c, color, kind)
		}
	}
}

func nRandomActions(b boardInterface, n int) {
	for n > 0 {
		randomAction(b)
		n--
	}
}

func randomInoffensiveMove(b *Board, player Color) Ply {
	var coords []coord

	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.IsOccupied(row, col) {
				continue
			}

			color, _ := b.Get(row, col)
			if color == player {
				coords = append(coords, coord{row, col})
			}
		}
	}

	if len(coords) == 0 {
		return Ply{}
	}

	randomCoord := coords[rand.Int()%len(coords)]
	srow, scol := randomCoord.row, randomCoord.col

	var drow, dcol byte
	for {
		drow, dcol = rn8(), rn8()
		if !b.IsOccupied(drow, dcol) {
			break
		}
	}

	return Ply{MoveInstruction(srow, scol, drow, dcol)}
}
