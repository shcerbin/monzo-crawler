# monzo-crawler

This is a simple recursive web crawler that takes a domain as input and returns a list of all the links on that domain.

app.go - handle all concurrency logic and save the results to the csv file.

crawler.go - sends requests and parse html response.

# Usage

Run the following command to run the web crawler
```bash
go run ./cmd/main.go
```

You can also specify the domain as an environment variable
```bash
DOMAIN=example.com go run ./cmd/main.go
```

# Highlights
- Concurrency is controlled by concurrentRequests variable
- The results are saved to a file as soon as new links are found, so as not to lose them if something goes wrong
- Everything fit into two files with logic: app.go and crawler.go
- The code is covered with integration tests that run in parallel
- 23s to crawl monzo.com domain
