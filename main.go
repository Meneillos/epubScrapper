package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Epub struct {
	XMLName          xml.Name `xml:"package"`
	Text             string   `xml:",chardata"`
	Version          string   `xml:"version,attr"`
	UniqueIdentifier string   `xml:"unique-identifier,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	Metadata         struct {
		Text       string `xml:",chardata"`
		Dc         string `xml:"dc,attr"`
		Opf        string `xml:"opf,attr"`
		Title      string `xml:"title"`
		Author     string `xml:"creator"`
		Subject    string `xml:"subject"`
		Identifier string `xml:"identifier"`
		Dates      []struct {
			Text  string `xml:",chardata"`
			Event string `xml:"event"`
		} `xml:"date"`
		Meta []struct {
			Name    string `xml:"name,attr"`
			Content string `xml:"content,attr"`
		} `xml:"meta"`
	} `xml:"metadata"`
}

func main() {
	var workpath string
	var err error
	if len(os.Args) > 1 {
		workpath, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		workpath, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	filter := workpath + "\\*.epub"
	epubList, err := filepath.Glob(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(epubList)
	epub, err := zip.OpenReader("Gideon_la_Novena_Tamsyn_Muir.epub")
	if err != nil {
		panic(err)
	}
	defer epub.Close()

	var file *zip.File
	for _, f := range epub.File {
		if f.Name == "OEBPS/content.opf" {
			file = f
			break
		}
	}

	content, err := file.Open()
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	_, err = io.Copy(&output, content)
	if err != nil {
		panic(err)
	}
	content.Close()

	data := new(Epub)
	err = xml.Unmarshal(output.Bytes(), data)
	if err != nil {
		panic(err)
	}
}
