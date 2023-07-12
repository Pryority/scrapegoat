package main

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/pryority/scrapegoat/pkg/scraper"
	"github.com/valyala/fasthttp"

	log "github.com/sirupsen/logrus"
)

var (
	// Command-line flags
	targetURL      = flag.String("url", "https://www.bestbuy.ca/en-ca/category/laptops-macbooks/20352", "Target website URL")
	proxyList      = flag.String("proxies", "http://193.233.202.75:8080,http://139.59.1.14:8080", "Comma-separated list of proxies")
	rotateInterval = flag.Duration("rotate-interval", 5*time.Second, "Interval between IP rotations")
	requestTimeout = flag.Duration("request-timeout", 10*time.Second, "Timeout for each request")
	numWorkers     = flag.Int("num-workers", 5, "Number of concurrent workers")
	logFile        = flag.String("log-file", "", "Path to the log file")

	scrapedProducts sync.Map
)

func main() {
	flag.Parse()

	// Configure logrus
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Set log output to a file if provided
	if *logFile != "" {
		logFile, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Start scraping loop
	var wg sync.WaitGroup
	workerChan := make(chan struct{}, *numWorkers)

	websites := []scraper.Website{
		{
			URL:           "https://www.bestbuy.ca/en-ca/category/laptops-macbooks/20352",
			NameSelector:  ".productItemName_3IZ3c",
			PriceSelector: ".price_2j8lL",
		},
		{
			URL:           "https://www.website2.com",
			NameSelector:  ".name-selector",
			PriceSelector: ".price-selector",
		},
	}

	for _, website := range websites {
		workerChan <- struct{}{}
		wg.Add(1)

		go func(website scraper.Website) {
			defer wg.Done()

			// Create a new fasthttp.Client with custom settings
			client := &fasthttp.Client{
				ReadTimeout:  *requestTimeout,
				WriteTimeout: *requestTimeout,
			}

			// Perform scraping using the current IP address
			err := scraper.ScrapeData(client, website.URL, website.NameSelector, website.PriceSelector)
			if err != nil {
				log.Error("Scraping error:", err)
			}

			// Rotate IP address
			time.Sleep(*rotateInterval)

			<-workerChan
		}(website)
	}

	wg.Wait()
}
