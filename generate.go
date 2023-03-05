package main

var offBoth = [2]int8{-1, +1}

func generateSimplePawnMoves(iss [][]instruction, b *board, row, col byte, color color) [][]instruction {
	drow := byte(int8(row) + forward[color])
	if drow >= 8 {
		return iss
	}
	crown := crowningRow[color] == drow
	for _, dir := range offBoth {
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
	for _, roff := range offBoth {
		for _, coff := range offBoth {
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
			if kind == pawnKind {
				iss = generateSimplePawnMoves(iss, b, row, col, color)
			} else {
				iss = generateSimpleKingMoves(iss, b, row, col, color)
			}
		}
	}

	return iss
}

// WARNING
// generating captures modifies the board!

func followPawnCaptures(iss [][]instruction, stack []instruction, b *board, row, col byte, color color) [][]instruction {
	sink := true

	for _, roff := range offBoth {
		for _, coff := range offBoth {
			drow, dcol := byte(int8(row)+2*roff), byte(int8(col)+2*coff)
			if drow >= 8 || dcol >= 8 || b.isOccupied(drow, dcol) {
				continue
			}
			mrow, mcol := byte(int8(row)+roff), byte(int8(col)+coff)
			if !b.isOccupied(mrow, mcol) {
				continue
			}
			mcolor, mkind := b.get(mrow, mcol)
			if mcolor == color {
				continue
			}

			sink = false

			// do
			stack = append(stack, makeMoveInstruction(row, col, drow, dcol))
			stack = append(stack, makeCaptureInstruction(mrow, mcol, mcolor, mkind))
			b.move(row, col, drow, dcol)
			b.clear(mrow, mcol)

			iss = followPawnCaptures(iss, stack, b, drow, dcol, color)

			// undo
			b.set(mrow, mcol, mcolor, mkind)
			b.move(drow, dcol, row, col)
			stack = stack[:len(stack)-2]

			// I could make another stack variable 'substack', copy the slice stack to it, append
			// the instructions to that one and pass it to the recursive call,
			// so I don't need to shrink the stack at the end,
			// BUT this can be less efficient if appending the substack grows the slice:
			// currently if the stack grows once it can use that leftover capacity in further recursive calls;
			// if we used a substack and it grew, it'd grow again in the next iteration of the loop
		}
	}

	if sink && stack != nil {
		isLen := len(stack)
		crown := row == crowningRow[color]
		if crown {
			isLen += 1
		}
		is := make([]instruction, isLen)
		copy(is, stack)
		if crown {
			is[isLen-1] = makeCrownInstruction(row, col)
		}
		iss = append(iss, is)
	}
	return iss
}

func generatePawnCaptureMoves(iss [][]instruction, b *board, row, col byte, color color) [][]instruction {
	return followPawnCaptures(iss, nil, b, row, col, color)
}

func generateKingCaptureMoves(iss [][]instruction, b *board, row, col byte, color color) [][]instruction {
	return iss
}

func generateCaptureMoves(iss [][]instruction, b *board) [][]instruction {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.isOccupied(row, col) {
				continue
			}

			color, kind := b.get(row, col)
			if kind == pawnKind {
				iss = generatePawnCaptureMoves(iss, b, row, col, color)
			} else {
				iss = generateKingCaptureMoves(iss, b, row, col, color)
			}
		}
	}
	return iss
}
