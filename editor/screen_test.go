package editor_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestScreen(t *testing.T) {
	t.Run("Update()", func(t *testing.T) {
		tests := []struct {
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
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Width: tt.givenWidth,
			}
			sc.Update(tt.givenBuffer)
			if diff := cmp.Diff(tt.wantRows, sc.Rows); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("Scroll()", func(t *testing.T) {
		tests := []struct {
			desc         string
			givenHeight  int
			givenVscroll int
			givenRows    []*editor.ScreenRow
			givenDiff    int
			wantVscroll  int
		}{
			{
				desc:         "scroll down",
				givenHeight:  2,
				givenVscroll: 0,
				givenRows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				givenDiff:   1,
				wantVscroll: 1,
			},
			{
				desc:         "scroll down",
				givenHeight:  2,
				givenVscroll: 2,
				givenRows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				givenDiff:   -1,
				wantVscroll: 1,
			},
			{
				desc:         "scroll down over bottom",
				givenHeight:  2,
				givenVscroll: 0,
				givenRows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				givenDiff:   3,
				wantVscroll: 2,
			},
			{
				desc:         "scroll up over top",
				givenHeight:  2,
				givenVscroll: 2,
				givenRows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				givenDiff:   -3,
				wantVscroll: 0,
			},
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Height:  tt.givenHeight,
				Rows:    tt.givenRows,
				Vscroll: tt.givenVscroll,
			}
			sc.Scroll(tt.givenDiff)
			if diff := cmp.Diff(tt.wantVscroll, sc.Vscroll); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("View()", func(t *testing.T) {
		tests := []struct {
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
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Height:  tt.givenHeight,
				Vscroll: tt.givenVscroll,
				Rows:    tt.givenRows,
			}
			got := sc.View()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("MoveCursorHorizontally()", func(t *testing.T) {
		tests := []struct {
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
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Rows: tt.givenRows,
				Cx:   tt.givenCx,
				Cy:   tt.givenCy,
			}
			sc.MoveCursorHorizontally(tt.givenDiff)
			if diff := cmp.Diff(fmt.Sprintf("%d,%d", tt.wantCx, tt.wantCy), fmt.Sprintf("%d,%d", sc.Cx, sc.Cy)); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("MoveCursorVertically()", func(t *testing.T) {
		tests := []struct {
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
			{
				desc: "go down and keep x on screen",
				givenRows: []*editor.ScreenRow{
					{Body: "あいう", Len: 3, ScreenXs: []int{0, 2, 4}},
					{Body: "abcdef", Len: 6, ScreenXs: []int{0, 1, 2, 3, 4, 5}},
				},
				givenCx:   1,
				givenCy:   0,
				givenDiff: 1,
				wantCx:    2,
				wantCy:    1,
			},
			{
				desc: "go down and keep x on screen (dest row is shorter than source row)",
				givenRows: []*editor.ScreenRow{
					{Body: "あいう", Len: 3, ScreenXs: []int{0, 2, 4}},
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
				},
				givenCx:   2,
				givenCy:   0,
				givenDiff: 1,
				wantCx:    1,
				wantCy:    1,
			},
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Rows: tt.givenRows,
				Cx:   tt.givenCx,
				Cy:   tt.givenCy,
			}
			sc.MoveCursorVertically(tt.givenDiff)
			if diff := cmp.Diff(fmt.Sprintf("%d,%d", tt.wantCx, tt.wantCy), fmt.Sprintf("%d,%d", sc.Cx, sc.Cy)); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})
}
