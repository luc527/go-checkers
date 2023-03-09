package main

var offBoth = [2]int8{-1, +1}

type captureRule bool

const (
	capturesMandatory    = captureRule(true)
	capturesNotMandatory = captureRule(false)
)

type bestRule bool

const (
	bestMandatory    = bestRule(true)
	bestNotMandatory = bestRule(false)
)

// ply: a single action a player can execute when it's their turn to play
type ply []instruction

func (p ply) String() string {
	return instructionsToString(p)
}

func generateSimplePawnPlies(ps []ply, b *board, row, col byte, color color) []ply {
	drow := byte(int8(row) + forward[color])
	if drow >= 8 {
		return ps
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
		ps = append(ps, ply(is))
	}
	return ps
}

func generateSimpleKingPlies(ps []ply, b *board, row, col byte, color color) []ply {
	for _, roff := range offBoth {
		for _, coff := range offBoth {
			dist := int8(1)
			for {
				drow, dcol := byte(int8(row)+dist*roff), byte(int8(col)+dist*coff)
				if drow >= 8 || dcol >= 8 || b.isOccupied(drow, dcol) {
					break
				}

				is := []instruction{makeMoveInstruction(row, col, drow, dcol)}
				ps = append(ps, ply(is))

				dist++
			}
		}
	}

	return ps
}

func generateSimplePlies(ps []ply, b *board, player color) []ply {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.isOccupied(row, col) {
				continue
			}

			color, kind := b.get(row, col)
			if color != player {
				continue
			}

			if kind == pawnKind {
				ps = generateSimplePawnPlies(ps, b, row, col, color)
			} else {
				ps = generateSimpleKingPlies(ps, b, row, col, color)
			}
		}
	}

	return ps
}

// WARNING
// generating captures modifies the board!

func followPawnCaptures(ps []ply, stack []instruction, b *board, row, col byte, color color) []ply {
	// sink: there are no more captures available from here
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

			ps = followPawnCaptures(ps, stack, b, drow, dcol, color)

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

	// stack is nil on the first call where no captures have been made yet
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
		ps = append(ps, ply(is))
	}
	return ps
}

func followKingCaptures(ps []ply, stack []instruction, b *board, row, col byte, player color) []ply {
	sink := true

	for _, roff := range offBoth {
		for _, coff := range offBoth {

			pastCapture := false // currently iterating through positions after the captured one

			// captured piece, if any (only one)
			var crow, ccol byte
			var ccolor color
			var ckind kind

			dist := int8(1)
			for {

				// iteration row and col, for lack of a better single-letter abbreviation
				irow, icol := byte(int8(row)+dist*roff), byte(int8(col)+dist*coff)

				if irow >= 8 || icol >= 8 {
					break
				}

				if b.isOccupied(irow, icol) {
					if pastCapture {
						break
					}
					icolor, ikind := b.get(irow, icol)
					if icolor != player {
						// this is the capture
						pastCapture = true
						crow, ccol = irow, icol
						ccolor, ckind = icolor, ikind
					}
				} else if pastCapture {
					// this is a destination
					sink = false

					// do
					stack = append(stack, makeMoveInstruction(row, col, irow, icol))
					stack = append(stack, makeCaptureInstruction(crow, ccol, ccolor, ckind))
					b.move(row, col, irow, icol)
					b.clear(crow, ccol)

					ps = followKingCaptures(ps, stack, b, irow, icol, player)

					// undo
					b.set(crow, ccol, ccolor, ckind)
					b.move(irow, icol, row, col)
					stack = stack[:len(stack)-2]
				}

				dist++
			}

		}
	}

	// same code as in followSimpleCaptures, except no crowning since the piece is already a king
	if sink && stack != nil {
		is := make([]instruction, len(stack))
		copy(is, stack)
		ps = append(ps, ply(is))
	}

	return ps
}

func generateCapturePlies(ps []ply, b *board, player color) []ply {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.isOccupied(row, col) {
				continue
			}

			color, kind := b.get(row, col)
			if color != player {
				continue
			}

			if kind == pawnKind {
				ps = followPawnCaptures(ps, nil, b, row, col, color)
			} else {
				ps = followKingCaptures(ps, nil, b, row, col, color)
			}
		}
	}
	return ps
}

func generatePlies(b *board, player color, captureRule captureRule, bestRule bestRule) []ply {

	ps := generateCapturePlies(nil, b, player)

	capturesMandatory := captureRule == capturesMandatory
	bestMandatory := bestRule == bestMandatory

	if len(ps) == 0 || (!capturesMandatory && !bestMandatory) {
		ps = generateSimplePlies(ps, b, player)
	}

	if len(ps) > 0 && bestMandatory {
		captureCountPerMove := make([]int, len(ps))
		mostCaptures := 0
		for k, p := range ps {
			captureCount := 0
			for _, i := range p {
				if i.t == captureInstruction {
					captureCount++
				}
			}
			captureCountPerMove[k] = captureCount
			if captureCount > mostCaptures {
				mostCaptures = captureCount
			}
		}

		var best []ply
		for k, p := range ps {
			if captureCountPerMove[k] == mostCaptures {
				best = append(best, p)
			}
		}
		ps = best
	}

	return ps
}
