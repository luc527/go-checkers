package main

import (
	"bytes"
	"fmt"
	"strings"
)

type instructionType byte

const (
	moveInstruction = instructionType(iota)
	captureInstruction
	crownInstruction
)

func (t instructionType) String() string {
	switch t {
	case moveInstruction:
		return "move"
	case captureInstruction:
		return "capture"
	case crownInstruction:
		return "crown"
	default:
		return "UNKNOWN"
	}
}

// instruction should contain all information needed to undo it
// that's the reason for storing the color and kind of the captured piece

type Instruction struct {
	t   instructionType
	row byte
	col byte
	// arbitraty data
	// row, col for moveInstruction
	// color, kind for captureInstruction
	// unused for crownInstruction
	d [2]byte
}

func (i Instruction) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "{%s (%d, %d)", i.t.String(), i.row, i.col)
	if i.t == moveInstruction {
		fmt.Fprintf(buf, " to (%d, %d)", i.d[0], i.d[1])
	} else if i.t == captureInstruction {
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

func MoveInstruction(sourceRow, sourceCol, destinationRow, destinationCol byte) Instruction {
	var i Instruction
	i.t = moveInstruction
	i.row, i.col = sourceRow, sourceCol
	i.d[0], i.d[1] = destinationRow, destinationCol
	return i
}

func CaptureInstruction(row, col byte, c Color, k Kind) Instruction {
	var i Instruction
	i.t = captureInstruction
	i.row, i.col = row, col
	i.d[0], i.d[1] = byte(c), byte(k)
	return i
}

func CrownInstruction(row, col byte) Instruction {
	var i Instruction
	i.t = crownInstruction
	i.row, i.col = row, col
	return i
}

func PerformInstructions(b *Board, is []Instruction) {
	for _, i := range is {
		switch i.t {
		case moveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			b.Move(fromRow, fromCol, toRow, toCol)
		case captureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := Color(i.d[0]), Kind(i.d[1])
			actualColor, actualKind := b.Get(row, col)

			// TODO return err instead of panicking?

			if capturedColor != actualColor || capturedKind != actualKind {
				panic(fmt.Sprintf(
					"performed capture instruction of %s %s on row %d %d but piece is a %s %s",
					capturedColor, capturedKind,
					row, col,
					actualColor, actualKind,
				))
			}
			b.Clear(row, col)
		case crownInstruction:
			b.Crown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}

// if you pass is to PerformInstructions
// the same is passed to UndoInstructions will undo the instructionList performed
// you don't need to reverse them
func UndoInstructions(b *Board, is []Instruction) {
	for k := len(is) - 1; k >= 0; k-- {
		i := is[k]
		switch i.t {
		case moveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			b.Move(toRow, toCol, fromRow, fromCol)
		case captureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := Color(i.d[0]), Kind(i.d[1])
			b.Set(row, col, capturedColor, capturedKind)
		case crownInstruction:
			b.Uncrown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}
