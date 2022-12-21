package game

var directionUp = []int8{-1}
var directionDown = []int8{1}
var directionsBoth = []int8{1, -1}

func GenerateSimpleMoves(board *Board, player Color) []SimpleMove {

	// TODO we generate moves much more often than read all board positions sequentially
	// would make more sense for the board to store a list of white positions with cells and a list of black positions with cells

	// this way we fill the cache with just what we need, instead of filling with the whole board

	var moves []SimpleMove

	for row := int8(0); row < 8; row++ {
		for col := int8(0); col < 8; col++ {

			if !board.Occupied(row, col) {
				continue
			}

			color, kind := board.Get(row, col)
			if color != player {
				continue
			}

			reach := int8(1)
			if kind == King {
				reach = 127 // unlimited reach
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
						destRow := row + distance*rowOffset
						destCol := col + distance*colOffset

						if !Inbounds(destRow, destCol) {
							break
						}

						if board.Occupied(destRow, destCol) {
							break
						}

						crown := IsCrowningRow(player, destRow)

						moves = append(moves, SimpleMove{move{coord{row, col}, coord{destRow, destCol}}, crown})
					}

				}
			}

		}
	}

	return moves
}
