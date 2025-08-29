// 代码生成时间: 2025-08-29 23:02:40
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

// FileSyncer 结构体定义文件同步器，包含源和目标路径
type FileSyncer struct {
    Source string
    Target string
}

// NewFileSyncer 创建并返回一个新的FileSyncer实例
func NewFileSyncer(source, target string) *FileSyncer {
    return &FileSyncer{
        Source: source,
        Target: target,
    }
}

// Sync 同步文件
func (fs *FileSyncer) Sync() error {
    // 检查源路径是否存在
    if _, err := os.Stat(fs.Source); os.IsNotExist(err) {
        return fmt.Errorf("source path does not exist: %s", fs.Source)
    }

    // 确保目标路径存在
    if err := os.MkdirAll(fs.Target, 0755); err != nil {
        return fmt.Errorf("failed to create target directory: %s", err)
    }

    // 同步文件
    filepath.Walk(fs.Source, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            // 创建目标目录
            targetPath := filepath.Join(fs.Target, path[len(fs.Source):])
            if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
                return fmt.Errorf("failed to create directory %s: %s", targetPath, err)
            }
        } else {
            // 同步文件
            src, err := os.Open(path)
            if err != nil {
                return fmt.Errorf("failed to open source file %s: %s", path, err)
            }
            defer src.Close()

            dst, err := os.Create(filepath.Join(fs.Target, path[len(fs.Source):]))
            if err != nil {
                return fmt.Errorf("failed to create target file %s: %s", filepath.Join(fs.Target, path[len(fs.Source):]), err)
            }
            defer dst.Close()

            if _, err := io.Copy(dst, src); err != nil {
                return fmt.Errorf("failed to copy file %s: %s", path, err)
            }
        }
        return nil
    })
    return nil
}

// Backup 备份文件
func (fs *FileSyncer) Backup() error {
    // 确保目标路径存在
    if err := os.MkdirAll(fs.Target, 0755); err != nil {
        return fmt.Errorf("failed to create target directory: %s", err)
    }

    // 备份文件
    filepath.Walk(fs.Source, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            // 创建目标目录
            targetPath := filepath.Join(fs.Target, path[len(fs.Source):])
            if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
                return fmt.Errorf("failed to create directory %s: %s", targetPath, err)
            }
        } else {
            // 备份文件
            src, err := os.Open(path)
            if err != nil {
                return fmt.Errorf("failed to open source file %s: %s", path, err)
            }
            defer src.Close()

            backupPath := filepath.Join(fs.Target, path[len(fs.Source):])
            if _, err := os.Stat(backupPath); os.IsNotExist(err) {
                // 创建备份文件
                dst, err := os.Create(backupPath)
                if err != nil {
                    return fmt.Errorf("failed to create backup file %s: %s", backupPath, err)
                }
                defer dst.Close()

                if _, err := io.Copy(dst, src); err != nil {
                    return fmt.Errorf("failed to copy file %s: %s", path, err)
                }
            } else if err != nil {
                return err
            }
        }
        return nil
    })
    return nil
}

func main() {
    // 示例使用FileSyncer
    sourcePath := "/path/to/source"
    targetPath := "/path/to/backup"

    fs := NewFileSyncer(sourcePath, targetPath)

    // 同步文件
    if err := fs.Sync(); err != nil {
        log.Fatalf("failed to sync files: %s", err)
    }
    fmt.Println("Files synced successfully")

    // 备份文件
    if err := fs.Backup(); err != nil {
        log.Fatalf("failed to backup files: %s", err)
    }
    fmt.Println("Files backed up successfully")
}
