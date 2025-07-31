// 代码生成时间: 2025-08-01 06:22:00
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "log"

    "github.com/gorilla/mux"
)

// Renamer defines the structure for file renaming
type Renamer struct {
    // Pattern defines the pattern for renaming files
    // e.g., "{index}_{name}"
    Pattern string
}

// NewRenamer creates a new Renamer instance with the given pattern
func NewRenamer(pattern string) *Renamer {
    return &Renamer{Pattern: pattern}
}

// Rename renames a file from its current name to the new name based on the pattern
func (r *Renamer) Rename(basePath, currentName, newName string) error {
    // Construct the full path for the old file
    oldPath := filepath.Join(basePath, currentName)
    // Construct the full path for the new file
    newPath := filepath.Join(basePath, newName)
    // Check if the old file exists
    if _, err := os.Stat(oldPath); os.IsNotExist(err) {
        return fmt.Errorf("file not found: %w", err)
    }
    // Attempt to rename the file
    if err := os.Rename(oldPath, newPath); err != nil {
        return fmt.Errorf("failed to rename file: %w", err)
    }
    return nil
}

// HandleRename is the HTTP handler for renaming files
func HandleRename(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Path      string `json:"path"`
        Pattern   string `json:"pattern"`
        OldNames  []string `json:"oldNames"`
        NewNames  []string `json:"newNames"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    
    renamer := NewRenamer(req.Pattern)
    
    for i, oldName := range req.OldNames {
        newName := req.NewNames[i]
        if err := renamer.Rename(req.Path, oldName, newName); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    
    // Send a success response
    fmt.Fprintf(w, "{"status": "success"}")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc('/rename', HandleRename).Methods("POST")

    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// Note: This code assumes the existence of a JSON decoding function and the necessary
// imports for the HTTP package, which are not included here for brevity.