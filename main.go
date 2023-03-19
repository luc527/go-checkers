package main

import (
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
)

func main() {
	f, perr := os.Create("cpu.pprof")
	if perr != nil {
		log.Fatal(perr)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	trials := 100_000
	for t := 0; t < trials; t++ {
		g := NewStandardGame(CapturesMandatory, BestNotMandatory)
		for !g.IsOver() {
			plies := g.Plies()
			random := plies[rand.Int()%len(plies)]
			g.DoPly(random)
		}
		for g.HasLastPly() {
			g.UndoLastPly()
		}
	}
}
