package main

import (
	"os"

	"github.com/shcerbin/monzo-crawler/app"
)

func main() {
	domain := os.Getenv("DOMAIN")

	if domain == "" {
		domain = "monzo.com"
	}

	app.Run(domain)
}
