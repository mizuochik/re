package editor

import (
	"unicode"
)

type Screen struct {
	Width   int
	Height  int
	Vscroll int
	Cx      int
	Cy      int
	Rows    []*ScreenRow
}

type ScreenRow struct {
	Body     string
	Len      int
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

func (s *Screen) Update(buffer []string) {
	var rs []*ScreenRow
	for _, row := range buffer {
		w := 0
		bi := 0
		l := 0
		var xs []int
		for i, c := range row {
			xs = append(xs, w)
			l++
			bw := w
			if c <= unicode.MaxASCII {
				w++
			} else {
				w += 2
			}
			if w > s.Width {
				rs = append(rs, &ScreenRow{
					Body:     row[bi:i],
					ScreenXs: append([]int(nil), xs[:len(xs)-1]...),
					Len:      l - 1,
				})
				bi = i
				l = 1
				w = w - bw
				xs = []int{0}
			}
		}
		rs = append(rs, &ScreenRow{
			Body:     row[bi:],
			ScreenXs: append([]int(nil), xs...),
			Len:      l,
		})
	}
	s.Rows = rs
}

func (s *Screen) Scroll(diff int) {
	v := s.Vscroll
	v += diff
	maxV := len(s.Rows) - s.Height
	minV := 0
	if v > maxV {
		v = maxV
	}
	if v < minV {
		v = minV
	}
	s.Vscroll = v
}

func (s *Screen) View() []*ScreenRow {
	bottom := s.Vscroll + s.Height
	if bottom > len(s.Rows) {
		bottom = len(s.Rows)
	}
	return s.Rows[s.Vscroll:bottom]
}

func (s *Screen) MoveCursorHorizontally(diff int) {
	if diff > 0 {
		rest := s.Rows[s.Cy].Len - s.Cx - 1
		if diff < rest {
			s.Cx += diff
		} else if s.Cy+1 >= len(s.Rows) {
			s.Cx = s.Rows[s.Cy].Len - 1
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
			if s.Rows[s.Cy].Len <= 0 {
				s.Cx = 0
			} else {
				s.Cx = s.Rows[s.Cy].Len - 1
			}
			s.MoveCursorHorizontally(diff + rest + 1)
		}
	}
}

func (s *Screen) MoveCursorVertically(diff int) {
	origSx := 0
	if s.Rows[s.Cy].Len > 0 {
		origSx = s.Rows[s.Cy].ScreenXs[s.Cx]
	}
	if diff > 0 {
		rest := len(s.Rows) - s.Cy - 1
		if diff < rest {
			s.Cy -= diff
		} else {
			s.Cy = len(s.Rows) - 1
		}
	} else if diff < 0 {
		rest := s.Cy
		if -diff < rest {
			s.Cy += diff
		} else {
			s.Cy = 0
		}
	}
	s.Cx = 0
	if s.Rows[s.Cy].Len > 0 {
		s.Cx = s.Rows[s.Cy].Len - 1
		for x, sx := range s.Rows[s.Cy].ScreenXs {
			if sx >= origSx {
				s.Cx = x
				break
			}
		}
	}
}

func (s *Screen) CursorPosition() (int, int) {
	x := s.Rows[s.Cy].ScreenXs[s.Cx]
	y := s.Cy - s.Vscroll
	return x, y
}
