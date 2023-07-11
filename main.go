package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	proxyList = []string{
		"http://193.233.202.75:8080",
		"http://139.59.1.14:8080",
	}

	// RotateIPInterval defines the interval between IP rotations
	RotateIPInterval = 5 * time.Second

	// Timeout for each request
	RequestTimeout = 30 * time.Second
)

func main() {
	// Create a new HTTP client with a custom transport
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(parseProxyURL(getRandomProxy())),
		},
		Timeout: RequestTimeout,
	}

	// Start scraping loop
	for {
		// Perform scraping using the current IP address
		err := scrapeData(client)
		if err != nil {
			log.Println("Scraping error:", err)
		}

		// Rotate IP address
		time.Sleep(RotateIPInterval)
		client.Transport.(*http.Transport).Proxy = http.ProxyURL(parseProxyURL(getRandomProxy()))
	}
}

// scrapeData performs the actual scraping
func scrapeData(client *http.Client) error {
	// Make a GET request to the target website
	response, err := client.Get("http://www.example.com")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Use goquery to parse the HTML response
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
	}

	// Extract data from the parsed HTML document
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	return nil
}

// getRandomProxy returns a random proxy URL from the proxyList
func getRandomProxy() string {
	rand.Seed(time.Now().UnixNano())
	return proxyList[rand.Intn(len(proxyList))]
}

// parseProxyURL parses the proxy URL string into a URL struct
func parseProxyURL(proxyURL string) *url.URL {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		log.Fatal(err)
	}
	return parsedURL
}
