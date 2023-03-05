package main

import "fmt"

// TODO the generateXxx procedures should take a color, duh

func main() {

	// TODO make some proper unit tests for this

	b := new(board)

	b.set(2, 2, whiteColor, pawnKind)
	b.set(3, 3, blackColor, pawnKind)

	fmt.Println("before:")
	fmt.Println(b)

	var iss [][]instruction
	iss = generateCaptureMoves(iss, b)

	fmt.Println("after (should be the same):")
	fmt.Println(b)
	fmt.Printf("instructions (%d):\n", len(iss))
	for _, is := range iss {
		fmt.Println(is)
	}
}
