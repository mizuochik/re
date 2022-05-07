package editor_test

import (
	"fmt"
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
			wantRows    []*editor.ScreenRow
		}{
			{
				desc:        "wraps long lines",
				givenWidth:  2,
				givenBuffer: []string{"abcd"},
				wantRows: []*editor.ScreenRow{
					{
						Body:     "ab",
						ScreenXs: []int{0, 1},
						Len:      2,
					},
					{
						Body:     "cd",
						ScreenXs: []int{0, 1},
						Len:      2,
					},
				},
			},
			{
				desc:        "considers non-ascii characters",
				givenWidth:  2,
				givenBuffer: []string{"あい"},
				wantRows: []*editor.ScreenRow{
					{
						Body:     "あ",
						ScreenXs: []int{0},
						Len:      1,
					},
					{
						Body:     "い",
						ScreenXs: []int{0},
						Len:      1,
					},
				},
			},
			{
				desc:        "considers ascii and non-ascii characters",
				givenWidth:  2,
				givenBuffer: []string{"aあいb"},
				wantRows: []*editor.ScreenRow{
					{
						Body:     "a",
						ScreenXs: []int{0},
						Len:      1,
					},
					{
						Body:     "あ",
						ScreenXs: []int{0},
						Len:      1,
					},
					{
						Body:     "い",
						ScreenXs: []int{0},
						Len:      1,
					},
					{
						Body:     "b",
						ScreenXs: []int{0},
						Len:      1,
					},
				},
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
			givenRows    []*editor.ScreenRow
			want         []*editor.ScreenRow
		}{
			{
				desc:         "no scroll",
				givenHeight:  2,
				givenVscroll: 0,
				givenRows: []*editor.ScreenRow{
					{Body: "a", ScreenXs: []int{0}, Len: 1},
					{Body: "b", ScreenXs: []int{0}, Len: 1},
					{Body: "c", ScreenXs: []int{0}, Len: 1},
				},
				want: []*editor.ScreenRow{
					{Body: "a", ScreenXs: []int{0}, Len: 1},
					{Body: "b", ScreenXs: []int{0}, Len: 1},
				},
			},
			{
				desc:         "scroll to bottom",
				givenHeight:  2,
				givenVscroll: 1,
				givenRows: []*editor.ScreenRow{
					{Body: "a", ScreenXs: []int{0}, Len: 1},
					{Body: "b", ScreenXs: []int{0}, Len: 1},
					{Body: "c", ScreenXs: []int{0}, Len: 1},
				},
				want: []*editor.ScreenRow{
					{Body: "b", ScreenXs: []int{0}, Len: 1},
					{Body: "c", ScreenXs: []int{0}, Len: 1},
				},
			},
			{
				desc:         "scroll to over bottom",
				givenHeight:  2,
				givenVscroll: 2,
				givenRows: []*editor.ScreenRow{
					{Body: "a", ScreenXs: []int{0}, Len: 1},
					{Body: "b", ScreenXs: []int{0}, Len: 1},
					{Body: "c", ScreenXs: []int{0}, Len: 1},
				},
				want: []*editor.ScreenRow{
					{Body: "c", ScreenXs: []int{0}, Len: 1},
				},
			},
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

	t.Run("MoveCursorHorizontally()", func(t *testing.T) {
		for _, c := range []struct {
			desc      string
			givenRows []*editor.ScreenRow
			givenCx   int
			givenCy   int
			givenDiff int
			wantCx    int
			wantCy    int
		}{
			{
				desc: "forward",
				givenRows: []*editor.ScreenRow{
					{Body: "abcd", Len: 4, ScreenXs: []int{0, 1, 2, 3}},
				},
				givenCx:   0,
				givenCy:   0,
				givenDiff: 2,
				wantCx:    2,
				wantCy:    0,
			},
			{
				desc: "back",
				givenRows: []*editor.ScreenRow{
					{Body: "abcd", Len: 4, ScreenXs: []int{0, 1, 2, 3}},
				},
				givenCx:   3,
				givenCy:   0,
				givenDiff: -2,
				wantCx:    1,
				wantCy:    0,
			},
			{
				desc: "forward to next line",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   0,
				givenCy:   0,
				givenDiff: 2,
				wantCx:    0,
				wantCy:    1,
			},
			{
				desc: "forward to end of screen",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   0,
				givenCy:   0,
				givenDiff: 4,
				wantCx:    1,
				wantCy:    1,
			},
			{
				desc: "back to before line",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   1,
				givenCy:   1,
				givenDiff: -2,
				wantCx:    1,
				wantCy:    0,
			},
			{
				desc: "back to start of screen",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   1,
				givenCy:   1,
				givenDiff: -4,
				wantCx:    0,
				wantCy:    0,
			},
		} {
			sc := &editor.Screen{
				Rows: c.givenRows,
				Cx:   c.givenCx,
				Cy:   c.givenCy,
			}
			sc.MoveCursorHorizontally(c.givenDiff)
			if diff := cmp.Diff(fmt.Sprintf("%d,%d", c.wantCx, c.wantCy), fmt.Sprintf("%d,%d", sc.Cx, sc.Cy)); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})

	t.Run("MoveCursorVertically()", func(t *testing.T) {
		for _, c := range []struct {
			desc      string
			givenRows []*editor.ScreenRow
			givenCx   int
			givenCy   int
			givenDiff int
			wantCx    int
			wantCy    int
		}{
			{
				desc: "go down",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   0,
				givenCy:   0,
				givenDiff: 1,
				wantCx:    0,
				wantCy:    1,
			},
			{
				desc: "go up",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   1,
				givenCy:   1,
				givenDiff: -1,
				wantCx:    1,
				wantCy:    0,
			},
			{
				desc: "go down over bottom",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   0,
				givenCy:   0,
				givenDiff: 2,
				wantCx:    0,
				wantCy:    1,
			},
			{
				desc: "go up over top",
				givenRows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   1,
				givenCy:   1,
				givenDiff: -2,
				wantCx:    1,
				wantCy:    0,
			},
		} {
			sc := &editor.Screen{
				Rows: c.givenRows,
				Cx:   c.givenCx,
				Cy:   c.givenCy,
			}
			sc.MoveCursorVertically(c.givenDiff)
			if diff := cmp.Diff(fmt.Sprintf("%d,%d", c.wantCx, c.wantCy), fmt.Sprintf("%d,%d", sc.Cx, sc.Cy)); diff != "" {
				t.Errorf("%s: %s", c.desc, diff)
			}
		}
	})
}
