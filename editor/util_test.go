package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestToControl(t *testing.T) {
	tests := []struct {
		given rune
		want  rune
	}{
		{given: 'A', want: '\x01'},
		{given: 'z', want: '\x1a'},
	}
	for _, tt := range tests {
		if diff := cmp.Diff(tt.want, editor.ToControl(tt.given)); diff != "" {
			t.Errorf("%c: %s", tt.given, diff)
		}
	}
}
