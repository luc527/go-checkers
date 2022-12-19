package main

import (
	"checkers/game"
	"fmt"
)

func main() {

	// board := game.InitialBoard()
	board := new(game.Board)
	board.Set(1, 0, game.White, game.Pawn)
	fmt.Println(board)

	moves := game.GenerateSimpleMoves(board, game.White)
	fmt.Println("simple moves:", moves)

	for _, move := range moves {
		fmt.Println("doing", move.String())
		move.Do(board)
		fmt.Println(board)
		fmt.Println("undoing", move.String())
		move.Undo(board)
		fmt.Println(board)
	}
}
