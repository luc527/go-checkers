package main

import (
	"fmt"
	"time"
)

func autoplay(captureRule CaptureRule, bestRule BestRule, whiteBot Minimax, blackBot Minimax) {
	g := NewStandardGame(captureRule, bestRule)

	var state GameResult
	start := time.Now()
	for {
		state = g.Result()
		if state.IsOver() {
			break
		}
		toPlay := whiteBot
		if g.ToPlay() == BlackColor {
			toPlay = blackBot
		}
		_, ply := toPlay.Search(g)
		g.DoPly(ply)
	}
	end := time.Now()

	duration := end.Sub(start)
	fmt.Println("Game duration:", duration)

	if state.HasWinner() {
		fmt.Println(state.Winner(), "wins!")
	} else {
		fmt.Println("Draw")
	}
}
