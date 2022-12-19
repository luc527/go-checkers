package game

import "fmt"

// First the simple implementation
// Then optimize it with bit fiddling and so on
// so we actually get to compare and see how much of an improvement we get from doing that

var directionUp = []int8{-1}
var directionDown = []int8{1}
var directionsBoth = []int8{1, -1}

type SimpleMove struct {
	FromRow uint8
	FromCol uint8
	ToRow   uint8
	ToCol   uint8
	Crowned bool
}

func (move SimpleMove) String() string {
	c := ""
	if move.Crowned {
		c = "crowned"
	}
	return fmt.Sprintf("(%v, %v) -> (%v, %v) %s", move.FromRow, move.FromCol, move.ToRow, move.ToCol, c)
}

func (move *SimpleMove) Do(board *Board) {
	color, kind := board.Take(move.FromRow, move.FromCol)
	toCrowningRow := (color == White && move.ToRow == 0) || (color == Black && move.ToRow == 7)
	doCrown := kind == Pawn && toCrowningRow
	if doCrown {
		move.Crowned = true
		kind = King
	}
	board.Set(move.ToRow, move.ToCol, color, kind)
}

func (move *SimpleMove) Undo(board *Board) {
	color, kind := board.Take(move.ToRow, move.ToCol)
	uncrown := move.Crowned
	if uncrown {
		kind = Pawn
	}
	board.Set(move.FromRow, move.FromCol, color, kind)
}

func GenerateSimpleMoves(board *Board, player Color) []SimpleMove {

	// TODO we generate moves much more often than read all board positions sequentially
	// would make more sense for the board to store a list of white positions with cells and a list of black positions with cells

	// this way we fill the cache with just what we need, instead of filling with the whole board

	// (although the whole board already probably fits, so I don't know how much of an improvement we would get)

	var moves []SimpleMove

	for row := uint8(0); row < 8; row++ {
		for col := uint8(0); col < 8; col++ {

			if !board.Occupied(row, col) {
				continue
			}

			color, kind := board.Get(row, col)
			if color != player {
				continue
			}

			reach := int8(1)
			if kind == King {
				reach = 127 // arbitrary large number for unlimited reach
			}

			rowOffsets := directionsBoth
			if kind != King {
				rowOffsets = directionUp
				if color != White {
					rowOffsets = directionDown
				}
			}

			for _, rowOffset := range rowOffsets {
				for _, colOffset := range directionsBoth {

					for distance := int8(1); distance <= reach; distance++ {
						destRow := uint8(int8(row) + distance*rowOffset)
						destCol := uint8(int8(col) + distance*colOffset)

						if destRow > 7 || destCol > 7 {
							break
						}

						if board.Occupied(destRow, destCol) {
							break
						}

						moves = append(moves, SimpleMove{row, col, destRow, destCol, false})
					}

				}
			}

		}
	}

	return moves
}
