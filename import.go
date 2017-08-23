package main

import "fmt"
import "log"
import "github.com/PuerkitoBio/goquery"

type Talk struct {
  title string
  audience string
  format string
  suggestedBy string
  presenter string
  source string
}

func ScrapeWiki(url string) {
  doc := LoadDocument(url)
  // The table for self volunteered talks
  ImportTable(doc.Find("table.confluenceTable:nth-child(1)"), "volunteer")
  // The table for suggested talks
  ImportTable(doc.Find("table.confluenceTable:nth-child(2)"), "suggestion")
}

func ImportTable(table *goquery.Selection, source string) []Talk {
  var talks []Talk
  table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
    talk := new(Talk)
    talk.title = row.Find("td:nth-child(1)").Text()
    talk.audience = row.Find("td:nth-child(2)").Text()
    talk.format = row.Find("td:nth-child(3)").Text()
    talk.suggestedBy = row.Find("td:nth-child(4)").Text()
    talk.presenter = row.Find("td:nth-child(5)").Text()
    talk.source = source
    if talk.title != "" {
      talks = append(talks, *talk)
      fmt.Printf("Title %d: %s, %s\n", i, talk.title, talk.audience)
    }
  })
  return talks
}

func LoadDocument(url string) *goquery.Document {
  doc, err := goquery.NewDocument(url)
  if err != nil {
    log.Fatal(err)
  }
  return doc
}

func main() {
  ScrapeWiki("https://wiki.duraspace.org/display/samvera/Suggestions+for+Samvera+Connect+2017+Program")
}
