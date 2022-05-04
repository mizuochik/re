package editor

import "unicode"

type Key struct {
	Value           rune
	EscapedSequence []rune
}

func (k Key) IsControl() bool {
	return !k.IsEscaped() && unicode.IsControl(k.Value)
}

func (k Key) IsEscaped() bool {
	return len(k.EscapedSequence) > 0
}
