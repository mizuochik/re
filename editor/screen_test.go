package editor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
			sc := &Screen{
				Width: c.givenWidth,
			}
			sc.Update(c.givenBuffer)
			if diff := cmp.Diff(c.wantRows, sc.Rows); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})
}
