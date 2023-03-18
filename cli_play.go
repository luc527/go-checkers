package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func play() {

	g := NewStandardGame(CapturesMandatory, BestNotMandatory)

	input := bufio.NewScanner(os.Stdin)

	// if the game has ended the player can still undo their last action, so we don't quit yet
	quit := false

	for !quit {
		fmt.Printf("It's %s's turn!\n", g.ToPlay())

		if g.IsOver() {
			if !g.HasWinner() {
				fmt.Println("It's a draw, no one wins")
			} else {
				fmt.Printf("The winner is %s!\n", g.Winner())
			}
		}

		plies := g.Plies()

		for i, p := range plies {
			fmt.Printf("[%2d]: %s\n", i, p.String())
		}
		if g.HasLastPly() {
			fmt.Println("[ u]: undo last move")
		}
		fmt.Println("[ q]: quit")

		fmt.Println(g.Board())

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
			g.UndoLastPly()
		} else if text == "q" {
			quit = true
		} else {
			i, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("Invalid move, try again (%v)\n", err)
				goto askForMove // Considered harmful!
			}
			if i < 0 || i >= len(plies) {
				fmt.Println("Invalid move, try again (out of bounds)")
				goto askForMove // Considered harmful!
			}
			g.DoPly(plies[i])
		}
	}
}
