// 代码生成时间: 2025-09-13 18:36:43
 * It is designed to be easily maintainable and extensible.
 */

package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "time"
)

// Constants for directory paths
const (
    sourceDir = "./src"
    destDir   = "./dest"
)

// FileSyncer defines the structure for file synchronization tasks
type FileSyncer struct {
    Source string
    Dest   string
}

// NewFileSyncer creates a new FileSyncer instance
func NewFileSyncer(src, dest string) *FileSyncer {
    return &FileSyncer{
        Source: src,
        Dest:   dest,
    }
}

// Sync syncs files from source to destination directory
func (fs *FileSyncer) Sync() error {
    err := filepath.Walk(fs.Source, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil // skip directories
        }

        relPath, err := filepath.Rel(fs.Source, path)
        if err != nil {
            return err
        }

        destPath := filepath.Join(fs.Dest, relPath)
        dir := filepath.Dir(destPath)
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            if err := os.MkdirAll(dir, 0755); err != nil {
                return err
            }
        }

        // Read the file from source directory
        srcFile, err := os.Open(path)
        if err != nil {
            return err
        }
        defer srcFile.Close()

        // Create the file in destination directory
        destFile, err := os.Create(destPath)
        if err != nil {
            return err
        }
        defer destFile.Close()

        // Copy the file content from source to destination
        if _, err := io.Copy(destFile, srcFile); err != nil {
            return err
        }

        // Update the file modification time
        if err := destFile.Chtimes(info.ModTime(), info.ModTime()); err != nil {
            return err
        }

        fmt.Printf("Synced: %s -> %s\
", path, destPath)
        return nil
    })

    if err != nil {
        return err
    }

    return nil
}

func main() {
    fs := NewFileSyncer(sourceDir, destDir)
    err := fs.Sync()
    if err != nil {
        log.Fatalf("Error syncing files: %v", err)
    }

    fmt.Println("File sync completed successfully.")
}