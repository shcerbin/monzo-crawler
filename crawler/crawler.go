package crawler

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func FindAllLinks(url string) map[string]struct{} {
	result := make(map[string]struct{})

	log.Printf("visit url: %s\n", url)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("GET url: %s, err: %s", url, err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Recursively visit nodes in the parse tree
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					result[a.Val] = struct{}{}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return result
}
