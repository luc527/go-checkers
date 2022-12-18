package game

func GenerateSimpleMoves(board *Board, player Color) []SimpleMove {
	moves := make([]SimpleMove, 0, 8)

	for bit, mask := uint8(0), uint64(1); bit < 64; bit, mask = bit+1, mask<<1 {

		if board.HasPiece&mask == 0 {
			continue
		}

		row, col := bit/8, bit%8
		color, kind := board.Get(row, col)
		if color != player {
			continue
		}

		if kind == Pawn {
			// more specialized code for the common case

			destRow := row - 1
			if color == Black {
				destRow = row + 1
			}

			// -1 overflows, so we don't need to check <0 (always false)
			if destRow > 7 {
				continue
			}

			destLeftCol := col - 1
			if destLeftCol < 8 && !board.Occupied(destRow, destLeftCol) {
				moves = append(moves, MakeSimpleMove(row, col, destRow, destLeftCol))
			}

			destRightCol := col + 1
			if destRightCol < 8 && !board.Occupied(destRow, destRightCol) {
				moves = append(moves, MakeSimpleMove(row, col, destRow, destRightCol))
			}
		} else { // King
			for _, rowOffset := range [2]int8{-1, 1} {
				for _, colOffset := range [2]int8{-1, 1} {
					destRow := uint8(int8(row) + rowOffset)
					destCol := uint8(int8(col) + colOffset)
					for destRow < 8 && destCol < 8 && !board.Occupied(destRow, destCol) {
						moves = append(moves, MakeSimpleMove(row, col, destRow, destCol))
						destRow = uint8(int8(destRow) + rowOffset)
						destCol = uint8(int8(destCol) + colOffset)
					}
				}
			}
		}

	}

	return moves
}
