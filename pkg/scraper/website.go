package scraper

// Website represents a website to be scraped
type Website struct {
	URL           string
	NameSelector  string
	PriceSelector string
}
