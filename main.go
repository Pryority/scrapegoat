package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

var (
	proxyList = []string{
		"http://193.233.202.75:8080",
		"http://139.59.1.14:8080",
	}

	// RotateIPInterval defines the interval between IP rotations
	RotateIPInterval = 5 * time.Second

	// Timeout for each request
	RequestTimeout = 10 * time.Second

	NumWorkers = 5
)

func main() {
	var wg sync.WaitGroup

	workerChan := make(chan struct{}, NumWorkers)

	// Start scraping loop
	for {
		workerChan <- struct{}{}

		wg.Add(1)

		go func() {
			defer wg.Done()

			// Create a new fasthttp.Client
			client := &fasthttp.Client{
				ReadTimeout:  RequestTimeout,
				WriteTimeout: RequestTimeout,
			}

			// Perform scraping using the current IP address
			err := scrapeData(client)
			if err != nil {
				log.Println("Scraping error:", err)
			}

			// Rotate IP address
			time.Sleep(RotateIPInterval)

			<-workerChan
		}()
	}
}

// scrapeData performs the actual scraping
func scrapeData(client *fasthttp.Client) error {
	// Make a GET request to the target website
	url := "https://www.bestbuy.ca/en-ca/category/laptops-macbooks/20352"
	statusCode, body, err := client.Get(nil, url)
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
		name := s.Find(".productItemName_3IZ3c").Text()
		price := s.Find(".price_2j8lL").Text()

		fmt.Println("Name:", name)
		fmt.Println("Price:", price)
		fmt.Println("========================")
	})

	return nil
}
