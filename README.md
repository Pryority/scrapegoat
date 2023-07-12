# ScrapeGOat

This is a Go program that performs web scraping to extract product names and prices from a target website. It uses the Goquery library for HTML parsing and the fasthttp library for making HTTP requests.

## Prerequisites

- Go 1.16 or higher installed

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/pryority/scrapegoat.git
   ```

2. Navigate to the project directory:

   ```bash
   cd scrapegoat
   ```

3. Build the program:

   ```bash
   go build
   ```

4. Run the program:

   ```bash
   ./scrapegoat
   ```

## Configuration

- The target website URL is set to "<https://www.bestbuy.ca/en-ca/category/laptops-macbooks/20352>". You can modify this URL in the `scrapeData` function of the `main.go` file.

- The program uses a rotating proxy list defined in the `proxyList` variable. You can modify the list or add more proxies as needed.

- The program is configured to scrape data with multiple workers. The number of workers can be adjusted by modifying the `NumWorkers` variable.

## Output

The program will scrape the target website periodically and print the product names and prices. It ensures that each product is scraped only once by keeping track of the already scraped products in the `scrapedProducts` map.

## Extending the Program

To extend the program to scrape additional websites or perform custom parsing, you can follow these steps:

1. Create a new package in the `pkg` directory for the website-specific code. For example, create a `website1` package for scraping and parsing website1.

2. Implement the necessary functions in the new package. For scraping, create a function that takes a client and performs the scraping logic. For parsing, create a function that takes a goquery.Document and extracts the desired data.

3. In the `main.go` file, import the new package and add a new `Website` struct to the `websites` slice. Provide the appropriate URL, parse function, and price by name function.

4. Run the program, and it will scrape and parse the new website in addition to the existing ones.

## License

This project is licensed under the [MIT License](LICENSE).
