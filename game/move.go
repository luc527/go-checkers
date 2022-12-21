package game

import (
	"fmt"
	"strings"
)

type coord struct {
	row int8
	col int8
}

type move struct {
	from coord
	to   coord
}

type SimpleMove struct {
	move
	crowned bool
}

type capture struct {
	coord
	color Color
	kind  Kind
}

type CaptureMove struct {
	from     coord
	sequence []coord
	captures []capture
	crowned  bool
}

func (coord coord) String() string {
	return fmt.Sprintf("(%v, %v)", coord.row, coord.col)
}

func (move SimpleMove) String() string {
	c := ""
	if move.crowned {
		c = "crowned"
	}
	return fmt.Sprintf("%v -> %v %s", move.from, move.to, c)
}

func (move *SimpleMove) Do(board *Board) {
	color, kind := board.Take(move.from.row, move.from.col)
	if move.crowned {
		kind = King
	}
	board.Set(move.to.row, move.to.col, color, kind)
}

func (move *SimpleMove) Undo(board *Board) {
	color, kind := board.Take(move.to.row, move.to.col)
	if move.crowned {
		kind = Pawn
	}
	board.Set(move.from.row, move.from.col, color, kind)
}

func (move CaptureMove) String() string {
	strsequence := make([]string, 0, 1+len(move.sequence))
	strsequence = append(strsequence, move.from.String())
	for _, coord := range move.sequence {
		strsequence = append(strsequence, coord.String())
	}
	sequenceJoined := strings.Join(strsequence, " -> ")

	strcaptures := make([]string, 0, len(move.captures))
	for _, coord := range move.captures {
		strcaptures = append(strcaptures, coord.String())
	}
	capturesJoined := strings.Join(strcaptures, ", ")

	crowned := ""
	if move.crowned {
		crowned = ", Crowned"
	}

	return fmt.Sprintf("Sequence %v, Captures %v%v", sequenceJoined, capturesJoined, crowned)
}

func (move *CaptureMove) Do(board *Board) {
	for _, capture := range move.captures {
		board.Clear(capture.row, capture.col)
	}

	color, kind := board.Take(move.from.row, move.from.col)
	if move.crowned {
		kind = King
	}
	final := move.sequence[len(move.sequence)-1]
	board.Set(final.row, final.col, color, kind)
}

func (move *CaptureMove) Undo(board *Board) {
	final := move.sequence[len(move.sequence)-1]
	color, kind := board.Take(final.row, final.col)
	if move.crowned {
		kind = Pawn
	}
	board.Set(move.from.row, move.from.col, color, kind)

	for _, capture := range move.captures {
		board.Set(capture.row, capture.col, capture.color, capture.kind)
	}
}
