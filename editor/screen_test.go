package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestScreen(t *testing.T) {
	t.Run("Update()", func(t *testing.T) {
		for _, c := range []struct {
			desc        string
			givenWidth  int
			givenBuffer []string
			wantRows    []string
		}{
			{
				desc:        "wraps long lines",
				givenWidth:  2,
				givenBuffer: []string{"abcd"},
				wantRows:    []string{"ab", "cd"},
			},
			{
				desc:        "considers non-ascii characters",
				givenWidth:  2,
				givenBuffer: []string{"あい"},
				wantRows:    []string{"あ", "い"},
			},
			{
				desc:        "considers ascii and non-ascii characters",
				givenWidth:  2,
				givenBuffer: []string{"aあいb"},
				wantRows:    []string{"a", "あい", "b"},
			},
		} {
			sc := &editor.Screen{
				Width: c.givenWidth,
			}
			sc.Update(c.givenBuffer)
			if diff := cmp.Diff(c.wantRows, sc.Rows); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})

	t.Run("View()", func(t *testing.T) {
		for _, c := range []struct {
			desc         string
			givenHeight  int
			givenVscroll int
			givenRows    []string
			want         []string
		}{
			{},
		} {
			sc := &editor.Screen{
				Height:  c.givenHeight,
				Vscroll: c.givenVscroll,
				Rows:    c.givenRows,
			}
			got := sc.View()
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})
}
