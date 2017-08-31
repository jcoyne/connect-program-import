# A program to scrape the Samvera Connect wiki into github issues

## Building
You must have go (https://golang.org/) installed.

Then fetch dependencies:
```
% go get -d
```

Build the binary:
```
% go build import.go
```

## Running
Visit https://github.com/settings/tokens and generate a personal token with
`public_repo` access. Paste the access token as the argument to the command as
seen here:

```
% ./import ACCESS_TOKEN
```
