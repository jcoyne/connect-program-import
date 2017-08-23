# A program to scrape the Samvera Connect wiki into github issues

## Building
You must have go (https://golang.org/) installed.

Then run these commands:
```
% go get github.com/google/go-github/github
% go get github.com/PuerkitoBio/goquery
% go build import.go
```

## Running
Visit https://github.com/settings/tokens and generate a personal token with
`public_repo` access. Paste the access token as the argument to the command as
seen here:

```
% ./import ACCESS_TOKEN
```
