package main

var dirBoth = [2]int8{-1, +1}

func generateSimplePawnMoves(b *board, row, col byte, color color, ch chan<- []ins) {
	drow := byte(int8(row) + forward[color])
	if drow >= 8 {
		return
	}
	crown := crowningRow[color] == drow
	for _, dir := range dirBoth {
		dcol := byte(int8(col) + dir)
		if dcol >= 8 || b.isOccupied(drow, dcol) {
			continue
		}
		var is []ins
		is = append(is, makeMoveIns(row, col, drow, dcol))
		if crown {
			is = append(is, makeCrownIns(drow, dcol))
		}
		ch <- is
	}
}

func generateSimpleKingMoves(b *board, row, col byte, color color, ch chan<- []ins) {
	for _, roff := range dirBoth {
		for _, coff := range dirBoth {
			for dist := int8(1); ; dist++ {
				drow, dcol := byte(int8(row)+dist*roff), byte(int8(col)+dist*coff)
				if drow > 8 || dcol > 8 || b.isOccupied(drow, dcol) {
					break
				}

				var is []ins
				is = append(is, makeMoveIns(row, col, drow, dcol))
				if drow == crowningRow[color] {
					is = append(is, makeCrownIns(drow, dcol))
				}
				ch <- is

				dist++
			}
		}
	}
}

func generateSimpleMoves(b *board, ch chan<- []ins) {
	for row := byte(0); row < 8; row++ {
		for col := byte(0); col < 8; col++ {
			if !b.isOccupied(row, col) {
				continue
			}

			color, kind := b.get(row, col)
			if kind == kPawn {
				generateSimplePawnMoves(b, row, col, color, ch)
			} else {
				generateSimpleKingMoves(b, row, col, color, ch)
			}
		}
	}
}

func callGenerateSimpleMoves(b *board) []insList {
	var iss []insList
	done := make(chan struct{})
	ch := make(chan []ins)
	go func() {
		generateSimpleMoves(b, ch)
		done <- struct{}{}
	}()
	// idk yet if generateSimpleMoves should close the channel
	for {
		select {
		case is := <-ch:
			iss = append(iss, is)
		case <-done:
			return iss
		}
	}
}
