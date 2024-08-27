package app

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shcerbin/monzo-crawler/crawler"
)

type Domain = string
type Link = string

var concurrentRequests = 19

func Run(domain Domain) {
	log.Printf("Checking domain: %s\n", domain)

	resultFileName := strings.ReplaceAll(domain, "/", "_") + "_result.csv"

	resultFile, err := os.Create(resultFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer resultFile.Close()

	writer := csv.NewWriter(resultFile)

	var (
		uniqueLinks    sync.Map
		lock           sync.Mutex
		requestLimiter = make(chan struct{}, concurrentRequests)
		timer          = time.Now()
		counter        atomic.Uint64
	)

	// Recursively visit links
	var f func(string)
	f = func(link string) {
		// would block if requestLimiter channel is already filled
		requestLimiter <- struct{}{}
		counter.Add(1)
		crawlerLinks := crawler.FindAllLinks(link)
		<-requestLimiter

		if len(crawlerLinks) == 0 {
			return
		}

		var newLinks []string

		lock.Lock()
		for crawlerLink := range crawlerLinks {
			parsedCrawlerLink, err := parseLink(crawlerLink, domain)
			if err != nil {
				continue
			}

			if _, ok := uniqueLinks.Load(parsedCrawlerLink); ok {
				continue
			}

			// store the visited URL
			uniqueLinks.Store(parsedCrawlerLink, struct{}{})
			newLinks = append(newLinks, parsedCrawlerLink)

			fmt.Printf(`
#########################################################################################
found new link: %s, request counter: %d, time: %s
#########################################################################################

`, parsedCrawlerLink, counter.Load(), time.Since(timer))

			// add a CSV record to the output file
			writer.Write([]string{parsedCrawlerLink})
		}
		writer.Flush()
		lock.Unlock()

		// after storing results, recursively visit the next new domain link
		for _, newLink := range newLinks {
			if !strings.HasPrefix(newLink, "https://"+domain) && !strings.HasPrefix(newLink, "http://"+domain) {
				continue
			}

			f(newLink)
		}
	}

	// start the recursive link checking
	parsedDomainLink, _ := parseLink("/", domain)
	f(parsedDomainLink)

	fmt.Printf(`total time: %s, request counter: %d, resultFileName: %s`,
		time.Since(timer), counter.Load(), resultFileName)
}

func parseLink(link Link, domain Domain) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		log.Printf("failed to parse link: %s, err: %s", link, err)
		return "", err
	}

	if u.Scheme == "" && strings.Contains(domain, ":") {
		u.Scheme = "http"
	} else {
		u.Scheme = "https"
	}

	if u.Host == "" {
		u.Host = strings.TrimPrefix(strings.TrimPrefix(domain, "http://"), "https://")
	}

	parsedLink := strings.TrimSuffix(u.String(), "/")

	// validate parsed link

	if parsedLink == "https://"+domain+"/" {
		log.Printf("skip host link: %s", parsedLink)
		return parsedLink, errors.New("host link")
	}

	if strings.HasPrefix(parsedLink, "tel:") {
		log.Printf("skip tel link: %s", parsedLink)
		return parsedLink, errors.New("tel link")
	}

	if strings.HasPrefix(parsedLink, "mailto:") {
		log.Printf("skip mailto link: %s", parsedLink)
		return parsedLink, errors.New("mailto link")
	}

	if strings.HasSuffix(parsedLink, ".pdf") {
		log.Printf("skip pdf link: %s", parsedLink)
		return parsedLink, errors.New("pdf link")
	}

	if strings.HasSuffix(parsedLink, ".mp3") {
		log.Printf("skip mp3 link: %s", parsedLink)
		return parsedLink, errors.New("mp3 link")
	}

	return parsedLink, nil
}
