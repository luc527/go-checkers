package main

import "testing"

type instructionParseTest struct {
	s       string
	wantErr bool
	want    Instruction
}

func runInstructionParseTests(t *testing.T, tests []instructionParseTest) {
	for _, test := range tests {
		got, err := parseInstruction(test.s)
		if test.wantErr {
			if err == nil {
				t.Errorf("expected encoded instruction %s to throw an error when parsed, but it didn't", test.s)
			} else {
				t.Logf("ok, got error %v", err)
			}
		} else if err != nil {
			t.Errorf("got unwatned error %v", err)
		} else if !got.Equals(test.want) {
			t.Errorf("want instruction %v, got %v", test.want, got)
		}
	}
}

func TestParseCrownInstruction(t *testing.T) {
	runInstructionParseTests(t, []instructionParseTest{
		{
			s:       "k",
			wantErr: true,
		},
		{
			s:       "k7",
			wantErr: true,
		},
		{
			s:       "k19",
			wantErr: true,
		},
		{
			s:    "k22",
			want: CrownInstruction(2, 2),
		},
		{
			s:    "k01",
			want: CrownInstruction(0, 1),
		},
		{
			s:    "k76",
			want: CrownInstruction(7, 6),
		},
	})
}

func placeholderCapture(row, col byte) Instruction {
	return CaptureInstruction(row, col, WhiteColor, KingKind)
}

func TestParseCaptureInstruction(t *testing.T) {
	runInstructionParseTests(t, []instructionParseTest{
		{
			s:       "cm1",
			wantErr: true,
		},
		{
			s:    "c71",
			want: placeholderCapture(7, 1),
		},
		{
			s:    "c00",
			want: placeholderCapture(0, 0),
		},
		{
			s:    "c36",
			want: placeholderCapture(3, 6),
		},
	})
}

func TestParseMoveInstruction(t *testing.T) {
	runInstructionParseTests(t, []instructionParseTest{
		{
			s:       "m12",
			wantErr: true,
		},
		{
			s:       "m12kk",
			wantErr: true,
		},
		{
			s:       "m2290",
			wantErr: true,
		},
		{
			s:    "m1234",
			want: MoveInstruction(1, 2, 3, 4),
		},
		{
			s:    "m0770",
			want: MoveInstruction(0, 7, 7, 0),
		},
		{
			s:    "m7016",
			want: MoveInstruction(7, 0, 1, 6),
		},
		{
			s:    "m0000",
			want: MoveInstruction(0, 0, 0, 0),
		},
	})
}

func TestParsePly(t *testing.T) {
	type plyParseTest struct {
		s       string
		wantErr bool
		want    Ply
	}

	tests := []plyParseTest{
		{
			s:       "",
			wantErr: true,
		},
		{
			s: "k12,m1234,c77",
			want: Ply{
				CrownInstruction(1, 2),
				MoveInstruction(1, 2, 3, 4),
				placeholderCapture(7, 7),
			},
		},
		{
			s:       "m7700;c12",
			wantErr: true,
		},
		{
			s: "m7700",
			want: Ply{
				MoveInstruction(7, 7, 0, 0),
			},
		},
		{
			s:       "m7700, c12, m7700",
			wantErr: true,
		},
	}

	for _, test := range tests {
		got, err := parsePly(test.s)
		if test.wantErr {
			if err == nil {
				t.Errorf("expected encoded ply %q to throw error, but it didn't", test.s)
			}
		} else if err != nil {
			t.Errorf("got unwanted error %v", err)
		} else if !sliceEq(got, test.want) {
			t.Errorf("%q: want ply %v got %v", test.s, test.want, got)
		}
	}
}
