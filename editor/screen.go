package editor

import "unicode"

type Screen struct {
	Width   int
	Rows    []string
	Vscroll int
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

func (s *Screen) View(Height int) []string {
	return nil
}
