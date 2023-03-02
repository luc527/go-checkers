package main

import (
	"math/rand"
)

func rn8() uint8 {
	return uint8(rand.Uint32() % 8)
}

func rnColor() color {
	return color(rand.Uint32() % 2)
}

func rnKind() kind {
	return kind(rand.Uint32() % 2)
}

func randomAction(b boardInterface) {
	x := rand.Uint32() % 5
	if x == 0 { // happens 1/5 of the times
		b.clear(rn8(), rn8())
	} else if x == 1 || x == 2 || x == 3 { // happens 3/5 of the times
		b.set(rn8(), rn8(), rnColor(), rnKind())
	} else { // x == 4, x == 5, happens 2/5 of the time

		// flip

		r, c := rn8(), rn8()
		if b.isOccupied(r, c) {
			color, kind := b.get(r, c)

			if color == kWhite {
				color = kBlack
			} else {
				color = kWhite
			}

			if kind == kKing {
				kind = kPawn
			} else {
				kind = kKing
			}

			b.set(r, c, color, kind)
		}
	}
}

func nRandomActions(b boardInterface, n int) {
	for n > 0 {
		randomAction(b)
		n--
	}
}

func rnInstruction(b *board) instruction {

	type coord struct {
		row, col byte
	}

	var occupied []coord

	anyFree := false

	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if b.isOccupied(row, col) {
				occupied = append(occupied, coord{row, col})
			} else {
				anyFree = true
			}
		}
	}

	rnOccupied := func() coord {
		return occupied[int(rand.Uint32())%len(occupied)]
	}

	rnFreeIfAny := func() coord {
		if !anyFree {
			return rnOccupied()
		}
		row, col := rn8(), rn8()
		for b.isOccupied(row, col) {
			row, col = rn8(), rn8()
		}
		return coord{row, col}
	}

	t := instructionType(rand.Uint32() % 3)
	switch t {
	case moveInstruction:
		src := rnOccupied()
		dst := rnFreeIfAny()
		return makeMoveInstruction(src.row, src.col, dst.row, dst.col)
	case captureInstruction:
		src := rnOccupied()
		capColor, capKind := b.get(src.row, src.col)
		return makeCaptureInstruction(src.row, src.col, capColor, capKind)
	case crownInstruction:
		c := rnOccupied()
		return makeCrownInstruction(c.row, c.col)
	default:
		panic("unreachable")
	}

}
