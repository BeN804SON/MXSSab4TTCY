// 代码生成时间: 2025-09-22 17:45:32
package main

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "log"

    // Import Gorilla Web toolkit (Mux) for URL routing
    "github.com/gorilla/mux"
)

// URLValidator defines the structure for the URL validation
type URLValidator struct {
    client *http.Client
}

// NewURLValidator creates a new instance of URLValidator
func NewURLValidator() *URLValidator {
    // Create a new HTTP client
    return &URLValidator{
        client: &http.Client{},
    }
}

// Validate checks if the URL is valid and reachable
func (v *URLValidator) Validate(u string) error {
    // Parse the URL to check its validity
    parsedURL, err := url.ParseRequestURI(u)
    if err != nil {
        return fmt.Errorf("invalid URL format: %w", err)
    }

    // Check for scheme and host presence in the URL
    if parsedURL.Scheme == "" || parsedURL.Host == "" {
        return fmt.Errorf("URL must include a scheme and host")
    }

    // Create an HTTP request to check if the URL is reachable
    req, err := http.NewRequest("HEAD", u, nil)
    if err != nil {
        return fmt.Errorf("failed to create HTTP request: %w", err)
    }

    // Set a short timeout for the request to not wait too long
    req.Timeout = 5 * time.Second

    // Send the request
    resp, err := v.client.Do(req)
    if err != nil {
        return fmt.Errorf("URL not reachable: %w", err)
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("URL returned a non-200 status code: %d", resp.StatusCode)
    }

    // If we reached this point, the URL is valid and reachable
    return nil
}

// setupRoutes sets up the routing for the application
func setupRoutes(r *mux.Router, validator *URLValidator) {
    r.HandleFunc("/validate", func(w http.ResponseWriter, req *http.Request) {
        // Retrieve the URL from the query parameters
        u := req.URL.Query().Get("url")
        if u == "" {
            http.Error(w, "URL parameter is required", http.StatusBadRequest)
            return
        }

        // Validate the URL
        if err := validator.Validate(u); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // If the URL is valid, return a success message
        fmt.Fprintf(w, "URL is valid and reachable: %s", u)
    })
}

func main() {
    // Create a new router
    r := mux.NewRouter()

    // Create a new URL validator
    validator := NewURLValidator()

    // Set up the routes
    setupRoutes(r, validator)

    // Start the server
    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}