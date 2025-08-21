// 代码生成时间: 2025-08-21 23:50:41
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

// Handler is a function type that handles HTTP requests.
type Handler func(http.ResponseWriter, *http.Request)

// handleNotFound handles 404 not found HTTP responses.
func handleNotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "404 Not Found: %s\
", r.URL.Path)
}

// handlePing is a simple handler that returns a 'pong' message.
func handlePing(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong\
")
}

// handleHealthCheck checks the health of the service and returns a status.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Service is up and running\
")
}

// NewRouter creates a new Gorilla Mux router and sets up the routes.
func NewRouter() *mux.Router {
    router := mux.NewRouter()
    // Handle 404 by providing a custom not found handler
    router.NotFoundHandler = http.HandlerFunc(handleNotFound)

    // Define routes
    router.HandleFunc("/ping", handlePing).Methods("GET")
    router.HandleFunc("/health", handleHealthCheck).Methods("GET\)

    return router
}

// main is the entry point of the program.
func main() {
    // Create a new router
    router := NewRouter()

    // Start the HTTP server
    fmt.Println("Starting the HTTP server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        fmt.Printf("Failed to start server: %s\
", err)
    }
}
