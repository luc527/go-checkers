package main

import (
	"checkers/game"
	"fmt"
)

func main() {

	b := game.InitialBoard()

	for col := uint8(0); col < 8; col++ {
		b.Clear(2, col)
	}

	for col := uint8(0); col < 8; col++ {
		b.Clear(5, col)
	}

	b.Set(3, 3, game.White, game.King)

	simpleMoves := game.GenerateSimpleMoves(b, game.White)

	fmt.Println("Simple moves", simpleMoves)

	for _, move := range simpleMoves {
		fmt.Printf("doing %v\n", move.String())
		move.Do(b)
		fmt.Println(b)
		fmt.Printf("undoing %v\n", move.String())
		move.Undo(b)
		fmt.Println(b)
	}
}
