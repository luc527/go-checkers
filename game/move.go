package game

import (
	"fmt"
	"strconv"
)

type SimpleMove uint16

// for testing, TODO remove?
func NewSimpleMove(fromRow, fromCol, toRow, toCol uint8) *SimpleMove {
	move := MakeSimpleMove(fromRow, fromCol, toRow, toCol)
	return &move
}

func MakeSimpleMove(fromRow, fromCol, toRow, toCol uint8) SimpleMove {
	move := SimpleMove(0)
	move.Set(fromRow, fromCol, toRow, toCol)
	return move
}

// can be used when the memory was previously allocated outside
func (move *SimpleMove) Set(fromRow, fromCol, toRow, toCol uint8) {
	// least significant bit unused
	fromBits := uint16(0)
	fromBits |= uint16(fromRow) << (8 + 5)
	fromBits |= uint16(fromCol) << (8 + 2)
	// one bit for crowned

	// two least significant bits unused
	toBits := uint16(0)
	toBits |= uint16(toRow) << 5
	toBits |= uint16(toCol) << 2

	bits := fromBits | toBits

	*move = SimpleMove(bits)
}

func (move *SimpleMove) Debug() {
	fromBits := (*move & 0xFF00) >> 8
	toBits := *move & 0xFF
	fmt.Printf("from: %s\n", strconv.FormatUint(uint64(fromBits), 2))
	fmt.Printf("to:   %s\n", strconv.FormatUint(uint64(toBits), 2))
}

func (move *SimpleMove) String() string {
	fromRow, fromCol, toRow, toCol := move.coords()
	crowned := move.crowned()
	return fmt.Sprintf("from (%d, %d) to (%d, %d) (crowned? %v)", fromRow, fromCol, toRow, toCol, crowned)
}

func (move SimpleMove) coords() (fromRow, fromCol, toRow, toCol uint8) {
	toBits := (move & 0xFF) >> 2
	toCol = uint8(toBits & 07)
	toRow = uint8(toBits&070) >> 3

	fromBits := (move & 0xFF00) >> (8 + 2)
	fromCol = uint8(fromBits & 07)
	fromRow = uint8(fromBits&070) >> 3

	return
}

func (move SimpleMove) crowned() bool {
	bit := (move >> (8 + 1)) & 1
	return bit != 0
}

func SimpleMoveDo(fromRow, fromCol, toRow, toCol uint8, board *Board) (crowned bool) {
	color, kind := board.Take(fromRow, fromCol)

	toCrowningRow := (color == White && toRow == 0) || (color == Black && toRow == 7)
	crowned = kind == Pawn && toCrowningRow
	if crowned {
		kind = King
	}

	board.Set(toRow, toCol, color, kind)
	return
}

func SimpleMoveUndo(fromRow, fromCol, toRow, toCol uint8, crowned bool, board *Board) {
	color, kind := board.Take(toRow, toCol)
	if crowned {
		kind = Pawn
	}
	board.Set(fromRow, fromCol, color, kind)
}

func (move *SimpleMove) Do(board *Board) {
	fromRow, fromCol, toRow, toCol := move.coords()
	crowned := SimpleMoveDo(fromRow, fromCol, toRow, toCol, board)
	if crowned {
		*move |= 1 << (8 + 1)
	}
}

func (move *SimpleMove) Undo(board *Board) {
	fromRow, fromCol, toRow, toCol := move.coords()
	crowned := move.crowned()
	SimpleMoveUndo(fromRow, fromCol, toRow, toCol, crowned, board)
}
