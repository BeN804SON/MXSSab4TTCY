// 代码生成时间: 2025-10-11 21:49:53
package main

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
    "github.com/gorilla/schema"
)

// CDNConfig holds the configuration settings for the CDN
type CDNConfig struct {
    Origin string `schema:"origin"` // Origin server URL
}

// CDNHandler is the handler for CDN content distribution
type CDNHandler struct {
    config CDNConfig
    decoder *schema.Decoder
}

// NewCDNHandler creates a new instance of CDNHandler
func NewCDNHandler(config CDNConfig) *CDNHandler {
    return &CDNHandler{
        config: config,
        decoder: schema.NewDecoder(),
    }
}

// ServeHTTP handles HTTP requests and serves content from the origin server
func (h *CDNHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    var params struct{
        URL string `schema:"url"` // URL to fetch content from
    }
    if err := h.decoder.Decode(&params, r.URL.Query()); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate the URL parameter
    if params.URL == "" {
        http.Error(w, "URL parameter is required", http.StatusBadRequest)
        return
    }

    // Fetch content from the origin server
    response, err := http.Get(h.config.Origin + params.URL)
    if err != nil {
        http.Error(w, "Failed to fetch content from origin server", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    // Copy response body to the client
    if _, err := io.Copy(w, response.Body); err != nil {
        http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
        return
    }
}

func main() {
    // Define the CDN configuration
    config := CDNConfig{Origin: "https://example.com/"}

    // Create a new router
    router := mux.NewRouter()

    // Register the CDN handler with the router
    router.Handle("/cdn/{url}", NewCDNHandler(config)).Methods("GET")

    // Start the server
    fmt.Println("Starting CDN server on port 8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        fmt.Printf("Failed to start server: %v", err)
    }
}
