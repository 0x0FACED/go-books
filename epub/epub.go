package epub

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

type OPF struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
}

type Metadata struct {
	Identifiers []Identifier `xml:"identifier"`
}

type Identifier struct {
	ID     string `xml:"id,attr"`
	Value  string `xml:",chardata"`
	Scheme string `xml:"scheme,attr"`
}

type EPUB struct {
	Metadata Metadata
}

func (epub EPUB) ISBN() string {
	for _, id := range epub.Metadata.Identifiers {
		if isISBN(id) {
			return id.Value
		}
	}
	return ""
}

func (epub EPUB) CleanISBN() string {
	for _, id := range epub.Metadata.Identifiers {
		if isISBN(id) {
			return cleanISBN(id.Value)
		}
	}
	return ""
}

func ParseEPUBFromReader(r io.Reader) (*EPUB, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	containerData, err := readFileFromZip(zipReader, "META-INF/container.xml")
	if err != nil {
		return nil, err
	}

	var container struct {
		Rootfiles struct {
			Rootfile struct {
				FullPath string `xml:"full-path,attr"`
			} `xml:"rootfile"`
		} `xml:"rootfiles"`
	}
	if err := xml.Unmarshal(containerData, &container); err != nil {
		return nil, err
	}

	opfData, err := readFileFromZip(zipReader, container.Rootfiles.Rootfile.FullPath)
	if err != nil {
		return nil, err
	}

	var opf OPF
	if err := xml.Unmarshal(opfData, &opf); err != nil {
		return nil, err
	}

	return &EPUB{Metadata: opf.Metadata}, nil
}

func ParseEPUB(epubPath string) (*EPUB, error) {
	zipReader, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	containerData, err := readFileFromZip(&zipReader.Reader, "META-INF/container.xml")
	if err != nil {
		return nil, err
	}

	var container struct {
		Rootfiles struct {
			Rootfile struct {
				FullPath string `xml:"full-path,attr"`
			} `xml:"rootfile"`
		} `xml:"rootfiles"`
	}
	if err := xml.Unmarshal(containerData, &container); err != nil {
		return nil, err
	}

	opfData, err := readFileFromZip(&zipReader.Reader, container.Rootfiles.Rootfile.FullPath)
	if err != nil {
		return nil, err
	}

	var opf OPF
	if err := xml.Unmarshal(opfData, &opf); err != nil {
		return nil, err
	}

	return &EPUB{Metadata: opf.Metadata}, nil
}

func readFileFromZip(zipReader *zip.Reader, filePath string) ([]byte, error) {
	normalizedPath := strings.ReplaceAll(filePath, "\\", "/")
	for _, f := range zipReader.File {
		if strings.ReplaceAll(f.Name, "\\", "/") == normalizedPath {
			file, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()
			return io.ReadAll(file)
		}
	}
	return nil, errors.New("file not found in EPUB: " + filePath)
}

func isISBN(id Identifier) bool {
	if strings.ToLower(id.Scheme) == "isbn" {
		return true
	}
	if strings.HasPrefix(strings.ToLower(id.Value), "urn:isbn:") {
		return true
	}
	if strings.Contains(strings.ToLower(id.ID), "isbn") {
		return true
	}
	return false
}

func cleanISBN(isbn string) string {
	isbn = strings.TrimPrefix(isbn, "urn:isbn:")
	isbn = strings.TrimPrefix(isbn, "URN:ISBN:")
	var builder strings.Builder
	for _, r := range isbn {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		} else if r == 'X' || r == 'x' {
			builder.WriteRune('X')
		}
	}
	return builder.String()
}
