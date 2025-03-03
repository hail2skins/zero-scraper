// Package scrape provides functionality to scrape news articles.
package scrape

import (
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ScrapeArticle fetches the article content and byline from a given URL using Colly.
// It returns the article content, byline (author information), and an error if one occurred.
func ScrapeArticle(url string) (string, string, error) {
	// articleContent will accumulate the article's text.
	var articleContent string
	// author will store a combined byline if present.
	var author string
	// authors is a slice to store individual author names, if found.
	var authors []string

	// Create a new Colly collector.
	// The collector handles HTTP requests, response parsing, and event callbacks.
	c := colly.NewCollector(
	// Optionally restrict domains by uncommenting and modifying the following:
	// colly.AllowedDomains("apnews.com"),
	)

	// Capture the authors from a div with class "Page-authors" (used by AP News for the byline).
	c.OnHTML(`div.Page-authors`, func(e *colly.HTMLElement) {
		// Extract the complete byline text.
		text := e.Text
		if text != "" {
			// Trim any surrounding white space.
			author = strings.TrimSpace(text)
		}
		// Look for individual <a> elements inside the byline (often each name is linked).
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			name := strings.TrimSpace(el.Text)
			if name != "" {
				// Append the name to the authors slice.
				authors = append(authors, name)
			}
		})
	})

	// This callback extracts text content from all <p> (paragraph) elements to capture the article content.
	c.OnHTML("p", func(e *colly.HTMLElement) {
		// Append the text of every paragraph along with a newline.
		articleContent += e.Text + "\n"
	})

	// Handle HTTP errors during scraping.
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %v at %s\n", err, r.Request.URL)
	})

	// Begin the scraping process by visiting the specified URL.
	err := c.Visit(url)
	if err != nil {
		return "", "", err
	}

	// If individual author names were found but the combined author text is empty, join them.
	if author == "" && len(authors) > 0 {
		author = strings.Join(authors, " and ")
	}

	// Return the scraped article content, byline, and any error (nil if none occurred).
	return articleContent, author, nil
}
