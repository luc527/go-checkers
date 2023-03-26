package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
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
		err = fmt.Errorf("invalid player: %s", s)
	}
	return
}

func parseCaptureRule(s string) (rule CaptureRule, err error) {
	if s == "mandatory" || s == "capturesmandatory" {
		rule = CapturesMandatory
	} else if s == "notmandatory" || s == "capturesnotmandatory" {
		rule = CapturesNotMandatory
	} else {
		err = fmt.Errorf("invalid capture rule %s", s)
	}
	return
}

func parseBestRule(s string) (rule BestRule, err error) {
	if s == "mandatory" || s == "bestmandatory" {
		rule = BestMandatory
	} else if s == "notmandatory" || s == "bestnotmandatory" {
		rule = BestNotMandatory
	} else {
		err = fmt.Errorf("invalid best rule %s", s)
	}
	return
}

func handleGeneratePlies(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var strBoard string
	var strPlayer string
	var strCaptureRule string
	var strBestRule string

	for param, values := range r.Form {
		if len(values) == 0 {
			continue
		}
		value := values[len(values)-1]
		switch param {
		case "board":
			strBoard = strings.ToLower(value)
		case "player":
			strPlayer = strings.ToLower(value)
		case "captureRule":
			strCaptureRule = strings.ToLower(value)
		case "bestRule":
			strBestRule = strings.ToLower(value)
		}
	}

	board, err := parseBoard(strBoard)
	if err != nil {
		badRequest(w, err)
		return
	}

	player, err := parsePlayer(strPlayer)
	if err != nil {
		badRequest(w, err)
		return
	}

	captureRule, err := parseCaptureRule(strCaptureRule)
	if err != nil {
		badRequest(w, err)
		return
	}

	bestRule, err := parseBestRule(strBestRule)
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

// func handleAI(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	for param, values := range r.Form {
// 		if len(values) == 0 {
// 			continue
// 		}
// 		value := values[len(values)-1]
// 		switch param {
// 		case "board":
// 		case "player":
// 		case "captureRule":
// 		case "bestRule":
// 		case "turnsSinceCapture":
// 		case "turnsSincePawnMove":
// 		case "turnsInSpecialEnding":
// 		}
// 	}

// 	g := new(Game)
// }

func runServer() {
	http.HandleFunc("/generatePlies", handleGeneratePlies)
	// http.HandleFunc("/ai", handleAI)
	fmt.Println("Running server at http://localhost:8080")
	log.Fatalln(http.ListenAndServe("localhost:8080", nil))
}
