package main

import (
	"os"
)

func main() {
	keyBuf := make([]byte, 1)
outer:
	for {
		_, err := os.Stdin.Read(keyBuf)
		if err != nil {
			panic(err)
		}
		switch keyBuf[0] {
		case 'q':
			break outer
		}
	}
}
