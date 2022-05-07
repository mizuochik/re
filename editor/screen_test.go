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
			desc     string
			width    int
			buffer   []string
			wantRows []*editor.ScreenRow
		}{
			{
				desc:   "wraps long lines",
				width:  2,
				buffer: []string{"abcd"},
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
				desc:   "considers non-ascii characters",
				width:  2,
				buffer: []string{"あい"},
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
				desc:   "considers ascii and non-ascii characters",
				width:  2,
				buffer: []string{"aあいb"},
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
				Width: tt.width,
			}
			sc.Update(tt.buffer)
			if diff := cmp.Diff(tt.wantRows, sc.Rows); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("Scroll()", func(t *testing.T) {
		tests := []struct {
			desc        string
			height      int
			vscroll     int
			rows        []*editor.ScreenRow
			diff        int
			wantVscroll int
		}{
			{
				desc:    "scroll down",
				height:  2,
				vscroll: 0,
				rows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				diff:        1,
				wantVscroll: 1,
			},
			{
				desc:    "scroll down",
				height:  2,
				vscroll: 2,
				rows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				diff:        -1,
				wantVscroll: 1,
			},
			{
				desc:    "scroll down over bottom",
				height:  2,
				vscroll: 0,
				rows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				diff:        3,
				wantVscroll: 2,
			},
			{
				desc:    "scroll up over top",
				height:  2,
				vscroll: 2,
				rows: []*editor.ScreenRow{
					{},
					{},
					{},
					{},
				},
				diff:        -3,
				wantVscroll: 0,
			},
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Height:  tt.height,
				Rows:    tt.rows,
				Vscroll: tt.vscroll,
			}
			sc.Scroll(tt.diff)
			if diff := cmp.Diff(tt.wantVscroll, sc.Vscroll); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("View()", func(t *testing.T) {
		tests := []struct {
			desc    string
			height  int
			vscroll int
			rows    []*editor.ScreenRow
			want    []*editor.ScreenRow
		}{
			{
				desc:    "no scroll",
				height:  2,
				vscroll: 0,
				rows: []*editor.ScreenRow{
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
				desc:    "scroll to bottom",
				height:  2,
				vscroll: 1,
				rows: []*editor.ScreenRow{
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
				desc:    "scroll to over bottom",
				height:  2,
				vscroll: 2,
				rows: []*editor.ScreenRow{
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
				Height:  tt.height,
				Vscroll: tt.vscroll,
				Rows:    tt.rows,
			}
			got := sc.View()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("MoveCursorHorizontally()", func(t *testing.T) {
		tests := []struct {
			desc   string
			rows   []*editor.ScreenRow
			cx     int
			cy     int
			diff   int
			wantCx int
			wantCy int
		}{
			{
				desc: "forward",
				rows: []*editor.ScreenRow{
					{Body: "abcd", Len: 4, ScreenXs: []int{0, 1, 2, 3}},
				},
				cx:     0,
				cy:     0,
				diff:   2,
				wantCx: 2,
				wantCy: 0,
			},
			{
				desc: "back",
				rows: []*editor.ScreenRow{
					{Body: "abcd", Len: 4, ScreenXs: []int{0, 1, 2, 3}},
				},
				cx:     3,
				cy:     0,
				diff:   -2,
				wantCx: 1,
				wantCy: 0,
			},
			{
				desc: "forward to next line",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     0,
				cy:     0,
				diff:   2,
				wantCx: 0,
				wantCy: 1,
			},
			{
				desc: "forward to end of screen",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     0,
				cy:     0,
				diff:   4,
				wantCx: 1,
				wantCy: 1,
			},
			{
				desc: "back to before line",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     1,
				cy:     1,
				diff:   -2,
				wantCx: 1,
				wantCy: 0,
			},
			{
				desc: "back to start of screen",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "cd", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     1,
				cy:     1,
				diff:   -4,
				wantCx: 0,
				wantCy: 0,
			},
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Rows: tt.rows,
				Cx:   tt.cx,
				Cy:   tt.cy,
			}
			sc.MoveCursorHorizontally(tt.diff)
			if diff := cmp.Diff(fmt.Sprintf("%d,%d", tt.wantCx, tt.wantCy), fmt.Sprintf("%d,%d", sc.Cx, sc.Cy)); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})

	t.Run("MoveCursorVertically()", func(t *testing.T) {
		tests := []struct {
			desc   string
			rows   []*editor.ScreenRow
			cx     int
			cy     int
			diff   int
			wantCx int
			wantCy int
		}{
			{
				desc: "go down",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     0,
				cy:     0,
				diff:   1,
				wantCx: 0,
				wantCy: 1,
			},
			{
				desc: "go up",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     1,
				cy:     1,
				diff:   -1,
				wantCx: 1,
				wantCy: 0,
			},
			{
				desc: "go down over bottom",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     0,
				cy:     0,
				diff:   2,
				wantCx: 0,
				wantCy: 1,
			},
			{
				desc: "go up over top",
				rows: []*editor.ScreenRow{
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
					{Body: "bc", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     1,
				cy:     1,
				diff:   -2,
				wantCx: 1,
				wantCy: 0,
			},
			{
				desc: "go down and keep x on screen",
				rows: []*editor.ScreenRow{
					{Body: "あいう", Len: 3, ScreenXs: []int{0, 2, 4}},
					{Body: "abcdef", Len: 6, ScreenXs: []int{0, 1, 2, 3, 4, 5}},
				},
				cx:     1,
				cy:     0,
				diff:   1,
				wantCx: 2,
				wantCy: 1,
			},
			{
				desc: "go down and keep x on screen (dest row is shorter than source row)",
				rows: []*editor.ScreenRow{
					{Body: "あいう", Len: 3, ScreenXs: []int{0, 2, 4}},
					{Body: "ab", Len: 2, ScreenXs: []int{0, 1}},
				},
				cx:     2,
				cy:     0,
				diff:   1,
				wantCx: 1,
				wantCy: 1,
			},
		}
		for _, tt := range tests {
			sc := &editor.Screen{
				Rows: tt.rows,
				Cx:   tt.cx,
				Cy:   tt.cy,
			}
			sc.MoveCursorVertically(tt.diff)
			if diff := cmp.Diff(fmt.Sprintf("%d,%d", tt.wantCx, tt.wantCy), fmt.Sprintf("%d,%d", sc.Cx, sc.Cy)); diff != "" {
				t.Errorf("%s: %s", tt.desc, diff)
			}
		}
	})
}
