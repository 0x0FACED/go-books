package gobooks_test

import (
	"testing"

	gobooks "github.com/0x0FACED/go-books"
	"github.com/stretchr/testify/assert"
)

// https://www.googleapis.com/books/v1/volumes?q=isbn:978-5-4461-1766-6
func TestFindByISBN(t *testing.T) {
	client := gobooks.NewGoogleBooksClient()

	resp, err := client.GetByISBN("978-5-4461-1766-6")

	assert.NoError(t, err)

	assert.NotNil(t, resp)
	// First author
	assert.Equal(t, "Эндрю Таненбаум", resp.Items[0].VolumeInfo.Authors[0])
	// Second author
	assert.Equal(t, "Ник Фимстер", resp.Items[0].VolumeInfo.Authors[1])
	// Third author
	assert.Equal(t, "Дэвид Уэзеролл", resp.Items[0].VolumeInfo.Authors[2])
	// ISBN_13
	assert.Equal(t, "9785446117666", resp.Items[0].VolumeInfo.IndustryIdentifiers[0].Identifier)
	// ISBN_10
	assert.Equal(t, "5446117662", resp.Items[0].VolumeInfo.IndustryIdentifiers[1].Identifier)
}
