// 代码生成时间: 2025-08-18 08:40:05
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "log"

    "github.com/gorilla/mux"
)

// UserType defines the type of user.
type UserType int

const (
    UserTypeAdmin UserType = iota
    UserTypeRegular
)

// User structure represents a user with their permissions.
type User struct {
    ID       string    `json:"id"`
    Username string    `json:"username"`
    Type     UserType  `json:"type"`
}

// userRouter is the router for user-related endpoints.
var userRouter = mux.NewRouter()

// userPermissions contains the permissions for each user type.
var userPermissions = map[UserType][]string{
    UserTypeAdmin:    {"create", "read", "update", "delete"},
    UserTypeRegular: {"read"},
}

// GetUserPermissions returns permissions for a given user type.
func GetUserPermissions(userType UserType) []string {
    return userPermissions[userType]
}

// CreateUser handles the creation of a new user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Define the request body structure.
    var newUser User
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Here you would add logic to save the user to a database.
    // For simplicity, we just write the user back to the response.
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newUser)
}

// GetUserPermissions handles requests to get a user's permissions.
func GetUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]
    // Here you would fetch the user from the database and check their permissions.
    // For simplicity, we assume the user is of UserTypeAdmin.
    userPermissions := GetUserPermissions(UserTypeAdmin)
    json.NewEncoder(w).Encode(userPermissions)
}

func main() {
    // Define the routes.
    userRouter.HandleFunc("/users", CreateUser).Methods("POST")
    userRouter.HandleFunc("/users/{id}/permissions", GetUserPermissionsHandler).Methods("GET")

    // Start the server.
    log.Println("Server starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", userRouter))
}
