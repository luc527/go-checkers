package main

var dirBoth = [2]int8{-1, +1}

func generateSimplePawnMoves(iss [][]instruction, b *board, row, col byte, color color) [][]instruction {
	drow := byte(int8(row) + forward[color])
	if drow >= 8 {
		return iss
	}
	crown := crowningRow[color] == drow
	for _, dir := range dirBoth {
		dcol := byte(int8(col) + dir)
		if dcol >= 8 || b.isOccupied(drow, dcol) {
			continue
		}
		var is []instruction
		is = append(is, makeMoveInstruction(row, col, drow, dcol))
		if crown {
			is = append(is, makeCrownInstruction(drow, dcol))
		}
		iss = append(iss, is)
	}
	return iss
}

func generateSimpleKingMoves(iss [][]instruction, b *board, row, col byte, color color) [][]instruction {
	for _, roff := range dirBoth {
		for _, coff := range dirBoth {
			dist := int8(1)
			for {
				drow, dcol := byte(int8(row)+dist*roff), byte(int8(col)+dist*coff)
				if drow >= 8 || dcol >= 8 || b.isOccupied(drow, dcol) {
					break
				}

				is := []instruction{makeMoveInstruction(row, col, drow, dcol)}
				iss = append(iss, is)

				dist++
			}
		}
	}

	return iss
}

// should always be called like append:
// iss = generateSimpleMoves(iss, b)
func generateSimpleMoves(iss [][]instruction, b *board) [][]instruction {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.isOccupied(row, col) {
				continue
			}

			color, kind := b.get(row, col)
			if kind == kPawn {
				iss = generateSimplePawnMoves(iss, b, row, col, color)
			} else {
				iss = generateSimpleKingMoves(iss, b, row, col, color)
			}
		}
	}

	return iss
}
