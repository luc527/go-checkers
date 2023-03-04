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
