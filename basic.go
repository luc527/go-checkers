package main

type color byte

type kind byte

const (
	kWhite = color(1)
	kBlack = color(0)
	kKing  = kind(1)
	kPawn  = kind(0)
)

func (c color) String() string {
	if c == kWhite {
		return "white"
	} else {
		return "black"
	}
}

func (k kind) String() string {
	if k == kKing {
		return "king"
	} else {
		return "pawn"
	}
}
