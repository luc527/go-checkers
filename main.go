package main

import (
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
)

// could do the following:
// run the game many times and gather statistics about the average number of plies as a function of piece count
// then make a map from piece count to the avg number of plies (maybe a little higher)
// and use that for the initial capacity of the []Ply at GeneratePlies

// another possible optimization
// GeneratePlies already iterates through the board and finds pieces
// it could also count the pieces so we don't have to do it as a separate process
// but this would make the logic in game.go a little confusing

func prof() {
	f, perr := os.Create("cpu4.pprof")
	if perr != nil {
		log.Fatal(perr)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// plyCounts := make(map[int][]int)

	trials := 100_000
	for t := 0; t < trials; t++ {
		g := NewStandardGame(CapturesMandatory, BestNotMandatory)
		for !g.ComputeState().IsOver() {
			plies := g.Plies()

			// pc := g.pieceCount
			// c := int(pc.BlackKings) + int(pc.BlackPawns) + int(pc.WhiteKings) + int(pc.WhitePawns)
			// plyCounts[c] = append(plyCounts[c], len(plies))

			random := plies[rand.Int()%len(plies)]
			g.DoPly(random)
		}
		for g.HasLastPly() {
			g.UndoLastPly()
		}
	}

	// for i := 0; i <= 26; i++ {
	// 	var avg, stddev float64

	// 	{
	// 		sum := 0
	// 		for _, c := range plyCounts[i] {
	// 			sum += c
	// 		}
	// 		avg = float64(sum) / float64(len(plyCounts[i]))
	// 	}

	// 	{
	// 		sum := 0.0
	// 		for _, c := range plyCounts[i] {
	// 			dev := float64(c) - avg
	// 			sq := dev * dev
	// 			sum += sq
	// 		}
	// 		variance := sum / float64(len(plyCounts[i]))
	// 		stddev = math.Sqrt(variance)
	// 	}

	// 	fmt.Printf("%3d -> avg %5.2g, stddev %5.2g, avg+1std %5.2g, avg+2std %5.2g, \n", i, avg, stddev, avg+stddev, avg+2*stddev)
	// }

}

func main() {
	runServer()
}
