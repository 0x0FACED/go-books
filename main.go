package main

import (
	"fmt"
	"os"

	"github.com/0x0FACED/go-books/epub"
	"github.com/0x0FACED/go-books/fb2"
)

func main() {
	file := "/home/podliva/Downloads/Tanenbaum_E__Fimster_N__Uezeroll_D_-_Kompyuternye_Seti_6-E_Izd__klassika_computer_science_-_2023.epub"
	file2 := "/home/podliva/Downloads/Достоевский_Федор_Вечный_муж_royallib_ru.fb2"
	epubBook, err := epub.ParseEPUB(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileEpub, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	epub2, err := epub.ParseEPUBFromReader(fileEpub)
	if err != nil {
		fmt.Println(err)
		return
	}

	parser := fb2.NewParser()

	fb2, err := parser.UnmarshalFromFile(file2, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fb2.ISBN())
	fmt.Println(epubBook.ISBN())
	fmt.Println(epub2.ISBN())
}
