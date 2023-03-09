package main

type gameState byte

const (
	playingState = gameState(iota)
	whiteWonState
	blackWonState
	drawState
)

// the part of the state that you need to explicitely remember in order to undo
type rememberedState struct {
	state   gameState
	plies   []ply
	lastPly ply
}

type game struct {
	captureRule
	bestRule
	rememberedState
	board   *board
	toPlay  color
	history []rememberedState
}

func newGame(captureRule captureRule, bestRule bestRule) *game {
	var g game
	g.board = new(board)
	placeInitialPieces(g.board)
	g.toPlay = whiteColor
	g.state = playingState
	g.captureRule = captureRule
	g.bestRule = bestRule
	g.plies = generatePlies(g.board, g.toPlay, captureRule, bestRule)
	g.lastPly = nil
	return &g
}

func (g *game) isOver() bool {
	return g.state != playingState
}

func (g *game) hasWinner() bool {
	return g.state == whiteWonState || g.state == blackWonState
}

func (g *game) winner() color {
	if g.state == whiteWonState {
		return whiteColor
	} else {
		return blackColor
	}
}

func (g *game) doPly(p ply) {
	performInstructions(g.board, p)
	g.toPlay = g.toPlay.opposite()

	// save current state in the history
	g.history = append(g.history, g.rememberedState)

	g.plies = nil
	g.lastPly = p

	count := g.board.pieceCount()
	whiteCount := count.whiteKings + count.whitePawns
	blackCount := count.blackKings + count.blackPawns

	if whiteCount == 0 {
		g.state = blackWonState
	} else if blackCount == 0 {
		g.state = whiteWonState
	}

	if g.isOver() {
		return
	}

	g.plies = generatePlies(g.board, g.toPlay, g.captureRule, g.bestRule)
	if len(g.plies) == 0 {
		if g.toPlay == whiteColor {
			g.state = blackWonState
		} else {
			g.state = whiteWonState
		}
	}

	// TODO draw detection
	// involves keeping more state variables in the rememberedState

	// TODO also, find a better name than 'rememberedState', it's a little annoying to type
}

func (g *game) undoLastPly() {
	if g.lastPly == nil {
		return
	}

	undoInstructions(g.board, g.lastPly)
	g.toPlay = g.toPlay.opposite()
	g.rememberedState = g.history[len(g.history)-1]

	g.history = g.history[:len(g.history)-1]
}
