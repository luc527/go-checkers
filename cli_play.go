package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func play() {

	// TODO no game state representation yet
	// ^ consider that the game state has to be able to undo as well
	// and that includes undoing a win, draw etc. in general rolling black every part of the state
	// also: need to detect draws.

	b := new(board)
	placeInitialPieces(b)

	captureRule := capturesMandatory
	bestRule := bestNotMandatory

	toPlay := whiteColor
	gameOver := false
	draw := false
	var winner color

	input := bufio.NewScanner(os.Stdin)

	for !gameOver {
		fmt.Printf("It's %s's turn!\n", toPlay)

		moves := generateMoves(b, toPlay, captureRule, bestRule)

		if len(moves) == 0 {
			gameOver = true
			winner = toPlay.opposite()
		}

		fmt.Println(b)

		for i, m := range moves {
			fmt.Printf("[%2d]: %s\n", i, instructionsToString(m))
		}

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
		i, err := strconv.Atoi(text)
		if err != nil {
			fmt.Printf("Invalid move, try again (%v)\n", err)
			goto askForMove // Considered harmful!
		}
		if i < 0 || i >= len(moves) {
			fmt.Println("Invalid move, try again")
			goto askForMove // Considered harmful!
		}

		moveTaken := moves[i]

		performInstructions(b, moveTaken)

		count := b.pieceCount()

		whites := count.whitePawns + count.whiteKings
		blacks := count.blackPawns + count.blackKings

		if whites == 0 {
			gameOver = true
			winner = blackColor
		} else if blacks == 0 {
			gameOver = true
			winner = whiteColor
		}

		toPlay = toPlay.opposite()
	}

	if gameOver {
		if draw {
			fmt.Println("It's a draw, no one wins")
		} else {
			fmt.Printf("The winner is %s!\n", winner)
		}
	} else {
		fmt.Println("Oops")
	}
}
