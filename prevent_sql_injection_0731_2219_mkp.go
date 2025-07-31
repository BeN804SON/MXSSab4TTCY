// 代码生成时间: 2025-07-31 22:19:36
package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/gorilla/mux"
)

// User model represents a user in the database
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex"`
    Password string
}

// SetupDatabase initializes the SQLite database connection
func SetupDatabase() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
        // Disables the default log output to avoid logging SQL statements
        SkipDefaultTransaction: true,
    })
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    
    // Migrate the schema
    db.AutoMigrate(&User{})
    return db
}

// CreateUser handles HTTP POST requests to create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    username := r.Form.Get("username")
    password := r.Form.Get("password")
    
    // Check for empty fields
    if username == "" || password == "" {
        fmt.Fprintf(w, "Username and password are required")
        return
    }
    
    // Hash the password before storing in the database
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Use GORM to prevent SQL injection
    db := SetupDatabase()
    if result := db.Where(&User{Username: username}).First(&User{}); result.Error == nil {
        fmt.Fprintf(w, "Username already exists")
        return
    } else if result.Error != gorm.ErrRecordNotFound {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    
    // Create a new user
    newUser := User{Username: username, Password: string(hashedPassword)}
    result := db.Create(&newUser)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "User created successfully")
}

func main() {
    router := mux.NewRouter()
    
    // Define routes
    router.HandleFunc("/users", CreateUser).Methods("POST")
    
    // Start the server
    log.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", router)
}