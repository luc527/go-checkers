package main

import (
	"bytes"
	"fmt"
	"strings"
)

type insType byte

const (
	moveIns = insType(iota)
	captureIns
	crownIns
)

func (t insType) String() string {
	switch t {
	case moveIns:
		return "move"
	case captureIns:
		return "capture"
	case crownIns:
		return "crown"
	default:
		return "UNKNOWN"
	}
}

// instruction should contain all information needed to undo it
// that's the reason for storing the color and kind of the captured piece
// when you'd just need the coordinate for actually removing the piece

type ins struct {
	t   insType
	row byte
	col byte
	// arbitraty data
	// row, col for moveIns
	// color, kind for captureIns
	// unused for crownIns
	d [2]byte
}

type insList []ins

func (i ins) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "{%s instruction (%d, %d)", i.t.String(), i.row, i.col)
	if i.t == moveIns {
		fmt.Fprintf(buf, " (%d, %d)", i.d[0], i.d[1])
	} else if i.t == captureIns {
		fmt.Fprintf(buf, " %s %s", color(i.d[0]), kind(i.d[1]))
	}
	buf.WriteRune('}')
	return buf.String()
}

func (is insList) String() string {
	ss := make([]string, 0, len(is))
	for _, i := range is {
		ss = append(ss, i.String())
	}
	return strings.Join(ss, ";")
}

func makeMoveIns(sourceRow, sourceCol, destinationRow, destinationCol byte) ins {
	var i ins
	i.t = moveIns
	i.row, i.col = sourceRow, sourceCol
	i.d[0], i.d[1] = destinationRow, destinationCol
	return i
}

func makeCaptureIns(row, col byte, c color, k kind) ins {
	var i ins
	i.t = captureIns
	i.row, i.col = row, col
	i.d[0], i.d[1] = byte(c), byte(k)
	return i
}

func makeCrownIns(row, col byte) ins {
	var i ins
	i.t = crownIns
	i.row, i.col = row, col
	return i
}

func performInstructions(b *board, is []ins) {
	for _, i := range is {
		switch i.t {
		case moveIns:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			movedColor, movedKind := b.get(fromRow, fromCol)
			b.clear(fromRow, fromCol)
			b.set(toRow, toCol, movedColor, movedKind)
		case captureIns:
			row, col := i.row, i.col
			capturedColor, capturedKind := color(i.d[0]), kind(i.d[1])
			actualColor, actualKind := b.get(row, col)

			// TODO return err instead of panicking

			if capturedColor != actualColor || capturedKind != actualKind {
				panic(fmt.Sprintf(
					"performed capture instruction of %s %s on row %d %d but piece is a %s %s",
					capturedColor, capturedKind,
					row, col,
					actualColor, actualKind,
				))
			}
			b.clear(row, col)
		case crownIns:
			b.crown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}

// if you pass is to performInstructions
// the same is passed to undoInstructions will undo the insList performed
// you don't need to reverse them
func undoInstructions(b *board, is []ins) {
	for k := len(is) - 1; k >= 0; k-- {
		i := is[k]
		switch i.t {
		case moveIns:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			movedColor, movedKind := b.get(toRow, toCol)
			b.clear(toRow, toCol)
			b.set(fromRow, fromCol, movedColor, movedKind)
		case captureIns:
			row, col := i.row, i.col
			capturedColor, capturedKind := color(i.d[0]), kind(i.d[1])
			b.set(row, col, capturedColor, capturedKind)
		case crownIns:
			b.uncrown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}
