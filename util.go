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

func showcasePlies(b *board, ps []ply) {
	for _, p := range ps {
		boards := make([]string, 0, len(p))
		for k := range p {
			performInstructions(b, p[k:k+1])
			boards = append(boards, b.String())
		}
		fmt.Println(p)
		sideBySide(boards)
		undoInstructions(b, p)
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

	iss := generateCapturePlies(nil, b, blackColor)
	showcasePlies(b, iss)
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

			iss := generatePlies(b, blackColor, capRule, bRule)
			showcasePlies(b, iss)
		}
	}
}

func sliceEq[T comparable](a []T, b []T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		y := b[i]
		if x != y {
			return false
		}
	}
	return true
}
