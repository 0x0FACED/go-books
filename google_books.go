package gobooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

type VolumeInfo struct {
	Title               string       `json:"title"`
	Authors             []string     `json:"authors"`
	IndustryIdentifiers []IndustryID `json:"industryIdentifiers"`
	Language            string       `json:"language"`
	ImageLinks          *ImageLinks  `json:"imageLinks,omitempty"`
}

type IndustryID struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail,omitempty"`
	Thumbnail      string `json:"thumbnail,omitempty"`
}

type Volume struct {
	ID         string     `json:"id"`
	SelfLink   string     `json:"selfLink"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type VolumesResponse struct {
	Kind       string   `json:"kind"`
	TotalItems int      `json:"totalItems"`
	Items      []Volume `json:"items"`
}

type GoogleBooksClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewGoogleBooksClient() *GoogleBooksClient {
	return &GoogleBooksClient{
		BaseURL:    "https://www.googleapis.com/books/v1/volumes?",
		HTTPClient: &http.Client{},
	}
}

func (c *GoogleBooksClient) GetByISBN(isbn string) (*VolumesResponse, error) {
	if !ValidateISBN(isbn) {
		return nil, errors.New("invalid isbn")
	}

	query := url.Values{}
	query.Add("q", fmt.Sprintf("isbn:%s", isbn))

	req := fmt.Sprintf("%s%s", c.BaseURL, query.Encode())

	resp, err := c.HTTPClient.Get(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result VolumesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

func ValidateISBN(isbn string) bool {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", "")

	length := len(cleaned)
	if length != 10 && length != 13 {
		return false
	}

	if length == 10 {
		return validateISBN10(cleaned)
	}
	return validateISBN13(cleaned)
}

func validateISBN10(isbn string) bool {
	if len(isbn) != 10 {
		return false
	}

	sum := 0
	for i := range 9 {
		char := rune(isbn[i])
		if !unicode.IsDigit(char) {
			return false
		}
		digit := int(char - '0')
		sum += digit * (10 - i)
	}

	lastChar := rune(isbn[9])
	var lastDigit int
	if unicode.IsDigit(lastChar) {
		lastDigit = int(lastChar - '0')
	} else if unicode.ToUpper(lastChar) == 'X' {
		lastDigit = 10
	} else {
		return false
	}

	sum += lastDigit

	return sum%11 == 0
}

func validateISBN13(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}

	sum := 0
	for i := range 13 {
		char := rune(isbn[i])
		if !unicode.IsDigit(char) {
			return false
		}
		digit := int(char - '0')

		weight := 1
		if i%2 == 1 {
			weight = 3
		}
		sum += digit * weight
	}

	return sum%10 == 0
}
