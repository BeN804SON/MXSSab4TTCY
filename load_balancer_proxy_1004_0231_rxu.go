// 代码生成时间: 2025-10-04 02:31:25
Features:
- Simple round-robin load balancing.
- WebSocket proxying to backend services.

Usage:
- Start the proxy with desired backends.
- Connect clients to the proxy using WebSockets.
- Proxy will forward messages to backend services and return responses to clients.
*/

package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/gorilla/websocket"
)

// Global variables to hold backends and current index
var backends = []string{"ws://backend1:port","ws://backend2:port"}
var currentIndex = 0

// Upgrader will upgrade an HTTP server connection to a WebSocket protocol connection.
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for simplicity
    },
}

// proxyHandler handles WebSocket requests and proxies them to the backend services.
func proxyHandler(w http.ResponseWriter, r *http.Request) {
    backendURL := backends[currentIndex]
    currentIndex = (currentIndex + 1) % len(backends) // Round-robin load balancing

    // Dial the backend WebSocket
    connBackend, _, err := websocket.DefaultDialer.Dial(backendURL, nil)
    if err != nil {
        log.Printf("Failed to connect to backend: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer connBackend.Close()

    // Upgrade the incoming request to a WebSocket connection
    connClient, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    defer connClient.Close()

    // Start proxying messages between client and backend
    go func() {
        for {
            _, message, err := connClient.ReadMessage()
            if err != nil {
                break
            }
            if err := connBackend.WriteMessage(websocket.TextMessage, message); err != nil {
                break
            }
        }
    }()

    for {
        _, message, err := connBackend.ReadMessage()
        if err != nil {
            break
        }
        if err := connClient.WriteMessage(websocket.TextMessage, message); err != nil {
            break
        }
    }
}

func main() {
    // Set up HTTP server to handle WebSocket requests
    http.HandleFunc("/ws", proxyHandler)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Starting proxy server on port %s", port)
    log.Fatal(http.ListenAndServe(":" + port, nil))
}
