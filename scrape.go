package main

import (
	"libgenscrape/views/components"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
)

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}
func GetLists(query string) []components.BookType {
	c := colly.NewCollector()
	baseUrl := "https://libgen.is"
	searchUrl := baseUrl + "/search.php?req=" + EncodeParam(query)
	var books []components.BookType
	c.OnHTML("table.c tbody", func(e *colly.HTMLElement) {
		// Extract data from HTML elements
		idCol := 0
		// sizeCol := 0
		extensionCol := 0
		// mirrorCol1 := 0
		mirrorCol2 := 0
		titleCol := 0
		e.ForEach("tr", func(rowNum int, row *colly.HTMLElement) {
			var book components.BookType
			var mirrors []string
			md5 := ""
			row.ForEach("td", func(colNum int, col *colly.HTMLElement) {
				text := col.Text
				if rowNum == 0 {

					switch text {
					case "ID":
						idCol = colNum
					case "Extension":
						extensionCol = colNum
					case "Mirrors":
						// mirrorCol1 = colNum
						mirrorCol2 = colNum + 1
					case "Title":
						titleCol = colNum
					}
				} else {
					switch colNum {
					case idCol:
						book.ID = text
					case extensionCol:
						book.Extension = text
					// case mirrorCol1:
					// 	mirrors = append(mirrors, col.ChildAttr("a", "href"))
					case mirrorCol2:
						mirrors = append(mirrors, col.ChildAttr("a", "href"))
					case titleCol:
						book.BookName = text
						md5 = strings.Split(col.ChildAttr("a[id]", "href"), "=")[1]
					}
				}

			})
			if rowNum != 0 {
				coverFolder, err := strconv.ParseInt(book.ID, 10, 32)
				if err != nil {
					log.Fatal(err)
				}
				if coverFolder < 1000 {
					coverFolder = 0
				} else {
					coverFolder = coverFolder / 1000 * 1000
				}
				book.ImgUrl = baseUrl + "/covers/" + strconv.Itoa(int(coverFolder)) + "/" + strings.ToLower(md5)
				book.Mirrors = mirrors
				books = append(books, book)
			}
		})
	})
	err := c.Visit(searchUrl)
	if err != nil {
		log.Info(err)
	}

	return books
}

func GetDownloadUrl(visitUrl string) string {
	c := colly.NewCollector()
	website := GetDomain(visitUrl)

	var downloadUrl string

	selector := "#download"

	if website == "libgen.li" {
		selector = "#main"
	}

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		downloadUrl = e.ChildAttr("a", "href")
	})
	err := c.Visit(visitUrl)
	if err != nil {
		log.Info(err)
	}
	if website == "libgen.li" {
		downloadUrl = "https://libgen.li/" + downloadUrl
	}
	return strings.TrimSpace(downloadUrl)
}
