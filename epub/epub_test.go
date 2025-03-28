package epub_test

import (
	"testing"

	"github.com/0x0FACED/go-books/epub"
)

func TestExtactISBN(t *testing.T) {
	file := "/home/podliva/Downloads/Murder_in_the_Fog-Dominic_Butler.epub"

	epub, err := epub.ParseEPUB(file)
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(epub.ISBN())
}
