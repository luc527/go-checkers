package main

import (
	"checkers/game"
	"fmt"
)

// TODO turn into proper tests
func testSimpleMove() {
	b := new(game.Board)

	b.Set(1, 0, game.White, game.Pawn)
	fmt.Println(b)

	{
		m := game.NewSimpleMove(1, 0, 0, 1)
		fmt.Println("doing", m)
		m.Do(b)
		fmt.Println(b)

		fmt.Println("undoing", m)
		m.Undo(b)
		fmt.Println(b)

		b.Set(5, 4, game.Black, game.Pawn)
		fmt.Println(b)
	}

	{
		m0 := game.NewSimpleMove(5, 4, 6, 5)
		fmt.Println("doing", m0)
		m0.Do(b)
		fmt.Println(b)

		m1 := game.NewSimpleMove(6, 5, 7, 4)
		fmt.Println("doing", m1)
		m1.Do(b)
		fmt.Println(b)

		fmt.Println("undoing", m1)
		m1.Undo(b)
		fmt.Println(b)

		fmt.Println("undoing", m0)
		m0.Undo(b)
		fmt.Println(b)
	}

}

func main() {
	testSimpleMove()
}
