package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type coordObject struct {
	Row byte `json:"row"`
	Col byte `json:"col"`
}

type instructionObject struct {
	Type        string       `json:"type"`
	Source      *coordObject `json:"source,omitempty"`
	Destination *coordObject `json:"destination,omitempty"`
}

func makeInstructionObject(i Instruction) instructionObject {
	o := instructionObject{
		Type: i.t.String(),
		Source: &coordObject{
			Row: i.row,
			Col: i.col,
		},
	}
	if i.t == moveInstruction {
		o.Destination = &coordObject{
			Row: i.d[0],
			Col: i.d[1],
		}
	}
	return o
}

func parseBoard(s string) (*Board, error) {
	if len(s)%4 != 0 {
		return nil, fmt.Errorf("invalid board %q", s)
	}

	s = strings.ToLower(s)
	b := new(Board)

	for len(s) > 0 {
		colorRune := s[0]
		kindRune := s[1]
		rowRune := s[2]
		colRune := s[3]

		var color Color
		switch colorRune {
		case 'w':
			color = WhiteColor
		case 'b':
			color = BlackColor
		default:
			return nil, fmt.Errorf("invalid piece color %q", colorRune)
		}

		var kind Kind
		switch kindRune {
		case 'p':
			kind = PawnKind
		case 'k':
			kind = KingKind
		default:
			return nil, fmt.Errorf("invalid piece kind %q", kindRune)
		}

		var row byte
		if rowRune < '0' || rowRune > '7' {
			return nil, fmt.Errorf("invalid row %q", rowRune)
		}
		row = byte(rowRune - '0')

		var col byte
		if colRune < '0' || colRune > '7' {
			return nil, fmt.Errorf("invalid col %q", colRune)
		}
		col = byte(colRune - '0')

		b.Set(row, col, color, kind)

		s = s[4:]
	}

	return b, nil
}

func parsePlayer(s string) (player Color, err error) {
	if s == "white" || s == "w" {
		player = WhiteColor
	} else if s == "black" || s == "b" {
		player = BlackColor
	} else {
		err = fmt.Errorf("invalid player: %q", s)
	}
	return
}

func parseCaptureRule(s string) (rule CaptureRule, err error) {
	if s == "mandatory" || s == "capturesmandatory" {
		rule = CapturesMandatory
	} else if s == "notmandatory" || s == "capturesnotmandatory" {
		rule = CapturesNotMandatory
	} else {
		err = fmt.Errorf("invalid capture rule %q", s)
	}
	return
}

func parseBestRule(s string) (rule BestRule, err error) {
	if s == "mandatory" || s == "bestmandatory" {
		rule = BestMandatory
	} else if s == "notmandatory" || s == "bestnotmandatory" {
		rule = BestNotMandatory
	} else {
		err = fmt.Errorf("invalid best rule %q", s)
	}
	return
}

func parseState(s string) (g GameState, err error) {
	switch s {
	case "playing":
		g = PlayingState
	case "whiteWon":
		g = WhiteWonState
	case "blackWon":
		g = BlackWonState
	case "draw":
		g = DrawState
	default:
		err = fmt.Errorf("invalid state %q", s)
	}
	return
}

func parsePly(s string) (ply Ply, err error) {
	stringInstructions := strings.Split(s, ",")
	ply = make([]Instruction, len(stringInstructions))
	for k, stringInstruction := range stringInstructions {
		var instruction Instruction
		instruction, err = parseInstruction(stringInstruction)
		if err != nil {
			return
		}
		ply[k] = instruction
	}
	return
}

func isCoord(r rune) bool {
	return r >= '0' && r <= '7'
}

// TODO test parsing stuff
// specially parseInstruction

// capture instruction will be returned incomplete
// it has no access to the board to tell what kind of piece
// was captured, and
func parseInstruction(s string) (instruction Instruction, err error) {
	if len(s) < 3 {
		err = fmt.Errorf("empty instruction %q", s)
		return
	}

	var t instructionType
	switch s[0] {
	case 'k':
		t = crownInstruction
	case 'c':
		t = captureInstruction
	case 'm':
		t = moveInstruction
	}

	r1 := rune(s[1])
	if !isCoord(r1) {
		err = fmt.Errorf("invalid instruction ([1] not row) %q", s)
		return
	}
	r2 := rune(s[2])
	if !isCoord(r2) || r2 < '0' || r2 > '7' {
		err = fmt.Errorf("invalid instruction ([2] not col) %q", s)
		return
	}

	row, col := byte(r1-'0'), byte(r2-'0')

	switch t {
	case moveInstruction:
		if len(s) != 5 {
			err = fmt.Errorf("invalid instruction (missing move destination) %q", s)
			return
		}
		d1 := rune(s[3])
		if !isCoord(d1) {
			err = fmt.Errorf("invalid instruction ([3] not row) %q", s)
			return
		}
		d2 := rune(s[4])
		if !isCoord(d2) {
			err = fmt.Errorf("invalid instruction ([3] not col) %q", s)
			return
		}
		drow, dcol := byte(d1-'0'), byte(d2-'0')

		instruction = MoveInstruction(row, col, drow, dcol)
	case captureInstruction:
		placeholderColor := WhiteColor
		placeholderKind := KingKind
		instruction = CaptureInstruction(row, col, placeholderColor, placeholderKind)
	case crownInstruction:
		instruction = CrownInstruction(row, col)
	}

	return
}

func handleGeneratePlies(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	board, err := parseBoard(strings.ToLower(r.Form.Get("board")))
	if err != nil {
		badRequest(w, err)
		return
	}

	player, err := parsePlayer(strings.ToLower(r.Form.Get("player")))
	if err != nil {
		badRequest(w, err)
		return
	}

	captureRule, err := parseCaptureRule(strings.ToLower(r.Form.Get("captureRule")))
	if err != nil {
		badRequest(w, err)
		return
	}

	bestRule, err := parseBestRule(strings.ToLower(r.Form.Get("bestRule")))
	if err != nil {
		badRequest(w, err)
		return
	}

	plies := GeneratePlies(make([]Ply, 0, 20), board, player, captureRule, bestRule)

	pliesArray := make([][]instructionObject, len(plies))
	for k, ply := range plies {
		pliesArray[k] = make([]instructionObject, len(ply))
		for l, ins := range ply {
			pliesArray[k][l] = makeInstructionObject(ins)
		}
	}

	body, err := json.Marshal(pliesArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(body))
}

func handleDoPly(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		serverError(w, err)
		return
	}

	form := r.Form

	board, err := parseBoard(form.Get("board"))
	if err != nil {
		badRequest(w, err)
		return
	}

	player, err := parsePlayer(strings.ToLower(form.Get("player")))
	if err != nil {
		badRequest(w, err)
		return
	}

	captureRule, err := parseCaptureRule(strings.ToLower(form.Get("captureRule")))
	if err != nil {
		badRequest(w, err)
		return
	}

	bestRule, err := parseBestRule(strings.ToLower(form.Get("bestRule")))
	if err != nil {
		badRequest(w, err)
		return
	}

	turnsSinceCapture, err := strconv.Atoi(form.Get("turnsSinceCapture"))
	if err != nil {
		badRequest(w, err)
		return
	}

	turnsSincePawnMove, err := strconv.Atoi(form.Get("turnsSincePawnMove"))
	if err != nil {
		badRequest(w, err)
		return
	}

	turnsInSpecialEnding, err := strconv.Atoi(form.Get("turnsInSpecialEnding"))
	if err != nil {
		badRequest(w, err)
		return
	}

	state, err := parseState(form.Get("state"))
	if err != nil {
		badRequest(w, err)
		return
	}

	ply, err := parsePly(form.Get("ply"))
	if err != nil {
		badRequest(w, err)
		return
	}

	g := new(Game)

	g.stagnantTurnsToDraw = 20 // default
	g.pieceCount = board.PieceCount()
	g.plies = g.generatePlies()

	g.board = board
	g.toPlay = player
	g.captureRule = captureRule
	g.bestRule = bestRule
	g.turnsSinceCapture = int16(turnsSinceCapture)
	g.turnsSincePawnMove = int16(turnsSincePawnMove)
	g.turnsInSpecialEnding = int16(turnsInSpecialEnding)

	if state.IsOver() {
		badRequest(w, fmt.Errorf("TODO(1)"))
	}

	g.DoPly(ply)

	badRequest(w, fmt.Errorf("TODO(2)"))
}

// TODO finish implementing /doPly
// TODO test /doPly

func runServer() {
	http.HandleFunc("/generatePlies", handleGeneratePlies)
	http.HandleFunc("/doPly", handleDoPly)
	fmt.Println("Running server at http://localhost:8080")
	log.Fatalln(http.ListenAndServe("localhost:8080", nil))
}
