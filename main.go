// Package main implements a simple web scraper using the Colly library.
// It extracts the text content of an article as well as author information (byline)
// from a specified news article URL.
package main

import (
	"flag"    // For command-line flag parsing.
	"fmt"     // For formatted I/O.
	"log"     // For logging errors and informational messages.
	"strings" // For string manipulation.

	"github.com/gocolly/colly/v2" // The Colly web scraping framework.
)

// scrapeArticle accepts a URL string and returns the scraped article content,
// the author/byline information, and an error (if any) from the scraping process.
func scrapeArticle(url string) (string, string, error) {
	// articleContent will accumulate the article's text.
	var articleContent string
	// author will store a combined byline if present.
	var author string
	// authors is a slice to store individual author names, if found.
	var authors []string

	// Create a new Colly collector.
	// The collector is the main component of Colly that handles HTTP requests,
	// response parsing, and event callbacks.
	c := colly.NewCollector(
	// Uncomment and modify AllowedDomains if you need to restrict the scraper to certain domains.
	// colly.AllowedDomains("apnews.com"),
	)

	// Capture the authors from a div with class "Page-authors"
	// "Page-authors" is a class that is used by AP News for the byline.
	// This callback is triggered when the collector visits an element matching the CSS selector.
	c.OnHTML(`div.Page-authors`, func(e *colly.HTMLElement) {
		// Extract the complete byline text.
		text := e.Text
		if text != "" {
			// Trim any surrounding white space.
			author = strings.TrimSpace(text)
		}
		// Also look for individual <a> elements inside the byline (often each name is linked).
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			// Trim the text content of each <a> element.
			name := strings.TrimSpace(el.Text)
			if name != "" {
				// Append the name to the authors slice.
				authors = append(authors, name)
			}
		})
	})

	// This callback extracts text content from all <p> (paragraph) elements.
	// It's a simple method to capture the article content.
	c.OnHTML("p", func(e *colly.HTMLElement) {
		// Append the text of every paragraph along with a newline.
		articleContent += e.Text + "\n"
	})

	// This callback handles HTTP errors that occur during scraping.
	c.OnError(func(r *colly.Response, err error) {
		// Log the error along with the URL that caused it.
		log.Printf("Error: %v at %s\n", err, r.Request.URL)
	})

	// Start the scraping process by visiting the specified URL.
	err := c.Visit(url)
	if err != nil {
		// If visiting the URL fails, return the error.
		return "", "", err
	}

	// If individual author names were found but the combined author text is empty,
	// join them together with " and " as a separator.
	if author == "" && len(authors) > 0 {
		author = strings.Join(authors, " and ")
	}

	// Return the scraped article content, author byline, and nil as the error.
	return articleContent, author, nil
}

// main is the entry point of the program.
func main() {
	// Define a command-line flag "-url" for the URL of the article to scrape.
	urlPtr := flag.String("url", "", "The URL of the news article to scrape")
	// Parse the command-line flags.
	flag.Parse()

	// If the URL flag is not provided, exit with an error.
	if *urlPtr == "" {
		log.Fatal("Please provide a URL using the -url flag")
	}

	// Call scrapeArticle with the provided URL.
	article, byline, err := scrapeArticle(*urlPtr)
	if err != nil {
		// Log and exit if scraping fails.
		log.Fatalf("Error scraping article: %v", err)
	}

	// If no article content was scraped, log a message.
	if article == "" {
		log.Println("No article content found.")
	} else {
		// Otherwise, output the scraped article content.
		fmt.Println("Scraped Article Content:")
		fmt.Println(article)
	}

	// Output the scraped author information (byline).
	if byline == "" {
		fmt.Println("No author information found.")
	} else {
		fmt.Println("Byline:", byline)
	}
}
