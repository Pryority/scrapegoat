package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pryority/scrapegoat/pkg/scraper"
	"github.com/valyala/fasthttp"
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

	if *logFile != "" {
		// Initialize log output to file
		logFile, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Parse proxy list
	// proxies := strings.Split(*proxyList, ",")

	// Start scraping loop
	var wg sync.WaitGroup
	workerChan := make(chan struct{}, *numWorkers)

	for {
		workerChan <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Create a new fasthttp.Client with custom settings
			client := &fasthttp.Client{
				ReadTimeout:  *requestTimeout,
				WriteTimeout: *requestTimeout,
			}

			// Perform scraping using the current IP address
			err := scraper.ScrapeData(client)
			if err != nil {
				log.Println("Scraping error:", err)
			}

			// Rotate IP address
			time.Sleep(*rotateInterval)

			<-workerChan
		}()
	}

	wg.Wait()
}

// scrapeCallback is called for each scraped item
func scrapeCallback(name, price string) {
	if _, ok := scrapedProducts.LoadOrStore(name, struct{}{}); !ok {
		fmt.Println("Name:", name)
		fmt.Println("Price:", price)
		fmt.Println("========================")
	}
}
