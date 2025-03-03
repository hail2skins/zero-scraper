// Package main is the entry point of the application.
// It parses command-line flags, calls the scraping function, and outputs the results.
package main

import (
	"flag" // For command-line flag parsing
	"fmt"  // For formatted I/O
	"log"  // For logging errors and informational messages

	"github.com/hail2skins/zero-scraper/internal/scrape" // Import the scrape package from the internal directory. Adjust the module path as necessary.
)

func main() {
	// Define a command-line flag '-url' for the URL of the article to scrape.
	urlPtr := flag.String("url", "", "The URL of the news article to scrape")

	// Parse the command-line flags.
	flag.Parse()

	// If the URL flag is not provided, log a fatal error and exit.
	if *urlPtr == "" {
		log.Fatal("Please provide a URL using the -url flag")
	}

	// Call the ScrapeArticle function from the scrape package.
	// This function returns the article content, the author/byline, and an error, if any.
	article, byline, err := scrape.ScrapeArticle(*urlPtr)
	if err != nil {
		log.Fatalf("Error scraping article: %v", err)
	}

	// Check if any article content was returned.
	if article == "" {
		log.Println("No article content found.")
	} else {
		// Otherwise, print the scraped article content to the console.
		fmt.Println("Scraped Article Content:")
		fmt.Println(article)
	}

	// Output the scraped author information (byline) if available.
	if byline == "" {
		fmt.Println("No author information found.")
	} else {
		fmt.Println("Byline:", byline)
	}
}
