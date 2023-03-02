// Compares the compact default board implementation (three uint64)
// with the more simple and direct implementation (matrix of pieces)

package main

// arguably interfaces shouldn't be named xxxInterface
// but this is just for the naive board comparison, so whatever
type boardInterface interface {
	isOccupied(row, col uint8) bool
	clear(row, col uint8)
	set(row, col uint8, c color, k kind)
	get(row, col uint8) (color, kind)
}

var _ boardInterface = &naiveBoard{}
var _ boardInterface = &board{}

type naivePiece struct {
	color
	kind
}

type naiveBoard struct {
	matrix [8][8]*naivePiece
}

func newEmptyNaiveBoard() *naiveBoard {
	return &naiveBoard{}
}

func (b *naiveBoard) isOccupied(row, col uint8) bool {
	return b.matrix[row][col] != nil
}

func (b *naiveBoard) clear(row, col uint8) {
	b.matrix[row][col] = nil
}

func (b *naiveBoard) set(row, col uint8, c color, k kind) {
	b.matrix[row][col] = &naivePiece{c, k}
}

func (b *naiveBoard) get(row, col uint8) (color, kind) {
	p := b.matrix[row][col]
	return p.color, p.kind
}
