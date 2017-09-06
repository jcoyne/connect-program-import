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

func scrapeWiki(url string) ([]Talk, error) {
	doc, err := loadDocument(url)
	if err != nil {
		return nil, err
	}
	return importTable(doc.Find("table.confluenceTable")), nil
}

func importTable(table *goquery.Selection) []Talk {
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

func loadDocument(url string) (*goquery.Document, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: import [ACCESS_TOKEN]\n")
	flag.PrintDefaults()
}

func getAccessToken() (string, error) {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return "", errors.New("access token is missing")
	}
	return flag.Args()[0], nil
}

func initClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func isValidFormat(format string) bool {
	switch format {
	case
		"Breakout",
		"Lightning talk",
		"Panel",
		"Plenary",
		"Presentation",
		"Unconference",
		"Workshop":
		return true
	}
	return false
}

func isValidAudience(audience string) bool {
	switch audience {
	case
		"All",
		"Developers",
		"Managers",
		"System Administrators",
		"Metadata":
		return true
	}
	return false
}

func newIssue(talk Talk) *github.IssueRequest {
	req := new(github.IssueRequest)
	req.Title = &talk.title
	body := fmt.Sprintf("%s", talk)
	req.Body = &body
	var labels []string

	if isValidFormat(talk.format) {
		labels = append(labels, talk.format)
	}

	if isValidAudience(talk.audience) {
		labels = append(labels, talk.audience)
	}

	if labels != nil {
		req.Labels = &labels
	}

	return req
}

func createIssues(token string, talks []Talk) error {
	ctx := context.Background()
	client := initClient(ctx, token)
	for _, talk := range talks {
		req := newIssue(talk)
		issue, _, err := client.Issues.Create(ctx, "jcoyne", "bl-test", req)
		if err != nil {
			return err
		}
		fmt.Printf("ISSUE %s", issue)
	}
	return nil
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR %s\n", err)
	os.Exit(1)
}

func main() {
	token, err := getAccessToken()
	if err != nil {
		os.Exit(1)
	}
	var url = "https://wiki.duraspace.org/display/samvera/Suggestions+for+Samvera+Connect+2017+Program"
	talks, err2 := scrapeWiki(url)
	if err2 != nil {
		handleError(err2)
	}
	if err = createIssues(token, talks); err != nil {
		handleError(err)
	}
}
