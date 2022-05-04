package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestToControl(t *testing.T) {
	for _, c := range []struct {
		given rune
		want  rune
	}{
		{given: 'A', want: '\x01'},
		{given: 'z', want: '\x1a'},
	} {
		if diff := cmp.Diff(c.want, editor.ToControl(c.given)); diff != "" {
			t.Errorf("%c: %s", c.given, diff)
		}
	}
}
