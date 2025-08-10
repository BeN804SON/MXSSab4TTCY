// 代码生成时间: 2025-08-10 19:41:12
package main

import (
    "net/http"
    "encoding/json"
    "log"
    "github.com/gorilla/mux"
)

// InventoryItem represents an item in the inventory
type InventoryItem struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Quantity    int    `json:"quantity"`
}

// InventoryManager is a simple struct to manage inventory
type InventoryManager struct {
    items map[string]InventoryItem
}

// NewInventoryManager creates a new InventoryManager
func NewInventoryManager() *InventoryManager {
    return &InventoryManager{
        items: make(map[string]InventoryItem),
    }
}

// AddItem adds a new item to the inventory
func (im *InventoryManager) AddItem(item InventoryItem) error {
    if _, exists := im.items[item.ID]; exists {
        return ErrItemExists
    }
    im.items[item.ID] = item
    return nil
}

// GetItem retrieves an item by its ID from the inventory
func (im *InventoryManager) GetItem(id string) (InventoryItem, error) {
    item, exists := im.items[id]
    if !exists {
        return InventoryItem{}, ErrItemNotFound
    }
    return item, nil
}

// UpdateItem updates an existing item in the inventory
func (im *InventoryManager) UpdateItem(id string, newQuantity int) error {
    if _, exists := im.items[id]; !exists {
        return ErrItemNotFound
    }
    im.items[id].Quantity = newQuantity
    return nil
}

// DeleteItem removes an item from the inventory
func (im *InventoryManager) DeleteItem(id string) error {
    if _, exists := im.items[id]; !exists {
        return ErrItemNotFound
    }
    delete(im.items, id)
    return nil
}

// ErrItemExists is an error indicating that an item with the same ID already exists
var ErrItemExists = errors.New("item with this ID already exists")

// ErrItemNotFound is an error indicating that an item with the specified ID was not found
var ErrItemNotFound = errors.New("item not found")

// router is the HTTP router for the inventory management system
var router = mux.NewRouter()

func main() {
    im := NewInventoryManager()
    // Add some initial items to the inventory
    im.AddItem(InventoryItem{ID: "1", Name: "Widget", Quantity: 100})
    im.AddItem(InventoryItem{ID: "2", Name: "Gadget", Quantity: 50})

    router.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(im.items)
    })
    
    router.HandleFunc("/inventory/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        item, err := im.GetItem(id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(item)
    })
    
    router.HandleFunc("/inventory/{id}/update", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        err := json.NewDecoder(r.Body).Decode(&InventoryItem{})
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
            return
        }
        // Update logic would go here
        // For now, just return the item
        item, err := im.GetItem(id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(item)
    })
    
    router.HandleFunc("/inventory/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]
        err := im.DeleteItem(id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": "item deleted"})
    })

    log.Println("Starting inventory management system on port 8080")
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}