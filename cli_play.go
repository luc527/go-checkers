package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func play() {

	g := newGame(capturesMandatory, bestNotMandatory)

	input := bufio.NewScanner(os.Stdin)

	quit := false

	for !quit {
		fmt.Printf("It's %s's turn!\n", g.toPlay)

		if g.isOver() {
			if !g.hasWinner() {
				fmt.Println("It's a draw, no one wins")
			} else {
				fmt.Printf("The winner is %s!\n", g.winner())
			}
		}

		for i, p := range g.plies {
			fmt.Printf("[%2d]: %s\n", i, p.String())
		}
		if len(g.history) > 0 {
			fmt.Println("[ u]: undo last move")
		}
		fmt.Println("[ q]: quit")

		fmt.Println(g.board)

	askForMove:
		fmt.Print("Your choice: ")
		if !input.Scan() {
			break
		}
		if err := input.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		text := input.Text()
		if text == "u" {
			g.undoLastPly()
		} else if text == "q" {
			quit = true
		} else {
			i, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("Invalid move, try again (%v)\n", err)
				goto askForMove // Considered harmful!
			}
			if i < 0 || i >= len(g.plies) {
				fmt.Println("Invalid move, try again")
				goto askForMove // Considered harmful!
			}
			g.doPly(g.plies[i])
		}
	}
}
