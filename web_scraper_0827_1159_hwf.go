// 代码生成时间: 2025-08-27 11:59:37
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "log"
    "golang.org/x/net/html"
)

// WebScraper is a struct that holds the URL to scrape.
type WebScraper struct {
    URL string
}

// NewWebScraper creates a new instance of WebScraper.
func NewWebScraper(url string) *WebScraper {
    return &WebScraper{URL: url}
}

// Scrape fetches the HTML content of the webpage and extracts text.
func (s *WebScraper) Scrape() (string, error) {
    // Send an HTTP request to the webpage
    resp, err := http.Get(s.URL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check if the response status code is 200 OK
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to retrieve webpage: status code %d", resp.StatusCode)
    }

    // Create a tokenizer to parse the HTML content
    tokenizer := html.NewTokenizer(resp.Body)
    var content strings.Builder

    // Iterate over the tokens in the HTML
    for {
        tt := tokenizer.Next()
        if tt == html.ErrorToken {
            // If there's an error, return it
            return "", tokenizer.Err()
        }
        if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
            t := tokenizer.Token()
            // Skip script and style tags to avoid extracting JavaScript and CSS
            if t.Data == "script" || t.Data == "style" {
                for {
                    tt = tokenizer.Next()
                    if tt == html.EndTagToken && tokenizer.Token().Data == "script" || tokenizer.Token().Data == "style" {
                        break
                    }
                }
                continue
            }
        }
        if tt == html.TextToken {
            text := tokenizer.Token().Data
            // Add the text to the content builder, excluding any script or style content
            content.WriteString(text)
        }
        if tt == html.EndTagToken && tokenizer.Token().Data == "html" {
            break
        }
    }
    return content.String(), nil
}

// SaveToFile saves the scraped content to a file.
func (s *WebScraper) SaveToFile(fileName string) error {
    content, err := s.Scrape()
    if err != nil {
        return err
    }
    // Create the file and write the content to it
    file, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = io.WriteString(file, content)
    if err != nil {
        return err
    }
    return nil
}

func main() {
    url := "http://example.com"
    scraper := NewWebScraper(url)
    fileName := "scrapped_content.txt"
    err := scraper.SaveToFile(fileName)
    if err != nil {
        log.Fatalf("An error occurred while scraping: %s", err)
    }
    fmt.Printf("Content saved to %s
", fileName)
}