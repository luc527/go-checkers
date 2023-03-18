// Compares the compact default board implementation (three uint64)
// with the more simple and direct implementation (matrix of pieces)

package main

// arguably interfaces shouldn't be named xxxInterface
// but this is just for the naive board comparison, so whatever
type boardInterface interface {
	IsOccupied(row, col uint8) bool
	Clear(row, col uint8)
	Set(row, col uint8, c Color, k Kind)
	Get(row, col uint8) (Color, Kind)
}

var _ boardInterface = &naiveBoard{}
var _ boardInterface = &Board{}

type naivePiece struct {
	Color
	Kind
}

type naiveBoard struct {
	matrix [8][8]*naivePiece
}

func newEmptyNaiveBoard() *naiveBoard {
	return &naiveBoard{}
}

func (b *naiveBoard) IsOccupied(row, col uint8) bool {
	return b.matrix[row][col] != nil
}

func (b *naiveBoard) Clear(row, col uint8) {
	b.matrix[row][col] = nil
}

func (b *naiveBoard) Set(row, col uint8, c Color, k Kind) {
	b.matrix[row][col] = &naivePiece{c, k}
}

func (b *naiveBoard) Get(row, col uint8) (Color, Kind) {
	p := b.matrix[row][col]
	return p.Color, p.Kind
}
