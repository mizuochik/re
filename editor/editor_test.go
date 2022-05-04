package editor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mizuochikeita/re/editor"
)

func TestEditor(t *testing.T) {
	t.Run("OpenFile()", func(t *testing.T) {
		e := editor.New()
		if err := e.OpenFile("sample.txt"); err != nil {
			t.Fatal(err)
		}
		want := []string{
			"hello world",
			"bye world",
		}
		if diff := cmp.Diff(want, e.Buffer); diff != "" {
			t.Error(diff)
		}
	})
}
