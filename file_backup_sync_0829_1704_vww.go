// 代码生成时间: 2025-08-29 17:04:59
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"
    "strings"

    "github.com/gorilla/mux"
)

// Config stores configuration settings
type Config struct {
    SourceDir  string
    DestinationDir string
    SyncInterval time.Duration
}

// BackupService handles backup and sync operations
type BackupService struct {
    config Config
}

// NewBackupService creates a new BackupService instance
func NewBackupService(config Config) *BackupService {
    return &BackupService{config: config}
}

// StartSync starts syncing files at the specified interval
func (s *BackupService) StartSync() {
    ticker := time.NewTicker(s.config.SyncInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := s.Sync(); err != nil {
                log.Printf("Error syncing files: %v", err)
            }
        }
    }
}

// Sync synchronizes files from source to destination directory
func (s *BackupService) Sync() error {
    // Get all files in source directory
    files, err := os.ReadDir(s.config.SourceDir)
    if err != nil {
        return fmt.Errorf("failed to read source directory: %w", err)
    }

    for _, file := range files {
        srcPath := filepath.Join(s.config.SourceDir, file.Name())
        destPath := filepath.Join(s.config.DestinationDir, file.Name())

        // Check if file already exists in destination
        if _, err := os.Stat(destPath); os.IsNotExist(err) {
            // Copy file to destination
            if err := s.copyFile(srcPath, destPath); err != nil {
                return fmt.Errorf("failed to copy file %s: %w", file.Name(), err)
            }
        } else if err != nil {
            return fmt.Errorf("failed to check file existence: %w", err)
        }
    }
    return nil
}

// copyFile copies a file from source to destination
func (s *BackupService) copyFile(src, dest string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("failed to open source file: %w", err)
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dest)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer destFile.Close()

    if _, err := io.Copy(destFile, sourceFile); err != nil {
        return fmt.Errorf("failed to copy file content: %w", err)
    }
    return nil
}

// HealthCheck returns a health check response
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service is up and running")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/health", HealthCheck).Methods("GET")

    config := Config{
        SourceDir: "/path/to/source",
        DestinationDir: "/path/to/destination",
        SyncInterval: 5 * time.Minute,
    }

    backupService := NewBackupService(config)
    go backupService.StartSync()

    log.Println("Starting file backup and sync service")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
