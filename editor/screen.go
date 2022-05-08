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
	var (
		rows     []*ScreenRow
		screenXs []int
	)
	for _, row := range buffer {
		rr := []rune(row)
		rr = append(rr, ' ') // Add a space expressing end of the row
		l := 0
		w := 0
		var nw int
		if rr[0] <= unicode.MaxASCII {
			nw = 1
		} else {
			nw = 2
		}
		for i := 0; i < len(rr); i++ {
			screenXs = append(screenXs, w)
			w = nw
			if i < len(rr)-1 {
				if rr[i+1] <= unicode.MaxASCII {
					nw = w + 1
				} else {
					nw = w + 2
				}
				if nw > s.Width {
					rows = append(rows, &ScreenRow{
						Len:      len(rr[l : i+1]),
						Body:     string(rr[l : i+1]),
						ScreenXs: screenXs,
					})
					l = i + 1
					nw = nw - w
					w = 0
					screenXs = nil
				}
			}
		}
		rows = append(rows, &ScreenRow{
			Len:      len(rr[l:]),
			Body:     string(rr[l:]),
			ScreenXs: screenXs,
		})
	}
	s.Rows = rows
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
			s.Cy += diff
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
