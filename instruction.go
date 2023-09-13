package go_checkers

import (
	"bytes"
	"fmt"
	"strings"
)

type InstructionType byte

const (
	MoveInstruction = InstructionType(iota)
	CaptureInstruction
	CrownInstruction
)

func (t InstructionType) String() string {
	switch t {
	case MoveInstruction:
		return "move"
	case CaptureInstruction:
		return "capture"
	case CrownInstruction:
		return "crown"
	default:
		return "UNKNOWN"
	}
}

// instruction should contain all information needed to undo it
// that's the reason for storing the color and kind of the captured piece

type Instruction struct {
	t   InstructionType
	row byte
	col byte
	// arbitraty data
	// destination row and col for MoveInstruction
	// color and kind for CaptureInstruction
	// unused for CrownInstruction
	d [2]byte
}

func (i Instruction) Equals(o Instruction) bool {
	return i.t == o.t &&
		i.row == o.row &&
		i.col == o.col &&
		i.d[0] == o.d[0] &&
		i.d[1] == o.d[1]
}

func (i Instruction) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "{%s (%d, %d)", i.t.String(), i.row, i.col)
	if i.t == MoveInstruction {
		fmt.Fprintf(buf, " to (%d, %d)", i.d[0], i.d[1])
	} else if i.t == CaptureInstruction {
		fmt.Fprintf(buf, " %s %s", Color(i.d[0]), Kind(i.d[1]))
	}
	buf.WriteRune('}')
	return buf.String()
}

func instructionsToString(is []Instruction) string {
	ss := make([]string, 0, len(is))
	for _, i := range is {
		ss = append(ss, i.String())
	}
	return strings.Join(ss, ";")
}

func MakeMoveInstruction(sourceRow, sourceCol, destinationRow, destinationCol byte) Instruction {
	var i Instruction
	i.t = MoveInstruction
	i.row, i.col = sourceRow, sourceCol
	i.d[0], i.d[1] = destinationRow, destinationCol
	return i
}

func MakeCaptureInstruction(row, col byte, c Color, k Kind) Instruction {
	var i Instruction
	i.t = CaptureInstruction
	i.row, i.col = row, col
	i.d[0], i.d[1] = byte(c), byte(k)
	return i
}

func MakeCrownInstruction(row, col byte) Instruction {
	var i Instruction
	i.t = CrownInstruction
	i.row, i.col = row, col
	return i
}

func PerformInstructions(b *Board, is []Instruction) error {
	for _, i := range is {
		switch i.t {
		case MoveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			b.Move(fromRow, fromCol, toRow, toCol)
		case CaptureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := Color(i.d[0]), Kind(i.d[1])
			actualColor, actualKind := b.Get(row, col)
			if capturedColor != actualColor || capturedKind != actualKind {
				return fmt.Errorf(
					"performed capture instruction of %s %s on row %d %d but piece is a %s %s",
					capturedColor, capturedKind,
					row, col,
					actualColor, actualKind,
				)
			}
			b.Clear(row, col)
		case CrownInstruction:
			b.Crown(i.row, i.col)
		default:
			return fmt.Errorf("invalid instruction type %s", i.t)
		}
	}
	return nil
}

// if you pass 'is' to PerformInstructions
// the same 'is' passed to UndoInstructions will undo the instructions performed
// you don't need to reverse them
func UndoInstructions(b *Board, is []Instruction) {
	for k := len(is) - 1; k >= 0; k-- {
		i := is[k]
		switch i.t {
		case MoveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			b.Move(toRow, toCol, fromRow, fromCol)
		case CaptureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := Color(i.d[0]), Kind(i.d[1])
			b.Set(row, col, capturedColor, capturedKind)
		case CrownInstruction:
			b.Uncrown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}
