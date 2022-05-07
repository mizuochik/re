package editor

import (
	"fmt"
	"os"
)

func ToControl(r rune) rune {
	return r & 0x1f
}

func Debugf(format string, a ...interface{}) {
	f, err := os.OpenFile("/tmp/re.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintf(f, format, a...)
	fmt.Fprintln(f)
}
