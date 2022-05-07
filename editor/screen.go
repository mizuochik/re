package editor

import (
	"unicode"
)

type Screen struct {
	Width   int
	Height  int
	Rows    []string
	Vscroll int
	Cx      int
	Cy      int
	Rows_   []*ScreenRow
}

type ScreenRow struct {
	Body     string
	ScreenXs []int
}

func (r *ScreenRow) UpdateXs() {
	var xs []int
	x := 0
	for _, c := range r.Body {
		xs = append(xs, x)
		if c <= unicode.MaxASCII {
			x++
		} else {
			x += 2
		}
	}
	r.ScreenXs = xs
}

func (s *Screen) Update_(buffer []string) {
	var rs []*ScreenRow
	for _, row := range buffer {
		w := 0
		l := 0
		var xs []int
		for i, c := range row {
			xs = append(xs, w)
			bw := w
			if c <= unicode.MaxASCII {
				w++
			} else {
				w += 2
			}
			if w > s.Width {
				rs = append(rs, &ScreenRow{
					Body:     row[l:i],
					ScreenXs: append([]int(nil), xs[:len(xs)-1]...),
				})
				l = i
				w = w - bw
				xs = []int{0}
			}
		}
		rs = append(rs, &ScreenRow{
			Body:     row[l:],
			ScreenXs: append([]int(nil), xs...),
		})
	}
	s.Rows_ = rs
}

func (s *Screen) Update(buffer []string) {
	var rs []string
	for _, row := range buffer {
		l := 0
		w := 0
		for i, c := range row {
			if c <= unicode.MaxASCII {
				w++
			} else {
				w += 2
			}
			if w > s.Width {
				rs = append(rs, row[l:i])
				l = i
				w = 0
			}
		}
		rs = append(rs, row[l:])
	}
	s.Rows = rs
}

func (s *Screen) Scroll(diff int) {
	s.Vscroll += diff
}

func (s *Screen) View() []string {
	bottom := s.Vscroll + s.Height
	if bottom > len(s.Rows) {
		bottom = len(s.Rows)
	}
	return s.Rows[s.Vscroll:bottom]
}

func (s *Screen) MoveCursorHorizontally(diff int) {
	if diff > 0 {
		rest := len(s.Rows[s.Cy]) - s.Cx - 1
		if diff < rest {
			s.Cx += diff
		} else if s.Cy+1 >= len(s.Rows) {
			s.Cx = len(s.Rows[s.Cy]) - 1
		} else {
			s.Cy++
			s.Cx = 0
			s.MoveCursorHorizontally(diff - rest - 1)
		}
	} else if diff < 0 {
		rest := s.Cx
		if -diff < rest {
			s.Cx += diff
		} else if s.Cy <= 0 {
			s.Cx = 0
		} else {
			s.Cy--
			if len(s.Rows[s.Cy]) <= 0 {
				s.Cx = 0
			} else {
				s.Cx = len(s.Rows[s.Cy]) - 1
			}
			s.MoveCursorHorizontally(diff + rest + 1)
		}
	}
}
