package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestToControl(t *testing.T) {
	tests := []struct {
		in   rune
		want rune
	}{
		{in: 'A', want: '\x01'},
		{in: 'z', want: '\x1a'},
	}
	for _, tt := range tests {
		if diff := cmp.Diff(tt.want, editor.ToControl(tt.in)); diff != "" {
			t.Errorf("%c: %s", tt.in, diff)
		}
	}
}
