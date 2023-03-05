package main

import "fmt"

func main() {
	b := new(board)

	b.set(5, 5, kWhite, kKing)
	b.set(2, 2, kBlack, kKing)
	b.set(0, 7, kWhite, kKing)

	fmt.Println(b)
}
