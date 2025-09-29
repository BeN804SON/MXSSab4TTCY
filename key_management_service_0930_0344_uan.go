// 代码生成时间: 2025-09-30 03:44:26
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

// Key represents the structure for a key
type Key struct {
    ID       string `json:"id"`
    KeyValue string `json:"keyValue"`
}

// CreateKeyResponse represents the response for creating a key
type CreateKeyResponse struct {
    Message string `json:"message"`
    Key     Key    `json:"key"`
}

// KeyManager is the struct for managing keys
type KeyManager struct {
    keys map[string]Key
}

// NewKeyManager creates a new instance of KeyManager
func NewKeyManager() *KeyManager {
    return &KeyManager{
        keys: make(map[string]Key),
    }
}

// GenerateKey generates a new unique key and stores it
func (km *KeyManager) GenerateKey() (Key, error) {
    id := generateID() // Assume generateID() is a function that generates a unique ID
    key := Key{
        ID:       id,
        KeyValue: "value" + id, // Simplified for demo purposes
    }
    km.keys[id] = key
    return key, nil
}

// GetKey retrieves a key by its ID
func (km *KeyManager) GetKey(id string) (Key, error) {
    key, exists := km.keys[id]
    if !exists {
        return Key{}, ErrKeyNotFound{ID: id}
    }
    return key, nil
}

// ErrKeyNotFound defines the error for a not found key
type ErrKeyNotFound struct {
    ID string
}

// Error returns the error message for a not found key
func (e ErrKeyNotFound) Error() string {
    return "Key not found with ID: " + e.ID
}

// CreateKeyHandler handles the creation of a new key
func CreateKeyHandler(km *KeyManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        key, err := km.GenerateKey()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(CreateKeyResponse{Message: "Key created successfully", Key: key})
    }
}

// GetKeyHandler handles the retrieval of a key
func GetKeyHandler(km *KeyManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := mux.Vars(r)["id"]
        key, err := km.GetKey(id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(key)
    }
}

func main() {
    r := mux.NewRouter()
    km := NewKeyManager()

    // Define routes
    r.HandleFunc("/keys", CreateKeyHandler(km)).Methods("POST")
    r.HandleFunc("/keys/{id}", GetKeyHandler(km)).Methods("GET")

    // Start the server
    log.Println("Starting the key management service on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
