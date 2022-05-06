package editor

import "unicode"

type Screen struct {
	Width   int
	Height  int
	Rows    []string
	Vscroll int
	Cx      int
	Cy      int
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
	var bottom int
	if s.Height < len(s.Rows) {
		bottom = s.Height + s.Vscroll
	} else {
		bottom = s.Height + len(s.Rows)
	}
	return s.Rows[s.Vscroll:bottom]
}
