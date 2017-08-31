package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Talk struct {
	title       string
	audience    string
	format      string
	suggestedBy string
	presenter   string
}

func (talk Talk) String() string {
	return fmt.Sprintf("%s\n\nSuggested by: %s\nPresenter: %s\nFormat: %s\nAudience: %s", talk.title, talk.suggestedBy, talk.presenter, talk.format, talk.audience)
}

func ScrapeWiki(url string) ([]Talk, error) {
	doc, err := LoadDocument(url)
	if err != nil {
		return nil, err
	}
	return ImportTable(doc.Find("table.confluenceTable")), nil
}

func ImportTable(table *goquery.Selection) []Talk {
	var talks []Talk
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		talk := &Talk{
			title:       row.Find("td:nth-child(1)").Text(),
			audience:    row.Find("td:nth-child(2)").Text(),
			format:      row.Find("td:nth-child(3)").Text(),
			suggestedBy: row.Find("td:nth-child(4)").Text(),
			presenter:   row.Find("td:nth-child(5)").Text(),
		}
		if talk.title != "" {
			talks = append(talks, *talk)
		}
	})
	return talks
}

func LoadDocument(url string) (*goquery.Document, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: import [ACCESS_TOKEN]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func GetAccessToken() (string, error) {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		return "", errors.New("access token is missing")
	}
	return args[0], nil
}

func InitClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func NewIssue(talk Talk) *github.IssueRequest {
	req := new(github.IssueRequest)
	req.Title = &talk.title
	body := fmt.Sprintf("%s", talk)
	req.Body = &body
	return req
}

func CreateIssues(token string, talks []Talk) error {
	ctx := context.Background()
	client := InitClient(ctx, token)
	for _, talk := range talks {
		req := NewIssue(talk)
		issue, _, err := client.Issues.Create(ctx, "jcoyne", "bl-test", req)
		if err != nil {
			return err
		}
		fmt.Printf("ISSUE %s", issue)
	}
	return nil
}

func HandleError(err error) {
	fmt.Printf("ERROR %s\n", err)
	os.Exit(1)
}

func main() {
	token, err := GetAccessToken()
	if err != nil {
		HandleError(err)
	}
	var url = "https://wiki.duraspace.org/display/samvera/Suggestions+for+Samvera+Connect+2017+Program"
	talks, err2 := ScrapeWiki(url)
	if err2 != nil {
		HandleError(err2)
	}
	if err = CreateIssues(token, talks); err != nil {
		HandleError(err)
	}
}
