// 代码生成时间: 2025-09-29 02:01:22
// autocomplete_service.go
package main

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
)

// AutocompleteService handles the search autocomplete functionality.
type AutocompleteService struct {
    // Storage for suggestions, could be a database or in-memory store in a real application.
    suggestions []string
}

// NewAutocompleteService creates a new instance of AutocompleteService.
func NewAutocompleteService() *AutocompleteService {
    return &AutocompleteService{
        suggestions: []string{"apple", "banana", "cherry", "date"},
    }
}

// Autocomplete searches for suggestions based on the given prefix.
func (service *AutocompleteService) Autocomplete(w http.ResponseWriter, r *http.Request) {
    var err error
    defer func() {
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }()

    vars := mux.Vars(r)
    prefix := vars["prefix"]

    if prefix == "" {
        err = errors.New("prefix parameter is required")
        return
    }

    // Search for suggestions that start with the given prefix.
    suggestions := []string{}
    for _, suggestion := range service.suggestions {
        if strings.HasPrefix(strings.ToLower(suggestion), strings.ToLower(prefix)) {
            suggestions = append(suggestions, suggestion)
        }
    }

    // Return the suggestions as JSON.
    w.Header().Set("Content-Type", "application/json")
    if err = json.NewEncoder(w).Encode(suggestions); err != nil {
        return
    }
}

// main function sets up the server and routes.
func main() {
    r := mux.NewRouter()
    service := NewAutocompleteService()

    // Define the route for the autocomplete service.
    r.HandleFunc("/autocomplete/{prefix}", service.Autocomplete).Methods("GET")

    // Start the server.
    http.ListenAndServe(":8080", r)
}