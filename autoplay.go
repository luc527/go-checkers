package main

import (
	"fmt"
	"time"
)

func autoplay(captureRule CaptureRule, bestRule BestRule, whiteBot Minimax, blackBot Minimax) {
	g := NewStandardGame(captureRule, bestRule)

	start := time.Now()
	for !g.IsOver() {
		toPlay := whiteBot
		if g.ToPlay() == BlackColor {
			toPlay = blackBot
		}
		_, ply := toPlay.Search(g, 0)
		g.DoPly(ply)
	}
	end := time.Now()

	duration := end.Sub(start)
	fmt.Println("Game duration:", duration)

	// gs := g.GenerateHistory()
	// for _, g := range gs {
	// 	fmt.Println(&g)
	// }

	if g.HasWinner() {
		fmt.Println(g.Winner(), "wins!")
	} else {
		fmt.Println("Draw")
	}
}
