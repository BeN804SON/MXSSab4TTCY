// 代码生成时间: 2025-10-03 02:39:22
package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "github.com/gorilla/mux"
)

// Desensitizer is a struct that holds the data to be desensitized
type Desensitizer struct {
    Data string `json:"data"`
}

// DesensitizedData is a struct that holds the desensitized data
type DesensitizedData struct {
    Desensitized string `json:"desensitized"`
}

// Desensitize is a function that takes a string and returns its desensitized version
func Desensitize(input string) string {
    // Simple desensitization: replace all digits with 'X'
    // This can be expanded to include more complex rules as needed
    return strings.ReplaceAll(input, "0123456789", "X")
}

// DesensitizeHandler is the HTTP handler function that processes the desensitization request
func DesensitizeHandler(w http.ResponseWriter, r *http.Request) {
    var d Desensitizer
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&d)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    desensitized := Desensitize(d.Data)
    responseData := DesensitizedData{Desensitized: desensitized}

    // Set the content type to JSON
    w.Header().Set("Content-Type", "application/json")
    // Marshal the responseData to JSON and write it to the response
    json.NewEncoder(w).Encode(responseData)
}

func main() {
    // Create a new router
    router := mux.NewRouter()

    // Define the route for the desensitization handler
    router.HandleFunc("/desensitize", DesensitizeHandler).Methods("POST")

    // Start the server on port 8080
    http.ListenAndServe(":8080", router)
}
