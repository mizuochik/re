package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestKey(t *testing.T) {
	t.Run("IsControl()", func(t *testing.T) {
		tests := []struct {
			desc string
			key  editor.Key
			want bool
		}{
			{desc: "NULL", key: editor.Key{Value: '\x00'}, want: true},
			{desc: "CUU", key: editor.Key{EscapedSequence: []rune{'A'}}, want: false},
		}
		for _, tt := range tests {
			got := tt.key.IsControl()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("IsEscaped()", func(t *testing.T) {
		tests := []struct {
			desc string
			key  editor.Key
			want bool
		}{
			{desc: "NULL", key: editor.Key{Value: '\x00'}, want: false},
			{desc: "CUU", key: editor.Key{EscapedSequence: []rune{'A'}}, want: true},
		}
		for _, tt := range tests {
			got := tt.key.IsEscaped()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})
}
