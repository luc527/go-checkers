package main

import (
	"strings"
	"testing"
)

func compareGeneratedInstructions(
	got []insList,
	want []insList,
) (
	extra []string,
	missing []string,
) {

	gotMap := make(map[string]bool)
	for _, is := range got {
		gotMap[is.String()] = true
	}

	wantMap := make(map[string]bool)
	for _, is := range want {
		wantMap[is.String()] = true
	}

	for is := range gotMap {
		if _, ok := wantMap[is]; !ok {
			extra = append(extra, is)
		}
	}

	for is := range wantMap {
		if _, ok := gotMap[is]; !ok {
			missing = append(missing, is)
		}
	}

	return
}

func insListsString(iss []insList) string {
	ss := make([]string, len(iss))
	for _, is := range iss {
		ss = append(ss, is.String())
	}
	return strings.Join(ss, "\n")
}

func TestSimplePawnMove(t *testing.T) {
	b := newEmptyBoard()
	b.set(1, 1, kWhite, kPawn)
	// no possibilities
	b.set(0, 0, kWhite, kPawn)
	// one possibility occupied
	b.set(0, 2, kWhite, kPawn)
	// against the (left and right) walls
	b.set(1, 7, kWhite, kPawn)
	b.set(3, 0, kWhite, kPawn)
	// crowning
	b.set(6, 6, kWhite, kPawn)

	// TODO test black

	movesGot := callGenerateSimpleMoves(b)
	movesWant := []insList{
		{makeMoveIns(1, 1, 2, 2)},
		{makeMoveIns(1, 1, 2, 0)},
		{makeMoveIns(0, 2, 1, 3)},
		{makeMoveIns(1, 7, 2, 6)},
		{makeMoveIns(3, 0, 4, 1)},
		{
			makeMoveIns(6, 6, 7, 7),
			makeCrownIns(7, 7),
		},
		{
			makeMoveIns(6, 6, 7, 5),
			makeCrownIns(7, 5),
		},
	}
	extra, missing := compareGeneratedInstructions(movesGot, movesWant)
	if len(extra) > 0 {
		t.Errorf("generated extra instruction lists:\n%s", strings.Join(extra, "\n"))
	}
	if len(missing) > 0 {
		t.Errorf("didn't generate these instruction lists:\n%s", strings.Join(missing, "\n"))
	}
}
