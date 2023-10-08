package core

import (
	"bytes"
	"fmt"
	"strings"
)

type InstructionType byte

const (
	MoveInstruction = InstructionType(iota)
	CaptureInstruction
	CrownInstruction
)

func (t InstructionType) String() string {
	switch t {
	case MoveInstruction:
		return "move"
	case CaptureInstruction:
		return "capture"
	case CrownInstruction:
		return "crown"
	default:
		return "UNKNOWN"
	}
}

// instruction should contain all information needed to undo it
// that's the reason for storing the color and kind of the captured piece

type Instruction struct {
	t   InstructionType
	row byte
	col byte
	// arbitraty data
	// destination row and col for MoveInstruction
	// color and kind for CaptureInstruction
	// unused for CrownInstruction
	d [2]byte
}

func (i Instruction) Equals(o Instruction) bool {
	return i.t == o.t &&
		i.row == o.row &&
		i.col == o.col &&
		(i.t == CrownInstruction ||
			(i.d[0] == o.d[0] &&
				i.d[1] == o.d[1]))
}

func (i Instruction) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "{%s (%d, %d)", i.t.String(), i.row, i.col)
	if i.t == MoveInstruction {
		fmt.Fprintf(buf, " to (%d, %d)", i.d[0], i.d[1])
	} else if i.t == CaptureInstruction {
		fmt.Fprintf(buf, " %s %s", Color(i.d[0]), Kind(i.d[1]))
	}
	buf.WriteRune('}')
	return buf.String()
}

var (
	colorChar = [...]byte{BlackColor: 'b', WhiteColor: 'w'}
	kindChar  = [...]byte{PawnKind: 'p', KingKind: 'k'}
)

func (i Instruction) SerializeInto(buf *bytes.Buffer) error {
	// Please, The Go Authors™, implement a 'try' keyword...
	switch i.t {
	case MoveInstruction:
		if err := buf.WriteByte('m'); err != nil {
			return err
		}
	case CaptureInstruction:
		if err := buf.WriteByte('c'); err != nil {
			return err
		}
	case CrownInstruction:
		if err := buf.WriteByte('k'); err != nil {
			return err
		}
	}

	if err := buf.WriteByte('0' + i.row); err != nil {
		return err
	}
	if err := buf.WriteByte('0' + i.col); err != nil {
		return err
	}

	if i.t == MoveInstruction {
		if err := buf.WriteByte('0' + i.d[0]); err != nil {
			return err
		}
		if err := buf.WriteByte('0' + i.d[1]); err != nil {
			return err
		}
	} else if i.t == CaptureInstruction {
		if err := buf.WriteByte(colorChar[i.d[0]]); err != nil {
			return err
		}
		if err := buf.WriteByte(kindChar[i.d[1]]); err != nil {
			return err
		}
	}
	return nil
}

func (i Instruction) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('"')
	if err := i.SerializeInto(&buf); err != nil {
		return nil, err
	}
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

func (i *Instruction) Unserialize(bs []byte) error {
	if len(bs) < 1 {
		return fmt.Errorf("unmarshal instruction: empty bytes")
	}
	tbyte := bs[0]
	switch tbyte {
	case 'm':
		i.t = MoveInstruction
	case 'c':
		i.t = CaptureInstruction
	case 'k':
		i.t = CrownInstruction
	default:
		return fmt.Errorf("unmarshal instruction: invalid type")
	}

	if len(bs) < 3 {
		return fmt.Errorf("unmarshal instruction: missing row and col")
	}
	row := bs[1] - '0'
	col := bs[2] - '0'
	if row > 7 || col > 7 {
		return fmt.Errorf("unmarshal instruction: out of bounds row or col")
	}
	i.row = row
	i.col = col

	if i.t == CrownInstruction {
		if len(bs) != 3 {
			return fmt.Errorf("unmarshal instruction: crown: trailing bytes")
		}
		return nil
	}
	if i.t == MoveInstruction {
		if len(bs) != 5 {
			return fmt.Errorf("unmarshal instruction: move: missing or trailing bytes")
		}
		drow := bs[3] - '0'
		dcol := bs[4] - '0'
		if drow > 7 || dcol > 7 {
			return fmt.Errorf("unmarshal instruction: move: out of bounds destination row or col")
		}
		i.d[0] = drow
		i.d[1] = dcol
		return nil
	}
	if i.t == CaptureInstruction {
		if len(bs) != 5 {
			return fmt.Errorf("unmarshal instruction: capture: missing or trailing bytes")
		}
		var color Color
		switch bs[3] {
		case 'b':
			color = BlackColor
		case 'w':
			color = WhiteColor
		default:
			return fmt.Errorf("unmarshal instruction: capture: invalid captured color")
		}
		var kind Kind
		switch bs[4] {
		case 'p':
			kind = PawnKind
		case 'k':
			kind = KingKind
		default:
			return fmt.Errorf("unmarshal instruction: capture: invalid captured kind")
		}
		i.d[0] = byte(color)
		i.d[1] = byte(kind)
		return nil
	}
	// Unreachable.
	return nil
}

func (i *Instruction) UnmarshalJSON(bs []byte) error {
	if len(bs) < 3 {
		return fmt.Errorf("unmarshal instruction: empty bytes")
	}
	if bs[0] != '"' && bs[len(bs)-1] != '"' {
		return fmt.Errorf("unmarshal instruction: not a string")
	}
	bs = bs[1 : len(bs)-1]
	return i.Unserialize(bs)
}

func instructionsToString(is []Instruction) string {
	ss := make([]string, 0, len(is))
	for _, i := range is {
		ss = append(ss, i.String())
	}
	return strings.Join(ss, ";")
}

func MakeMoveInstruction(sourceRow, sourceCol, destinationRow, destinationCol byte) Instruction {
	var i Instruction
	i.t = MoveInstruction
	i.row, i.col = sourceRow, sourceCol
	i.d[0], i.d[1] = destinationRow, destinationCol
	return i
}

func MakeCaptureInstruction(row, col byte, c Color, k Kind) Instruction {
	var i Instruction
	i.t = CaptureInstruction
	i.row, i.col = row, col
	i.d[0], i.d[1] = byte(c), byte(k)
	return i
}

func MakeCrownInstruction(row, col byte) Instruction {
	var i Instruction
	i.t = CrownInstruction
	i.row, i.col = row, col
	return i
}

func PerformInstructions(b *Board, is []Instruction) error {
	for _, i := range is {
		switch i.t {
		case MoveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			b.Move(fromRow, fromCol, toRow, toCol)
		case CaptureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := Color(i.d[0]), Kind(i.d[1])
			actualColor, actualKind := b.Get(row, col)
			if capturedColor != actualColor || capturedKind != actualKind {
				return fmt.Errorf(
					"performed capture instruction of %s %s on row %d %d but piece is a %s %s",
					capturedColor, capturedKind,
					row, col,
					actualColor, actualKind,
				)
			}
			b.Clear(row, col)
		case CrownInstruction:
			b.Crown(i.row, i.col)
		default:
			return fmt.Errorf("invalid instruction type %s", i.t)
		}
	}
	return nil
}

// if you pass 'is' to PerformInstructions
// the same 'is' passed to UndoInstructions will undo the instructions performed
// you don't need to reverse them
func UndoInstructions(b *Board, is []Instruction) {
	for k := len(is) - 1; k >= 0; k-- {
		i := is[k]
		switch i.t {
		case MoveInstruction:
			fromRow, fromCol := i.row, i.col
			toRow, toCol := i.d[0], i.d[1]
			b.Move(toRow, toCol, fromRow, fromCol)
		case CaptureInstruction:
			row, col := i.row, i.col
			capturedColor, capturedKind := Color(i.d[0]), Kind(i.d[1])
			b.Set(row, col, capturedColor, capturedKind)
		case CrownInstruction:
			b.Uncrown(i.row, i.col)
		default:
			panic(fmt.Sprintf("Invalid instruction type %s", i.t))
		}
	}
}
