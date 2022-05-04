package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestKey(t *testing.T) {
	t.Run("IsControl()", func(t *testing.T) {
		cases := []struct {
			desc  string
			given editor.Key
			want  bool
		}{
			{desc: "NULL", given: editor.Key{Value: '\x00'}, want: true},
			{desc: "CUU", given: editor.Key{EscapedSequence: []rune{'A'}}, want: false},
		}
		for _, c := range cases {
			got := c.given.IsControl()
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})

	t.Run("IsEscaped()", func(t *testing.T) {
		cases := []struct {
			desc  string
			given editor.Key
			want  bool
		}{
			{desc: "NULL", given: editor.Key{Value: '\x00'}, want: false},
			{desc: "CUU", given: editor.Key{EscapedSequence: []rune{'A'}}, want: true},
		}
		for _, c := range cases {
			got := c.given.IsEscaped()
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})
}
