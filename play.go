package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/luc527/go_checkers/core"
	mm "github.com/luc527/go_checkers/minimax"
)

func main() {
	results := make(chan core.GameResult)

	sema := make(chan struct{}, 8)

	var wg sync.WaitGroup

	for t := 1; t <= 20; t++ {
		wg.Add(1)
		go func(t int) {
			sema <- struct{}{}
			res := play(t, results)
			fmt.Printf("%d done, result is %v\n", t, res)
			wg.Done()
			<-sema
		}(t)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	count := make(map[core.GameResult]int)
	for res := range results {
		count[res]++
	}

	for k, v := range count {
		fmt.Println(k, ": ", v)
	}

}

func play(id int, results chan<- core.GameResult) core.GameResult {
	var aiWhite mm.Searcher = mm.TimeLimitedSearcher{
		ToMax:     core.WhiteColor,
		Heuristic: mm.WeightedCountHeuristic,
		TimeLimit: 4000 * time.Millisecond,
	}

	var aiBlack mm.Searcher = mm.TimeLimitedSearcher{
		ToMax:     core.BlackColor,
		Heuristic: mm.WeightedCountHeuristic,
		TimeLimit: 400 * time.Millisecond,
	}

	g := core.NewStandardGame()
	var res core.GameResult
	i := 1
	for {
		res = g.Result()
		if res.Over() {
			break
		}

		var ai mm.Searcher
		if g.WhiteToPlay() {
			ai = aiWhite
		} else {
			ai = aiBlack
		}

		// t0 := time.Now()
		ply := ai.Search(g)
		// dt := time.Since(t0)

		if _, err := g.DoPly(ply); err != nil {
			fmt.Println(g.Board())
			fmt.Println(ply)
			break
		}
		// fmt.Printf("  %2d (%3d): %v, %v\n", i, id, g.ToPlay(), dt)
		i++
	}
	results <- res
	return res
}
