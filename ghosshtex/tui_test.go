package ghosshtex

import (
	"log"
	"os"
	"testing"
)

func TestTUI(t *testing.T) {
	tui := NewEditorTUI(os.Stdout)
	if _, err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}
