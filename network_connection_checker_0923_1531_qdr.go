// 代码生成时间: 2025-09-23 15:31:35
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gorilla/mux"
)

// NetworkChecker is a struct that holds the route for checking network connection
type NetworkChecker struct {
    *mux.Router
}

// NewNetworkChecker initializes a new NetworkChecker
func NewNetworkChecker() *NetworkChecker {
    router := mux.NewRouter().StrictSlash(true)
    return &NetworkChecker{Router: router}
}

// CheckConnection is the handler function for checking network connection
func (nc *NetworkChecker) CheckConnection(w http.ResponseWriter, r *http.Request) {
    // Define the URL to check connection
    checkURL := "https://www.google.com"
    response, err := http.Get(checkURL)
    
    // Check for errors
    if err != nil {
        fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err.Error())
        return
    }
    defer response.Body.Close()
    
    // Check the HTTP status code for a successful connection
    if response.StatusCode != http.StatusOK {
        fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"Non-successful status code: %d\"}", response.StatusCode)
        return
    }
    
    // If the connection is successful
    fmt.Fprintf(w, "{\"status\":\"success\"}")
}

func main() {
    // Create a new instance of NetworkChecker
    nc := NewNetworkChecker()
    
    // Define the route for checking network connection
    nc.HandleFunc("/check-connection", nc.CheckConnection).Method("GET")
    
    // Start the server
    fmt.Println("Starting network connection checker on port 8080")
    err := http.ListenAndServe(":8080