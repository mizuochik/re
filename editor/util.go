package editor

func ToControl(r rune) rune {
	return r & 0x1f
}
