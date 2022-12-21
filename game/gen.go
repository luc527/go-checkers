package game

import "fmt"

var directionUp = []int8{-1}
var directionDown = []int8{1}
var directionsBoth = []int8{1, -1}

type captureWithDestinations struct {
	capture      capture
	destinations []coord
}

func GenerateSimpleMoves(board *Board, player Color) []SimpleMove {

	// TODO we generate moves much more often than read all board positions sequentially
	// would make more sense for the board to store a list of white positions with cells and a list of black positions with cells

	// this way we fill the cache with just what we need, instead of filling with the whole board

	var moves []SimpleMove

	for row := int8(0); row < 8; row++ {
		for col := int8(0); col < 8; col++ {

			if !board.Occupied(row, col) {
				continue
			}

			color, kind := board.Get(row, col)
			if color != player {
				continue
			}

			reach := int8(1)
			if kind == King {
				reach = 127 // unlimited reach
			}

			rowOffsets := directionsBoth
			if kind != King {
				rowOffsets = directionUp
				if color != White {
					rowOffsets = directionDown
				}
			}

			for _, rowOffset := range rowOffsets {
				for _, colOffset := range directionsBoth {

					for distance := int8(1); distance <= reach; distance++ {
						destRow := row + distance*rowOffset
						destCol := col + distance*colOffset

						if !Inbounds(destRow, destCol) {
							break
						}

						if board.Occupied(destRow, destCol) {
							break
						}

						crown := IsCrowningRow(player, destRow)

						moves = append(moves, SimpleMove{move{coord{row, col}, coord{destRow, destCol}}, crown})
					}

				}
			}

		}
	}

	return moves
}

// public only for testing
func captureDestinationsFrom(board *Board, row, col int8) []captureWithDestinations {
	if !Inbounds(row, col) {
		panic(fmt.Sprintf("out of bounds (%v, %v)", row, col))
	}

	if !board.Occupied(row, col) {
		panic(fmt.Sprintf("empty cell (%v, %v)", row, col))
	}

	var captures []captureWithDestinations

	color, kind := board.Get(row, col)

	if kind == Pawn {
		for _, rowOffset := range directionsBoth {
			for _, colOffset := range directionsBoth {
				capRow, capCol := row+rowOffset, col+colOffset
				if !Inbounds(capRow, capCol) || !board.Occupied(capRow, capCol) {
					break
				}
				capColor, capKind := board.Get(capRow, capCol)
				if capColor == color {
					break
				}
				destRow, destCol := row+2*rowOffset, col+2*colOffset
				if !Inbounds(destRow, destCol) || board.Occupied(destRow, destCol) {
					break
				}

				capture := capture{coord{capRow, capCol}, capColor, capKind}
				destination := coord{destRow, destCol}
				withDestination := captureWithDestinations{capture, []coord{destination}}
				captures = append(captures, withDestination)
			}
		}
	} else {
		for _, rowOffset := range directionsBoth {
			for _, colOffset := range directionsBoth {

				didFindCapture := false
				var foundCapture capture
				var destinations []coord

				distance := int8(1)
				for {
					curRow, curCol := row+distance*rowOffset, col+distance*colOffset
					if !Inbounds(curRow, curCol) {
						break
					}

					if didFindCapture {
						if board.Occupied(curRow, curCol) {
							break
						} else {
							destinations = append(destinations, coord{curRow, curCol})
						}
					} else if board.Occupied(curRow, curCol) {
						capColor, capKind := board.Get(curRow, curCol)
						if capColor == color {
							break
						} else {
							didFindCapture = true
							foundCapture = capture{coord{curRow, curCol}, capColor, capKind}
						}
					}

					distance++
				}

				if didFindCapture && len(destinations) > 0 {
					withDestinations := captureWithDestinations{foundCapture, destinations}
					captures = append(captures, withDestinations)
				}
			}
		}
	}

	return captures
}

//
// TESTS (unstructured)
//

func TestCaptureDestinationsFrom() {
	board := new(Board)

	board.Set(4, 3, White, King)
	board.Set(2, 1, Black, Pawn)
	board.Set(5, 2, Black, King)
	board.Set(7, 6, Black, Pawn)
	fmt.Println(board)

	cds := captureDestinationsFrom(board, 4, 3)

	for _, cd := range cds {
		cap := cd.capture
		fmt.Printf("capture %v %v on (%v, %v)\n", cap.color, cap.kind, cap.coord.row, cap.coord.col)

		fmt.Println("destinations")
		for _, dest := range cd.destinations {
			fmt.Printf("  (%v, %v)\n", dest.row, dest.col)
		}
	}
}
