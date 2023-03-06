package main

import "fmt"

func main() {
	b := new(board)
	placeInitialPieces(b)
	fmt.Println(b)
}
