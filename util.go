package main

import (
	"fmt"
	"strings"
)

func beforeEachLine(indent string, s string) string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	indented := make([]string, 0, len(lines))
	for _, line := range lines {
		indented = append(indented, indent+line)
	}
	return strings.Join(indented, "\n")
}

func sideBySide(boards []string) {
	var rowsPerBoard [][]string

	maxRows := -1

	for _, board := range boards {
		rows := strings.Split(board, "\n")
		if len(rows) > maxRows {
			maxRows = len(rows)
		}
		rowsPerBoard = append(rowsPerBoard, rows)
	}

	for row := 0; row < maxRows; row++ {
		for _, board := range rowsPerBoard {
			fmt.Print(board[row])
			fmt.Print("  ")
		}
		fmt.Println()
	}

}

func showcaseMoves(b *board, iss [][]instruction) {
	for _, is := range iss {
		boards := make([]string, 0, len(is))
		for k := range is {
			performInstructions(b, is[k:k+1])
			boards = append(boards, b.String())
		}
		fmt.Println(is)
		sideBySide(boards)
		undoInstructions(b, is)
		fmt.Println()
	}
}

func example0() {
	b := new(board)

	b.set(0, 5, blackColor, kingKind)
	b.set(2, 3, whiteColor, pawnKind)
	b.set(4, 3, whiteColor, pawnKind)
	b.set(4, 5, whiteColor, pawnKind)
	b.set(6, 5, whiteColor, pawnKind)
	b.set(5, 2, whiteColor, pawnKind)

	fmt.Println(b)

	iss := generateCaptureMoves(nil, b, blackColor)
	showcaseMoves(b, iss)
}

func example1() {

	b := new(board)

	//same board as example0
	b.set(0, 5, blackColor, kingKind)
	b.set(2, 3, whiteColor, pawnKind)
	b.set(4, 3, whiteColor, pawnKind)
	b.set(4, 5, whiteColor, pawnKind)
	b.set(6, 5, whiteColor, pawnKind)
	b.set(5, 2, whiteColor, pawnKind)

	fmt.Println(b)

	captureRules := []captureRule{capturesMandatory, capturesNotMandatory}
	bestRules := []bestRule{bestMandatory, bestNotMandatory}

	for _, capRule := range captureRules {
		for _, bRule := range bestRules {
			fmt.Println()

			capString := " NOT "
			if capRule == capturesMandatory {
				capString = " ARE "
			}

			bString := " NOT "
			if bRule == bestMandatory {
				bString = " ARE "
			}

			fmt.Printf("==================%s===================\n", "=====")
			fmt.Printf("========= captures%smandatory =========\n", capString)
			fmt.Printf("=========     best%smandatory =========\n", bString)
			fmt.Printf("==================%s===================\n", "=====")

			iss := generateMoves(b, blackColor, capRule, bRule)
			showcaseMoves(b, iss)
		}
	}
}
