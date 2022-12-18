package main

import (
	"checkers/board"
	"fmt"
)

func main() {
	b := board.InitialBoard()
	b.Set(0, 0, board.Black, board.King)
	b.Set(7, 7, board.White, board.King)
	fmt.Println(b)
}
