package main

import "fmt"

// I guess this is inspired by https://www.computerenhance.com/p/clean-code-horrible-performance
// a little edgy, I know, but I wanted to try it out

type instructionType uint8

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

type instruction struct {
	t   instructionType
	row byte
	col byte
	// arbitraty data
	// row, col for moveInstruction
	// color, kind for captureInstruction
	// unused for crownInstruction
	d [2]byte
}

func makeMoveInstruction(sourceRow, sourceCol, destinationRow, destinationCol byte) instruction {
	var i instruction
	i.t = moveInstruction
	i.row, i.col = sourceRow, sourceCol
	i.d[0], i.d[1] = destinationRow, destinationCol
	return i
}

func makeCaptureInstruction(row, col byte, c color, k kind) instruction {
	var i instruction
	i.t = captureInstruction
	i.row, i.col = row, col
	i.d[0], i.d[1] = byte(c), byte(k)
	return i
}

func makeCrownInstruction(row, col byte) instruction {
	var i instruction
	i.t = crownInstruction
	i.row, i.col = row, col
	return i
}

func performInstructions(b *board, is []instruction) {
	for _, i := range is {
		switch i.t {
		case moveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			movedColor, movedKind := b.get(fromRow, fromCol)
			b.clear(fromRow, fromCol)
			b.set(toRow, toCol, movedColor, movedKind)
		case captureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := color(i.d[0]), kind(i.d[1])
			actualColor, actualKind := b.get(row, col)
			if capturedColor != actualColor || capturedKind != actualKind {
				panic(fmt.Sprintf(
					"performed capture instruction of %s %s on row %d %d but piece is a %s %s",
					capturedColor, capturedKind,
					row, col,
					actualColor, actualKind,
				))
			}
			b.clear(row, col)
		case crownInstruction:
			b.crown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}

// if you pass is to performInstructions
// the same is passed to undoInstructions will undo the instructions performed
// you don't need to reverse them
func undoInstructions(b *board, is []instruction) {
	for k := len(is) - 1; k >= 0; k-- {
		i := is[k]
		switch i.t {
		case moveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			movedColor, movedKind := b.get(toRow, toCol)
			b.clear(toRow, toCol)
			b.set(fromRow, fromCol, movedColor, movedKind)
		case captureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := color(i.d[0]), kind(i.d[1])
			b.set(row, col, capturedColor, capturedKind)
		case crownInstruction:
			b.uncrown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}
