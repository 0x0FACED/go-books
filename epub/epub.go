package epub

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

type EPUB struct {
	Metadata Metadata `xml:"metadata"`
}

type Metadata struct {
	ISBN string `xml:"identifier"`
}

func (epub EPUB) ISBN() string {
	return epub.Metadata.ISBN
}

func readFileFromZIP(epubPath, filePath string) ([]byte, error) {
	zipReader, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.Name == filePath {
			file, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			return io.ReadAll(file)
		}
	}
	return nil, errors.New("file not found in EPUB")
}

func ParseEPUBFromReader(r io.Reader) (*EPUB, error) {
	var epub EPUB

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))

	err = decoder.Decode(&epub)
	if err != nil {
		return nil, err
	}

	return &epub, nil
}

func ParseEPUB(epubPath string) (*EPUB, error) {
	containerData, err := readFileFromZIP(epubPath, "META-INF/container.xml")
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

	opfData, err := readFileFromZIP(epubPath, container.Rootfiles.Rootfile.FullPath)
	if err != nil {
		return nil, err
	}

	var epub EPUB
	if err := xml.Unmarshal(opfData, &epub); err != nil {
		return nil, err
	}

	return &epub, nil
}
