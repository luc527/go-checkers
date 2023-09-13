package checkers

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
