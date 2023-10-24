package core

import (
	"bytes"
	"strings"
)

// TODO maybe there's some way to optimize this using the bit masks

// both offSets
var offBoth = [2]int8{-1, +1}

// some types to avoid boolean blindness
// since we don't have proper sum types
// e.g. NewCustomGame(true, false, ...) true or false what?

type CaptureRule bool

const (
	// 'capturesMandatory' means that if you have captures available, you must perform one of the captures, not the simple plies
	CapturesMandatory    = CaptureRule(true)
	CapturesNotMandatory = CaptureRule(false)
)

type BestRule bool

const (
	// 'BestMandatory' means that you must perform the best capture available (the one that captures the most pieces)
	BestMandatory    = BestRule(true)
	BestNotMandatory = BestRule(false)
)

// Ply: a single action a player can execute when it's their turn to play
type Ply []Instruction

func (p Ply) String() string {
	return instructionsToString(p)
}

func (p Ply) countCaptures() int {
	c := 0
	for _, i := range p {
		if i.t == CaptureInstruction {
			c++
		}
	}
	return c
}

func (p Ply) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	sep := ""
	buf.WriteByte('"')
	for _, i := range p {
		if _, err := buf.WriteString(sep); err != nil {
			return nil, err
		}
		if err := i.SerializeInto(&buf); err != nil {
			return nil, err
		}
		sep = ","
	}
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

func (p *Ply) UnmarshalJSON(bs []byte) error {
	s := string(bs)
	s = strings.Trim(s, " \"\t\n\r")
	if len(s) == 0 {
		*p = Ply{}
		return nil
	}
	sis := strings.Split(s, ",")
	for _, si := range sis {
		i := &Instruction{}
		if err := i.Unserialize([]byte(si)); err != nil {
			return err
		}
		*p = append(*p, *i)
	}
	return nil
}

func (p Ply) Equals(q Ply) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	if len(p) != len(q) {
		return false
	}
	for i, ins := range p {
		if !ins.Equals(q[i]) {
			return false
		}
	}
	return true
}

func (p Ply) Copy() Ply {
	is := make([]Instruction, len(p))
	copy(is, p)
	return is
}

func CopyPlies(ps []Ply) []Ply {
	a := make([]Ply, len(ps))
	for i, p := range ps {
		a[i] = p.Copy()
	}
	return a
}

// This function doesn't really do what is says. For the two ply lists to be equal,
// they not only need to contain the same plies, but also in the same order.
// Although it's mostly used in tests and it works fine there.
func PliesEquals(ps []Ply, qs []Ply) bool {
	if ps == nil && qs == nil {
		return true
	}
	if ps == nil || qs == nil {
		return false
	}
	if len(ps) != len(qs) {
		return false
	}
	for i, p := range ps {
		if !p.Equals(qs[i]) {
			return false
		}
	}
	return true
}

// the generateSimplePawnPlies and followPawnCaptures procedures
// are special cases of the same procedures for king pieces
// that have the distance bound to 1 (simple move) or 2 (capture)
// and direction bounded by forward[color] (simple move)
// you COULD make something like maxDist := king ? 100 : 2 and use the same procedure
// but my hope is that these pawn-specific versions are faster

// simple plies are ones not involving any captures, where the piece just moves

func generateSimplePawnPlies(ps []Ply, b *Board, row, col byte, color Color) []Ply {
	drow := byte(int8(row) + forward[color])
	if drow >= 8 {
		return ps
	}
	crown := crowningRow[color] == drow
	for _, dir := range offBoth {
		dcol := byte(int8(col) + dir)
		if dcol >= 8 || b.IsOccupied(drow, dcol) {
			continue
		}
		var is []Instruction
		is = append(is, MakeMoveInstruction(row, col, drow, dcol))
		if crown {
			is = append(is, MakeCrownInstruction(drow, dcol))
		}
		ps = append(ps, Ply(is))
	}
	return ps
}

func generateSimpleKingPlies(ps []Ply, b *Board, row, col byte, color Color) []Ply {
	for _, roff := range offBoth {
		for _, coff := range offBoth {
			dist := int8(1)
			for {
				drow, dcol := byte(int8(row)+dist*roff), byte(int8(col)+dist*coff)
				if drow >= 8 || dcol >= 8 || b.IsOccupied(drow, dcol) {
					break
				}

				is := []Instruction{MakeMoveInstruction(row, col, drow, dcol)}
				ps = append(ps, Ply(is))

				dist++
			}
		}
	}

	return ps
}

func generateSimplePlies(ps []Ply, b *Board, player Color) []Ply {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.IsOccupied(row, col) {
				continue
			}

			color, kind := b.Get(row, col)
			if color != player {
				continue
			}

			if kind == PawnKind {
				ps = generateSimplePawnPlies(ps, b, row, col, color)
			} else {
				ps = generateSimpleKingPlies(ps, b, row, col, color)
			}
		}
	}

	return ps
}

// after calling generateCapturePlies the board is the same as before calling it,
// but the procedures alter the board in order to generate the captures
// (we need to do a tree search in order to generate all possibilities of
// sequential captures, and we do that by backtracking)

func followPawnCaptures(ps []Ply, stack []Instruction, b *Board, row, col byte, color Color) []Ply {
	// sink: there are no more captures available from here
	sink := true

	for _, roff := range offBoth {
		for _, coff := range offBoth {
			drow, dcol := byte(int8(row)+2*roff), byte(int8(col)+2*coff)
			if drow >= 8 || dcol >= 8 || b.IsOccupied(drow, dcol) {
				continue
			}
			mrow, mcol := byte(int8(row)+roff), byte(int8(col)+coff)
			if !b.IsOccupied(mrow, mcol) {
				continue
			}
			mcolor, mkind := b.Get(mrow, mcol)
			if mcolor == color {
				continue
			}

			sink = false

			// do
			stack = append(stack, MakeMoveInstruction(row, col, drow, dcol))
			stack = append(stack, MakeCaptureInstruction(mrow, mcol, mcolor, mkind))
			b.Move(row, col, drow, dcol)
			b.Clear(mrow, mcol)

			ps = followPawnCaptures(ps, stack, b, drow, dcol, color)

			// undo
			b.Set(mrow, mcol, mcolor, mkind)
			b.Move(drow, dcol, row, col)
			stack = stack[:len(stack)-2]
		}
	}

	// TODO cleanup, someMeaningfulName := stack != nil, then if sink && someMeaningfulName
	// stack is nil at the first call when no captures have been made yet
	if sink && stack != nil {
		isLen := len(stack)
		crown := row == crowningRow[color]
		if crown {
			isLen += 1
		}
		is := make([]Instruction, isLen)
		copy(is, stack)
		if crown {
			is[isLen-1] = MakeCrownInstruction(row, col)
		}
		ps = append(ps, Ply(is))
	}
	return ps
}

func followKingCaptures(ps []Ply, stack []Instruction, b *Board, row, col byte, player Color) []Ply {
	sink := true

	for _, roff := range offBoth {
		for _, coff := range offBoth {

			pastCapture := false // currently iterating through positions after the captured one

			// captured piece, if any (only one)
			var crow, ccol byte
			var ccolor Color
			var ckind Kind

			dist := int8(1)
			for {

				// i for [i]teration row and col, for lack of a better single-letter abbreviation
				// TODO implement by irow, icol = irow+roff, icol+coff? would it be any more efficient?
				irow, icol := byte(int8(row)+dist*roff), byte(int8(col)+dist*coff)

				if irow >= 8 || icol >= 8 {
					break
				}

				if b.IsOccupied(irow, icol) {
					if pastCapture {
						break
					}
					icolor, ikind := b.Get(irow, icol)
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
					stack = append(stack, MakeMoveInstruction(row, col, irow, icol))
					stack = append(stack, MakeCaptureInstruction(crow, ccol, ccolor, ckind))
					b.Move(row, col, irow, icol)
					b.Clear(crow, ccol)

					ps = followKingCaptures(ps, stack, b, irow, icol, player)

					// undo
					b.Set(crow, ccol, ccolor, ckind)
					b.Move(irow, icol, row, col)
					stack = stack[:len(stack)-2]
				}

				dist++
			}

		}
	}

	// same code as in followSimpleCaptures, except no crowning since the piece is already a king
	if sink && stack != nil {
		is := make([]Instruction, len(stack))
		copy(is, stack)
		ps = append(ps, Ply(is))
	}

	return ps
}

func generateCapturePlies(ps []Ply, b *Board, player Color) []Ply {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.IsOccupied(row, col) {
				continue
			}

			color, kind := b.Get(row, col)
			if color != player {
				continue
			}

			if kind == PawnKind {
				ps = followPawnCaptures(ps, nil, b, row, col, color)
			} else {
				ps = followKingCaptures(ps, nil, b, row, col, color)
			}
		}
	}
	return ps
}

func GeneratePlies(ps []Ply, b *Board, player Color) []Ply {
	ps = generateCapturePlies(ps, b, player)
	if len(ps) == 0 {
		ps = generateSimplePlies(ps, b, player)
	}
	return ps
}
