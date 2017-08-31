package main

import "fmt"
import "flag"
import "log"
import "os"
import "github.com/PuerkitoBio/goquery"
import "golang.org/x/oauth2"
import "context"
import "github.com/google/go-github/github"

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

func ScrapeWiki(url string) []Talk {
	doc := LoadDocument(url)
	return ImportTable(doc.Find("table.confluenceTable"))
}

func ImportTable(table *goquery.Selection) []Talk {
	var talks []Talk
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		talk := new(Talk)
		talk.title = row.Find("td:nth-child(1)").Text()
		talk.audience = row.Find("td:nth-child(2)").Text()
		talk.format = row.Find("td:nth-child(3)").Text()
		talk.suggestedBy = row.Find("td:nth-child(4)").Text()
		talk.presenter = row.Find("td:nth-child(5)").Text()
		if talk.title != "" {
			talks = append(talks, *talk)
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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: import [ACCESS_TOKEN]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func GetAccessToken() string {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Access token is missing.")
		os.Exit(1)
	}
	return args[0]
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

func CreateIssues(token string, talks []Talk) {
	ctx := context.Background()
	client := InitClient(ctx, token)
	for _, talk := range talks {
		req := NewIssue(talk)
		issue, _, err := client.Issues.Create(ctx, "jcoyne", "bl-test", req)
		if err != nil {
			fmt.Printf("ERROR %s", err)
			os.Exit(1)
		}
		fmt.Printf("ISSUE %s", issue)
	}
}

func main() {
	token := GetAccessToken()
	talks := ScrapeWiki("https://wiki.duraspace.org/display/samvera/Suggestions+for+Samvera+Connect+2017+Program")
	CreateIssues(token, talks)
}
