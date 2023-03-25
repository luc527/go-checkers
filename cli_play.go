package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// TODO make server, then make a web client
// try to reuse most of the pin3-checkers client

func play() {

	g := NewStandardGame(CapturesMandatory, BestNotMandatory)

	input := bufio.NewScanner(os.Stdin)

	// if the game has ended the player can still undo their last action, so we don't quit yet
	quit := false

	mm := Minimax{
		ToMaximize: BlackColor,
		Heuristic:  WeightedCountHeuristic,
		Cutoff:     7,
	}

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

		if g.ToPlay() == WhiteColor {
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
				g.UndoLastPly() // undo AI
				g.UndoLastPly() // undo player
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
		} else {
			fmt.Println("Waiting fo the AI's choice...")
			_, ply := mm.Search(g, 0)
			fmt.Println("Choice:", ply)
			g.DoPly(ply)
		}

	}

	fmt.Print("\nHistory\n")
	gs := g.GenerateHistory()
	for _, g := range gs {
		fmt.Println(&g)
	}
}
