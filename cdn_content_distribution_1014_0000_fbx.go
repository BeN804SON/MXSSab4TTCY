// 代码生成时间: 2025-10-14 00:00:26
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/gorilla/mux"
)

// FileServer serves files from a directory
type FileServer struct {
    BasePath string
}

// NewFileServer creates a new FileServer instance
func NewFileServer(base string) *FileServer {
    return &FileServer{BasePath: base}
}

// ServeHTTP serves files from the configured base path
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Check if request path is valid
    if !isValidPath(r.URL.Path) {
        http.Error(w, "Invalid path", http.StatusNotFound)
        return
    }

    // Resolve file path from the request path
    filePath := filepath.Join(fs.BasePath, r.URL.Path)
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Check if it's a directory and serve index.html if available, otherwise 403
    if fileInfo.IsDir() {
        indexFilePath := filepath.Join(filePath, "index.html")
        if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
            http.Error(w, "Directory listing is not allowed", http.StatusForbidden)
            return
        }
        filePath = indexFilePath
    }

    // Serve the file
    http.ServeFile(w, r, filePath)
}

// isValidPath checks if the path is valid and avoids directory traversal attacks
func isValidPath(path string) bool {
    // Convert the path to an absolute path
    absolutePath, err := filepath.Abs(path)
    if err != nil {
        return false
    }

    // Check if the absolute path is within the base path
    return filepath.HasPrefix(absolutePath, filepath.Clean(fs.BasePath))
}

func main() {
    // Define the directory to serve
    directory := "./static"
    fs := NewFileServer(directory)

    // Create a new router
    router := mux.NewRouter()

    // Handle file serving
    router.PathPrefix("/").Handler(http.StripPrefix("/", http.HandlerFunc(fs.ServeHTTP)))

    // Define the port and start the server
    port := ":8080"
    log.Printf("Serving files from '%s' on port %s", directory, port)
    if err := http.ListenAndServe(port, router); err != nil {
        log.Fatalf("Unable to start server: %s", err)
    }
}
