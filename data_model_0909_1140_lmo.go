// 代码生成时间: 2025-09-09 11:40:19
package main

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
)

// User represents the data model for a user
type User struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"createdAt"`
}

// newUserResponse is the response structure for a newly created user
type newUserResponse struct {
    ID        string `json:"id"`
    CreatedAt string `json:"createdAt"`
}

// newUserRequest is the request structure for creating a new user
type newUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// CreateUser handles the creation of a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Parse the request body into newUserRequest
    var newUser newUserRequest
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Create a new user ID (in a real-world scenario, you should use a unique ID generator)
    userID := "user-" + strings.ReplaceAll(newUser.Name, " ", "")

    // Create a timestamp for when the user was created (in a real-world scenario, use time.Now())
    createdAt := "2023-04-01T12:00:00Z"

    // Create a new user object
    user := User{
        ID:        userID,
        Name:      newUser.Name,
        Email:     newUser.Email,
        CreatedAt: createdAt,
    }

    // Marshal the user object into JSON
    response, err := json.Marshal(newUserResponse{
        ID:        userID,
        CreatedAt: createdAt,
    })
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Write the response back to the client
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(response)
}

// main function to start the server
func main() {
    r := mux.NewRouter()
    r.HandleFunc("/users", CreateUser).Methods("POST")

    // Start the server
    http.ListenAndServe(":8080", r)
}