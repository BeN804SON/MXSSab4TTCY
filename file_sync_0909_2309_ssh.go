// 代码生成时间: 2025-09-09 23:09:48
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    "github.com/gorilla/mux"
)

// FileSyncer 结构体，用于保存源目录和目标目录的路径
type FileSyncer struct {
    SourceDir  string
    DestinationDir string
}

// Sync 同步文件的方法
func (fs *FileSyncer) Sync() error {
    return filepath.Walk(fs.SourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 如果是目录，则跳过
        if info.IsDir() {
            return nil
        }

        // 构建目标路径
        relativePath, err := filepath.Rel(fs.SourceDir, path)
        if err != nil {
            return err
        }
        destPath := filepath.Join(fs.DestinationDir, relativePath)

        // 创建目标目录
        if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
            return err
        }

        // 复制文件
        if _, err := CopyFile(destPath, path); err != nil {
            return err
        }

        return nil
    })
}

// CopyFile 复制文件的方法
func CopyFile(destPath, srcPath string) (int64, error) {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return 0, err
    }
    defer srcFile.Close()

    destFile, err := os.Create(destPath)
    if err != nil {
        return 0, err
    }
    defer destFile.Close()

    return io.Copy(destFile, srcFile)
}

// NewRouter 创建 Gorilla Mux 路由器
func NewRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/sync", SyncHandler).Methods("GET")
    return r
}

// SyncHandler 同步文件的处理函数
func SyncHandler(w http.ResponseWriter, r *http.Request) {
    fs := FileSyncer{
        SourceDir:  "/path/to/source",
        DestinationDir: "/path/to/destination",
    }
    if err := fs.Sync(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintln(w, "Files synced successfully.")
}

func main() {
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8080", router))
}

// 请注意，为了使代码运行，需要将 `/path/to/source` 和 `/path/to/destination` 替换为实际的路径，并导入 `io` 包。
