package scraper

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var scrapedProducts = make(map[string]struct{})

// ScrapeData performs the actual scraping for a given website
func ScrapeData(client *fasthttp.Client, targetURL string, nameSelector string, priceSelector string) error {
	// Make a GET request to the target website
	statusCode, body, err := client.Get(nil, targetURL)
	if err != nil {
		return err
	}

	if statusCode != fasthttp.StatusOK {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}

	// Use goquery to parse the HTML response
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Extract data from the parsed HTML document
	doc.Find(".productLine_2N9kG").Each(func(i int, s *goquery.Selection) {
		name := s.Find(nameSelector).Text()
		price := s.Find(priceSelector).Text()

		if _, ok := scrapedProducts[name]; !ok {
			scrapedProducts[name] = struct{}{}
			log.Println(name)
			log.Println(price)
			log.Println("========================")
		}
	})

	return nil
}
