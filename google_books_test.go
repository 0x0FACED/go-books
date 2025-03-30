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

func TestValidateISBN(t *testing.T) {
	testCases := []struct {
		name    string
		isbn    string
		isValid bool
	}{
		{"Valid ISBN-10", "5170815808", true},
		{"Valid ISBN-13", "9785170815807", true},
		{"Valid ISBN-13 with hyphens", "978-5-4461-1766-6", true},
		{"Valid ISBN-10 without hyphens", "5446117662", true},
		{"Invalid ISBN-10 checksum", "1234567890", false},
		{"Invalid ISBN-13 checksum", "9781234567890", false},
		{"Invalid ISBN empty string", "", false},
		{"Invalid ISBN-13 with wrong checksum", "978-5-4461-1766-5", false},
		{"Valid ISBN-10 with hyphens", "0-306-40615-2", true},
		{"Valid ISBN-10 without hyphens", "0306406152", true},
		{"Valid ISBN-13 matching ISBN-10", "9780306406157", true},
		{"Valid ISBN-10 with X check digit", "156881111X", true},
		{"Invalid ISBN-10 with lowercase x", "156881111x", true},
		{"Invalid characters in ISBN", "978-5-4461-1766-X", false},
		{"Too short", "123", false},
		{"Too long", "9785170815807123", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := gobooks.ValidateISBN(tc.isbn)
			if got != tc.isValid {
				t.Errorf("ValidateISBN(%q) = %v, expected %v", tc.isbn, got, tc.isValid)
			}
		})
	}
}
