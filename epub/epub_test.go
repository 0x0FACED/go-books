package epub_test

import (
	"os"
	"testing"

	"github.com/0x0FACED/go-books/epub"
	"github.com/stretchr/testify/assert"
)

func TestParseFromReader(t *testing.T) {
	path := "/home/test.epub"
	f, err := os.Open(path)

	assert.NoError(t, err)

	book, err := epub.ParseEPUBFromReader(f)

	assert.NoError(t, err)

	assert.NotNil(t, book)
	assert.Equal(t, "978-5-4461-1766-6", book.ISBN())
	assert.Equal(t, "9785446117666", book.CleanISBN())
}
