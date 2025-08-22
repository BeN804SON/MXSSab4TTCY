// 代码生成时间: 2025-08-22 19:32:41
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "github.com/gorilla/mux"
)

// Message represents the structure for a message
type Message struct {
    Content string `json:"content"`
}

// NewMessage creates a new message
func NewMessage(content string) *Message {
    return &Message{Content: content}
}

// MessageHandler handles the incoming message and sends a notification
func MessageHandler(w http.ResponseWriter, r *http.Request) {
    var msg Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Logic to send notification, e.g., to a message queue or a service
    // For demonstration, we'll just log the message
    fmt.Println("Received message: ", msg.Content)

    // Respond with a 200 OK status
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "message received"})
}

func main() {
    // Create a new router
    router := mux.NewRouter()

    // Register the MessageHandler for POST requests to "/message"
    router.HandleFunc("/message", MessageHandler).Methods("POST")

    // Start the server
    fmt.Println("Starting message notification system on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        fmt.Println("Failed to start server: ", err)
    }
}
