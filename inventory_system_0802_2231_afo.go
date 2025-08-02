// 代码生成时间: 2025-08-02 22:31:15
package main

import (
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

// InventoryItem represents an item in the inventory
type InventoryItem struct {
    ID           int    `json:"id"`
    Description  string `json:"description"`
    Quantity     int    `json:"quantity"`
}

// Inventory is the collection of inventory items
type Inventory struct {
    items map[int]InventoryItem
}

// NewInventory creates a new inventory
func NewInventory() *Inventory {
    return &Inventory{items: make(map[int]InventoryItem)}
}

// AddItem adds a new item to the inventory
func (inv *Inventory) AddItem(item InventoryItem) error {
    if _, exists := inv.items[item.ID]; exists {
        return fmt.Errorf("item with ID %d already exists", item.ID)
    }
    inv.items[item.ID] = item
    return nil
}

// GetItem retrieves an item from the inventory by ID
func (inv *Inventory) GetItem(id int) (InventoryItem, error) {
    item, exists := inv.items[id]
    if !exists {
        return InventoryItem{}, fmt.Errorf("item with ID %d not found", id)
    }
    return item, nil
}

// UpdateItem updates an existing item in the inventory
func (inv *Inventory) UpdateItem(id int, item InventoryItem) error {
    if _, exists := inv.items[id]; !exists {
        return fmt.Errorf("item with ID %d not found", id)
    }
    inv.items[id] = item
    return nil
}

// DeleteItem removes an item from the inventory by ID
func (inv *Inventory) DeleteItem(id int) error {
    if _, exists := inv.items[id]; !exists {
        return fmt.Errorf("item with ID %d not found", id)
    }
    delete(inv.items, id)
    return nil
}

// InventoryHandler handles HTTP requests for inventory operations
type InventoryHandler struct {
    inventory *Inventory
}

// AddItemHandler adds a new item to the inventory
func (handler *InventoryHandler) AddItemHandler(w http.ResponseWriter, r *http.Request) {
    var item InventoryItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := handler.inventory.AddItem(item); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

// GetItemHandler retrieves an item from the inventory
func (handler *InventoryHandler) GetItemHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    item, err := handler.inventory.GetItem(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(item)
}

// UpdateItemHandler updates an existing item in the inventory
func (handler *InventoryHandler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    var item InventoryItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := handler.inventory.UpdateItem(id, item); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

// DeleteItemHandler removes an item from the inventory
func (handler *InventoryHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := handler.inventory.DeleteItem(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func main() {
    router := mux.NewRouter()
    inventory := NewInventory()
    handler := &InventoryHandler{inventory: inventory}

    router.HandleFunc("/inventory", handler.AddItemHandler).Methods("POST")
    router.HandleFunc("/inventory/{id}", handler.GetItemHandler).Methods("GET")
    router.HandleFunc("/inventory/{id}", handler.UpdateItemHandler).Methods("PUT")
    router.HandleFunc("/inventory/{id}", handler.DeleteItemHandler).Methods("DELETE")

    http.ListenAndServe(":8080", router)
}
